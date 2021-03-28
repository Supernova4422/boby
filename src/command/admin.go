package command

// ImAdminTrigger is a trigger to use for an ImAdmin command.
const ImAdminTrigger = "imadmin"

// SetAdminTrigger is a trigger to use for an SetAdmin command.
const SetAdminTrigger = "setadmin"

// UnsetAdminTrigger is a trigger to use for an UnsetAdmin command.
const UnsetAdminTrigger = "unsetadmin"

// IsAdminTrigger is a trigger to use for an IsAdmin command.
const IsAdminTrigger = "isadmin"

// Repo is a URL to this project's repository. Useful for showing with help information.
const Repo = "https://github.com/BKrajancic/boby"

// AdminCommands returns an array of commands for handling admins.
func AdminCommands() []Command {
	return []Command{
		{
			Trigger:    ImAdminTrigger,
			Parameters: []Parameter{},
			Exec:       ImAdmin,
			Help:       "Check if the sender is an admin.",
			HelpInput:  "[@role or @user]",
		},

		{
			Trigger: IsAdminTrigger,
			Parameters: []Parameter{
				{
					Name:        "User_or_Role",
					Description: "User/Role To check if is an admin",
					Type:        "string",
				},
			},
			Exec:      CheckAdmin,
			Help:      "Check if a role/user is admin.",
			HelpInput: "[@role or @user]",
		},

		{
			Trigger: "setadmin",
			Parameters: []Parameter{
				{
					Name:        "User_or_Role",
					Description: "User/Role to set as an admin",
					Type:        "string",
				},
			},
			Exec:      SetAdmin,
			Help:      "Set a role or user as an admin, therefore giving them all permissions for this bot. Users/Roles with any of the following server permissions are automatically treated as admin: 'Administrator', 'Manage Server', 'Manage Webhooks.'",
			HelpInput: "[@role or @user]",
		},

		{
			Trigger: UnsetAdminTrigger,
			Parameters: []Parameter{
				{
					Name:        "User_or_Role",
					Description: "User/Role to unset as an admin",
					Type:        "string",
				},
			},
			Exec:      UnsetAdmin,
			Help:      "Unset a role or user as an admin, therefore giving them usual permissions.",
			HelpInput: "[@role or @user]",
		},

		{
			Trigger: "setprefix",
			Parameters: []Parameter{
				{
					Name:        "prefix",
					Description: "This feature should no longer work.",
					Type:        "string",
				},
			},
			Exec:      SetPrefix,
			Help:      "Set the prefix of all commands of this bot, for this server.",
			HelpInput: "[word]",
		},
	}
}
