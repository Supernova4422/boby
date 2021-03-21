package storage

import (
	"bytes"
	"fmt"
	"sync"
	"testing"

	"github.com/BKrajancic/boby/m/v2/src/service"
)

// A TruncatableBuffer wraps a bytes.buffer object to fulfil the interface of TruncatableWriter.
type TruncatableBuffer struct {
	buffer *bytes.Buffer
}

func (t TruncatableBuffer) Truncate(n int64) error {
	t.buffer.Truncate(int(n))
	return nil
}

func (t TruncatableBuffer) Read(p []byte) (int, error) {
	return t.buffer.Read(p)
}

func (t TruncatableBuffer) Write(b []byte) (n int, err error) {
	return t.buffer.Write(b)
}

func (t TruncatableBuffer) Sync() error {
	return nil
}

// Seek is supposed to set the offset for following read/write, but here it does nothing.
//
// In GobStorage's usage, we expect that only seek(0, 0) is ever used, this is how
// a byte buffer typically works, so we just do nothing.
func (t TruncatableBuffer) Seek(offset int64, whence int) (n int64, err error) {
	if offset != 0 || whence != 0 {
		return 0, fmt.Errorf("byte buffers have no use for a seek method")
	}

	return 0, nil
}

// A TruncatableBufferError returns errors for each function.
type TruncatableBufferErrorProne struct{}

func (t TruncatableBufferErrorProne) Truncate(n int64) error {
	return fmt.Errorf("expecting an error")
}

func (t TruncatableBufferErrorProne) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("expecting an error")
}

func (t TruncatableBufferErrorProne) Write(b []byte) (n int, err error) {
	return 0, fmt.Errorf("expecting an error")
}

func (TruncatableBufferErrorProne) Seek(offset int64, whence int) (n int64, err error) {
	return 0, fmt.Errorf("expecting an error")
}

func (TruncatableBufferErrorProne) Sync() (err error) {
	return fmt.Errorf("expecting an error")
}

func TestGobSetGetGuildValue(t *testing.T) {
	bytesOut := bytes.NewBuffer([]byte{})
	writer := TruncatableBuffer{bytesOut}
	storage := GobStorage{
		TempStorage: GetTempStorage(),
		writer:      writer,
		mutex:       &sync.Mutex{},
	}

	guild := service.Guild{ServiceID: "0", GuildID: "0"}
	k0 := "k0"
	v0 := "v0"
	k1 := "k1"
	v1 := "v1"
	k2 := "k2"
	v2 := "v2"

	storage.SetGuildValue(guild, k0, v0)
	storage.SetGuildValue(guild, k1, v1)
	storage.SetGuildValue(guild, k2, v2)
	storage, err := LoadFromBuffer(writer)
	if err != nil {
		t.Fail()
	}

	valueOut, ok := storage.GetGuildValue(guild, k0)
	if ok == false || valueOut != v0 {
		t.Fail()
	}

	valueOut, ok = storage.GetGuildValue(guild, k1)
	if ok == false || valueOut != v1 {
		t.Fail()
	}

	valueOut, ok = storage.GetGuildValue(guild, k2)
	if ok == false || valueOut != v2 {
		t.Fail()
	}
}

func TestGobSetGetUserValue(t *testing.T) {
	bytesOut := bytes.NewBuffer([]byte{})
	writer := TruncatableBuffer{bytesOut}
	storage := GobStorage{
		TempStorage: GetTempStorage(),
		writer:      writer,
		mutex:       &sync.Mutex{},
	}

	user := service.User{ServiceID: "0", Name: "0"}
	k0 := "k0"
	v0 := "v0"
	k1 := "k1"
	v1 := "v1"
	k2 := "k2"
	v2 := "v2"

	storage.SetUserValue(user, k0, v0)
	storage.SetUserValue(user, k1, v1)
	storage.SetUserValue(user, k2, v2)
	storage, err := LoadFromBuffer(writer)
	if err != nil {
		t.Fail()
	}

	valueOut, ok := storage.GetUserValue(user, k0)
	if ok == false || valueOut != v0 {
		t.Fail()
	}

	valueOut, ok = storage.GetUserValue(user, k1)
	if ok == false || valueOut != v1 {
		t.Fail()
	}

	valueOut, ok = storage.GetUserValue(user, k2)
	if ok == false || valueOut != v2 {
		t.Fail()
	}
}

