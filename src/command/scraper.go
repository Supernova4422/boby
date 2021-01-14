package command

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"regexp"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/storage"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/utils"
)

// ScraperConfig is a struct that can be turned into a usable scraper.
type ScraperConfig struct {
	Trigger       string
	Capture       string
	TitleTemplate string // Title template that will be replaced by regex captures (using %s).
	TitleCapture  string // Regex captures for title replacement.
	URL           string // A url to scrape from, can contain one "%s" which is replaced with the first capture group.
	ReplyCapture  string // Regular expression used to parse a webpage.
	Help          string // Help message to display
}

// GetScraperConfigs returns a set of ScraperConfig by reading a file.
// If a file doesn't exist at the given filepath, an example is made in its place,
// and an error is returned.
func GetScraperConfigs(reader io.Reader) ([]ScraperConfig, error) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var config []ScraperConfig
	return config, json.Unmarshal(bytes, &config)
}

// GetScraper returns a webscraper command from a config, using HTTP to get a html.
func (config ScraperConfig) GetScraper() (Command, error) {
	return config.GetScraperWithHTMLGetter(utils.HTMLGetWithHTTP)
}

// GetScraperWithHTMLGetter makes a scraper from a config.
func (config ScraperConfig) GetScraperWithHTMLGetter(htmlGetter HTMLGetter) (Command, error) {
	webpageCapture := regexp.MustCompile(config.ReplyCapture)
	titleCapture := regexp.MustCompile(config.TitleCapture)

	curry := func(sender service.Conversation, user service.User, msg [][]string, storage *storage.Storage, sink func(service.Conversation, service.Message)) {
		scraper(config.URL,
			webpageCapture,
			config.TitleTemplate,
			titleCapture,
			sender,
			user,
			msg,
			storage,
			sink,
			htmlGetter,
		)
	}
	regex, err := regexp.Compile(config.Capture)
	if err != nil {
		return Command{}, err
	}
	return Command{
		Trigger: config.Trigger,
		Pattern: regex,
		Exec:    curry,
		Help:    config.Help,
	}, nil
}

// Return the received message
func scraper(urlTemplate string, webpageCapture *regexp.Regexp, titleTemplate string, titleCapture *regexp.Regexp, sender service.Conversation, user service.User, msg [][]string, storage *storage.Storage, sink func(service.Conversation, service.Message), htmlGetter HTMLGetter) {
	substitutions := strings.Count(urlTemplate, "%s")
	url := urlTemplate
	if substitutions > 0 {
		if msg == nil || len(msg) == 0 || len(msg[0]) < substitutions {
			sink(sender, service.Message{Description: "An error when building the url."})
			return
		}
		for _, capture := range msg[0][1:] {
			url = fmt.Sprintf(url, capture)
		}
	}

	htmlReader, err := htmlGetter(url)
	if err == nil {
		defer htmlReader.Close()
		body, err := ioutil.ReadAll(htmlReader)
		if err == nil {
			// Create a regular expression to find comments
			bodyS := string(body)

			matches := webpageCapture.FindAllStringSubmatch(bodyS, -1)
			titleMatches := titleCapture.FindAllStringSubmatch(bodyS, -1)

			if matches != nil {
				allCaptures := make([]string, len(matches))
				for i, captures := range matches {
					allCaptures[i] = strings.Join(captures[1:], " ")
				}

				reply := fmt.Sprintf("%s.\n\nRead more at: %s", strings.Join(allCaptures, " "), url)
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

				sink(sender, service.Message{
					Title:       replyTitle,
					Description: reply,
					URL:         url,
				})
			} else {
				sink(sender, service.Message{Description: "Could not extract data from the webpage."})
			}
		} else {
			sink(sender, service.Message{Description: "An error occurred when processing the webpage."})
		}
	} else {
		sink(sender, service.Message{
			Description: "An error occurred retrieving the webpage.",
			URL:         url,
		})
	}
}
