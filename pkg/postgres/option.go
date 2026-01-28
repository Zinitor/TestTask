package postgres

import (
	"time"

	"gorm.io/gorm/logger"
)

const (
	_defaultMaxPoolSize     = 1
	_defaultMaxIdleConns    = 2
	_defaultConnTimeout     = 5 * time.Second
	_defaultMaxConnLifetime = 30 * time.Minute
)

// Option - functional options pattern.
type Option func(*options)

type options struct {
	maxPoolSize     int
	maxIdleConns    int
	connAttempts    int
	maxConnLifetime time.Duration
	connTimeout     time.Duration
	gormLogger      logger.Interface
}

// defaultOptions - sane defaults.
func defaultOptions() *options {
	return &options{
		maxPoolSize:  _defaultMaxPoolSize,
		maxIdleConns: _defaultMaxIdleConns,
		connTimeout:  _defaultConnTimeout,
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

// MaxConnLifetime sets maximum connection lifetime
func MaxConnLifetime(lifetime time.Duration) Option {
	return func(o *options) { o.maxConnLifetime = lifetime }
}

// GormLogger configures GORM's internal logger
func GormLogger(l logger.Interface) Option {
	return func(o *options) { o.gormLogger = l }
}
