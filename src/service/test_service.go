package service

import (
	"fmt"
)

const service_id = "CLID"

type Cli_Service_Subject struct {
	// messages and users are co-indexed
	messages []string
	users    []User

	observers []Service_Observer
}

func (cli *Cli_Service_Subject) Id() string {
	return service_id
}

func (cli *Cli_Service_Subject) Register(observer Service_Observer) {
	cli.observers = append(cli.observers, observer)
}

func (cli *Cli_Service_Subject) AddMessage(user User, message string) {
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

type Cli_Service_Sender struct{}

func (cli Cli_Service_Sender) SendMessage(sender User, msg string) {
	fmt.Println(sender.Name+":", msg)
}

func (cli Cli_Service_Sender) OnMessage(sender User, msg string) {
	cli.SendMessage(sender, msg)
}

func (cli Cli_Service_Sender) Id() string {
	return service_id
}
