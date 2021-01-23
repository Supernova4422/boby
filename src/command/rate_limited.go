package command

import (
	"fmt"
	"sort"
	"time"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/storage"
)

// RateLimitConfig can be used to ensure a command is not used excessively.
type RateLimitConfig struct {
	TimesPerInterval   int    // How many times can it be used per interval?
	SecondsPerInterval int64  // How long is an interval, in seconds.
	Body               string // Reply when limit is reached.
	ID                 string // An ID used for storage purposes.
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
// If SecondsPerInterval and TimesPerInterval are both 0, this function returns
// the given command.
func (r RateLimitConfig) GetRateLimitedCommand(command Command) Command {
	if r.SecondsPerInterval == 0 && r.TimesPerInterval == 0 {
		return command
	}

	curry := func(sender service.Conversation, user service.User, msg [][]string, storage *storage.Storage, sink func(service.Conversation, service.Message)) {
		now := time.Now().Unix()
		history := []int64{}
		val, err := (*storage).GetUserValue(user, r.ID)
		if err == nil {
			var ok bool = false

			// (HACK) JSONStorage has a nasty side effect that large numbers get
			// unmarshalled as float64. In the future, this should be received
			// from storage as int64.
			switch val.(type) {
			case []int64:
				history, ok = val.([]int64)
			case []float64:
				for _, number := range val.([]float64) {
					history = append(history, int64(number))
				}
				ok = true
			case []interface{}:
				ok = true
				for _, number := range val.([]interface{}) {
					switch number.(type) {
					case int64:
						history = append(history, number.(int64))
					case float64:
						history = append(history, int64(number.(float64)))
					default:
						ok = false
						break
					}
				}
			}

			if ok == false {
				panic(fmt.Errorf("Interface type wasn't usable"))
			}
		}

		if r.rateLimited(now, history) {
			remaining := r.timeRemaining(now, history)
			var footer string
			if remaining > 60*60 {
				footer = fmt.Sprintf(
					"%d Hours and %d seconds remaining",
					(remaining/60)/60,
					remaining/60,
				)
			} else if remaining > 60 {
				footer = fmt.Sprintf("%d Minutes remaining", remaining/60)
			} else {
				footer = fmt.Sprintf("%d Seconds remaining", remaining)
			}

			sink(
				sender,
				service.Message{
					Title:       "Please try again later.",
					Description: r.Body,
					Footer:      footer,
				},
			)
		} else {
			go (*storage).SetUserValue(user, r.ID, append(history, now))
			command.Exec(sender, user, msg, storage, sink)
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
