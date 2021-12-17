package command

import (
	"strings"
	"testing"
	"time"

	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/service/demoservice"
	"github.com/BKrajancic/boby/m/v2/src/storage"
	"github.com/google/go-cmp/cmp"
)

func Repeater(sender service.Conversation, user service.User, msg []interface{}, storage *storage.Storage, sink func(service.Conversation, service.Message)) {
	sink(sender, service.Message{Description: msg[0].(string)})
}

func TestCleanHistory(t *testing.T) {
	history := []int64{
		0, 1, 2, 3, 4, 6, 7, 8,
	}

	r := RateLimitConfig{
		SecondsPerInterval: 5,
	}

	newHistory := r.cleanHistory(10, history)
	if cmp.Equal(newHistory, []int64{6, 7, 8}) == false {
		t.Fail()
	}
}

func TestRateLimited(t *testing.T) {
	history := []int64{
		0, 1, 2, 3, 4, 6, 7, 8, 9,
	}

	r := RateLimitConfig{
		SecondsPerInterval: 5,
		TimesPerInterval:   2,
	}

	if r.rateLimited(10, history) == false {
		t.Fail()
	}

	if r.rateLimited(20, history) {
		t.Fail()
	}
}

func TestTimeRemaining(t *testing.T) {
	history := []int64{
		0, 1, 2, 3, 4, 6, 7, 8, 9,
	}

	r := RateLimitConfig{
		SecondsPerInterval: 5,
		TimesPerInterval:   2,
	}

	remaining := r.timeRemaining(10, history)
	if remaining.Seconds() != 3 {
		t.Fail()
	}

	timeRemaining := r.timeRemaining(20, history)
	if timeRemaining.Seconds() != 0 {
		t.Fail()
	}
}

func TestTimeRemainingUnordered(t *testing.T) {
	history := []int64{
		8, 6, 4, 7, 2, 9, 1, 0, 3,
	}

	r := RateLimitConfig{
		SecondsPerInterval: 5,
		TimesPerInterval:   2,
	}

	remaining := r.timeRemaining(10, history)
	if remaining.Seconds() != 3 {
		t.Fail()
	}
}

func TestTimeRemainingTooFew(t *testing.T) {
	history := []int64{8}

	r := RateLimitConfig{
		SecondsPerInterval: 5,
		TimesPerInterval:   2,
	}

	if r.rateLimited(10, history) {
		t.Fail()
	}
}

func TestRateLimitedCommand(t *testing.T) {
	demoSender := demoservice.DemoSender{}
	// Message to repeat.
	testConversation := service.Conversation{
		ServiceID:      "0",
		ConversationID: "0",
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}
	testCmd := "repeat"

	replyCommand := Command{
		Trigger:    testCmd,
		Parameters: []Parameter{{Type: "string"}},
		Exec:       Repeater,
		Help:       "Help",
	}

	limitMsg := "You hit the limit"
	rateLimitConfig := RateLimitConfig{
		TimesPerInterval:   2,
		SecondsPerInterval: 2,
		Body:               limitMsg,
		ID:                 "cmd",
	}

	tempStorage := storage.GetTempStorage()
	var _storage storage.Storage = &tempStorage

	rateLimitedCommand := rateLimitConfig.GetRateLimitedCommand(replyCommand)
	replyMsg := "Hello"
	msg := []interface{}{replyMsg}

	rateLimitedCommand.Exec(
		testConversation, testSender,
		msg, &_storage, demoSender.SendMessage,
	)

	resultMessage, _ := demoSender.PopMessage()
	if resultMessage.Description != replyMsg {
		t.Fail()
	}

	for i := 0; i < 20; i++ {
		rateLimitedCommand.Exec(
			testConversation, testSender,
			msg, &_storage, demoSender.SendMessage,
		)
	}

	resultMessage, _ = demoSender.PopMessage()
	if resultMessage.Description == limitMsg {
		t.Fail()
	}
}

