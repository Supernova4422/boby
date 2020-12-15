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
	Command        string           // Regular expression which triggers this scraper. Can contain capture groups.
	Title_selector Selector_Capture // Regex captures for title replacement.
	Url            string           // A url to scrape from, can contain one "%s" which is replaced with the first capture group.
	Reply_selector Selector_Capture
	Help           string // Help message to display
}

type Selector_Capture struct {
	Template        string   // Message template to be filled out.
	Selectors       []string // What captures to use to fill out the template
	Handle_multiple string   // How to handle multiple captures. "Random" or "First."
}

// Match all selectors and fill out template. Then using handle_multiple decide which to use.
func Selector_Capture_To_String(doc goquery.Document, selector_capture Selector_Capture) string {

	var max_length int64 = math.MaxInt64
	all_captures := make([](*(goquery.Selection)), len(selector_capture.Selectors))
	for i, selector := range selector_capture.Selectors {
		capture := doc.Find(selector)
		all_captures[i] = capture
		capture_length := int64(capture.Length())

		if capture_length < max_length {
			max_length = capture_length
		}
	}

	reply := selector_capture.Template

	if max_length == 0 {
		return "There was an error retrieving information from the webpage."
	}

	max_length--

	var index int = 0
	if selector_capture.Handle_multiple == "Random" {
		rand.Seed(time.Now().UnixNano())
		index = int(rand.Int63n(max_length))
	} else if selector_capture.Handle_multiple == "Last" {
		index = int(max_length)
	}

	tmp := make([]interface{}, len(selector_capture.Selectors))
	for i, selector := range all_captures {
		selector_i := selector.Slice(int(index), int(index)+1)
		tmp[i] = strings.TrimSpace(selector_i.Text())
	}

	reply = fmt.Sprintf(reply, tmp...)
	return reply
}

// Get a usable scraper.
func GetGoqueryScraper(config GoQueryScraperConfig) (Command, error) {

	curry := func(sender service.Conversation, user service.User, msg [][]string, sink func(service.Conversation, service.Message)) {
		goquery_scraper(
			config,
			sender,
			user,
			msg,
			sink,
		)
	}
	regex, err := regexp.Compile(config.Command)
	if err != nil {
		return Command{}, err
	}
	return Command{
		Pattern: regex,
		Exec:    curry,
		Help:    config.Help,
	}, nil
}

// Return the received message
func goquery_scraper(goQueryScraperConfig GoQueryScraperConfig, sender service.Conversation, user service.User, msg [][]string, sink func(service.Conversation, service.Message)) {
	substitutions := strings.Count(goQueryScraperConfig.Url, "%s")
	url := goQueryScraperConfig.Url
	if (substitutions > 0) && (msg == nil || len(msg) == 0 || len(msg[0]) < substitutions) {
		sink(sender, service.Message{Description: "An error when building the url."})
		return
	}

	for _, capture := range msg[0][1:] {
		url = fmt.Sprintf(url, capture)
	}

	res, err := http.Get(url)
	defer res.Body.Close()
	if err == nil {
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err == nil {
			title := Selector_Capture_To_String(
				*doc,
				goQueryScraperConfig.Title_selector,
			)

			reply := fmt.Sprintf(
				"%s\n\nRead more at: %s",
				Selector_Capture_To_String(*doc, goQueryScraperConfig.Reply_selector),
				url,
			)

			sink(sender, service.Message{
				Title:       title,
				Description: reply,
				Url:         url,
			})
		} else {
			sink(sender, service.Message{Description: "An error occurred when processing the webpage."})
		}
	} else {
		sink(sender, service.Message{
			Description: "An error occurred retrieving the webpage.",
			Url:         url,
		})
	}
}

// Given a filepath, returns a ScraperConfig.
// A file doesn't exist, an example is made in its place, and an error is returned.
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
	config := []GoQueryScraperConfig{GoQueryScraperConfig{}}
	bytes, err := json.Marshal(config)

	if err != nil {
		return errors.New("Unable to create example JSON")
	}

	file, err := os.Create(filepath)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to create file: %s", filepath))
	}
	defer file.Close()

	_, err = file.Write(bytes)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to write to file: %s", filepath))
	}
	return errors.New(fmt.Sprintf("File %s did not exist, an example has been writen.", filepath))
}
