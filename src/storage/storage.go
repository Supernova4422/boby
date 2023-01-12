package storage

import "github.com/BKrajancic/boby/m/v2/src/service"

// A Storage implementation can be used to store and load data.
// Each implementation must be thread safe, that is multiple threads can call multiple
// functions of a storage implementation and expect no concurrency issue to occur.
type Storage interface {
	GetGuildValue(guild service.Guild, key string) (interface{}, bool)
	SetGuildValue(guild service.Guild, key string, value interface{}) error
	SetDefaultGuildValue(key string, value interface{}) error

	GetUserValue(user service.User, key string) (interface{}, bool)
	SetUserValue(user service.User, key string, value interface{}) error
	SetDefaultUserValue(key string, value interface{}) error

	IsAdmin(guild service.Guild, UserID string) bool
	SetAdmin(guild service.Guild, UserID string) error
	UnsetAdmin(guild service.Guild, UserID string) error

	SetGlobalValue(key string, value interface{}) error
	GetGlobalValue(key string) (interface{}, bool)
}
