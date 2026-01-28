package chats

import (
	"ReviewTask/internal/service"
)

func TransformMessageToSchema(msg *service.Message) *Message {
	return &Message{
		ID:        msg.ID,
		Text:      msg.Text,
		CreatedAt: msg.CreatedAt,
	}
}

// TransformToSchema converts service model to handler response schema
func TransformToSchema(serviceChat *service.Chat) *Chat {
	if serviceChat == nil {
		return nil
	}

	messages := make([]Message, len(serviceChat.Messages))
	for i, msg := range serviceChat.Messages {
		messages[i] = Message{
			ID:        msg.ID,
			Text:      msg.Text,
			CreatedAt: msg.CreatedAt,
		}
	}

	return &Chat{
		ID:        serviceChat.ID,
		Title:     serviceChat.Title,
		Messages:  messages,
		CreatedAt: serviceChat.CreatedAt,
	}
}
