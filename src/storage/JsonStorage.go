package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
)

// JSONStorage is an implementation of Storage that uses json and files for data storage.
type JSONStorage struct {
	TempStorage TempStorage
	filepath    string
	mutex       *sync.Mutex // Lock when calling any public function.
}

// LoadFromFile will load a JSON from a filepath.
func LoadFromFile(filepath string) (JSONStorage, error) {
	config := JSONStorage{
		filepath: filepath,
		mutex:    &sync.Mutex{},
	}

	// Make an empty file.
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		err := config.SaveToFile()
		if err != nil {
			return config, fmt.Errorf("There was an error making an example file")
		}
	}

	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Unable to read file: %s", filepath)
		return config, nil
	}

	err = json.Unmarshal(bytes, &config)
	return config, err
}

// SaveToFile saves JSONStorage's state to a file, which can be reloaded later using LoadFromFile.
func (j *JSONStorage) SaveToFile() error {
	bytes, err := json.Marshal(j)

	file, err := os.Create(j.filepath)
	if err != nil {
		return fmt.Errorf("Unable to create file: %s", j.filepath)
	}
	defer file.Close()

	_, err = file.Write(bytes)
	if err != nil {
		return fmt.Errorf("Unable to write to file: %s", j.filepath)
	}
	return nil
}

// GetValue retrieves the value for key, for a Guild.
// Returns an error if the key doesn't exist or can't be retrieved.
func (j *JSONStorage) GetValue(guild service.Guild, key string) (string, error) {
	return j.TempStorage.GetValue(guild, key)
}

// SetValue sets the value for key, for a Guild.
func (j *JSONStorage) SetValue(guild service.Guild, key string, value string) {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	j.TempStorage.SetValue(guild, key, value)
	j.SaveToFile()
}

// IsAdmin returns true if userID has been set using SetAdmin.
func (j *JSONStorage) IsAdmin(guild service.Guild, userID string) bool {
	return j.TempStorage.IsAdmin(guild, userID)
}

// SetAdmin sets a userID as an admin for a guild.
func (j *JSONStorage) SetAdmin(guild service.Guild, userID string) {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	j.TempStorage.SetAdmin(guild, userID)
	j.SaveToFile()
}

// UnsetAdmin removes userID as an admin for a guild.
func (j *JSONStorage) UnsetAdmin(guild service.Guild, userID string) {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	j.TempStorage.UnsetAdmin(guild, userID)
	j.SaveToFile()
}
