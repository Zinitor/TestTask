package service

import "ReviewTask/pkg/logger"

type chatRepository interface {
	Get()
	Add()
	Delete()
}
type messageRepository interface {
	Get()
	Add()
	Delete()
}

type ChatsService struct {
	logger      logger.Interface
	chatsRepo   chatRepository
	messageRepo messageRepository
}

func NewChatsService(logger logger.Interface, chatsRepo chatRepository, messageRepo messageRepository) ChatsService {
	return ChatsService{
		logger:      logger,
		chatsRepo:   chatsRepo,
		messageRepo: messageRepo,
	}
}

func (s ChatsService) Create() error {
	return nil
}
func (s ChatsService) GetMessages() ([]string, error) {
	var arr []string
	return arr, nil
}
func (s ChatsService) Get() error {
	return nil
}