func TestRateLimitedCommandDisaster(t *testing.T) {
	demoSender := demoservice.DemoSender{}
	// Message to repeat.
	testConversation := service.Conversation{
		ServiceID:      "0",
		ConversationID: "0",
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}
	testCmd := "repeat"

	replyCommand := Command{
		Trigger:    testCmd,
		Parameters: []Parameter{{Type: "string"}},
		Exec:       Repeater,
		Help:       "Help",
	}

	limitMsg := "You hit the limit"
	rateLimitID := "cmd"
	rateLimitConfig := RateLimitConfig{
		TimesPerInterval:   2,
		SecondsPerInterval: 2,
		Body:               limitMsg,
		ID:                 rateLimitID,
	}

	tempStorage := storage.GetTempStorage()
	var _storage storage.Storage = &tempStorage
	_storage.SetUserValue(testSender, rateLimitID, 0)

	rateLimitedCommand := rateLimitConfig.GetRateLimitedCommand(replyCommand)
	replyMsg := "Hello"
	msg := []interface{}{replyMsg}

	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()

	// If this doesn't panic, the test fails.
	rateLimitedCommand.Exec(
		testConversation, testSender,
		msg, &_storage, demoSender.SendMessage,
	)
}

func TestRateLimitedCommandDisasterGlobal(t *testing.T) {
	demoSender := demoservice.DemoSender{}
	// Message to repeat.
	testConversation := service.Conversation{
		ServiceID:      "0",
		ConversationID: "0",
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}
	testCmd := "repeat"

	replyCommand := Command{
		Trigger:    testCmd,
		Parameters: []Parameter{{Type: "string"}},
		Exec:       Repeater,
		Help:       "Help",
	}

	limitMsg := "You hit the limit"
	rateLimitID := "cmd"
	rateLimitConfig := RateLimitConfig{
		TimesPerInterval:   2,
		SecondsPerInterval: 2,
		Body:               limitMsg,
		ID:                 rateLimitID,
		Global:             true,
	}

	tempStorage := storage.GetTempStorage()
	var _storage storage.Storage = &tempStorage
	_storage.SetGlobalValue(rateLimitID, 0)

	rateLimitedCommand := rateLimitConfig.GetRateLimitedCommand(replyCommand)
	replyMsg := "Hello"
	msg := []interface{}{replyMsg}

	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()

	// If this doesn't panic, the test fails.
	rateLimitedCommand.Exec(
		testConversation, testSender,
		msg, &_storage, demoSender.SendMessage,
	)
}

func TestRateLimitedCommandWithGobStorage(t *testing.T) {
	demoSender := demoservice.DemoSender{}
	// Message to repeat.
	testConversation := service.Conversation{
		ServiceID:      "0",
		ConversationID: "0",
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}
	testCmd := "repeat"

	replyCommand := Command{
		Trigger:    testCmd,
		Parameters: []Parameter{{Type: "string"}},
		Exec:       Repeater,
		Help:       "Help",
	}

	limitMsg := "You hit the limit"
	rateLimitConfig := RateLimitConfig{
		TimesPerInterval:   2,
		SecondsPerInterval: 2,
		Body:               limitMsg,
		ID:                 "cmd",
	}

	tempStorage := storage.GetTempStorage()
	var _storage storage.Storage = &tempStorage

	rateLimitedCommand := rateLimitConfig.GetRateLimitedCommand(replyCommand)
	replyMsg := "Hello"
	msg := []interface{}{replyMsg}

	rateLimitedCommand.Exec(
		testConversation, testSender,
		msg, &_storage, demoSender.SendMessage,
	)

	resultMessage, _ := demoSender.PopMessage()
	if resultMessage.Description != replyMsg {
		t.Fail()
	}

	for i := 0; i < 20; i++ {
		rateLimitedCommand.Exec(
			testConversation, testSender,
			msg, &_storage, demoSender.SendMessage,
		)
	}

	resultMessage, _ = demoSender.PopMessage()
	if resultMessage.Description == limitMsg {
		t.Fail()
	}
}

