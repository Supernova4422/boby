package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
)

type JsonStorage struct {
	Admins   map[string]map[string][]string
	filepath string
}

func (j *JsonStorage) GetValue(guild service.Guild, key string)               {}
func (j *JsonStorage) SetValue(guild service.Guild, key string, value string) {}

// IsAdmin returns true if ID is an admin.
func (j *JsonStorage) IsAdmin(guild service.Guild, ID string) bool {
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
func (j *JsonStorage) SetAdmin(guild service.Guild, ID string) {
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
	j.SaveToFile()
}

// UnsetAdmin ensures ID can't be used as an admin.
func (j *JsonStorage) UnsetAdmin(guild service.Guild, ID string) {
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
	j.SaveToFile()
}

func (j *JsonStorage) SaveToFile() error {
	filepath := "storage.json"
	bytes, err := json.Marshal(j)

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("Unable to create file: %s", filepath)
	}
	defer file.Close()

	_, err = file.Write(bytes)
	if err != nil {
		return fmt.Errorf("Unable to write to file: %s", filepath)
	}
	return nil
}

func LoadFromFile(filepath string) (JsonStorage, error) {
	config := JsonStorage{filepath: filepath}
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		err := config.SaveToFile()
		if err != nil {
			return config, fmt.Errorf("There was an error making an example file")
		}

		return config, fmt.Errorf("File not found")
	}

	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Unable to read file: %s", filepath)
		return config, nil
	}

	err = json.Unmarshal(bytes, &config)
	return config, err
}
