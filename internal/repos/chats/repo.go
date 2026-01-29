package chats

import (
	"ReviewTask/internal/service"
	"ReviewTask/pkg/logger"
	"ReviewTask/pkg/postgres"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	logger logger.Interface
	db     *postgres.DB
}

func NewRepository(logger logger.Interface, db *postgres.DB) Repository {
	return Repository{
		logger: logger,
		db:     db,
	}
}

func (r Repository) Add(ctx context.Context, title string) (*service.Chat, error) {
	if title == "" {
		return nil, errors.New("chat title cannot be empty")
	}

	chat := Chat{Title: title}

	if err := r.db.WithContext(ctx).Create(&chat).Error; err != nil {
		r.logger.Error("failed to create chat in DB: %v | title: %s", err, title)
		return nil, fmt.Errorf("database error: %w", err)
	}

	chatIDStr := strconv.FormatUint(uint64(chat.ID), 10)
	r.logger.Info("chat created successfully | id: %s, title: %s", chatIDStr, title)

	return TransformChatToDTO(&chat), nil
}

func (r Repository) Get(ctx context.Context, chatID int, limit int, lastID int) (*service.Chat, error) {
	if chatID <= 0 {
		return nil, fmt.Errorf("invalid chat ID: %d", chatID)
	}

	if limit < 0 {
		return nil, fmt.Errorf("invalid message limit: %d", limit)
	}

	var chat Chat
	chatIDUint := uint(chatID)

	err := r.db.WithContext(ctx).Preload("Messages", func(db *gorm.DB) *gorm.DB {
		query := db.Order("created_at DESC").Limit(limit)
		if lastID > 0 {
			query = query.Where("id < ?", lastID)
		}
		return query
	}).First(&chat, chatIDUint).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("chat not found: %w", gorm.ErrRecordNotFound)
		}
		r.logger.Error("failed to get chat", "error", err, "chat_id", chatID)
		return nil, fmt.Errorf("database error: %w", err)
	}

	return TransformChatToDTO(&chat), nil
}

func (r Repository) Delete(ctx context.Context, chatID int) error {
	if chatID <= 0 {
		return fmt.Errorf("invalid chat ID: %d", chatID)
	}

	chatIDUint := uint(chatID)

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Message{}).
			Where("chat_id = ? AND deleted_at IS NULL", chatIDUint).
			Update("deleted_at", time.Now()).Error; err != nil {
			r.logger.Error("failed to soft-delete messages: %v | chat_id: %d", err, chatID)
			return fmt.Errorf("database error: %w", err)
		}

		result := tx.Delete(&Chat{}, chatIDUint)
		if result.Error != nil {
			r.logger.Error("failed to delete chat: %v | id: %d", result.Error, chatID)
			return fmt.Errorf("database error: %w", result.Error)
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("chat not found: %w", gorm.ErrRecordNotFound)
		}

		r.logger.Info("chat and %d messages soft-deleted successfully | id: %d",
			result.RowsAffected, chatID)
		return nil
	})
}

func (r Repository) AddMessage(ctx context.Context, chatID int, text string) (*service.Message, error) {
	if chatID <= 0 {
		return nil, errors.New("invalid chat ID: must be positive integer")
	}
	if text == "" {
		return nil, errors.New("message text cannot be empty")
	}
	chatIDUint := uint(chatID)

	if err := r.db.WithContext(ctx).Select("id").First(&Chat{}, chatIDUint).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Warn("chat not found or soft-deleted | id: %d", chatID)
			return nil, fmt.Errorf("chat not found: %w", gorm.ErrRecordNotFound)
		}
		r.logger.Error("database error checking chat existence: %v | id: %d", err, chatID)
		return nil, fmt.Errorf("database error: %w", err)
	}

	message := Message{
		Text:   text,
		ChatID: chatIDUint,
	}

	if err := r.db.WithContext(ctx).Create(&message).Error; err != nil {
		r.logger.Error("failed to insert message: %v | chat_id: %d", err, chatID)
		return nil, fmt.Errorf("failed to save message: %w", err)
	}

	return TransformMessageToDTO(&message), nil
}
