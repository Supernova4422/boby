package command

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"strings"

	"math"
	"net/url"

	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/storage"
	"github.com/BKrajancic/boby/m/v2/src/utils"
	"github.com/PuerkitoBio/goquery"
)

// GoQueryScraperConfig can be turned into a scraper that uses GoQuery.
type GoQueryScraperConfig struct {
	Title         string          // When sending a post, what should the title be.
	Trigger       string          // Word which triggers this command to activate.
	Parameters    []Parameter     // How to capture words.
	TitleSelector SelectorCapture // The output message's title.
	ErrorURL      string          // A url to show only when there is an error.
	URL           string          // A url to scrape from, can contain one "%s" which is replaced with the first capture group.
	URLSuffix     string          // When adding a URL to a message, this string is appended. This is useful for including referral links.
	ReplySelector SelectorCapture // The output message's body text.
	Fields        []GoQueryFieldCapture
	Help          string // Help message to display.
	HelpInput     string // Help message to display for input following command.
	HideURL       bool   // When true, a result returns no URL. Use with caution, attribution is often required.
}

// GoQueryFieldCapture is used to have a selector capture for a pair of selectors.
type GoQueryFieldCapture struct {
	Title       SelectorCapture
	Description SelectorCapture
}

// SelectorCapture will fill out a template string using webpage content selected with goquery.
type SelectorCapture struct {
	Template        string              // Message template to be filled out. Every %s in a template is replaced with results of selectors.
	Selectors       []string            // What goquery captures are used to fill out the template.
	Replacements    []map[string]string // String replacements for each entry in selectors.
	FullReplacement map[string]string   // String replacement that takes place on the completed selector.
	HandleMultiple  string              // How to handle multiple captures. "Random" or "First."
}

// A HTMLGetter returns a url and buffer based on a string.
type HTMLGetter = func(string) (url string, out io.ReadCloser, err error)

// selectorCaptureToString matches all selectors and fill out template.
// Then using HandleMultiple decide which to use.
func (s SelectorCapture) selectorCaptureToString(doc goquery.Document) (string, error) {
	if len(s.Selectors) == 0 || !strings.Contains(s.Template, "%s") {
		return s.Template, nil
	}

	var maxLength int64 = math.MaxInt64
	allCaptures := make([](*(goquery.Selection)), len(s.Selectors))
	for i, selector := range s.Selectors {
		capture := doc.Find(selector)
		allCaptures[i] = capture
		captureLength := int64(capture.Length())

		if captureLength < maxLength {
			maxLength = captureLength
		}
	}

	// if maxLength == 0 {
	// 	return "", fmt.Errorf("There was an error retrieving information from the webpage.")
	// }
	maxLength--

	var index int = 0
	if maxLength > 0 {
		if s.HandleMultiple == "Random" {
			index = int(rand.Int63n(maxLength))
		} else if s.HandleMultiple == "Last" {
			index = int(maxLength)
		}
	}

	tmp := make([]interface{}, len(s.Selectors))
	for i, selector := range allCaptures {
		val := ""
		if index < (*selector).Length() {
			val = strings.TrimSpace(selector.Slice(int(index), int(index)+1).Text())
			if i < len(s.Replacements) {
				for search, replace := range s.Replacements[i] {
					if strings.Contains(val, search) {
						val = strings.ReplaceAll(val, search, replace)
						break
					}
				}
			}
		}
		tmp[i] = val
	}

	reply := fmt.Sprintf(s.Template, tmp...)
	for search, replace := range s.FullReplacement {
		if strings.Contains(reply, search) {
			reply = strings.ReplaceAll(reply, search, replace)
		}
	}

	return reply, nil
}

// Command returns a webscraper Command from a config.
func (g GoQueryScraperConfig) Command() (Command, error) {
	return g.CommandWithHTMLGetter(utils.HTMLGetWithHTTP)
}

