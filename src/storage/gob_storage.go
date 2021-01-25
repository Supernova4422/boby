package storage

import (
	"bytes"
	"encoding/gob"
	"io"
	"sync"

	"github.com/BKrajancic/boby/m/v2/src/service"
)

// GobStorage is an implementation of Storage, saveing to a file using the Gob format.
type GobStorage struct {
	TempStorage TempStorage
	writer      TruncatableWriter
	mutex       *sync.Mutex // Lock when calling any public function.
}

// LoadFromBuffer will load a Gob from a buffer.
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
func (g *GobStorage) SaveToFile() error {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)

	if err := enc.Encode(g); err != nil {
		return err
	}

	if err := g.writer.Truncate(0); err != nil {
		return err
	}

	if _, err := g.writer.Seek(0, io.SeekStart); err != nil {
		return err
	}

	if _, err := g.writer.Write(buffer.Bytes()); err != nil {
		return err
	}

	return g.writer.Sync()
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
