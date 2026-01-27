package repo

import (
	"ReviewTask/pkg/logger"

	"gorm.io/gorm"
)

type ChatsRepository struct {
	logger *logger.Interface
	db     *gorm.DB
}

func NewChatsRepository(logger *logger.Interface, db *gorm.DB) ChatsRepository {
	return ChatsRepository{
		logger: logger,
		db:     db,
	}

}

func (r *ChatsRepository) Create(user *Chat) error {
	return r.db.Create(user).Error
}

func (r *ChatsRepository) FindAll() ([]*Chat, error) {
	var users []*Chat
	err := r.db.Find(&users).Error
	return users, err
}

func (r *ChatsRepository) FindByID(id uint) (*Chat, error) {
	var user Chat
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *ChatsRepository) Update(user *Chat) error {
	return r.db.Save(user).Error
}

func (r *ChatsRepository) Delete(id uint) error {
	return r.db.Delete(&Chat{}, id).Error
}
