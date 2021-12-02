package storage

import (
	"sync"
	"testing"

	"github.com/BKrajancic/boby/m/v2/src/service"
)

func TestSetGetValue(t *testing.T) {
	storage := TempStorage{mutex: &sync.Mutex{}}
	guild := service.Guild{ServiceID: "0", GuildID: "0"}
	key := "key"
	value := "value"
	storage.SetGuildValue(guild, key, value)
	valueOut, ok := storage.GetGuildValue(guild, key)
	valueStr := valueOut.(string)
	if ok == false || valueStr != value {
		t.Fail()
	}
}

func TestGetDefaultValue(t *testing.T) {
	storage := TempStorage{mutex: &sync.Mutex{}}
	key := "key"
	value := "value"
	storage.SetDefaultGuildValue(key, value)

	valueOut, ok := storage.GetGuildValue(
		service.Guild{ServiceID: "0", GuildID: "0"},
		key,
	)

	if ok == false || valueOut != value {
		t.Fail()
	}
}

func TestGetDefaultUserValue(t *testing.T) {
	storage := TempStorage{mutex: &sync.Mutex{}}
	key := "key"
	value := "value"
	storage.SetDefaultUserValue(key, value)
	storage.SetUserValue(service.User{ServiceID: "0", Name: "1"}, "key", "value")

	valueOut, ok := storage.GetUserValue(
		service.User{ServiceID: "0", Name: "0"},
		key,
	)

	if ok == false || valueOut != value {
		t.Fail()
	}
}

func TestGetDefaultUserValue2(t *testing.T) {
	storage := TempStorage{mutex: &sync.Mutex{}}
	key := "key"
	value := "value"
	storage.SetDefaultUserValue(key, value)
	user := service.User{ServiceID: "0", Name: "1"}
	storage.SetUserValue(user, "key2", "value")

	valueOut, ok := storage.GetUserValue(user, key)

	if ok == false || valueOut != value {
		t.Fail()
	}
}

func TestGetValue(t *testing.T) {
	storage := TempStorage{mutex: &sync.Mutex{}}
	guild := service.Guild{ServiceID: "0", GuildID: "0"}
	key := "key"
	_, ok := storage.GetGuildValue(guild, key)
	if ok == true {
		t.Fail()
	}
}

func TestGetValueMissingButHasService(t *testing.T) {
	storage := TempStorage{mutex: &sync.Mutex{}}
	guild := service.Guild{ServiceID: "0", GuildID: "0"}
	key1 := "key1"
	key2 := "key2"
	value := "value"
	storage.SetGuildValue(guild, key1, value)
	_, ok := storage.GetGuildValue(guild, key2)
	if ok {
		t.Fail()
	}
}

func TestGetValueDifferentGuilds(t *testing.T) {
	storage := TempStorage{mutex: &sync.Mutex{}}
	serviceID := "0"
	guild1 := service.Guild{ServiceID: serviceID, GuildID: "0"}
	key1 := "key1"
	key2 := "key2"
	value := "value"
	storage.SetGuildValue(guild1, key1, value)
	guild2 := service.Guild{ServiceID: serviceID, GuildID: "1"}
	_, ok := storage.GetGuildValue(guild2, key2)
	if ok {
		t.Fail()
	}
}

func TestSetAdmin(t *testing.T) {
	storage := TempStorage{mutex: &sync.Mutex{}}
	guild := service.Guild{ServiceID: "0", GuildID: "0"}
	user := "20"
	storage.SetAdmin(guild, user)
	if storage.IsAdmin(guild, user) == false {
		t.Fail()
	}
}

func TestUnsetAdmin(t *testing.T) {
	storage := TempStorage{mutex: &sync.Mutex{}}
	guild := service.Guild{ServiceID: "0", GuildID: "0"}
	user := "20"
	storage.SetAdmin(guild, user)
	if storage.IsAdmin(guild, user) == false {
		t.Fail()
	}

	storage.UnsetAdmin(guild, user)
	if storage.IsAdmin(guild, user) {
		t.Fail()
	}
}

func TestUnsetAdminWhenMultipleAdmins(t *testing.T) {
	storage := TempStorage{mutex: &sync.Mutex{}}
	guild := service.Guild{ServiceID: "0", GuildID: "0"}
	user1 := "20"
	user2 := "21"
	storage.SetAdmin(guild, user1)
	storage.SetAdmin(guild, user2)
	if storage.IsAdmin(guild, user1) == false || storage.IsAdmin(guild, user2) == false {
		t.Fail()
	}

	storage.UnsetAdmin(guild, user1)

	if storage.IsAdmin(guild, user1) {
		t.Fail()
	}
	if storage.IsAdmin(guild, user2) == false {
		t.Fail()
	}
}

func TestSetAdminDifferentGuilds(t *testing.T) {
	storage := TempStorage{mutex: &sync.Mutex{}}
	serviceID := "0"
	guild1 := service.Guild{ServiceID: serviceID, GuildID: "0"}
	guild2 := service.Guild{ServiceID: serviceID, GuildID: "1"}
	user := "20"

	storage.SetAdmin(guild1, user)
	if storage.IsAdmin(guild1, user) == false {
		t.Fail()
	}
	if storage.IsAdmin(guild2, user) {
		t.Fail()
	}

	storage.SetAdmin(guild2, user)
	if storage.IsAdmin(guild2, user) == false {
		t.Fail()
	}
}

func TestUnsetAdminDisaster(t *testing.T) {
	guild := service.Guild{ServiceID: "0", GuildID: "0"}
	storage := TempStorage{mutex: &sync.Mutex{}}
	storage.SetAdmin(guild, "Test")
	storage.UnsetAdmin(guild, "Test")

	storage.GuildValues[guild.ServiceID][guild.GuildID][AdminKey] = 0
	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()

	storage.UnsetAdmin(guild, "Test")
}

func TestSetAdminDisaster(t *testing.T) {
	guild := service.Guild{ServiceID: "0", GuildID: "0"}
	storage := TempStorage{mutex: &sync.Mutex{}}
	storage.SetAdmin(guild, "Test")
	storage.UnsetAdmin(guild, "Test")

	storage.GuildValues[guild.ServiceID][guild.GuildID][AdminKey] = 0
	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()

	storage.SetAdmin(guild, "Test")
}

func TestGetGlobalValueBeforeSet(t *testing.T) {
	storage := TempStorage{mutex: &sync.Mutex{}}
	_, ok := storage.GetGlobalValue("FAIL")
	if ok {
		t.Fail()
	}
}

func TestGetGlobalValueSetandGet(t *testing.T) {
	storage := TempStorage{mutex: &sync.Mutex{}}
	key := "key"
	value := "value"
	storage.SetGlobalValue(key, value)
	result, ok := storage.GetGlobalValue(key)
	if !ok {
		t.Fail()
	}

	if value != result {
		t.Fail()
	}
}