// CommandWithHTMLGetter makes a scraper Command from a config, retrieving HTML pages using HTMLGetter.
func (g GoQueryScraperConfig) CommandWithHTMLGetter(htmlGetter HTMLGetter) (Command, error) {
	curry := func(sender service.Conversation, user service.User, msg []interface{}, storage *storage.Storage, sink func(service.Conversation, service.Message) error) error {
		return g.onMessage(
			sender,
			user,
			msg,
			storage,
			sink,
			htmlGetter,
		)
	}

	return Command{
		Trigger:    g.Trigger,
		Parameters: g.Parameters,
		Exec:       curry,
		Help:       g.Help,
		HelpInput:  g.HelpInput,
	}, nil
}

// onMessage processes the request, and sends out messages.
func (g GoQueryScraperConfig) onMessage(sender service.Conversation, user service.User, msg []interface{}, storage *storage.Storage, sink func(service.Conversation, service.Message) error, htmlGetter HTMLGetter) error {
	substitutions := strings.Count(g.URL, "%s")
	if (substitutions > 0) && (len(msg) == 0 || len(msg) < substitutions) {
		return sink(
			sender,
			service.Message{
				Description: "An error occurred when building the url.",
			})
	}

	fields := make([]service.MessageField, 0)
	msgURL := g.URL

	for _, word := range msg {
		msgURL = fmt.Sprintf(msgURL, url.PathEscape(word.(string)))
	}

	redirect, htmlReader, err := htmlGetter(msgURL)
	if err == nil {
		defer htmlReader.Close()
	} else {
		return sink(
			sender,
			service.Message{
				Title:       "Error",
				Description: "An error occurred retrieving the webpage.",
				URL:         msgURL,
			},
		)
	}

	doc, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		return sink(
			sender,
			service.Message{
				Title:       msgURL,
				Description: "An error occurred when processing the webpage.",
				URL:         g.ErrorURL,
			},
		)
	}

	if doc.Text() == "" {
		captures := []string{}
		for _, item := range msg {
			captures = append(captures, item.(string))
		}

		return sink(
			sender,
			service.Message{
				Title:       "Error",
				Description: fmt.Sprintf("No result was found for \"%s\"", strings.Join(captures, " ")),
				URL:         g.ErrorURL,
			},
		)
	}

	title, err1 := g.TitleSelector.selectorCaptureToString(*doc)
	value, err2 := g.ReplySelector.selectorCaptureToString(*doc)
	if err1 == nil && err2 == nil && title != "" && value != "" {
		if g.HideURL {
			redirect = ""
		}

		fields = append(fields, service.MessageField{
			Field: title,
			Value: value,
			URL:   redirect + g.URLSuffix,
		})
	}

	for _, field := range g.Fields {
		fieldTitle, err1 := field.Title.selectorCaptureToString(*doc)
		value, err2 := field.Description.selectorCaptureToString(*doc)
		if err1 == nil && err2 == nil && fieldTitle != "" && value != "" {
			fields = append(fields,
				service.MessageField{
					Field:  fieldTitle,
					Value:  value,
					Inline: true,
				},
			)
		}
	}

	if len(fields) == 0 {
		fields = append(fields, service.MessageField{
			Field: "Error",
			Value: "No result was found",
			URL:   g.ErrorURL,
		})
	}

	replyMsg := service.Message{
		Title:       fields[0].Field,
		Description: fields[0].Value,
		URL:         fields[0].URL,
	}

	if len(fields) > 1 {
		replyMsg.Fields = fields[1:]
	}

	return sink(sender, replyMsg)
}

// GetGoqueryScraperConfigs retrieves an array of GoQueryScraperConfig by parsing JSON from a buffer.
// If a file doesn't exist, an example is made in its place, and an error is returned.
func GetGoqueryScraperConfigs(reader io.Reader) ([]GoQueryScraperConfig, error) {
	bytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var config []GoQueryScraperConfig
	return config, json.Unmarshal(bytes, &config)
}
