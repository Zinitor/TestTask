package logger

import (
	"context"
	"time"

	gormlogger "gorm.io/gorm/logger"
)

// GormAdapter bridges your logger to GORM's interface (ignores context).
type GormAdapter struct {
	l *Logger
}

func NewGormAdapter(l *Logger) *GormAdapter {
	return &GormAdapter{l: l}
}

func (a *GormAdapter) LogMode(level gormlogger.LogLevel) gormlogger.Interface { return a }
func (a *GormAdapter) Info(ctx context.Context, msg string, args ...any) {
	a.l.Info(msg, args...)
}
func (a *GormAdapter) Warn(ctx context.Context, msg string, args ...any) {
	a.l.Warn(msg, args...)
}
func (a *GormAdapter) Error(ctx context.Context, msg string, args ...any) {
	a.l.Error(msg, args...)
}

func (a *GormAdapter) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	if err != nil {
		a.l.Error("SQL ERROR", "elapsed", elapsed, "rows", rows, "sql", sql, "error", err)
	} else if elapsed > 200*time.Millisecond {
		a.l.Info("SLOW SQL", "elapsed", elapsed, "rows", rows, "sql", sql)
	}
}
