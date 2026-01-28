package chats

import (
	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model

	Title    string    `gorm:"size:200;not null"`
	Messages []Message `gorm:"foreignKey:ChatID"`
}

type Message struct {
	gorm.Model

	Text   string `gorm:"size:200;not null"`
	ChatID uint
}