func TestGobSetGetDefaultUserValue(t *testing.T) {
	bytesOut := bytes.NewBuffer([]byte{})
	writer := TruncatableBuffer{bytesOut}
	storage := GobStorage{
		TempStorage: GetTempStorage(),
		writer:      writer,
		mutex:       &sync.Mutex{},
	}

	k0 := "k0"
	v0 := "v0"

	storage.SetDefaultUserValue(k0, v0)
	storage, err := LoadFromBuffer(writer)
	if err != nil {
		t.Fail()
	}

	valueOut, ok := storage.GetUserValue(
		service.User{ServiceID: "0", Name: "0"},
		k0,
	)
	if ok == false || valueOut != v0 {
		t.Fail()
	}
}

func TestGobSetGetDefaultGuildValue(t *testing.T) {
	bytesOut := bytes.NewBuffer([]byte{})
	writer := TruncatableBuffer{bytesOut}
	storage := GobStorage{
		TempStorage: GetTempStorage(),
		writer:      writer,
		mutex:       &sync.Mutex{},
	}

	k0 := "k0"
	v0 := "v0"
	storage.SetDefaultGuildValue(k0, v0)
	storage, err := LoadFromBuffer(writer)

	if err != nil {
		t.Fail()
	}

	valueOut, ok := storage.GetGuildValue(
		service.Guild{ServiceID: "0", GuildID: "0"},
		k0,
	)
	if ok == false || valueOut != v0 {
		t.Fail()
	}
}
func TestGobGetUserValue(t *testing.T) {
	bytesOut := bytes.NewBuffer([]byte{})
	writer := TruncatableBuffer{bytesOut}
	storage := GobStorage{
		TempStorage: GetTempStorage(),
		writer:      writer,
		mutex:       &sync.Mutex{},
	}

	guild := service.Guild{ServiceID: "0", GuildID: "0"}
	key := "key"
	_, ok := storage.GetGuildValue(guild, key)
	if ok {
		t.Fail()
	}
}

