package storage

import (
	"sync"

	"github.com/BKrajancic/boby/m/v2/src/service"
)

// AdminKey is the key used in storage to store if someone is an Admin.
const AdminKey = "Admin"

// TempStorage implements the Storage interface, but data is lost on destruction.
type TempStorage struct {
	UserValues         map[string]map[string]map[string]interface{}
	DefaultUserValues  map[string]interface{}
	GuildValues        map[string]map[string]map[string]interface{}
	DefaultGuildValues map[string]interface{}
	GlobalValues       map[string]interface{}
	mutex              *sync.Mutex
}

// GetTempStorage returns a TempStorage.
func GetTempStorage() TempStorage {
	return TempStorage{mutex: &sync.Mutex{}}
}

// GetGuildValue retrieves the value for key, for a Guild.
// Returns an error if the key doesn't exist or can't be retrieved.
func (t *TempStorage) GetGuildValue(guild service.Guild, key string) (interface{}, bool) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	serviceMap, ok := t.GuildValues[guild.ServiceID]
	if ok == false {
		val, ok := t.DefaultGuildValues[key]
		return val, ok
	}

	keyMap, ok := serviceMap[guild.GuildID]
	if ok == false {
		val, ok := t.DefaultGuildValues[key]
		return val, ok
	}

	val, ok := keyMap[key]
	if ok == false {
		val, ok := t.DefaultGuildValues[key]
		return val, ok
	}

	return val, ok
}

// SetGuildValue sets the value for key, for a Guild.
func (t *TempStorage) SetGuildValue(guild service.Guild, key string, val interface{}) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

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

// SetDefaultGuildValue sets the default value for key, for all Guilds.
func (t *TempStorage) SetDefaultGuildValue(key string, val interface{}) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if t.DefaultGuildValues == nil {
		t.DefaultGuildValues = make(map[string]interface{})
	}

	t.DefaultGuildValues[key] = val
}

// SetDefaultUserValue sets the default value for key, for all Users.
func (t *TempStorage) SetDefaultUserValue(key string, val interface{}) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if t.DefaultUserValues == nil {
		t.DefaultUserValues = make(map[string]interface{})
	}

	t.DefaultUserValues[key] = val
}

// GetUserValue retrieves the value for key, for a User.
// Returns an error if the key doesn't exist or can't be retrieved.
func (t *TempStorage) GetUserValue(user service.User, key string) (interface{}, bool) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	serviceMap, ok := t.UserValues[user.ServiceID]
	if !ok {
		val, ok := t.DefaultUserValues[key]
		return val, ok
	}

	keyMap, ok := serviceMap[user.Name]
	if !ok {
		val, ok := t.DefaultUserValues[key]
		return val, ok
	}

	val, ok := keyMap[key]
	if !ok {
		val, ok := t.DefaultUserValues[key]
		return val, ok
	}

	return val, ok
}

// SetUserValue sets the value for key, for a Guild.
func (t *TempStorage) SetUserValue(user service.User, key string, val interface{}) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

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
	if val, ok := t.GetGuildValue(guild, AdminKey); ok {
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
	if val, ok := t.GetGuildValue(guild, AdminKey); ok {
		currentAdmins, ok := val.([]string)
		if ok {
			t.SetGuildValue(guild, AdminKey, append(currentAdmins, ID))
		} else {
			panic(ok)
		}
	} else {
		t.SetGuildValue(guild, AdminKey, []string{ID})
	}
}

// UnsetAdmin removes userID as an admin for a guild.
func (t *TempStorage) UnsetAdmin(guild service.Guild, ID string) {
	newAdmins := []string{}
	if val, ok := t.GetGuildValue(guild, AdminKey); ok {
		currentAdmins, ok := val.([]string)
		if !ok {
			panic(ok)
		}
		for _, adminID := range currentAdmins {
			if adminID != ID {
				newAdmins = append(newAdmins, adminID)
			}
		}
	}
	t.SetGuildValue(guild, AdminKey, newAdmins)
}

// SetGlobalValue sets a value that applies to globally.
func (t *TempStorage) SetGlobalValue(key string, value interface{}) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if t.GlobalValues == nil {
		t.GlobalValues = make(map[string]interface{})
	}

	t.GlobalValues[key] = value
}

// GetGlobalValue sets a value that applies to globally.
func (t *TempStorage) GetGlobalValue(key string) (interface{}, bool) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if t.GlobalValues == nil {
		t.GlobalValues = make(map[string]interface{})
	}

	keyMap, ok := t.GlobalValues[key]
	return keyMap, ok
}
