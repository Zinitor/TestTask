package postgres

import (
	"time"

	"gorm.io/gorm/logger"
)

// Option - functional options pattern.
type Option func(*options)

type options struct {
	maxPoolSize  int
	maxIdleConns int
	connAttempts int
	connTimeout  time.Duration
	gormLogger   logger.Interface
}

// defaultOptions - sane defaults.
func defaultOptions() *options {
	return &options{
		maxPoolSize:  10,
		maxIdleConns: 2,
		connTimeout:  5 * time.Second,
		gormLogger:   logger.Default,
	}
}

// MaxPoolSize -.
func MaxPoolSize(size int) Option {
	return func(o *options) {
		o.maxPoolSize = size
	}
}

// ConnAttempts -.
func ConnAttempts(attempts int) Option {
	return func(o *options) {
		o.connAttempts = attempts
	}
}

// ConnTimeout -.
func ConnTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.connTimeout = timeout
	}
}

// MaxIdleConns sets maximum idle connections
func MaxIdleConns(n int) Option {
	return func(o *options) { o.maxIdleConns = n }
}

// GormLogger configures GORM's internal logger
func GormLogger(l logger.Interface) Option {
	return func(o *options) { o.gormLogger = l }
}
