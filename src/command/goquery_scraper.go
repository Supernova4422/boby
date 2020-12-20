package command

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	"math"
	"net/http"
	"regexp"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
	"github.com/PuerkitoBio/goquery"
)

type GoQueryScraperConfig struct {
	Trigger       string          // Word which triggers this command to activate.
	Capture       string          // How to capture words.
	TitleSelector SelectorCapture // Regex captures for title replacement.
	URL           string          // A url to scrape from, can contain one "%s" which is replaced with the first capture group.
	ReplySelector SelectorCapture
	Help          string // Help message to display
}

// SelectorCapture is a method to capture from a webpage.
type SelectorCapture struct {
	Template       string   // Message template to be filled out.
	Selectors      []string // What captures to use to fill out the template
	HandleMultiple string   // How to handle multiple captures. "Random" or "First."
}

// Match all selectors and fill out template. Then using HandleMultiple decide which to use.
func selectorCaptureToString(doc goquery.Document, selectorCapture SelectorCapture) string {
	var maxLength int64 = math.MaxInt64
	allCaptures := make([](*(goquery.Selection)), len(selectorCapture.Selectors))
	for i, selector := range selectorCapture.Selectors {
		capture := doc.Find(selector)
		allCaptures[i] = capture
		captureLength := int64(capture.Length())

		if captureLength < maxLength {
			maxLength = captureLength
		}
	}

	reply := selectorCapture.Template

	if maxLength == 0 {
		return "There was an error retrieving information from the webpage."
	}

	maxLength--

	var index int = 0
	if maxLength > 0 {
		if selectorCapture.HandleMultiple == "Random" {
			rand.Seed(time.Now().UnixNano())
			index = int(rand.Int63n(maxLength))
		} else if selectorCapture.HandleMultiple == "Last" {
			index = int(maxLength)
		}
	}

	tmp := make([]interface{}, len(selectorCapture.Selectors))
	for i, selector := range allCaptures {
		selectorIndex := selector.Slice(int(index), int(index)+1)
		tmp[i] = strings.TrimSpace(selectorIndex.Text())
	}

	reply = fmt.Sprintf(reply, tmp...)
	return reply
}

// GetGoqueryScraper converts a config to a scraper.
func GetGoqueryScraper(config GoQueryScraperConfig) (Command, error) {

	curry := func(sender service.Conversation, user service.User, msg [][]string, sink func(service.Conversation, service.Message)) {
		goqueryScraper(
			config,
			sender,
			user,
			msg,
			sink,
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
func goqueryScraper(goQueryScraperConfig GoQueryScraperConfig, sender service.Conversation, user service.User, msg [][]string, sink func(service.Conversation, service.Message)) {
	substitutions := strings.Count(goQueryScraperConfig.URL, "%s")
	if (substitutions > 0) && (msg == nil || len(msg) == 0 || len(msg[0]) < substitutions) {
		sink(sender, service.Message{Description: "An error when building the url."})
		return
	}

	for _, capture := range msg {
		url := goQueryScraperConfig.URL

		for _, word := range capture[1:] {
			url = fmt.Sprintf(url, word)
		}

		res, err := http.Get(url)
		if err == nil {
			defer res.Body.Close()
			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err == nil {
				if doc.Text() == "" {
					sink(sender, service.Message{
						Title:       "Webpage not found.",
						Description: "Webpage not found at: " + url,
						URL:         url,
					})
				} else {
					title := selectorCaptureToString(
						*doc,
						goQueryScraperConfig.TitleSelector,
					)

					reply := fmt.Sprintf(
						"%s\n\nRead more at: %s",
						selectorCaptureToString(*doc, goQueryScraperConfig.ReplySelector),
						url,
					)

					sink(sender, service.Message{
						Title:       title,
						Description: reply,
						URL:         url,
					})
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
}

// GetGoqueryScraperConfigs retrieves an array of GoQueryScraperConfig from a json file.
// If a file doesn't exist, an example is made in its place, and an error is returned.
func GetGoqueryScraperConfigs(filepath string) ([]GoQueryScraperConfig, error) {
	var config []GoQueryScraperConfig
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return config, makeExampleGoQueryScraperConfig(filepath)
	}

	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Unable to read file: %s", filepath)
		return config, nil
	}

	json.Unmarshal(bytes, &config)
	return config, nil
}

func makeExampleGoQueryScraperConfig(filepath string) error {
	config := []GoQueryScraperConfig{{}}
	bytes, err := json.Marshal(config)

	if err != nil {
		return errors.New("Unable to create example JSON")
	}

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("Unable to create file: %s", filepath)
	}
	defer file.Close()

	_, err = file.Write(bytes)
	if err != nil {
		return fmt.Errorf("Unable to write to file: %s", filepath)
	}
	return fmt.Errorf("File %s did not exist, an example has been writen", filepath)
}
