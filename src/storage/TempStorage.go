package storage

import (
	"fmt"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
)

type UserKey struct {
	User      string
	ServiceID string
}

// TempStorage implements the Storage interface, but data is lost on destruction.
type TempStorage struct {
	UserValues  map[UserKey]map[string]string
	GuildValues map[service.Guild]map[string]string
	Admins      map[string]map[string][]string
}

// GetGuildValue retrieves the value for key, for a Guild.
// Returns an error if the key doesn't exist or can't be retrieved.
func (t *TempStorage) GetGuildValue(guild service.Guild, key string) (val string, err error) {
	err = fmt.Errorf("No value for key %s", key)
	if _, ok := t.GuildValues[guild]; ok == false {
		return
	}

	val, ok := t.GuildValues[guild][key]
	if ok {
		err = nil
	}
	return
}

// SetGuildValue sets the value for key, for a Guild.
func (t *TempStorage) SetGuildValue(guild service.Guild, key string, value string) {
	_, ok := t.GuildValues[guild]
	if ok == false {
		t.GuildValues = make(map[service.Guild]map[string]string)
	}
	t.GuildValues[guild][key] = value
}

// GetUserValue retrieves the value for key, for a User.
// Returns an error if the key doesn't exist or can't be retrieved.
func (t *TempStorage) GetUserValue(serviceID string, user string, key string) (val string, err error) {
	err = fmt.Errorf("No value for key %s", key)
	userKey := UserKey{User: user, ServiceID: serviceID}
	if _, ok := t.UserValues[userKey]; ok == false {
		return
	}

	val, ok := t.UserValues[userKey][key]
	if ok {
		err = nil
	}
	return
}

// SetUserValue sets the value for key, for a Guild.
func (t *TempStorage) SetUserValue(serviceID string, user string, key string, value string) {
	userKey := UserKey{User: user, ServiceID: serviceID}
	_, ok := t.UserValues[userKey]
	if ok == false {
		t.UserValues = make(map[UserKey]map[string]string)
	}
	t.UserValues[userKey][key] = value
}

// IsAdmin returns true if ID is an admin.
func (t *TempStorage) IsAdmin(guild service.Guild, ID string) bool {
	serviceAdmins, ok := t.Admins[guild.ServiceID]
	if ok {
		admins, ok := serviceAdmins[guild.GuildID]
		if ok {
			for _, id := range admins {
				if id == ID {
					return true
				}
			}
		}
	}
	return false
}

// SetAdmin sets a userID as an admin for a guild.
func (t *TempStorage) SetAdmin(guild service.Guild, ID string) {
	_, ok := t.Admins[guild.ServiceID]
	if ok {
		_, ok := t.Admins[guild.ServiceID][guild.GuildID]
		if ok {
			t.Admins[guild.ServiceID][guild.GuildID] = append(t.Admins[guild.ServiceID][guild.GuildID], ID)
		} else {
			t.Admins[guild.ServiceID] = make(map[string][]string)
		}
	} else {
		if t.Admins == nil {
			t.Admins = make(map[string]map[string][]string)
			t.Admins[guild.ServiceID] = make(map[string][]string)
		}
		t.Admins[guild.ServiceID][guild.GuildID] = []string{ID}
	}
}

// UnsetAdmin removes userID as an admin for a guild.
func (t *TempStorage) UnsetAdmin(guild service.Guild, ID string) {
	_, ok := t.Admins[guild.ServiceID]
	if ok {
		_, ok := t.Admins[guild.ServiceID][guild.GuildID]
		if ok {
			newAdmins := []string{}
			for _, admin := range t.Admins[guild.ServiceID][guild.GuildID] {
				if admin != ID {
					newAdmins = append(newAdmins, admin)
				}
			}
			t.Admins[guild.ServiceID][guild.GuildID] = newAdmins
		}
	}
}
