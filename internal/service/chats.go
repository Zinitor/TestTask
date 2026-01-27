package service

import "ReviewTask/pkg/logger"

type ChatsRepository interface {
	get()
	add()
	delete()
}

type ChatsService struct {
	logger    logger.Interface
	chatsRepo ChatsRepository
}

func NewChatsService(logger logger.Interface, chatsRepo ChatsRepository) ChatsService {
	return ChatsService{
		logger:    logger,
		chatsRepo: chatsRepo,
	}
}
