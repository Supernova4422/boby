package command

import (
	"strconv"
	"time"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/storage"
)

type RateLimitConfig struct {
	TimesPerInterval   int    // How many times can it be used per interval?
	SecondsPerInterval int    // How long is an interval, in seconds.
	Body               string // Reply when limit is reached.
}

// GetRateLimitedCommand wraps around a command to make it rate limited.
func (r RateLimitConfig) GetRateLimitedCommand(command Command, commandDelayID string) Command {
	curry := func(sender service.Conversation, user service.User, msg [][]string, storage *storage.Storage, sink func(service.Conversation, service.Message)) {
		val, err := (*storage).GetUserValue(sender.ServiceID, user.ID, commandDelayID)
		if err != nil {
			now := time.Now().Unix()
			(*storage).SetUserValue(sender.ServiceID, user.ID, commandDelayID, string(now))
		} else {
			previous, err := strconv.Atoi(val)
			if err == nil {
				previousUse := time.Unix(int64(previous), 0)
			}
		}
	}
}
