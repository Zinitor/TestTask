package chats

import (
	"context"
	"fmt"
	"testing"

	"ReviewTask/pkg/postgres"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// mockLogger implements logger.Interface for testing
type mockLogger struct {
	entries []any
}

func (m *mockLogger) Info(format string, args ...any) {
	m.entries = append(m.entries, "INFO: "+format)
}
func (m *mockLogger) Warn(format string, args ...any) {
	m.entries = append(m.entries, "WARN: "+format)
}
func (m *mockLogger) Error(format any, args ...any) {
	m.entries = append(m.entries, format)
}
func (m *mockLogger) Debug(format any, args ...any) {}
func (m *mockLogger) Fatal(format any, args ...any) {}

func setupTestDB(tb testing.TB) *postgres.DB {
	tb.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(tb, err)

	err = db.AutoMigrate(&Chat{}, &Message{})
	require.NoError(tb, err)

	return &postgres.DB{DB: db}
}

func TestRepository_Add(t *testing.T) {
	tests := []struct {
		name          string
		title         string
		expectError   bool
		expectedTitle string
	}{
		{
			name:        "empty title",
			title:       "",
			expectError: true,
		},
		{
			name:          "valid title",
			title:         "Test Chat",
			expectError:   false,
			expectedTitle: "Test Chat",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := NewRepository(&mockLogger{}, db)

			chat, err := repo.Add(context.Background(), tt.title)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, chat)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, chat)
				assert.Equal(t, tt.expectedTitle, chat.Title)
				assert.Greater(t, chat.ID, uint(0))
			}
		})
	}
}

func TestRepository_Get(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(&mockLogger{}, db)

	chat := Chat{Title: "Test Chat"}
	result := db.Create(&chat)
	require.NoError(t, result.Error)

	msg1 := Message{Text: "Hello", ChatID: chat.ID}
	msg2 := Message{Text: "World", ChatID: chat.ID}
	db.Create(&msg1)
	db.Create(&msg2)

	tests := []struct {
		name        string
		chatID      int
		limit       int
		lastID      int
		expectError bool
		expectCount int
	}{
		{
			name:        "invalid chat ID",
			chatID:      -1,
			limit:       10,
			expectError: true,
		},
		{
			name:        "chat not found",
			chatID:      99999,
			limit:       10,
			expectError: true,
		},
		{
			name:        "valid chat with messages",
			chatID:      int(chat.ID),
			limit:       10,
			expectError: false,
			expectCount: 2,
		},
		{
			name:        "pagination with lastID",
			chatID:      int(chat.ID),
			limit:       1,
			lastID:      int(msg2.ID),
			expectError: false,
			expectCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.Get(context.Background(), tt.chatID, tt.limit, tt.lastID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result.Messages, tt.expectCount)
			}
		})
	}
}

func TestRepository_Delete(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(*postgres.DB) uint
		chatID      int
		expectError bool
	}{
		{
			name:        "invalid chat ID",
			setup:       func(_ *postgres.DB) uint { return 0 },
			chatID:      -5,
			expectError: true,
		},
		{
			name:        "chat not found",
			setup:       func(_ *postgres.DB) uint { return 0 },
			chatID:      99999,
			expectError: true,
		},
		{
			name: "successful soft delete",
			setup: func(db *postgres.DB) uint {
				chat := Chat{Title: "To Delete"}
				db.Create(&chat)
				msg := Message{Text: "msg", ChatID: chat.ID}
				db.Create(&msg)
				return chat.ID
			},
			chatID:      0,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := NewRepository(&mockLogger{}, db)

			if tt.setup != nil {
				chatID := tt.setup(db)
				if tt.chatID == 0 && chatID > 0 {
					tt.chatID = int(chatID)
				}
			}

			err := repo.Delete(context.Background(), tt.chatID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				var chat Chat
				db.Unscoped().First(&chat, tt.chatID)
				assert.NotNil(t, chat.DeletedAt)

				var msg Message
				db.Unscoped().Where("chat_id = ?", tt.chatID).First(&msg)
				assert.NotNil(t, msg.DeletedAt)
			}
		})
	}
}

func TestRepository_AddMessage(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(&mockLogger{}, db)

	chat := Chat{Title: "Test Chat"}
	db.Create(&chat)

	tests := []struct {
		name        string
		chatID      int
		text        string
		expectError bool
	}{
		{
			name:        "invalid chat ID",
			chatID:      -1,
			text:        "Hi",
			expectError: true,
		},
		{
			name:        "empty message text",
			chatID:      int(chat.ID),
			text:        "",
			expectError: true,
		},
		{
			name:        "chat not found",
			chatID:      99999,
			text:        "Hi",
			expectError: true,
		},
		{
			name:        "valid message",
			chatID:      int(chat.ID),
			text:        "Hello world",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := repo.AddMessage(context.Background(), tt.chatID, tt.text)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, msg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, msg)
				assert.Equal(t, tt.text, msg.Text)
			}
		})
	}
}

// although this benchmark isn't representative of actual db
// it's still an enough measurement tool
func BenchmarkRepository_Get_WithPreload(b *testing.B) {
	db := setupTestDB(b)
	repo := NewRepository(&mockLogger{}, db)

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
