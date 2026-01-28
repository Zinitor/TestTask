// Package postgres implements postgres connection.
package postgres

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB -.
type DB struct {
	*gorm.DB
}

// New -.
func New(dsn string, opts ...Option) (*DB, error) {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: o.gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("gorm open failed: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB failed: %w", err)
	}
	sqlDB.SetMaxOpenConns(o.maxPoolSize)
	sqlDB.SetMaxIdleConns(o.maxIdleConns)
	sqlDB.SetConnMaxLifetime(o.maxConnLifetime)

	return &DB{DB: db}, nil
}

// Close -.
func (d *DB) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}
	return sqlDB.Close()
}
