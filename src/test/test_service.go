package test

import (
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
)

const service_id = "CLI"

type Cli_Service_Subject struct {
	// messages and users are co-indexed
	messages []string
	users    []service.User

	observers []service.Service_Observer
}

func (cli *Cli_Service_Subject) Id() string {
	return service_id
}

func (cli *Cli_Service_Subject) Register(observer service.Service_Observer) {
	cli.observers = append(cli.observers, observer)
}

func (cli *Cli_Service_Subject) AddMessage(user service.User, message string) {
	cli.messages = append(cli.messages, message)
	cli.users = append(cli.users, user)
}

func (cli *Cli_Service_Subject) Run() {
	if len(cli.messages) != len(cli.users) {
		panic("users and messages should have the same length because the arrays are co-indexed (i.e. user[0] sends message [0]).")
	}

	for i := 0; i < len(cli.messages); i++ {
		user := cli.users[i]
		msg := cli.messages[i]
		for _, observer := range cli.observers {
			observer.OnMessage(user, msg)
		}
	}
}

type Cli_Service_Sender struct {
	messages []string
	senders  []service.User
}

func (cli *Cli_Service_Sender) SendMessage(sender service.User, message string) {
	cli.messages = append(cli.messages, message)
	cli.senders = append(cli.senders, sender)
}

func (cli *Cli_Service_Sender) IsEmpty() bool {
	return len(cli.messages) == 0
}
func (cli *Cli_Service_Sender) PopMessage() (message string, sender service.User) {
	message = cli.messages[0]
	sender = cli.senders[0]
	cli.messages = cli.messages[1:]
	cli.senders = cli.senders[1:]
	return
}

func (cli Cli_Service_Sender) Id() string {
	return service_id
}
