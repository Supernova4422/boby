package command

import (
	"errors"

	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/storage"
)

// CreateError will return an error.
func CreateError(sender service.Conversation, user service.User, msg []interface{}, storage *storage.Storage, sink func(service.Conversation, service.Message) error) error {
	return errors.New("error created for testing purposes")
}
