package command

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/storage"
)

type JSONGetterConfig struct {
	Trigger     string
	Capture     string
	Title       SelectorCapture
	Captures    []SelectorCapture
	URL         string
	Help        string
	Description string
}

type JSONGetter = func(string) (out io.ReadCloser, err error)

// GetWebScraper returns a webscraper command from a config.
func (j JSONGetterConfig) GetWebScraper(jsonGetter JSONGetter) (Command, error) {
	curry := func(sender service.Conversation, user service.User, msg [][]string, storage *storage.Storage, sink func(service.Conversation, service.Message)) {
		jsonGetterFunc(
			j,
			sender,
			user,
			msg,
			storage,
			sink,
			jsonGetter,
		)
	}

	regex, err := regexp.Compile(j.Capture)
	if err != nil {
		return Command{}, err
	}

	return Command{
		Trigger: j.Trigger,
		Pattern: regex,
		Exec:    curry,
		Help:    j.Help,
	}, nil
}

// Return the received message
func jsonGetterFunc(config JSONGetterConfig, sender service.Conversation, user service.User, msg [][]string, storage *storage.Storage, sink func(service.Conversation, service.Message), jsonGetter JSONGetter) {
	substitutions := strings.Count(config.URL, "%s")
	if (substitutions > 0) && (msg == nil || len(msg) == 0 || len(msg[0]) < substitutions) {
		sink(sender, service.Message{Description: "An error occurred when building the url."})
		return
	}

	fields := make([]service.MessageField, 0)
	for _, capture := range msg {
		msgURL := config.URL

		for _, word := range capture[1:] {
			msgURL = fmt.Sprintf(msgURL, url.PathEscape(word))
		}

		jsonReader, err := jsonGetter(msgURL)
		if err == nil {
			defer jsonReader.Close()
			buf := new(strings.Builder)
			_, err := io.Copy(buf, jsonReader)
			if err == nil {
				dict := make(map[string]string)
				err := json.Unmarshal([]byte(buf.String()), &dict)
				if err == nil {
					for _, capture := range config.Captures {
						val, err := capture.ToStringWithMap(dict)
						if err == nil {
							fields = append(fields,
								service.MessageField{
									Value: val,
								})
						}
					}

					title, err := config.Title.ToStringWithMap(dict)
					if err == nil {
						sink(sender, service.Message{
							Title:       title,
							Fields:      fields,
							Description: config.Description,
						})
					}
				}
			}
		}
	}

}
