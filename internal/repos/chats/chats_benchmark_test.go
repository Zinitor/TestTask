//go:build bench

package chats

import (
	"context"
	"fmt"
	"testing"
)

// although this benchmark isn't representative of actual db
// it's still an enough measurement tool
func BenchmarkRepository_Get(b *testing.B) {
	db := setupTestDB(b)
	fmt.Println(">>> BEFORE repo creation")
	repo := NewRepository(&mockLogger{}, db)
	fmt.Println(">>> AFTER repo creation")

	chat := Chat{Title: "Benchmark Chat"}
	db.Create(&chat)
	for i := range 1000 {
		db.Create(&Message{Text: fmt.Sprintf("msg %d", i), ChatID: chat.ID})
	}

	b.ResetTimer()
	for b.Loop() {
		repo.Get(context.Background(), int(chat.ID), 50, 0)
	}
}
