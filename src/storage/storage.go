package storage

import "github.com/BKrajancic/FLD-Bot/m/v2/src/service"

type Storage interface {
	GetValue(guild service.Guild, key string)
	SetValue(guild service.Guild, key string, value string)
	IsAdmin(guild service.Guild, UserID string) bool
	SetAdmin(guild service.Guild, UserID string)
	UnsetAdmin(guild service.Guild, UserID string)
}