func TestRateLimitedCommandMinute(t *testing.T) {
	demoSender := demoservice.DemoSender{}
	// Message to repeat.
	testConversation := service.Conversation{
		ServiceID:      "0",
		ConversationID: "0",
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}
	testCmd := "repeat"

	replyCommand := Command{
		Trigger:    testCmd,
		Parameters: []Parameter{{Type: "string"}},
		Exec:       Repeater,
		Help:       "Help",
	}

	limitMsg := "You hit the limit"
	rateLimitConfig := RateLimitConfig{
		TimesPerInterval:   1,
		SecondsPerInterval: 60 * 5,
		Body:               limitMsg,
		ID:                 "cmd",
	}

	tempStorage := storage.GetTempStorage()
	var _storage storage.Storage = &tempStorage

	rateLimitedCommand := rateLimitConfig.GetRateLimitedCommand(replyCommand)
	replyMsg := "Hello"
	msg := []interface{}{replyMsg}

	rateLimitedCommand.Exec(
		testConversation, testSender,
		msg, &_storage, demoSender.SendMessage,
	)

	resultMessage, _ := demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, replyMsg) {
		t.Fail()
	}

	for i := 0; i < 3; i++ {
		rateLimitedCommand.Exec(
			testConversation, testSender,
			msg, &_storage, demoSender.SendMessage,
		)
	}

	resultMessage, _ = demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, limitMsg) {
		t.Fail()
	}
}

func TestRateLimitedCommandHour(t *testing.T) {
	demoSender := demoservice.DemoSender{}
	// Message to repeat.
	testConversation := service.Conversation{
		ServiceID:      "0",
		ConversationID: "0",
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}
	testCmd := "repeat"

	replyCommand := Command{
		Trigger:    testCmd,
		Parameters: []Parameter{{Type: "string"}},
		Exec:       Repeater,
		Help:       "Help",
	}

	limitMsg := "You hit the limit"
	rateLimitConfig := RateLimitConfig{
		TimesPerInterval:   1,
		SecondsPerInterval: 60 * 60 * 2,
		Body:               limitMsg,
		ID:                 "cmd",
	}

	tempStorage := storage.GetTempStorage()
	var _storage storage.Storage = &tempStorage

	rateLimitedCommand := rateLimitConfig.GetRateLimitedCommand(replyCommand)
	replyMsg := "Hello"
	msg := []interface{}{replyMsg}

	rateLimitedCommand.Exec(
		testConversation, testSender,
		msg, &_storage, demoSender.SendMessage,
	)

	resultMessage, _ := demoSender.PopMessage()
	if resultMessage.Description != replyMsg {
		t.Fail()
	}

	for i := 0; i < 3; i++ {
		rateLimitedCommand.Exec(
			testConversation, testSender,
			msg, &_storage, demoSender.SendMessage,
		)
	}

	resultMessage, _ = demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, limitMsg) {
		t.Fail()
	}
}

func TestRateLimitedUseless(t *testing.T) {
	testCmd := "repeat"

	replyCommand := Command{
		Trigger:    testCmd,
		Parameters: []Parameter{{Type: "string"}},
		Exec:       Repeater,
		Help:       "Help",
	}

	limitMsg := "You hit the limit"
	rateLimitConfig := RateLimitConfig{
		TimesPerInterval:   0,
		SecondsPerInterval: 0,
		Body:               limitMsg,
		ID:                 "cmd",
	}

	rateLimitedCommand := rateLimitConfig.GetRateLimitedCommand(replyCommand)
	if replyCommand.Help != rateLimitedCommand.Help {
		t.Fail()
	}

	rateLimitedCommandInfo := rateLimitConfig.GetRateLimitedCommandInfo(replyCommand)
	demoSender := demoservice.DemoSender{}
	tempStorage := storage.GetTempStorage()
	var _storage storage.Storage = &tempStorage
	rateLimitedCommandInfo.Exec(
		service.Conversation{}, service.User{},
		[]interface{}{"string"}, &_storage, demoSender.SendMessage,
	)

	resultMessage, _ := demoSender.PopMessage()
	if resultMessage.Title != "This command has unlimited usage." {
		t.Fail()
	}
}

func TestRateLimitedNotUseless(t *testing.T) {
	testCmd := "repeat"

	replyCommand := Command{
		Trigger:    testCmd,
		Parameters: []Parameter{{Type: "string"}},
		Exec:       Repeater,
		Help:       "Help",
	}

	limitMsg := "You hit the limit"
	rateLimitConfig := RateLimitConfig{
		TimesPerInterval:   2,
		SecondsPerInterval: 2,
		Body:               limitMsg,
		ID:                 "cmd",
	}

	rateLimitedCommand := rateLimitConfig.GetRateLimitedCommand(replyCommand)
	if replyCommand.Help == rateLimitedCommand.Help {
		t.Fail()
	}
}

