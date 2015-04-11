package main

import (
	"time"
)

// message: Generic chat message
type message struct {
	Sender    string
	Body      string
	SentTime  time.Time
	AvatarURL string
}
