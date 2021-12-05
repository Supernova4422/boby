package storage

import "github.com/BKrajancic/boby/m/v2/src/service"

// A Storage implementation can be used to store and load data.
// Each implementation must be thread safe, that is multiple threads can call multiple
// functions of a storage implementation and expect no concurrency issue to occur.
type Storage interface {
	GetGuildValue(guild service.Guild, key string) (interface{}, bool)
	SetGuildValue(guild service.Guild, key string, value interface{})
	SetDefaultGuildValue(key string, value interface{})

	GetUserValue(user service.User, key string) (interface{}, bool)
	SetUserValue(user service.User, key string, value interface{})
	SetDefaultUserValue(key string, value interface{})

	IsAdmin(guild service.Guild, UserID string) bool
	SetAdmin(guild service.Guild, UserID string)
	UnsetAdmin(guild service.Guild, UserID string)

	SetGlobalValue(key string, value interface{})
	GetGlobalValue(key string) (interface{}, bool)
}
