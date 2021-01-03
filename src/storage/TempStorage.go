package storage

import (
	"fmt"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
)

type TempStorage struct {
	Admins map[string]map[string][]string
	Values map[string]map[string]map[string]string
}

func (t *TempStorage) GetValue(guild service.Guild, key string) (string, error) {
	_, ok := t.Values[guild.ServiceId]
	if ok == false {
		return "", fmt.Errorf("No entries for ServiceID %s", guild.ServiceId)
	}

	_, ok = t.Values[guild.ServiceId][guild.GuildID]
	if ok == false {
		return "", fmt.Errorf("No entries for GuildID %s", guild.GuildID)
	}

	val, ok := t.Values[guild.ServiceId][guild.GuildID][key]
	if ok == false {
		return "", fmt.Errorf("No value for key %s", key)
	}
	return val, nil
}

func (t *TempStorage) SetValue(guild service.Guild, key string, value string) {
	_, ok := t.Values[guild.ServiceId]
	if ok == false {
		t.Values = make(map[string]map[string]map[string]string)
		t.Values[guild.ServiceId] = make(map[string]map[string]string)
	}

	_, ok = t.Values[guild.ServiceId][guild.GuildID]
	if ok == false {
		t.Values[guild.ServiceId][guild.GuildID] = make(map[string]string)
	}

	t.Values[guild.ServiceId][guild.GuildID][key] = value
}

// IsAdmin returns true if ID is an admin.
func (j *TempStorage) IsAdmin(guild service.Guild, ID string) bool {
	serviceAdmins, ok := j.Admins[guild.ServiceId]
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
func (j *TempStorage) SetAdmin(guild service.Guild, ID string) {
	_, ok := j.Admins[guild.ServiceId]
	if ok {
		_, ok := j.Admins[guild.ServiceId][guild.GuildID]
		if ok {
			j.Admins[guild.ServiceId][guild.GuildID] = append(j.Admins[guild.ServiceId][guild.GuildID], ID)
		} else {
			j.Admins[guild.ServiceId] = make(map[string][]string)
		}
	} else {
		if j.Admins == nil {
			j.Admins = make(map[string]map[string][]string)
			j.Admins[guild.ServiceId] = make(map[string][]string)
		}
		j.Admins[guild.ServiceId][guild.GuildID] = []string{ID}
	}
}

// UnsetAdmin ensures ID can't be used as an admin.
func (j *TempStorage) UnsetAdmin(guild service.Guild, ID string) {
	_, ok := j.Admins[guild.ServiceId]
	if ok {
		_, ok := j.Admins[guild.ServiceId][guild.GuildID]
		if ok {
			newAdmins := []string{}
			for _, admin := range j.Admins[guild.ServiceId][guild.GuildID] {
				if admin != ID {
					newAdmins = append(newAdmins, admin)
				}
			}
			j.Admins[guild.ServiceId][guild.GuildID] = newAdmins
		}
	}
}
