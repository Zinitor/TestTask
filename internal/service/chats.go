package service

import (
	"ReviewTask/internal"
	"ReviewTask/pkg/logger"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type chatRepository interface {
	Get(chatID int, limit int, lastID int) (*Chat, error)
	Add(title string) (*Chat, error)
	Delete(int) error
	AddMessage(chatID int, text string) (*Message, error)
}

type ChatsService struct {
	logger    logger.Interface
	chatsRepo chatRepository
}

func NewChatsService(logger logger.Interface, chatsRepo chatRepository) ChatsService {
	return ChatsService{
		logger:    logger,
		chatsRepo: chatsRepo,
	}
}

func (s ChatsService) Create(title string) (*Chat, error) {
	chatDTO, err := s.chatsRepo.Add(title)
	if err != nil {
		s.logger.Error("chats service: failed to create chat: %v", err)

		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, fmt.Errorf("%w: chat with this title already exists", internal.ErrInvalidInput)
		}
		return nil, fmt.Errorf("database error: %w", internal.ErrInternal)
	}

	return chatDTO, nil
}
func (s ChatsService) Get(chatID int, limit int, lastID int) (*Chat, error) {
	if chatID <= 0 {
		return nil, fmt.Errorf("%w: chat ID must be positive", internal.ErrInvalidInput)
	}

	chatDTO, err := s.chatsRepo.Get(chatID, limit, lastID)
	if err != nil {
		s.logger.Error("chats service: failed to get chat: %v | id: %d", err, chatID)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, internal.ErrChatNotFound
		}
		return nil, fmt.Errorf("database error: %w", internal.ErrInternal)
	}

	return chatDTO, nil
}

func (s ChatsService) AddMessage(chatID int, text string) (*Message, error) {
	messageDTO, err := s.chatsRepo.AddMessage(chatID, text)
	if err != nil {
		s.logger.Error("chats service: failed to add message: %v | chat_id: %d", err, chatID)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, internal.ErrChatNotFound
		}
		return nil, fmt.Errorf("database error: %w", internal.ErrInternal)
	}
	return messageDTO, nil
}

func (s ChatsService) Delete(chatID int) error {
	if chatID <= 0 {
		return fmt.Errorf("%w: chat ID must be positive", internal.ErrInvalidInput)
	}

	if err := s.chatsRepo.Delete(chatID); err != nil {
		s.logger.Error("chats service: failed to delete chat: %v | id: %d", err, chatID)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return internal.ErrChatNotFound
		}
		return fmt.Errorf("database error: %w", internal.ErrInternal)
	}

	return nil
}
