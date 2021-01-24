package routine

import (
	"time"

	"github.com/BKrajancic/boby/m/v2/src/service"
)

// Heartbeat sends a message every now and again using parameter route.
// This is only useful for testing purposes.
func Heartbeat(delay time.Duration, destination service.Conversation, msg service.Message, route func(service.Conversation, service.Message)) {
	for range time.Tick(delay) {
		route(destination, msg)
	}
}
