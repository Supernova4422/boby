package command

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/storage"
)

// JSONGetterConfig can be used to extract from JSON into a message.
type JSONGetterConfig struct {
	Trigger     string
	Capture     string
	Title       FieldCapture
	Captures    []JSONCapture
	URL         string
	Help        string
	HelpInput   string
	Description string
	Grouped     bool
	Delay       int
	Token       TokenMaker
	RateLimit   RateLimitConfig
}

// A TokenMaker is useful for creating a token that may be part of an API request.
type TokenMaker struct {
	Type    string // Can be MD5
	Prefix  string // When calculating a token, what should be prepended
	Postfix string // When calculating a token, what should be appended
	Size    int    // Take the first 'Size' characters from the result.
}

// MakeToken will make a token from this token maker.
// If t.type is "MD5", prefixes 'prefix' to the input, postfixes 'postfix' to the input,
// takes the 'size' number of characters from the MD5 sum.
func (t TokenMaker) MakeToken(input string) (out string) {
	if t.Type == "MD5" {
		fullString := []byte(t.Prefix + input + t.Postfix)
		return fmt.Sprintf("%x", md5.Sum(fullString))[:t.Size]
	}
	return ""
}

// A FieldCapture represents a template to be filled out by selectors.
// TODO Improve documentation here.
type FieldCapture struct {
	Template  string   // Message template to be filled out.
	Selectors []string // What captures to use to fill out the template
}

// ToStringWithMap uses a map to fill out the template.
// HandleMultiple is completely ignored.
// If a key is missing, it is skipped.
func (f FieldCapture) ToStringWithMap(dict map[string]interface{}) (out string, err error) {
	out = f.Template
	for _, selector := range f.Selectors {
		val, ok := dict[selector]
		if ok {
			out = fmt.Sprintf(out, val.(string))
		}
	}
	return out, err
}

// JSONCapture is a pair of FieldCapture to represent a title, body pair in a message.
type JSONCapture struct {
	Title FieldCapture
	Body  FieldCapture
}

// JSONGetter will accept a string and provide a reader. This could be a file, a webpage, who cares!
type JSONGetter = func(string) (out io.ReadCloser, err error)

// Command uses the config to make a Command that processes messages.
func (j JSONGetterConfig) Command(jsonGetter JSONGetter) (Command, error) {
	curry := func(sender service.Conversation, user service.User, msg [][]string, storage *storage.Storage, sink func(service.Conversation, service.Message)) {
		j.jsonGetterFunc(
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

// jsonGetterFunc processes a message.
func (j JSONGetterConfig) jsonGetterFunc(sender service.Conversation, user service.User, msg [][]string, storage *storage.Storage, sink func(service.Conversation, service.Message), jsonGetter JSONGetter) {
	substitutions := strings.Count(j.URL, "%s")
	if (substitutions > 0) && (msg == nil || len(msg) == 0 || len(msg[0]) < substitutions) {
		sink(sender, service.Message{Description: "An error occurred when building the url."})
		return
	}

	fields := make([]service.MessageField, 0)
	for _, capture := range msg {
		msgURL := j.URL

		replacements := strings.Count(msgURL, "%s")

		for i, word := range capture {
			if i == replacements {
				break
			}
			msgURL = strings.Replace(msgURL, "%s", url.PathEscape(word), 1)
		}

		msgURL += j.Token.MakeToken(strings.Join(capture, ""))
		fmt.Print(msgURL)
		jsonReader, err := jsonGetter(msgURL)
		if err == nil {
			defer jsonReader.Close()
			buf, err := ioutil.ReadAll(jsonReader)
			if err == nil {
				dict := make(map[string]interface{})
				err := json.Unmarshal(buf, &dict)
				if err == nil {
					for _, capture := range j.Captures {
						body, err := capture.Body.ToStringWithMap(dict)
						if err == nil && strings.Contains(body, "%s") == false {
							title, err := capture.Title.ToStringWithMap(dict)
							if err == nil {
								fields = append(fields, service.MessageField{Field: title, Value: body})
							}
						}
					}
					if j.Grouped {
						title, err := j.Title.ToStringWithMap(dict)
						if err == nil {
							sink(sender, service.Message{
								Title:       title,
								Fields:      fields,
								Description: j.Description,
							})
						}
					} else {
						for _, field := range fields {
							sink(sender, service.Message{
								Title:       field.Field,
								Description: field.Value,
							})

							time.Sleep(time.Duration(j.Delay) * time.Second)
						}
					}
				}
			}
		}
	}
}
