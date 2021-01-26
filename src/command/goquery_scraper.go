package command

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"

	"math"
	"net/url"
	"regexp"

	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/storage"
	"github.com/BKrajancic/boby/m/v2/src/utils"
	"github.com/PuerkitoBio/goquery"
)

// GoQueryScraperConfig can be turned into a scraper that uses GoQuery.
type GoQueryScraperConfig struct {
	Title         string          // When sending a post, what should the title be.
	Trigger       string          // Word which triggers this command to activate.
	Capture       string          // How to capture words.
	TitleSelector SelectorCapture // The output message's title.
	URL           string          // A url to scrape from, can contain one "%s" which is replaced with the first capture group.
	ReplySelector SelectorCapture // The output message's body text.
	Help          string          // Help message to display.
	HelpInput     string          // Help message to display for input following command.
}

// SelectorCapture will fill out a template string using webpage content selected with goquery.
type SelectorCapture struct {
	Template       string              // Message template to be filled out. Every %s in a template is replaced with results of selectors.
	Selectors      []string            // What goquery captures are used to fill out the template.
	Replacements   []map[string]string // String replacements for each entry in selectors.
	HandleMultiple string              // How to handle multiple captures. "Random" or "First."
}

// A HTMLGetter returns a url and buffer based on a string.
type HTMLGetter = func(string) (url string, out io.ReadCloser, err error)

// selectorCaptureToString matches all selectors and fill out template.
// Then using HandleMultiple decide which to use.
func (s SelectorCapture) selectorCaptureToString(doc goquery.Document) string {
	var maxLength int64 = math.MaxInt64
	allCaptures := make([](*(goquery.Selection)), len(s.Selectors))
	if len(s.Selectors) == 0 {
		maxLength = 0
	} else {
		for i, selector := range s.Selectors {
			capture := doc.Find(selector)
			allCaptures[i] = capture
			captureLength := int64(capture.Length())

			if captureLength < maxLength {
				maxLength = captureLength
			}
		}
	}

	reply := s.Template

	if maxLength == 0 && strings.Contains(reply, "%s") {
		return "There was an error retrieving information from the webpage."
	} else if maxLength > 0 {
		maxLength--

		var index int = 0
		if maxLength > 0 {
			if s.HandleMultiple == "Random" {
				rand.Seed(time.Now().UnixNano())
				index = int(rand.Int63n(maxLength))
			} else if s.HandleMultiple == "Last" {
				index = int(maxLength)
			}
		}

		tmp := make([]interface{}, len(s.Selectors))
		for i, selector := range allCaptures {
			selectorIndex := selector.Slice(int(index), int(index)+1)
			val := strings.TrimSpace(selectorIndex.Text())
			if i < len(s.Replacements) {
				for search, replace := range s.Replacements[i] {
					if strings.Contains(val, search) {
						val = strings.ReplaceAll(val, search, replace)
						break
					}
				}
			}
			tmp[i] = val
		}

		reply = fmt.Sprintf(reply, tmp...)
	}
	return reply
}

// Command returns a webscraper Command from a config.
func (g GoQueryScraperConfig) Command() (Command, error) {
	return g.CommandWithHTMLGetter(utils.HTMLGetWithHTTP)
}

// CommandWithHTMLGetter makes a scraper Command from a config, retrieving HTML pages using HTMLGetter.
func (g GoQueryScraperConfig) CommandWithHTMLGetter(htmlGetter HTMLGetter) (Command, error) {
	curry := func(sender service.Conversation, user service.User, msg [][]string, storage *storage.Storage, sink func(service.Conversation, service.Message)) {
		g.onMessage(
			sender,
			user,
			msg,
			storage,
			sink,
			htmlGetter,
		)
	}

	regex, err := regexp.Compile(g.Capture)
	if err != nil {
		return Command{}, err
	}

	return Command{
		Trigger:   g.Trigger,
		Pattern:   regex,
		Exec:      curry,
		Help:      g.Help,
		HelpInput: g.HelpInput,
	}, nil
}

// onMessage processes the request, and sends out messages.
func (g GoQueryScraperConfig) onMessage(sender service.Conversation, user service.User, msg [][]string, storage *storage.Storage, sink func(service.Conversation, service.Message), htmlGetter HTMLGetter) {
	substitutions := strings.Count(g.URL, "%s")
	if (substitutions > 0) && (msg == nil || len(msg) == 0 || len(msg[0]) < substitutions) {
		sink(sender, service.Message{Description: "An error occurred when building the url."})
		return
	}

	fields := make([]service.MessageField, 0)
	for _, capture := range msg {
		msgURL := g.URL

		for _, word := range capture {
			msgURL = fmt.Sprintf(msgURL, url.PathEscape(word))
		}

		redirect, htmlReader, err := htmlGetter(msgURL)
		if err == nil {
			defer htmlReader.Close()
			doc, err := goquery.NewDocumentFromReader(htmlReader)
			if err == nil {
				if doc.Text() == "" {
					fields = append(fields, service.MessageField{
						Field: msgURL,
						Value: fmt.Sprintf("Webpage not found at %s", redirect),
						URL:   "",
					})
				} else {
					fields = append(fields, service.MessageField{
						Field: g.TitleSelector.selectorCaptureToString(*doc),
						Value: g.ReplySelector.selectorCaptureToString(*doc),
						URL:   redirect,
					})
				}
			} else {
				fields = append(fields, service.MessageField{
					Field: msgURL,
					Value: "An error occurred when processing the webpage.",
				})
			}
		} else {
			fields = append(fields, service.MessageField{
				Field: "Error",
				Value: "An error occurred retrieving the webpage.",
				URL:   msgURL,
			})
		}
	}

	sink(sender, service.Message{
		Title:       g.Title,
		Description: "",
		Fields:      fields,
	})
}

// GetGoqueryScraperConfigs retrieves an array of GoQueryScraperConfig by parsing JSON from a buffer.
// If a file doesn't exist, an example is made in its place, and an error is returned.
func GetGoqueryScraperConfigs(reader io.Reader) ([]GoQueryScraperConfig, error) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var config []GoQueryScraperConfig
	return config, json.Unmarshal(bytes, &config)
}
