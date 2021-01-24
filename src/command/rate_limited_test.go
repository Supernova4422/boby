package command

import (
	"regexp"
	"testing"

	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/service/demoservice"
	"github.com/BKrajancic/boby/m/v2/src/storage"
	"github.com/google/go-cmp/cmp"
)

func Repeater(sender service.Conversation, user service.User, msg [][]string, storage *storage.Storage, sink func(service.Conversation, service.Message)) {
	sink(sender, service.Message{Description: msg[0][1]})
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
	if remaining != 3 {
		t.Fail()
	}

	if r.timeRemaining(20, history) != 0 {
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
	if remaining != 3 {
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
		Trigger: testCmd,
		Pattern: regexp.MustCompile("(.*)"),
		Exec:    Repeater,
		Help:    "Help",
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
	msg := [][]string{{"", replyMsg}}

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
		Trigger: testCmd,
		Pattern: regexp.MustCompile("(.*)"),
		Exec:    Repeater,
		Help:    "Help",
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
	msg := [][]string{{"", replyMsg}}

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
		Trigger: testCmd,
		Pattern: regexp.MustCompile("(.*)"),
		Exec:    Repeater,
		Help:    "Help",
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
	msg := [][]string{{"", replyMsg}}

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
