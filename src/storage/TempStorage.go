package storage

import (
	"fmt"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
)

// TempStorage implements the Storage interface, but data is lost on destruction.
type TempStorage struct {
	Admins map[string]map[string][]string
	Values map[string]map[string]map[string]string
}

// GetValue retrieves the value for key, for a Guild.
// Returns an error if the key doesn't exist or can't be retrieved.
func (t *TempStorage) GetValue(guild service.Guild, key string) (string, error) {
	_, ok := t.Values[guild.ServiceID]
	if ok == false {
		return "", fmt.Errorf("No entries for ServiceID %s", guild.ServiceID)
	}

	_, ok = t.Values[guild.ServiceID][guild.GuildID]
	if ok == false {
		return "", fmt.Errorf("No entries for GuildID %s", guild.GuildID)
	}

	val, ok := t.Values[guild.ServiceID][guild.GuildID][key]
	if ok == false {
		return "", fmt.Errorf("No value for key %s", key)
	}
	return val, nil
}

// SetValue sets the value for key, for a Guild.
func (t *TempStorage) SetValue(guild service.Guild, key string, value string) {
	_, ok := t.Values[guild.ServiceID]
	if ok == false {
		t.Values = make(map[string]map[string]map[string]string)
		t.Values[guild.ServiceID] = make(map[string]map[string]string)
	}

	_, ok = t.Values[guild.ServiceID][guild.GuildID]
	if ok == false {
		t.Values[guild.ServiceID][guild.GuildID] = make(map[string]string)
	}

	t.Values[guild.ServiceID][guild.GuildID][key] = value
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
