package internal

import "errors"

var (
	ErrChatNotFound   = errors.New("chat not found")
	ErrInvalidInput   = errors.New("invalid input")
	ErrInternal       = errors.New("internal error")
	ErrMessageTooLong = errors.New("message text exceeds maximum length")
	ErrMessageEmpty   = errors.New("message text cannot be empty")
)