func TestRateLimitedCommandMultiUser(t *testing.T) {
	demoSender := demoservice.DemoSender{}
	// Message to repeat.
	testConversation := service.Conversation{
		ServiceID:      "0",
		ConversationID: "0",
	}
	badSender := service.User{Name: "bad_User", ServiceID: demoSender.ID()}
	goodSender := service.User{Name: "good_User", ServiceID: demoSender.ID()}
	testCmd := "repeat"

	replyCommand := Command{
		Trigger:    testCmd,
		Parameters: []Parameter{{Type: "string"}},
		Exec:       Repeater,
		Help:       "Help",
	}

	limitMsg := "You hit the limit"
	rateLimitConfig := RateLimitConfig{
		TimesPerInterval:   1,
		SecondsPerInterval: 60,
		Body:               limitMsg,
		ID:                 "cmd",
	}

	tempStorage := storage.GetTempStorage()
	var _storage storage.Storage = &tempStorage

	rateLimitedCommand := rateLimitConfig.GetRateLimitedCommand(replyCommand)
	replyMsg := "Hello"
	msg := []interface{}{replyMsg}

	rateLimitedCommand.Exec(
		testConversation, badSender,
		msg, &_storage, demoSender.SendMessage,
	)

	resultMessage, _ := demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, replyMsg) {
		t.Fail()
	}

	// Hit the limit
	rateLimitedCommand.Exec(
		testConversation, badSender,
		msg, &_storage, demoSender.SendMessage,
	)

	// Not limited
	rateLimitedCommand.Exec(
		testConversation, goodSender,
		msg, &_storage, demoSender.SendMessage,
	)

	// Ensure spamming user fails
	resultMessage, _ = demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, limitMsg) {
		t.Fail()
	}

	// Ensure good user fails
	resultMessage, _ = demoSender.PopMessage()
	if strings.HasPrefix(resultMessage.Description, limitMsg) {
		t.Fail()
	}
}

func TestRateLimitedCommandGlobal(t *testing.T) {
	demoSender := demoservice.DemoSender{}
	// Message to repeat.
	testConversation := service.Conversation{
		ServiceID:      "0",
		ConversationID: "0",
	}
	badSender := service.User{Name: "bad_User", ServiceID: demoSender.ID()}
	goodSender := service.User{Name: "good_User", ServiceID: demoSender.ID()}
	testCmd := "repeat"

	replyCommand := Command{
		Trigger:    testCmd,
		Parameters: []Parameter{{Type: "string"}},
		Exec:       Repeater,
		Help:       "Help",
	}

	limitMsg := "You hit the limit"
	rateLimitConfig := RateLimitConfig{
		TimesPerInterval:   1,
		SecondsPerInterval: 60,
		Body:               limitMsg,
		ID:                 "cmd",
		Global:             true,
	}

	tempStorage := storage.GetTempStorage()
	var _storage storage.Storage = &tempStorage

	rateLimitedCommand := rateLimitConfig.GetRateLimitedCommand(replyCommand)
	replyMsg := "Hello"
	msg := []interface{}{replyMsg}

	// First message should be fine
	rateLimitedCommand.Exec(
		testConversation, badSender,
		msg, &_storage, demoSender.SendMessage,
	)

	resultMessage, _ := demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, replyMsg) {
		t.Fail()
	}

	// Hit the limit
	rateLimitedCommand.Exec(
		testConversation, badSender,
		msg, &_storage, demoSender.SendMessage,
	)

	// Not limited
	rateLimitedCommand.Exec(
		testConversation, goodSender,
		msg, &_storage, demoSender.SendMessage,
	)

	// Ensure spamming user fails
	resultMessage, _ = demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, limitMsg) {
		t.Fail()
	}

	// Since this is a global command, a good user also gets punished.
	resultMessage, _ = demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, limitMsg) {
		t.Fail()
	}
}

