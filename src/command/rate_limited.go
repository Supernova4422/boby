package command

import (
	"fmt"
	"sort"
	"time"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/storage"
)

type RateLimitConfig struct {
	TimesPerInterval   int    // How many times can it be used per interval?
	SecondsPerInterval int64  // How long is an interval, in seconds.
	Body               string // Reply when limit is reached.
}

// rateLimited returns true if a message should be rate limited.
// now and History are expected to be in unix time.
func (r RateLimitConfig) rateLimited(now int64, history []int64) bool {
	return r.timeRemaining(now, r.cleanHistory(now, history)) > 0
}

// cleanHistory will remove items from the history that are outside the interval.
func (r RateLimitConfig) cleanHistory(now int64, history []int64) []int64 {
	newHistory := make([]int64, 0)
	cutoff := now - r.SecondsPerInterval
	for _, timestamp := range history {
		if timestamp > cutoff {
			newHistory = append(newHistory, timestamp)
		}
	}

	return newHistory
}

// timeRemaining returns the number of seconds until rate limiting is off.
func (r RateLimitConfig) timeRemaining(now int64, history []int64) int64 {
	history = r.cleanHistory(now, history)
	historyLen := len(history)
	if historyLen < r.TimesPerInterval {
		return 0
	}

	sort.Slice(history, func(i, j int) bool { return history[i] < history[j] })
	return (history[historyLen-r.TimesPerInterval] + r.SecondsPerInterval) - now
}

// GetRateLimitedCommand wraps around a command to make it rate limited.
func (r RateLimitConfig) GetRateLimitedCommand(command Command, commandDelayID string) Command {
	curry := func(sender service.Conversation, user service.User, msg [][]string, storage *storage.Storage, sink func(service.Conversation, service.Message)) {
		now := time.Now().Unix()
		history := []int64{}
		val, err := (*storage).GetUserValue(sender.ServiceID, user.ID, commandDelayID)
		if err == nil {
			var ok bool
			if history, ok = val.([]int64); ok == false {
				panic(fmt.Errorf("Interface type wasn't expected"))
			}
		}

		if r.rateLimited(now, history) {
			sink(
				service.Conversation{},
				service.Message{
					Title:       "Please try again later.",
					Description: r.Body,
					Footer:      fmt.Sprintf("%d Minutes remaining", now/60),
				},
			)
		} else {
			command.Exec(sender, user, msg, storage, sink)
			(*storage).SetUserValue(sender.ServiceID, user.ID, commandDelayID, append(history, now))
		}
	}

	return Command{
		Trigger: command.Trigger,
		Pattern: command.Pattern,
		Exec:    curry,
		Help: fmt.Sprintf(
			"%s. Can only be used %d times every %d seconds.",
			command.Help,
			r.TimesPerInterval,
			r.SecondsPerInterval,
		),
	}
}
