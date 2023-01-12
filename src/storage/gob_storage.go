package storage

import (
	"bytes"
	"encoding/gob"
	"io"
	"sync"

	"github.com/BKrajancic/boby/m/v2/src/service"
)

// GobStorage is an implementation of Storage, saving to a file using the Gob format.
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
func (g *GobStorage) GetGuildValue(guild service.Guild, key string) (interface{}, bool) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	return g.TempStorage.GetGuildValue(guild, key)
}

// SetGuildValue sets the value for key, for a Guild.
func (g *GobStorage) SetGuildValue(guild service.Guild, key string, value interface{}) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	err := g.TempStorage.SetGuildValue(guild, key, value)
	if err != nil {
		return err
	}
	return g.SaveToFile()
}

// SetDefaultGuildValue sets the default value for key, for all Guilds.
func (g *GobStorage) SetDefaultGuildValue(key string, value interface{}) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	err := g.TempStorage.SetDefaultGuildValue(key, value)
	if err != nil {
		return err
	}

	return g.SaveToFile()
}

// GetUserValue retrieves the value for key, for a User.
// Returns an error if the key doesn't exist or can't be retrieved.
func (g *GobStorage) GetUserValue(user service.User, key string) (interface{}, bool) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	return g.TempStorage.GetUserValue(user, key)
}

// SetUserValue sets the value for key, for a User.
func (g *GobStorage) SetUserValue(user service.User, key string, val interface{}) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	err := g.TempStorage.SetUserValue(user, key, val)
	if err != nil {
		return err
	}

	return g.SaveToFile()
}

// SetDefaultUserValue sets the default value for key, for all Users.
func (g *GobStorage) SetDefaultUserValue(key string, val interface{}) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	err := g.TempStorage.SetDefaultUserValue(key, val)
	if err != nil {
		return err
	}

	return g.SaveToFile()
}

// IsAdmin returns true if userID has been set using SetAdmin.
func (g *GobStorage) IsAdmin(guild service.Guild, userID string) bool {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	return g.TempStorage.IsAdmin(guild, userID)
}

// SetAdmin sets a userID as an admin for a guild.
func (g *GobStorage) SetAdmin(guild service.Guild, userID string) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	err := g.TempStorage.SetAdmin(guild, userID)
	if err != nil {
		return err
	}
	return g.SaveToFile()
}

// UnsetAdmin removes userID as an admin for a guild.
func (g *GobStorage) UnsetAdmin(guild service.Guild, userID string) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	err := g.TempStorage.UnsetAdmin(guild, userID)
	if err != nil {
		return err
	}

	return g.SaveToFile()
}

// SetGlobalValue sets a value that applies to globally.
func (g *GobStorage) SetGlobalValue(key string, value interface{}) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	err := g.TempStorage.SetGlobalValue(key, value)
	if err != nil {
		return err
	}

	return g.SaveToFile()
}

// GetGlobalValue sets a value that applies to globally.
func (g *GobStorage) GetGlobalValue(key string) (interface{}, bool) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	return g.TempStorage.GetGlobalValue(key)
}
