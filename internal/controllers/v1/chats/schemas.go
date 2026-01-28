package chats

import "time"

type CreateChatRequest struct {
	Title string `json:"title" validate:"required,min=1,max=200"`
}

type CreateMessageRequest struct {
	Text string `json:"text" validate:"required,min=1,max=500"`
}

type Message struct {
	ID        uint      `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type Chat struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Messages  []Message `json:"messages"`
	CreatedAt time.Time `json:"created_at"`
}
