package command

import (
	"fmt"
	"log"
	"math"
	"sort"
	"strconv"
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
	Global             bool
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
func (r RateLimitConfig) timeRemaining(now int64, history []int64) time.Duration {
	history = r.cleanHistory(now, history)
	historyLen := len(history)
	if historyLen < r.TimesPerInterval {
		return 0
	}
	sort.Slice(history, func(i, j int) bool { return history[i] < history[j] })
	seconds := (history[historyLen-r.TimesPerInterval] + r.SecondsPerInterval) - now
	result, err := time.ParseDuration(strconv.FormatInt(seconds, 10) + "s")

	if err != nil {
		log.Panic(err)
	}

	return result
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
		now, history := r.GetRateLimitHistory(storage, user)
		if r.rateLimited(now, history) {
			remaining := r.timeRemaining(now, history)
			remainingAsString := r.timeRemainingToString(remaining) + " remaining"
			sink(
				sender,
				service.Message{
					Title:       "Please try again later.",
					Description: r.Body + " " + remainingAsString,
				},
			)
		} else {
			r.SetRateLimitHistory(append(history, now), storage, user)
			command.Exec(sender, user, msg, storage, sink)
		}
	}

	durationAsStr, _ := time.ParseDuration(strconv.FormatInt(r.SecondsPerInterval, 10) + "s")
	interval := r.timeRemainingToString(durationAsStr)
	rateLimitedCommand.Help = fmt.Sprintf(
		"%s. Can only be used %d times every %s",
		command.Help,
		r.TimesPerInterval,
		interval,
	)

	return rateLimitedCommand
}

// GetRateLimitedCommandInfo returns a command that reports the rate limiting of another command.
// If SecondsPerInterval and TimesPerInterval are both 0, this function returns
// the given command.
func (r RateLimitConfig) GetRateLimitedCommandInfo(command Command) Command {
	if r.SecondsPerInterval == 0 && r.TimesPerInterval == 0 {
		rateLimitedCommand := command
		rateLimitedCommand.Exec = func(sender service.Conversation, user service.User, msg []interface{}, storage *storage.Storage, sink func(service.Conversation, service.Message)) {
			sink(
				sender,
				service.Message{
					Title: "This command has unlimited usage.",
				},
			)
		}
		return rateLimitedCommand
	}

	command.Trigger = "info" + command.Trigger
	command.Parameters = []Parameter{}
	command.Help = "Get info about this command"

	rateLimitedCommand := command
	rateLimitedCommand.Exec = func(sender service.Conversation, user service.User, msg []interface{}, storage *storage.Storage, sink func(service.Conversation, service.Message)) {
		now, history := r.GetRateLimitHistory(storage, user)
		durationAsStr, _ := time.ParseDuration(strconv.FormatInt(r.SecondsPerInterval, 10) + "s")
		history = r.cleanHistory(now, history)
		interval := r.timeRemainingToString(durationAsStr)
		if r.rateLimited(now, history) {
			sink(
				sender,
				service.Message{
					Title:       "Command is currently rate limited",
					Description: fmt.Sprintf("%d/%d per %s", len(history), r.TimesPerInterval, interval) + " remaining.",
				},
			)
		} else {
			sink(
				sender,
				service.Message{
					Title:       "Currently not rate limited",
					Description: fmt.Sprintf("%d/%d per %s", len(history), r.TimesPerInterval, interval) + " remaining.",
				},
			)
		}
	}

	return rateLimitedCommand
}

// GetRateLimitHistory returns the history of usages for this command by a user, according to storage.
// 'now' refers to the current time, this is in unix time with seconds precisionn.
func (r RateLimitConfig) GetRateLimitHistory(storage *storage.Storage, user service.User) (now int64, history []int64) {
	if r.Global {
		if val, ok := (*storage).GetGlobalValue(r.ID); ok {
			history, ok = val.([]int64)
			if !ok {
				log.Panic("interface type wasn't usable")
			}
		}
	} else {
		if val, ok := (*storage).GetUserValue(user, r.ID); ok {
			history, ok = val.([]int64)
			if !ok {
				log.Panic("interface type wasn't usable")
			}
		}
	}

	now = time.Now().Unix()
	history = r.cleanHistory(now, history)
	r.SetRateLimitHistory(history, storage, user)

	return now, history
}

// SetRateLimitHistory sets the history.
func (r RateLimitConfig) SetRateLimitHistory(history []int64, storage *storage.Storage, user service.User) {
	if r.Global {
		(*storage).SetGlobalValue(r.ID, history)
	} else {
		(*storage).SetUserValue(user, r.ID, history)
	}
}

func (r RateLimitConfig) timeRemainingToString(remaining time.Duration) string {
	var countdown string
	if remaining.Hours() >= 24 {
		countdown = fmt.Sprintf("%.2f Days", math.Ceil(remaining.Hours()/24))
	} else if remaining.Hours() >= 1 {
		countdown = fmt.Sprintf("%.2f Hours", math.Ceil(remaining.Hours()))
	} else if remaining.Minutes() >= 1 {
		countdown = fmt.Sprintf("%.2f Minutes", math.Ceil(remaining.Minutes()))
	} else {
		countdown = fmt.Sprintf("%.2f Seconds", math.Ceil(remaining.Seconds()))
	}

	return countdown
}
