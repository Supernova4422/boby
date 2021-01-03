package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
)

type JsonStorage struct {
	TempStorage TempStorage
	filepath    string
	mutex       *sync.Mutex // Lock when calling any public function.
}

func LoadFromFile(filepath string) (JsonStorage, error) {
	config := JsonStorage{
		filepath: filepath,
	}
	config.mutex.Unlock()
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

func (j *JsonStorage) GetValue(guild service.Guild, key string) (string, error) {
	return j.TempStorage.GetValue(guild, key)
}
func (j *JsonStorage) SetValue(guild service.Guild, key string, value string) {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	j.TempStorage.SetValue(guild, key, value)
	j.SaveToFile()
}
func (j *JsonStorage) IsAdmin(guild service.Guild, userID string) bool {
	return j.TempStorage.IsAdmin(guild, userID)
}
func (j *JsonStorage) SetAdmin(guild service.Guild, userID string) {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	j.TempStorage.SetAdmin(guild, userID)
	j.SaveToFile()
}
func (j *JsonStorage) UnsetAdmin(guild service.Guild, userID string) {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	j.TempStorage.UnsetAdmin(guild, userID)
	j.SaveToFile()
}
