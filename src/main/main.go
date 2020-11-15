package main

import (
	"github.com/BKrajancic/FLD-Bot/m/v2/src/fld_bot"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
)

func main() {
	simple_bot := fld_bot.Simple_Bot{}

	// Create the source
	test_service_subject := service.Cli_Service_Subject{}
	test_service_subject.Register(simple_bot)

	// Create the destination
	test_service_sender := service.Cli_Service_Sender{}
	simple_bot.Register(test_service_sender)

	// Add Messages
	test_service_subject.AddMessage(
		service.User{Name: "Test_User", Id: test_service_subject.Id()},
		"Test1",
	)

	test_service_subject.Run()
}
