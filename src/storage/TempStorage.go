package storage

import (
	"fmt"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
)

const adminKey = "Admin"

// TempStorage implements the Storage interface, but data is lost on destruction.
type TempStorage struct {
	UserValues  map[string]map[string]map[string]interface{}
	GuildValues map[string]map[string]map[string]interface{}
}

// GetGuildValue retrieves the value for key, for a Guild.
// Returns an error if the key doesn't exist or can't be retrieved.
func (t *TempStorage) GetGuildValue(guild service.Guild, key string) (val interface{}, err error) {
	err = fmt.Errorf("No value for key %s", key)
	if serviceMap, ok := t.GuildValues[guild.ServiceID]; ok == true {
		if guildMap, ok := serviceMap[guild.GuildID]; ok {
			val, ok = guildMap[key]
			if ok {
				err = nil
			}
		}
	}
	return
}

// SetGuildValue sets the value for key, for a Guild.
func (t *TempStorage) SetGuildValue(guild service.Guild, key string, val interface{}) {
	if t.GuildValues == nil {
		t.GuildValues = make(map[string]map[string]map[string]interface{})
	}

	if t.GuildValues[guild.ServiceID] == nil {
		t.GuildValues[guild.ServiceID] = make(map[string]map[string]interface{})
	}

	if t.GuildValues[guild.ServiceID][guild.GuildID] == nil {
		t.GuildValues[guild.ServiceID][guild.GuildID] = make(map[string]interface{})
	}

	t.GuildValues[guild.ServiceID][guild.GuildID][key] = val
}

// GetUserValue retrieves the value for key, for a User.
// Returns an error if the key doesn't exist or can't be retrieved.
func (t *TempStorage) GetUserValue(user service.User, key string) (val interface{}, err error) {
	err = fmt.Errorf("No value for key %s", key)
	if serviceMap, ok := t.UserValues[user.ServiceID]; ok {
		if userMap, ok := serviceMap[user.Name]; ok {
			val, ok = userMap[key]
			if ok {
				err = nil
			}
		}
	}
	return
}

// SetUserValue sets the value for key, for a Guild.
func (t *TempStorage) SetUserValue(user service.User, key string, val interface{}) {
	if t.UserValues == nil {
		t.UserValues = make(map[string]map[string]map[string]interface{})
	}

	if t.UserValues[user.ServiceID] == nil {
		t.UserValues[user.ServiceID] = make(map[string]map[string]interface{})
	}

	if t.UserValues[user.ServiceID][user.Name] == nil {
		t.UserValues[user.ServiceID][user.Name] = make(map[string]interface{})
	}

	t.UserValues[user.ServiceID][user.Name][key] = val
}

// IsAdmin returns true if ID is an admin.
func (t *TempStorage) IsAdmin(guild service.Guild, ID string) bool {
	if val, err := t.GetGuildValue(guild, adminKey); err == nil {
		if admins, ok := val.([]string); ok {
			for _, adminID := range admins {
				if adminID == ID {
					return true
				}
			}
		}
	}
	return false
}

// SetAdmin sets a userID as an admin for a guild.
func (t *TempStorage) SetAdmin(guild service.Guild, ID string) {
	if val, err := t.GetGuildValue(guild, adminKey); err == nil {
		currentAdmins, ok := val.([]string)
		if ok {
			t.SetGuildValue(guild, adminKey, append(currentAdmins, ID))
		} else {
			panic(ok)
		}
	} else {
		t.SetGuildValue(guild, adminKey, []string{ID})
	}
}

// UnsetAdmin removes userID as an admin for a guild.
func (t *TempStorage) UnsetAdmin(guild service.Guild, ID string) {
	newAdmins := []string{}
	if val, err := t.GetGuildValue(guild, adminKey); err == nil {
		currentAdmins, ok := val.([]string)
		if ok == false {
			panic(ok)
		}
		for _, adminID := range currentAdmins {
			if adminID != ID {
				newAdmins = append(newAdmins, adminID)
			}
		}
	}
	t.SetGuildValue(guild, adminKey, newAdmins)
}
