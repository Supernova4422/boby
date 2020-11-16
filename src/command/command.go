package command

import (
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
)

type Command func(service.User, [][]string, func(service.User, string))
