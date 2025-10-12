package command

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/storage"
)

// JSONGetterConfig can be used to extract from JSON into a message.
type JSONGetterConfig struct {
	Trigger     string          // What a message must begin with to trigger this command.
	Parameters  []Parameter     // Capture is a regexp, that is used to capture everything following 'trigger.'
	Message     JSONCapture     // The primary title and body of a message.
	Fields      []JSONCapture   // A message is composed of several fields. Captures is used to make fields of a message.
	Grouped     bool            // If true, only a single message is sent, if false each entry in .
	URL         string          // URL to retrieve a JSON from.
	URLSelector string          // Selector to make into the URL for the title.
	Help        string          // Message shown when help command is used.
	HelpInput   string          // Message shown used to explain what expected user input is following trigger.
	Delay       int             // If grouped is false, what is the delay between each message sent.
	Token       TokenMaker      // Often an API requires a calculated API, Token is used to help create a token and append to a URL prior to requests.
	RateLimit   RateLimitConfig // RateLimit places a limit on how frequently a user can send messages.
}

// MessagesFromJSON accepts a dict (which usually represents a JSON) and returns a sequence of messages based on the configuration.
func (j JSONGetterConfig) MessagesFromJSON(dict map[string]interface{}) (messages []service.Message) {
	messages = make([]service.Message, 0)

	fields := make([]service.MessageField, 0)
	for _, field := range j.Fields {
		if field, err := field.MessageField(dict); err == nil {
			if field.Field != "" && field.Value != "" {
				fields = append(fields, field)
			}
		}
	}

	if j.Grouped {
		if body, err := j.Message.MessageField(dict); err == nil {
			messages = append(messages, service.Message{
				Title:       body.Field,
				Description: body.Value,
				URL:         body.URL,
				Fields:      fields,
			})
		}
	} else {
		if body, err := j.Message.MessageField(dict); err == nil {
			if body.Field != "" || body.Value != "" {
				messages = append(messages, service.Message{
					Title:       body.Field,
					Description: body.Value,
					URL:         body.URL,
				})
			}
		}

		for _, field := range fields {
			messages = append(messages, service.Message{
				Title:       field.Field,
				Description: field.Value,
				URL:         field.URL,
			})
		}
	}
	return
}

// JSONCapture is a pair of FieldCapture to represent a title, body pair in a message.
type JSONCapture struct {
	Title 	    FieldCapture
	Body  	    FieldCapture
	URLSelector string
}

// MessageField uses a dict (which is a usually a reading of a JSON file), to create a MessageField.
// Returns an error if %s is present in either the title or description even after replacements are made.
func (j JSONCapture) MessageField(dict map[string]interface{}) (field service.MessageField, err error) {
	body, err := j.Body.ToStringWithMap(dict)
	if err != nil {
		body = j.Body.ErrorMsg
	}

	title, err := j.Title.ToStringWithMap(dict)
	if err != nil {
		title = j.Title.ErrorMsg
	}

	url := ""

	if val, ok := dict[j.URLSelector]; ok && val != "" {
		url = val.(string)
	} 
	
	// TODO: Temporary work around for interpreting "'", a better solution is needed.
	field = service.MessageField{
		Field: strings.ReplaceAll(title, "&#39;", "'"),
		Value: strings.ReplaceAll(body, "&#39;", "'"),
		URL: url,
		Inline: true,
	}

	return field, nil
}

// A FieldCapture represents a template to be filled out by selectors.
type FieldCapture struct {
	Template  string   // Message template to be filled out. Use %s to denote text to be replaced.
	Selectors []string // What captures to use to fill out the template
	ErrorMsg  string   // If the template has any %s remaining, replace the entire msg with this msg.
}

// ToStringWithMap uses a map to fill out the template.
// If a key is missing, it is skipped.
func (f FieldCapture) ToStringWithMap(dict map[string]interface{}) (out string, err error) {
	split := strings.SplitAfter(f.Template, "%s")
	expectReplacments := len(split) - 1

	replacements := 0
	for _, selector := range f.Selectors {
		if val, ok := dict[selector]; ok && val != "" {
			out += fmt.Sprintf(split[replacements], val.(string))
			replacements++
		} else {
			break
		}

		if replacements == expectReplacments {
			break
		}
	}

	out += split[len(split)-1]

	if replacements < expectReplacments {
		return "", fmt.Errorf("there was an error processing results")
	}

	return out, nil
}

// JSONGetter will accept a string and provide a reader. This could be a file, a webpage, who cares!
type JSONGetter = func(string) (out io.ReadCloser, err error)

// Command uses the config to make a Command that processes messages.
func (j JSONGetterConfig) Command(jsonGetter JSONGetter) (Command, error) {
	curry := func(sender service.Conversation, user service.User, msg []interface{}, storage *storage.Storage, sink func(service.Conversation, service.Message) error) error {
		return j.jsonGetterFunc(
			sender,
			user,
			msg,
			storage,
			sink,
			jsonGetter,
		)
	}

	return Command{
		Trigger:    j.Trigger,
		Parameters: j.Parameters,
		Exec:       curry,
		Help:       j.Help,
		HelpInput:  j.HelpInput,
	}, nil
}

// jsonGetterFunc processes a message.
func (j JSONGetterConfig) jsonGetterFunc(sender service.Conversation, user service.User, msg []interface{}, storage *storage.Storage, sink func(service.Conversation, service.Message) error, jsonGetter JSONGetter) error {
	substitutions := strings.Count(j.URL, "%s")
	noCapture := len(msg) == 0

	if (substitutions > 0) && (noCapture || len(msg) < substitutions) {
		return sink(sender, service.Message{Description: "An error occurred when building the url."})
	}

	if noCapture {
		msg = []interface{}{""}
	}

	msgURL := j.URL

	replacements := strings.Count(msgURL, "%s")

	// HACK: noCapture is pretty hacky.
	if !noCapture {
		for i, word := range msg {
			if i == replacements {
				break
			}
			msgURL = strings.Replace(msgURL, "%s", url.PathEscape(word.(string)), 1)
		}

		output := []string{}
		for _, item := range msg {
			output = append(output, item.(string))
		}

		msgURL += j.Token.MakeToken(strings.Join(output, ""))
	}

	if jsonReader, err := jsonGetter(msgURL); err == nil {
		defer jsonReader.Close()
		if buf, err := io.ReadAll(jsonReader); err == nil {
			dict := make(map[string]interface{})
			if err := json.Unmarshal(buf, &dict); err == nil {
				for _, msg := range j.MessagesFromJSON(dict) {
					err := sink(sender, msg)
					if err != nil {
						return err
					}
					time.Sleep(time.Duration(j.Delay) * time.Second)
				}
			}
		}
	}
	return nil
}

// A TokenMaker is useful for creating a token that may be part of an API request.
type TokenMaker struct {
	Type    string // Can be MD5
	Prefix  string // When calculating a token, what should be prepended
	Postfix string // When calculating a token, what should be appended
	Size    int    // Take the first 'Size' characters from the result.
	Suffix  string // String to append after the token
}

// MakeToken will make a token from this token maker.
// If t.type is "MD5", prefixes 'prefix' to the input, postfixes 'postfix' to the input,
// takes the 'size' number of characters from the MD5 sum.
func (t TokenMaker) MakeToken(input string) (out string) {
	if t.Type == "MD5" {
		fullString := []byte(t.Prefix + input + t.Postfix)
		return fmt.Sprintf("%x", md5.Sum(fullString))[:t.Size] + t.Suffix
	}
	return t.Suffix
}