func TestGobGetValueMissingButHasService(t *testing.T) {
	bytesOut := bytes.NewBuffer([]byte{})
	writer := TruncatableBuffer{bytesOut}
	storage := GobStorage{
		TempStorage: GetTempStorage(),
		writer:      writer,
		mutex:       &sync.Mutex{},
	}

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

func TestGobGetValueDifferentGuilds(t *testing.T) {
	bytesOut := bytes.NewBuffer([]byte{})
	writer := TruncatableBuffer{bytesOut}
	storage := GobStorage{
		TempStorage: GetTempStorage(),
		writer:      writer,
		mutex:       &sync.Mutex{},
	}

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

func TestGobSetAdmin(t *testing.T) {
	bytesOut := bytes.NewBuffer([]byte{})
	writer := TruncatableBuffer{bytesOut}
	storage := GobStorage{
		TempStorage: GetTempStorage(),
		writer:      writer,
		mutex:       &sync.Mutex{},
	}

	guild := service.Guild{ServiceID: "0", GuildID: "0"}
	user := "20"
	storage.SetAdmin(guild, user)
	if storage.IsAdmin(guild, user) == false {
		t.Fail()
	}
}

func TestGobUnsetAdmin(t *testing.T) {
	bytesOut := bytes.NewBuffer([]byte{})
	writer := TruncatableBuffer{bytesOut}
	storage := GobStorage{
		TempStorage: GetTempStorage(),
		writer:      writer,
		mutex:       &sync.Mutex{},
	}

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

func TestGobUnsetAdminWhenMultipleAdmins(t *testing.T) {
	bytesOut := bytes.NewBuffer([]byte{})
	writer := TruncatableBuffer{bytesOut}
	storage := GobStorage{
		TempStorage: GetTempStorage(),
		writer:      writer,
		mutex:       &sync.Mutex{},
	}

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

func TestGobSetAdminDifferentGuilds(t *testing.T) {
	bytesOut := bytes.NewBuffer([]byte{})
	writer := TruncatableBuffer{bytesOut}
	storage := GobStorage{
		TempStorage: GetTempStorage(),
		writer:      writer,
		mutex:       &sync.Mutex{},
	}

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

func TestBadLoad(t *testing.T) {
	writer := TruncatableBufferErrorProne{}
	_, err := LoadFromBuffer(writer)
	if err == nil {
		t.Fail()
	}
}

func TestBadSave(t *testing.T) {
	storage := GobStorage{
		TempStorage: GetTempStorage(),
		writer:      TruncatableBufferErrorProne{},
		mutex:       &sync.Mutex{},
	}

	if storage.SaveToFile() == nil {
		t.Fail()
	}
}

// A TruncatableBufferError returns errors for each function.
type TruncatableBufferErrorOnWrite struct{}

func (t TruncatableBufferErrorOnWrite) Truncate(n int64) error {
	return nil
}

func (t TruncatableBufferErrorOnWrite) Read(p []byte) (int, error) {
	return 0, nil
}

func (t TruncatableBufferErrorOnWrite) Write(b []byte) (n int, err error) {
	return 0, fmt.Errorf("expecting an error")
}

func (TruncatableBufferErrorOnWrite) Seek(offset int64, whence int) (n int64, err error) {
	return 0, nil
}

func (TruncatableBufferErrorOnWrite) Sync() (err error) {
	return nil
}

func TestBadWrite(t *testing.T) {
	storage := GobStorage{
		TempStorage: GetTempStorage(),
		writer:      TruncatableBufferErrorOnWrite{},
		mutex:       &sync.Mutex{},
	}

	if storage.SaveToFile() == nil {
		t.Fail()
	}
}

// A TruncatableBufferOnSeek returns errors for each function.
type TruncatableBufferErrorOnSeek struct{}

func (t TruncatableBufferErrorOnSeek) Truncate(n int64) error {
	return nil
}

func (t TruncatableBufferErrorOnSeek) Read(p []byte) (int, error) {
	return 0, nil
}

func (t TruncatableBufferErrorOnSeek) Write(b []byte) (n int, err error) {
	return 0, nil
}

func (TruncatableBufferErrorOnSeek) Seek(offset int64, whence int) (n int64, err error) {
	return 0, fmt.Errorf("expecting an error")
}

func (TruncatableBufferErrorOnSeek) Sync() (err error) {
	return nil
}

func TestBadSeek(t *testing.T) {
	storage := GobStorage{
		TempStorage: GetTempStorage(),
		writer:      TruncatableBufferErrorOnSeek{},
		mutex:       &sync.Mutex{},
	}

	if storage.SaveToFile() == nil {
		t.Fail()
	}
}

func TestGobTypeChange(t *testing.T) {
	bytesOut := bytes.NewBuffer([]byte{})
	writer := TruncatableBuffer{bytesOut}
	storage := GobStorage{
		TempStorage: GetTempStorage(),
		writer:      writer,
		mutex:       &sync.Mutex{},
	}

	user := service.User{ServiceID: "0", Name: "0"}
	k0 := "k0"
	var v0 int64 = 1000000000000000000
	storage.SetUserValue(user, k0, v0)
	storage, err := LoadFromBuffer(writer)
	if err != nil {
		t.Fail()
	}

	valueOut, _ := storage.GetUserValue(user, k0)
	switch valueOut.(type) {
	case int64:
		break
	case float64:
		t.Fail() // Usually happens.
	default:
		t.Fail()
	}
}
