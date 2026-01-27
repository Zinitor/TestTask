package chats

import (
	"ReviewTask/internal/repos/messages"
	"time"

	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model

	ID        int    `gorm:"primaryKey"`
	Title     string `gorm:"size:200"`
	CreatedAt time.Time
	Messages  []messages.Message `gorm:"foreignKey:UserId"`
}
