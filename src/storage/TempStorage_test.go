package storage

import (
	"testing"

	"github.com/BKrajancic/boby/m/v2/src/service"
)

func TestSetGetValue(t *testing.T) {
	storage := TempStorage{}
	guild := service.Guild{ServiceID: "0", GuildID: "0"}
	key := "key"
	value := "value"
	storage.SetGuildValue(guild, key, value)
	valueOut, err := storage.GetGuildValue(guild, key)
	if err != nil || valueOut != value {
		t.Fail()
	}
}

func TestGetValue(t *testing.T) {
	storage := TempStorage{}
	guild := service.Guild{ServiceID: "0", GuildID: "0"}
	key := "key"
	_, err := storage.GetGuildValue(guild, key)
	if err == nil {
		t.Fail()
	}
}

func TestGetValueMissingButHasService(t *testing.T) {
	storage := TempStorage{}
	guild := service.Guild{ServiceID: "0", GuildID: "0"}
	key1 := "key1"
	key2 := "key2"
	value := "value"
	storage.SetGuildValue(guild, key1, value)
	_, err := storage.GetGuildValue(guild, key2)
	if err == nil {
		t.Fail()
	}
}

func TestGetValueDifferentGuilds(t *testing.T) {
	storage := TempStorage{}
	serviceID := "0"
	guild1 := service.Guild{ServiceID: serviceID, GuildID: "0"}
	key1 := "key1"
	key2 := "key2"
	value := "value"
	storage.SetGuildValue(guild1, key1, value)
	guild2 := service.Guild{ServiceID: serviceID, GuildID: "1"}
	_, err := storage.GetGuildValue(guild2, key2)
	if err == nil {
		t.Fail()
	}
}

func TestSetAdmin(t *testing.T) {
	storage := TempStorage{}
	guild := service.Guild{ServiceID: "0", GuildID: "0"}
	user := "20"
	storage.SetAdmin(guild, user)
	if storage.IsAdmin(guild, user) == false {
		t.Fail()
	}
}

func TestUnsetAdmin(t *testing.T) {
	storage := TempStorage{}
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
	storage := TempStorage{}
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
	storage := TempStorage{}
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
	storage := TempStorage{}
	storage.SetAdmin(guild, "Test")
	storage.UnsetAdmin(guild, "Test")

	storage.GuildValues[guild.ServiceID][guild.GuildID]["Admin"] = 0
	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()

	storage.UnsetAdmin(guild, "Test")
}

func TestSetAdminDisaster(t *testing.T) {
	guild := service.Guild{ServiceID: "0", GuildID: "0"}
	storage := TempStorage{}
	storage.SetAdmin(guild, "Test")
	storage.UnsetAdmin(guild, "Test")

	storage.GuildValues[guild.ServiceID][guild.GuildID]["Admin"] = 0
	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()

	storage.SetAdmin(guild, "Test")
}
