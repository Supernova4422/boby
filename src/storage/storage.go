package storage

import "github.com/BKrajancic/FLD-Bot/m/v2/src/service"

// A Storage implementation can be used to store and load data.
// Each implementation must be thread safe, that is multiple threads can call multiple
// functions of a storage implementation and expect no concurrency issue to occur.
type Storage interface {
	GetValue(guild service.Guild, key string) (string, error)
	SetValue(guild service.Guild, key string, value string)
	IsAdmin(guild service.Guild, UserID string) bool
	SetAdmin(guild service.Guild, UserID string)
	UnsetAdmin(guild service.Guild, UserID string)
}