func TestRateLimitedInfoCommandSecond(t *testing.T) {
	demoSender := demoservice.DemoSender{}
	// Message to repeat.
	testConversation := service.Conversation{
		ServiceID:      "0",
		ConversationID: "0",
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}
	testCmd := "repeat"

	replyCommand := Command{
		Trigger:    testCmd,
		Parameters: []Parameter{{Type: "string"}},
		Exec:       Repeater,
		Help:       "Help",
	}

	limitMsg := "You hit the limit"
	rateLimitConfig := RateLimitConfig{
		TimesPerInterval:   1,
		SecondsPerInterval: 1,
		Body:               limitMsg,
		ID:                 "cmd",
	}

	tempStorage := storage.GetTempStorage()
	var _storage storage.Storage = &tempStorage

	rateLimitedCommand := rateLimitConfig.GetRateLimitedCommand(replyCommand)
	rateLimitedInfoCommand := rateLimitConfig.GetRateLimitedCommandInfo(replyCommand)
	replyMsg := "Hello"
	msg := []interface{}{replyMsg}

	rateLimitedInfoCommand.Exec(
		testConversation, testSender,
		msg, &_storage, demoSender.SendMessage,
	)
	resultMessage, _ := demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, "0/1 per 1.00 Seconds remaining.") {
		t.Fail()
	}

	for i := 0; i < 3; i++ {
		rateLimitedCommand.Exec(
			testConversation, testSender,
			msg, &_storage, demoSender.SendMessage,
		)
		demoSender.PopMessage()
	}

	rateLimitedInfoCommand.Exec(
		testConversation, testSender,
		msg, &_storage, demoSender.SendMessage,
	)

	resultMessage, _ = demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, "1/1 per 1.00 Seconds remaining.") {
		t.Fail()
	}
}

func TestRateLimitedInfoCommandMinute(t *testing.T) {
	demoSender := demoservice.DemoSender{}
	// Message to repeat.
	testConversation := service.Conversation{
		ServiceID:      "0",
		ConversationID: "0",
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}
	testCmd := "repeat"

	replyCommand := Command{
		Trigger:    testCmd,
		Parameters: []Parameter{{Type: "string"}},
		Exec:       Repeater,
		Help:       "Help",
	}

	limitMsg := "You hit the limit"
	rateLimitConfig := RateLimitConfig{
		TimesPerInterval:   1,
		SecondsPerInterval: 60,
		Body:               limitMsg,
		ID:                 "cmd",
	}

	tempStorage := storage.GetTempStorage()
	var _storage storage.Storage = &tempStorage

	rateLimitedCommand := rateLimitConfig.GetRateLimitedCommand(replyCommand)
	rateLimitedInfoCommand := rateLimitConfig.GetRateLimitedCommandInfo(replyCommand)
	replyMsg := "Hello"
	msg := []interface{}{replyMsg}

	rateLimitedInfoCommand.Exec(
		testConversation, testSender,
		msg, &_storage, demoSender.SendMessage,
	)
	resultMessage, _ := demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, "0/1 per 1.00 Minutes remaining.") {
		t.Fail()
	}

	for i := 0; i < 3; i++ {
		rateLimitedCommand.Exec(
			testConversation, testSender,
			msg, &_storage, demoSender.SendMessage,
		)
		demoSender.PopMessage()
	}

	rateLimitedInfoCommand.Exec(
		testConversation, testSender,
		msg, &_storage, demoSender.SendMessage,
	)

	resultMessage, _ = demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, "1/1 per 1.00 Minutes remaining.") {
		t.Fail()
	}
}

func TestRateLimitedInfoCommandHour(t *testing.T) {
	demoSender := demoservice.DemoSender{}
	// Message to repeat.
	testConversation := service.Conversation{
		ServiceID:      "0",
		ConversationID: "0",
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}
	testCmd := "repeat"

	replyCommand := Command{
		Trigger:    testCmd,
		Parameters: []Parameter{{Type: "string"}},
		Exec:       Repeater,
		Help:       "Help",
	}

	limitMsg := "You hit the limit"
	rateLimitConfig := RateLimitConfig{
		TimesPerInterval:   1,
		SecondsPerInterval: 60 * 60,
		Body:               limitMsg,
		ID:                 "cmd",
	}

	tempStorage := storage.GetTempStorage()
	var _storage storage.Storage = &tempStorage

	rateLimitedCommand := rateLimitConfig.GetRateLimitedCommand(replyCommand)
	rateLimitedInfoCommand := rateLimitConfig.GetRateLimitedCommandInfo(replyCommand)
	replyMsg := "Hello"
	msg := []interface{}{replyMsg}

	rateLimitedInfoCommand.Exec(
		testConversation, testSender,
		msg, &_storage, demoSender.SendMessage,
	)
	resultMessage, _ := demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, "0/1 per 1.00 Hours remaining.") {
		t.Fail()
	}

	for i := 0; i < 3; i++ {
		rateLimitedCommand.Exec(
			testConversation, testSender,
			msg, &_storage, demoSender.SendMessage,
		)
		demoSender.PopMessage()
	}

	rateLimitedInfoCommand.Exec(
		testConversation, testSender,
		msg, &_storage, demoSender.SendMessage,
	)

	resultMessage, _ = demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, "1/1 per 1.00 Hours remaining.") {
		t.Fail()
	}
}

