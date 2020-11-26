package command

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"net/http"
	"regexp"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
)

type ScraperConfig struct {
	Command        string // Regular expression which triggers this scraper. Can contain capture groups.
	Title_template string // Title template that will be replaced by regex captures (using %s).
	Title_capture  string // Regex captures for title replacement.
	Url            string // A url to scrape from, can contain one "%s" which is replaced with the first capture group.
	Reply_capture  string // Regular expression used to parse a webpage.
	Help           string // Help message to display
}

// Given a filepath, returns a ScraperConfig.
// A file doesn't exist, an example is made in its place, and an error is returned.
func GetScraperConfigs(filepath string) ([]ScraperConfig, error) {
	var config []ScraperConfig
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return config, makeExampleScraperConfig(filepath)
	}

	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Unable to read file: %s", filepath)
		return config, nil
	}

	json.Unmarshal(bytes, &config)
	return config, nil
}

func makeExampleScraperConfig(filepath string) error {
	config := []ScraperConfig{ScraperConfig{}}
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

// Get a usable scraper.
func GetScraper(config ScraperConfig) (Command, error) {
	webpage_capture := regexp.MustCompile(config.Reply_capture)
	title_capture := regexp.MustCompile(config.Title_capture)

	curry := func(sender service.Conversation, user service.User, msg [][]string, sink func(service.Conversation, service.Message)) {
		scraper(config.Url,
			webpage_capture,
			config.Title_template,
			title_capture,
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
func scraper(url_template string, webpage_capture *regexp.Regexp, title_template string, title_capture *regexp.Regexp, sender service.Conversation, user service.User, msg [][]string, sink func(service.Conversation, service.Message)) {
	substitutions := strings.Count(url_template, "%s")
	url := url_template
	if (substitutions > 0) && (msg == nil || len(msg) == 0 || len(msg[0]) < substitutions) {
		sink(sender, service.Message{Description: "An error when building the url."})
		return
	}

	for _, capture := range msg[0][1:] {
		url = fmt.Sprintf(url, capture)
	}

	response, err := http.Get(url)
	if err == nil {
		// Read response data in to memory
		body, err := ioutil.ReadAll(response.Body)

		if err == nil {
			// Create a regular expression to find comments
			body_s := string(body)

			matches := webpage_capture.FindAllStringSubmatch(body_s, -1)
			title_matches := title_capture.FindAllStringSubmatch(body_s, -1)

			if matches != nil {
				all_captures := make([]string, len(matches))
				for i, captures := range matches {
					all_captures[i] = strings.Join(captures[1:], " ")
				}

				reply := fmt.Sprintf("%s.\n\nRead more at: %s", strings.Join(all_captures, " "), url)
				reply_title := title_template

				for _, captures := range title_matches {
					for _, capture_group := range captures[1:] {
						reply_title = fmt.Sprintf(reply_title, capture_group)
					}
				}

				sink(sender, service.Message{
					Title:       reply_title,
					Description: reply,
					Url:         url,
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
			Url:         url,
		})
	}
}
