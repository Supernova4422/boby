package command

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/storage"
)

// OxfordDictionaryConfig can be turned into a scraper that uses GoQuery.
type OxfordDictionaryConfig struct {
	AppId              string
	AppKey             string
	Trigger            string
	HelpText           string
	HelpInput          string
	SourceLanguage     string
	TargetLanguage     string
	TimesPerInterval   int
	SecondsPerInterval int
	Body               string
	ID                 string
}

// GetOxfordConfigs retrieves an array of OxfordDictionaryConfig by parsing JSON from a buffer.
// If a file doesn't exist, an example is made in its place, and an error is returned.
func GetOxfordConfigs(reader io.Reader) ([]OxfordDictionaryConfig, error) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var config []OxfordDictionaryConfig
	return config, json.Unmarshal(bytes, &config)
}

func (o *OxfordDictionaryConfig) Command() (Command, Command, error) {
	curry := func(sender service.Conversation, user service.User, msg []interface{}, storage *storage.Storage, sink func(service.Conversation, service.Message)) {
		url := fmt.Sprintf("https://od-api.oxforddictionaries.com/api/v2/translations/%s/%s/%s?strictMatch=false", o.SourceLanguage, o.TargetLanguage, url.PathEscape(msg[0].(string)))
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println(fmt.Sprintf("Oxford API error 1: %s", err))
			sink(
				sender,
				service.Message{
					Title: "An error occured processing your request",
				},
			)
			return
		}

		req.Header.Set("app_id", o.AppId)
		req.Header.Set("app_key", o.AppKey)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println(fmt.Sprintf("Oxford API error 2: %s", err))
			sink(
				sender,
				service.Message{
					Title: "An error occured processing your request",
				},
			)
			return
		}

		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}

		var dict OxfordTranslateResponseStruct
		err = json.Unmarshal(buf, &dict)
		if err != nil {
			return
		}
		if len(dict.Results) == 0 || len(dict.Results[0].LexicalEntries) == 0 {
			sink(sender, service.Message{
				Title:       "Unable to find word",
				Description: "Please check the spelling, and try again.",
				URL:         "https://www.oxfordlearnersdictionaries.com/",
			})
			return
		}

		results := []string{}
		for _, result := range dict.Results {
			for _, lexicalEntry := range result.LexicalEntries {
				for _, entry := range lexicalEntry.Entries {
					for _, sense := range entry.Senses {
						for _, translation := range sense.Translations {
							text := translation.Text
							unique := true
							for _, resultsEntry := range results {
								if text == resultsEntry {
									unique = false
									break
								}
							}

							if unique {
								results = append(results, text)
							}
						}
					}
				}
			}
		}

		val := fmt.Sprintf("[%s] ", dict.Results[0].LexicalEntries[0].LexicalCategory.Text)
		for _, resultsEntry := range results {
			val += resultsEntry + "; "
		}

		sink(sender, service.Message{
			Title:       "Translation: " + msg[0].(string),
			Description: val,
			URL:         "https://www.oxfordlearnersdictionaries.com/",
		})
	}

	cmd := Command{
		Trigger: o.Trigger,
		Parameters: []Parameter{{
			Type:        "string",
			Name:        "word",
			Description: "word to translate",
		}},
		Exec:      curry,
		Help:      o.HelpText,
		HelpInput: o.HelpInput,
	}

	config := RateLimitConfig{
		TimesPerInterval:   o.TimesPerInterval,
		SecondsPerInterval: int64(o.SecondsPerInterval),
		Body:               o.Body,
		ID:                 o.ID,
		Global:             true,
	}

	actual := config.GetRateLimitedCommand(cmd)
	info := config.GetRateLimitedCommandInfo(cmd)
	return actual, info, nil
}

type ResponseStruct struct {
	ID       string `json:"id"`
	Metadata struct {
		Operation string `json:"operation"`
		Provider  string `json:"provider"`
		Schema    string `json:"schema"`
	} `json:"metadata"`
	Results []struct {
		ID             string `json:"id"`
		Language       string `json:"language"`
		LexicalEntries []struct {
			Entries []struct {
				Etymologies    []string `json:"etymologies"`
				Pronunciations []struct {
					AudioFile        string   `json:"audioFile,omitempty"`
					Dialects         []string `json:"dialects"`
					PhoneticNotation string   `json:"phoneticNotation"`
					PhoneticSpelling string   `json:"phoneticSpelling"`
				} `json:"pronunciations"`
				Senses []struct {
					Definitions []string `json:"definitions"`
					Examples    []struct {
						Text string `json:"text"`
					} `json:"examples"`
					ID              string `json:"id"`
					SemanticClasses []struct {
						ID   string `json:"id"`
						Text string `json:"text"`
					} `json:"semanticClasses"`
					ShortDefinitions []string `json:"shortDefinitions"`
					Subsenses        []struct {
						Definitions []string `json:"definitions"`
						Examples    []struct {
							Text string `json:"text"`
						} `json:"examples"`
						ID              string `json:"id"`
						SemanticClasses []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
						} `json:"semanticClasses"`
						ShortDefinitions []string `json:"shortDefinitions"`
					} `json:"subsenses,omitempty"`
					Synonyms []struct {
						Language string `json:"language"`
						Text     string `json:"text"`
					} `json:"synonyms"`
					ThesaurusLinks []struct {
						EntryID string `json:"entry_id"`
						SenseID string `json:"sense_id"`
					} `json:"thesaurusLinks"`
				} `json:"senses"`
			} `json:"entries"`
			Language        string `json:"language"`
			LexicalCategory struct {
				ID   string `json:"id"`
				Text string `json:"text"`
			} `json:"lexicalCategory"`
			Phrases []struct {
				ID   string `json:"id"`
				Text string `json:"text"`
			} `json:"phrases"`
			Text string `json:"text"`
		} `json:"lexicalEntries"`
		Type string `json:"type"`
		Word string `json:"word"`
	} `json:"results"`
	Word string `json:"word"`
}
