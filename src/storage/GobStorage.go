package storage

import (
	"bytes"
	"encoding/gob"
	"io"
	"sync"

	"github.com/BKrajancic/boby/m/v2/src/service"
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

// GobStorage is an implementation of Storage that uses Gob and files for data storage.
type GobStorage struct {
	TempStorage TempStorage
	writer      TruncatableWriter
	mutex       *sync.Mutex // Lock when calling any public function.
}

// LoadFromBuffer will load a Gob from a filepath.
func LoadFromBuffer(t TruncatableWriter) (config GobStorage, err error) {
	enc := gob.NewDecoder(t)
	if err := enc.Decode(&config); err != nil {
		return config, err
	}

	config.mutex = &sync.Mutex{}
	config.writer = t

	config.TempStorage.mutex = &sync.Mutex{}
	return config, nil
}

// SetWriter sets the writer of a gobstorage to t.
func (g *GobStorage) SetWriter(t TruncatableWriter) {
	g.writer = t
}

// SaveToFile saves GobStorage's state to a file, which can be reloaded later using LoadFromFile.
func (j *GobStorage) SaveToFile() error {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)

	if err := enc.Encode(j); err != nil {
		return err
	}

	if err := j.writer.Truncate(0); err != nil {
		return err
	}

	if _, err := j.writer.Seek(0, io.SeekStart); err != nil {
		return err
	}

	if _, err := j.writer.Write(buffer.Bytes()); err != nil {
		return err
	}

	return j.writer.Sync()
}

// GetGuildValue retrieves the value for key, for a Guild.
// Returns an error if the key doesn't exist or can't be retrieved.
func (j *GobStorage) GetGuildValue(guild service.Guild, key string) (interface{}, error) {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	return j.TempStorage.GetGuildValue(guild, key)
}

// SetGuildValue sets the value for key, for a Guild.
func (j *GobStorage) SetGuildValue(guild service.Guild, key string, value interface{}) {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	j.TempStorage.SetGuildValue(guild, key, value)
	j.SaveToFile()
}

// GetUserValue retrieves the value for key, for a Guild.
// Returns an error if the key doesn't exist or can't be retrieved.
func (j *GobStorage) GetUserValue(user service.User, key string) (interface{}, error) {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	return j.TempStorage.GetUserValue(user, key)
}

// SetUserValue sets the value for key, for a Guild.
func (j *GobStorage) SetUserValue(user service.User, key string, val interface{}) {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	j.TempStorage.SetUserValue(user, key, val)
	j.SaveToFile()
}

// IsAdmin returns true if userID has been set using SetAdmin.
func (j *GobStorage) IsAdmin(guild service.Guild, userID string) bool {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	return j.TempStorage.IsAdmin(guild, userID)
}

// SetAdmin sets a userID as an admin for a guild.
func (j *GobStorage) SetAdmin(guild service.Guild, userID string) {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	j.TempStorage.SetAdmin(guild, userID)
	j.SaveToFile()
}

// UnsetAdmin removes userID as an admin for a guild.
func (j *GobStorage) UnsetAdmin(guild service.Guild, userID string) {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	j.TempStorage.UnsetAdmin(guild, userID)
	j.SaveToFile()
}