func TestRateLimitedInfoCommandDays(t *testing.T) {
	demoSender := demoservice.DemoSender{}
	// Message to repeat.
	testConversation := service.Conversation{
		ServiceID:      "0",
		ConversationID: "0",
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}
	testCmd := "repeat"

	replyCommand := Command{
		Trigger:    testCmd,
		Parameters: []Parameter{{Type: "string"}},
		Exec:       Repeater,
		Help:       "Help",
	}

	limitMsg := "You hit the limit"
	rateLimitConfig := RateLimitConfig{
		TimesPerInterval:   1,
		SecondsPerInterval: 60 * 60 * 24,
		Body:               limitMsg,
		ID:                 "cmd",
	}

	tempStorage := storage.GetTempStorage()
	var _storage storage.Storage = &tempStorage

	rateLimitedCommand := rateLimitConfig.GetRateLimitedCommand(replyCommand)
	rateLimitedInfoCommand := rateLimitConfig.GetRateLimitedCommandInfo(replyCommand)
	replyMsg := "Hello"
	msg := []interface{}{replyMsg}

	rateLimitedInfoCommand.Exec(
		testConversation, testSender,
		msg, &_storage, demoSender.SendMessage,
	)
	resultMessage, _ := demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, "0/1 per 1.00 Days remaining.") {
		t.Fail()
	}

	for i := 0; i < 3; i++ {
		rateLimitedCommand.Exec(
			testConversation, testSender,
			msg, &_storage, demoSender.SendMessage,
		)
		demoSender.PopMessage()
	}

	rateLimitedInfoCommand.Exec(
		testConversation, testSender,
		msg, &_storage, demoSender.SendMessage,
	)

	resultMessage, _ = demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, "1/1 per 1.00 Days remaining.") {
		t.Fail()
	}
}

func TestRateLimitedInfoCommandSecondOld(t *testing.T) {
	demoSender := demoservice.DemoSender{}
	// Message to repeat.
	testConversation := service.Conversation{
		ServiceID:      "0",
		ConversationID: "0",
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}
	testCmd := "repeat"

	replyCommand := Command{
		Trigger:    testCmd,
		Parameters: []Parameter{{Type: "string"}},
		Exec:       Repeater,
		Help:       "Help",
	}

	limitMsg := "You hit the limit"
	rateLimitConfig := RateLimitConfig{
		TimesPerInterval:   1,
		SecondsPerInterval: 1,
		Body:               limitMsg,
		ID:                 "cmd",
	}

	tempStorage := storage.GetTempStorage()
	var _storage storage.Storage = &tempStorage

	rateLimitedCommand := rateLimitConfig.GetRateLimitedCommand(replyCommand)
	rateLimitedInfoCommand := rateLimitConfig.GetRateLimitedCommandInfo(replyCommand)
	replyMsg := "Hello"
	msg := []interface{}{replyMsg}

	rateLimitedInfoCommand.Exec(
		testConversation, testSender,
		msg, &_storage, demoSender.SendMessage,
	)
	resultMessage, _ := demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, "0/1 per 1.00 Seconds remaining.") {
		t.Fail()
	}

	for i := 0; i < 3; i++ {
		rateLimitedCommand.Exec(
			testConversation, testSender,
			msg, &_storage, demoSender.SendMessage,
		)
		demoSender.PopMessage()
	}

	rateLimitedInfoCommand.Exec(
		testConversation, testSender,
		msg, &_storage, demoSender.SendMessage,
	)

	resultMessage, _ = demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, "1/1 per 1.00 Seconds remaining.") {
		t.Fail()
	}

	time.Sleep(2 * time.Second)

	rateLimitedInfoCommand.Exec(
		testConversation, testSender,
		msg, &_storage, demoSender.SendMessage,
	)

	resultMessage, _ = demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, "0/1 per 1.00 Seconds remaining.") {
		t.Fail()
	}
}
