package messages

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model

	ID        int    `gorm:"primaryKey"`
	Text      string `gorm:"size:200"`
	CreatedAt time.Time
}
