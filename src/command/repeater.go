package command

import "github.com/BKrajancic/FLD-Bot/m/v2/src/service"

// Return the received message
func Repeater(sender service.User, msg [][]string, sink func(service.User, string)) {
	sink(sender, msg[0][1])
}
