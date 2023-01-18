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
		},

		{
			Trigger: IsAdminTrigger + "user",
			Parameters: []Parameter{
				{
					Name:        "user",
					Description: "User to check if is an admin",
					Type:        "user",
				},
			},
			Exec: CheckAdmin,
			Help: "Check if a user is admin.",
		},

		{
			Trigger: SetAdminTrigger + "user",
			Parameters: []Parameter{
				{
					Name:        "user",
					Description: "User to set as an admin",
					Type:        "user",
				},
			},
			Exec: SetAdmin,
			Help: "Set a user as an admin, therefore giving them all permissions for this bot. Users/Roles with any of the following server permissions are automatically treated as admin: 'Administrator', 'Manage Server', 'Manage Webhooks.'",
		},

		{
			Trigger: UnsetAdminTrigger + "user",
			Parameters: []Parameter{
				{
					Name:        "user",
					Description: "User/Role to unset as an admin",
					Type:        "user",
				},
			},
			Exec:      UnsetAdmin,
			Help:      "Unset a role or user as an admin, therefore giving them usual permissions.",
			HelpInput: "[@role or @user]",
		},

		{
			Trigger: IsAdminTrigger + "role",
			Parameters: []Parameter{
				{
					Name:        "role",
					Description: "Role To check if is an admin",
					Type:        "role",
				},
			},
			Exec: CheckAdmin,
			Help: "Check if a role is admin.",
		},

		{
			Trigger: SetAdminTrigger + "role",
			Parameters: []Parameter{
				{
					Name:        "role",
					Description: "Role to set as an admin",
					Type:        "role",
				},
			},
			Exec: SetAdmin,
			Help: "Set a role as an admin, therefore giving them all permissions for this bot. Users/Roles with any of the following server permissions are automatically treated as admin: 'Administrator', 'Manage Server', 'Manage Webhooks.'",
		},

		{
			Trigger: UnsetAdminTrigger + "role",
			Parameters: []Parameter{
				{
					Name:        "role",
					Description: "Role to unset as an admin",
					Type:        "role",
				},
			},
			Exec: UnsetAdmin,
			Help: "Unset a role as an admin, therefore giving them usual permissions.",
		},
		{
			Trigger: "setprefix",
			Parameters: []Parameter{
				{
					Name:        "prefix",
					Description: "Set the prefix of commands, for this server.",
					Type:        "string",
				},
			},
			Exec:      SetPrefix,
			Help:      "Set the prefix of all commands of this bot, for this server.",
			HelpInput: "[word]",
		},
		{
			Trigger:    "error",
			Parameters: []Parameter{},
			Exec:       CreateError,
			Help:       "For testing purposes. Creates an error.",
		},
	}
}
