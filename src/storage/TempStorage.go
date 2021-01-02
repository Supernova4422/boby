package storage

import (
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
)

type TempStorage struct {
	Admins map[string]map[string][]string
}

func (t *TempStorage) GetValue(guild service.Guild, key string)               {}
func (t *TempStorage) SetValue(guild service.Guild, key string, value string) {}

// IsAdmin returns true if ID is an admin.
func (t *TempStorage) IsAdmin(guild service.Guild, ID string) bool {
	serviceAdmins, ok := t.Admins[guild.ServiceId]
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

// SetAdmin saves an ID as an admin.
func (t *TempStorage) SetAdmin(guild service.Guild, ID string) {
	_, ok := t.Admins[guild.ServiceId]
	if ok {
		_, ok := t.Admins[guild.ServiceId][guild.GuildID]
		if ok {
			t.Admins[guild.ServiceId][guild.GuildID] = append(t.Admins[guild.ServiceId][guild.GuildID], ID)
		} else {
			t.Admins[guild.ServiceId] = make(map[string][]string)
		}
	} else {
		if t.Admins == nil {
			t.Admins = make(map[string]map[string][]string)
			t.Admins[guild.ServiceId] = make(map[string][]string)
			t.Admins[guild.ServiceId][guild.GuildID] = make([]string, 0)
		}
		t.Admins[guild.ServiceId][guild.GuildID] = append(t.Admins[guild.ServiceId][guild.GuildID], ID)
	}
}

// UnsetAdmin ensures ID can't be used as an admin.
func (t *TempStorage) UnsetAdmin(guild service.Guild, ID string) {
	_, ok := t.Admins[guild.ServiceId]
	if ok {
		_, ok := t.Admins[guild.ServiceId][guild.GuildID]
		if ok {
			newAdmins := []string{}
			for _, admin := range t.Admins[guild.ServiceId][guild.GuildID] {
				if admin != ID {
					newAdmins = append(newAdmins, admin)
				}
			}
			t.Admins[guild.ServiceId][guild.GuildID] = newAdmins
		}
	}
}
