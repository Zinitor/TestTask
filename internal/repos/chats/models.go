package repo

import "time"

type Chat struct {
	ID        int    `gorm:"primaryKey"`
	Title     string `gorm:"size:200"`
	CreatedAt time.Time
}
