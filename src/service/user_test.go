package service

import (
	"testing"
)

func TestSimple(t *testing.T) {
	user := User{"0", "0"}
	out := user.ToString()
	if out != "0,0" {
		t.Fail()
	}
}
