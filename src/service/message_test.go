package service

import (
	"testing"
)

func TestMessageEmpty(t *testing.T) {
	fields := []MessageField{{}, {}}
	msg := Message{Fields: fields}
	msg.ToString()
}
