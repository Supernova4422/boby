package command

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/storage"
)

// RateLimitConfig is a wrapper around a Command, that ensures a command is not used excessively.
type RateLimitConfig struct {
	TimesPerInterval   int    // How many times can it be used per interval?
	SecondsPerInterval int64  // How long is an interval, in seconds.
	Body               string // Reply when limit is reached.
	ID                 string // An ID used for storage purposes.
}

// rateLimited returns true if a message should be rate limited.
// now and history are expected to be in unix time.
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
// If SecondsPerInterval and TimesPerInterval are both 0, this function returns
// the given command.
//
// TODO: This function is not thread safe, calling it twice from two threads will get two
// versions of history, then writing the new history will only include on message. This is
// ignored for now since it hasn't been an issue yet.
func (r RateLimitConfig) GetRateLimitedCommand(command Command) Command {
	if r.SecondsPerInterval == 0 && r.TimesPerInterval == 0 {
		return command
	}

	rateLimitedCommand := command
	rateLimitedCommand.Exec = func(sender service.Conversation, user service.User, msg []interface{}, storage *storage.Storage, sink func(service.Conversation, service.Message)) {
		now := time.Now().Unix()
		history := []int64{}

		if val, ok := (*storage).GetUserValue(user, r.ID); ok {
			history, ok = val.([]int64)
			if !ok {
				panic(fmt.Errorf("interface type wasn't usable"))
			}
		}

		if r.rateLimited(now, history) {
			remaining := r.timeRemaining(now, history)
			var countdown string
			if remaining > 60*60 {
				countdown = fmt.Sprintf(
					"%d Hours and %d seconds remaining",
					(remaining/60)/60,
					remaining/60,
				)
			} else if remaining > 60 {
				countdown = fmt.Sprintf("%d Minutes remaining", remaining/60)
			} else {
				countdown = fmt.Sprintf("%d Seconds remaining", remaining)
			}

			sink(
				sender,
				service.Message{
					Title:       "Please try again later.",
					Description: strings.Join([]string{r.Body, countdown}, "\n"),
				},
			)
		} else {
			(*storage).SetUserValue(user, r.ID, append(history, now))
			command.Exec(sender, user, msg, storage, sink)
		}
	}

	rateLimitedCommand.Help = fmt.Sprintf(
		"%s. Can only be used %d times every %d seconds.",
		command.Help,
		r.TimesPerInterval,
		r.SecondsPerInterval,
	)

	return rateLimitedCommand
}
