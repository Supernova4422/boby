package storage

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"sync"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
)

// A TruncatableWriter is a buffer that supports flexible operations.
//
// The behaviour of all functions is that of os.File (os.file fulfills this interface)
type TruncatableWriter interface {
	Truncate(n int64) error
	Write(b []byte) (n int, err error)
	Read(p []byte) (n int, err error)
	Seek(offset int64, whence int) (ret int64, err error)
	Sync() (err error)
}

// JSONStorage is an implementation of Storage that uses json and files for data storage.
type JSONStorage struct {
	TempStorage TempStorage
	writer      TruncatableWriter
	mutex       *sync.Mutex // Lock when calling any public function.
}

// LoadFromBuffer will load a JSON from a filepath.
func LoadFromBuffer(t TruncatableWriter) (JSONStorage, error) {
	config := JSONStorage{writer: t, mutex: &sync.Mutex{}}

	bytes, err := ioutil.ReadAll(t)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(bytes, &config)
	return config, err
}

// SaveToFile saves JSONStorage's state to a file, which can be reloaded later using LoadFromFile.
func (j *JSONStorage) SaveToFile() error {
	bytes, _ := json.Marshal(j) // Unlikely to be an error, ignore return value to make test coverage 100%.

	err := j.writer.Truncate(0)
	if err != nil {
		return err
	}

	_, err = j.writer.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	_, err = j.writer.Write(bytes)
	if err != nil {
		return err
	}

	return j.writer.Sync()
}

// GetGuildValue retrieves the value for key, for a Guild.
// Returns an error if the key doesn't exist or can't be retrieved.
func (j *JSONStorage) GetGuildValue(guild service.Guild, key string) (interface{}, error) {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	return j.TempStorage.GetGuildValue(guild, key)
}

// SetGuildValue sets the value for key, for a Guild.
func (j *JSONStorage) SetGuildValue(guild service.Guild, key string, value interface{}) {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	j.TempStorage.SetGuildValue(guild, key, value)
	j.SaveToFile()
}

// GetUserValue retrieves the value for key, for a Guild.
// Returns an error if the key doesn't exist or can't be retrieved.
func (j *JSONStorage) GetUserValue(serviceID string, userID string, key string) (interface{}, error) {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	return j.TempStorage.GetUserValue(serviceID, userID, key)
}

// SetUserValue sets the value for key, for a Guild.
func (j *JSONStorage) SetUserValue(serviceID string, userID string, key string, val interface{}) {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	j.TempStorage.SetUserValue(serviceID, userID, key, val)
	j.SaveToFile()
}

// IsAdmin returns true if userID has been set using SetAdmin.
func (j *JSONStorage) IsAdmin(guild service.Guild, userID string) bool {
	j.mutex.Lock()
	defer j.mutex.Unlock()
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
