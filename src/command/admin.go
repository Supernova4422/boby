package command

import (
	"regexp"

	"github.com/BKrajancic/boby/m/v2/src/storage"
)

const ImAdminTrigger = "imadmin"
const SetAdminTrigger = "setadmin"
const UnsetAdminTrigger = "unsetadmin"
const IsAdminTrigger = "isadmin"

// Repo is a URL to this project's repository. Useful for showing with help information.
const Repo = "https://github.com/BKrajancic/boby"

// AdminCommands returns an array of commands for handling admins.
func AdminCommands(storage *storage.Storage) []Command {
	return []Command{
		{
			Trigger:   ImAdminTrigger,
			Pattern:   regexp.MustCompile("(.*)"),
			Exec:      ImAdmin,
			Help:      "Check if the sender is an admin.",
			HelpInput: "[@role or @user]",
			Storage:   storage,
		},

		{
			Trigger:   IsAdminTrigger,
			Pattern:   regexp.MustCompile("(.*)"),
			Exec:      CheckAdmin,
			Help:      "Check if a role or user is an admin.",
			HelpInput: "[@role or @user]",
			Storage:   storage,
		},

		{
			Trigger:   "setadmin",
			Pattern:   regexp.MustCompile("(.*)"),
			Exec:      SetAdmin,
			Help:      "Set a role or user as an admin, therefore giving them all permissions for this bot. Users/Roles with any of the following server permissions are automatically treated as admin: 'Administrator', 'Manage Server', 'Manage Webhooks.'",
			HelpInput: "[@role or @user]",
			Storage:   storage,
		},

		{
			Trigger:   UnsetAdminTrigger,
			Pattern:   regexp.MustCompile("(.*)"),
			Exec:      UnsetAdmin,
			Help:      "Unset a role or user as an admin, therefore giving them usual permissions.",
			HelpInput: "[@role or @user]",
			Storage:   storage,
		},

		{
			Trigger:   "setprefix",
			Pattern:   regexp.MustCompile("(.*)"),
			Exec:      SetPrefix,
			Help:      "Set the prefix of all commands of this bot, for this server.",
			HelpInput: "[word]",
			Storage:   storage,
		},
	}
}
