package service

import "time"

type Message struct {
	ID        uint
	Text      string
	CreatedAt time.Time
}

type Chat struct {
	ID        uint
	Title     string
	Messages  []Message
	CreatedAt time.Time
}
