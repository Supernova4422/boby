// This is not for configuring tests, this is for testing configs.

package test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"

	"github.com/BKrajancic/boby/m/v2/src/config"
	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/service/demoservice"
	"github.com/BKrajancic/boby/m/v2/src/storage"
	"github.com/google/go-cmp/cmp"
)

type ConfigTest struct {
	Input  string
	Expect [][]service.Message
}

// Check if msg is in a list.
func MessageInList(msg service.Message, list []service.Message) bool {
	// Checks each character one by one for breakpoint debugging.
	for _, msg2 := range list {
		if msg.Title == msg2.Title {
			if msg.URL == msg2.URL {
				good := true
				descLength := len(msg.Description)
				if len(msg2.Description) < descLength {
					descLength = len(msg.Description)
				}

				for i := 0; i < descLength; i++ {
					char1 := msg.Description[i]
					char2 := msg2.Description[i]
					if char1 != char2 {
						good = false
						break
					}
				}
				if good && cmp.Equal(msg, msg2) {
					return true
				}
			}
		}
	}
	return false
}

func GetTestInputs(filepath string) ([]ConfigTest, error) {
	var configTests []ConfigTest
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Printf("Unable to read file: %s", filepath)
		return configTests, nil
	}

	json.Unmarshal(bytes, &configTests)
	return configTests, nil
}

func TestConfig(t *testing.T) {
	configTests := "config_tests.json"
	if len(os.Args) > 3 {
		configDir := os.Args[len(os.Args)-1]
		inputFp := path.Join(configDir, configTests)
		_, err := os.Stat(inputFp)
		if err != nil {
			t.Log("Configuration file was not used for this test.")
		} else {
			t.Log("Configuration file was used for this test.")
			tempStorage := storage.GetTempStorage()
			var _storage storage.Storage = &tempStorage
			commands, err := config.ConfiguredBot(configDir, &_storage)
			_storage.SetDefaultGuildValue("prefix", "!")

			if err == nil {
				demoService := demoservice.DemoService{
					ServiceID: demoservice.ServiceID,
					Storage:   &_storage,
				}

				demoSender := demoservice.DemoSender{
					ServiceID: demoservice.ServiceID,
				}

				for i := range commands {
					commands[i].AddSender(&demoSender)
					types := []string{}
					for _, commandParameter := range commands[i].Parameters {
						types = append(types, commandParameter.Type)
					}
					demoService.Register(commands[i].Trigger, types, commands[i].Exec, commands[i].RouteByID)
				}

				inputTest, _ := GetTestInputs(inputFp)

				testConversation := service.Conversation{
					ServiceID:      demoSender.ID(),
					ConversationID: "0",
				}

				testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}
				results := make([]service.Message, 0)
				for _, input := range inputTest {
					demoService.AddMessage(testConversation, testSender, input.Input)
					demoService.Run()
					for _, expect := range input.Expect {
						if demoSender.IsEmpty() {
							t.Errorf("No responses on msg: %s", input.Input)
							t.Fail()
						} else {
							resultMessage, _ := demoSender.PopMessage()
							results = append(results, resultMessage)
							if !MessageInList(resultMessage, expect) {
								t.Errorf("Failed on msg: %s", input.Input)
								msg, _ := json.Marshal(resultMessage)
								t.Log(string(msg))
								msgExpect, _ := json.Marshal(expect)
								t.Log(string(msgExpect))
								t.Fail()
							}
						}

					}
					if demoSender.IsEmpty() == false {
						t.Errorf("Too many responses from: %s", input.Input)
						t.Fail()
					}
					// msg, _ := json.Marshal(results)
					// t.Log(string(msg))
				}
			} else {
				t.Fail()
			}
		}
	} else {
		here, err := os.Getwd()
		if os.IsNotExist(err) {
			here = "[ERROR, USE ABSOLUTE PATH]"
		}

		t.Log("When running go test, add arguments: '-args <dir>' where " +
			"<dir> is the same directory used when running the main " +
			"program (relative to " + here + ")." +
			"If a file '" + configTests + "' is present, it can be " +
			"used to ensure that all the configuration files are valid, " +
			"and produce the expected output.")

		// TODO add information on the format of config_tests.

		t.Log("The argument is not present, but no error will be raised.")
	}
}
