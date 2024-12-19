package command

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"

	"regexp"

	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/storage"
	"github.com/BKrajancic/boby/m/v2/src/utils"
)

// RegexpScraperConfig is a struct that can be made into a command.
// That Command will process a HTML based on a regexp, to send responses.
type RegexpScraperConfig struct {
	Trigger       string
	Parameters    []Parameter
	TitleTemplate string // Title template that will be replaced by regex captures (using %s).
	TitleCapture  string // Regex captures for title replacement.
	URL           string // A url to scrape from, can contain one "%s" which is replaced with the first capture group.
	ReplyCapture  string // Regular expression used to parse a webpage.
	Help          string // Help message to display
	HelpInput     string // Help message to display for input following command
}

// GetRegexpScraperConfigs returns a set of RegexScraperConfig by reading a file.
// If a file doesn't exist at the given filepath, an example is made in its place,
// and an error is returned.
func GetRegexpScraperConfigs(reader io.Reader) ([]RegexpScraperConfig, error) {
	bytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var config []RegexpScraperConfig
	return config, json.Unmarshal(bytes, &config)
}

// Command returns a webscraper command from a config, using HTTP to get a html.
func (r RegexpScraperConfig) Command() (Command, error) {
	return r.CommandWithHTMLGetter(utils.HTMLGetWithHTTP)
}

// CommandWithHTMLGetter makes a scraper from a config.
func (r RegexpScraperConfig) CommandWithHTMLGetter(htmlGetter HTMLGetter) (Command, error) {
	webpageCapture := regexp.MustCompile(r.ReplyCapture)
	titleCapture := regexp.MustCompile(r.TitleCapture)

	curry := func(sender service.Conversation, user service.User, msg []interface{}, storage *storage.Storage, sink func(service.Conversation, service.Message) error) error {
		return scraper(r.URL,
			webpageCapture,
			r.TitleTemplate,
			titleCapture,
			sender,
			user,
			msg,
			storage,
			sink,
			htmlGetter,
		)
	}

	return Command{
		Trigger:    r.Trigger,
		Parameters: r.Parameters,
		Exec:       curry,
		Help:       r.Help,
		HelpInput:  r.HelpInput,
	}, nil
}

// scraper returns the received message
func scraper(urlTemplate string, webpageCapture *regexp.Regexp, titleTemplate string, titleCapture *regexp.Regexp, sender service.Conversation, user service.User, msg []interface{}, storage *storage.Storage, sink func(service.Conversation, service.Message) error, htmlGetter HTMLGetter) error {
	substitutions := strings.Count(urlTemplate, "%s")
	urlPage := urlTemplate
	if substitutions > 0 {
		if len(msg) == 0 || len(msg) < substitutions {
			return sink(sender, service.Message{Description: "An error when building the url."})
		}
		for _, capture := range msg {
			urlPage = fmt.Sprintf(urlPage, url.PathEscape(capture.(string)))
		}
	}

	_, htmlReader, err := htmlGetter(urlPage)
	if err != nil {
		return sink(sender, service.Message{
			Description: "An error occurred retrieving the webpage.",
			URL:         urlPage,
		})
	}

	defer htmlReader.Close()
	body, err := io.ReadAll(htmlReader)
	if err != nil {
		return sink(sender, service.Message{Description: "An error occurred when processing the webpage."})
	}

	// Create a regular expression to find comments
	bodyS := string(body)
	matches := webpageCapture.FindAllStringSubmatch(bodyS, -1)
	titleMatches := titleCapture.FindAllStringSubmatch(bodyS, -1)
	if matches == nil {
		return sink(sender, service.Message{Description: "Could not extract data from the webpage."})
	}
	allCaptures := make([]string, len(matches))
	for i, captures := range matches {
		allCaptures[i] = strings.Join(captures[1:], " ")
	}

	reply := strings.Join(allCaptures, "\n")
	replyTitle := titleTemplate

	if strings.Contains(replyTitle, "%s") {
		titleCaptures := ""
		for _, captures := range titleMatches {
			for _, captureGroup := range captures[1:] {
				if titleCaptures != "" {
					titleCaptures += " "
				}
				titleCaptures += captureGroup
			}
		}
		replyTitle = fmt.Sprintf(replyTitle, titleCaptures)
	}

	return sink(sender, service.Message{
		Title:       replyTitle,
		Description: reply,
		URL:         urlPage,
	})
}
