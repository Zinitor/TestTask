package chats

import "ReviewTask/internal/service"

// TransformMessageToDTO converts repo model to service dto.
func TransformMessageToDTO(msg *Message) *service.Message {
	return &service.Message{
		ID:        msg.ID,
		Text:      msg.Text,
		CreatedAt: msg.CreatedAt,
	}
}

// TransformChatToDTO converts repo model to service dto.
func TransformChatToDTO(repoChat *Chat) *service.Chat {
	if repoChat == nil {
		return nil
	}

	messages := make([]service.Message, len(repoChat.Messages))
	for i, msg := range repoChat.Messages {
		messages[i] = service.Message{
			ID:        msg.ID,
			Text:      msg.Text,
			CreatedAt: msg.CreatedAt,
		}
	}

	return &service.Chat{
		ID:        repoChat.ID,
		Title:     repoChat.Title,
		CreatedAt: repoChat.CreatedAt,
		Messages:  messages,
	}
}
