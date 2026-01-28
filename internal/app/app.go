// Package app configures and runs application.
package app

import (
	"ReviewTask/config"
	_ "ReviewTask/docs" // register swagger spec at init time
	v1 "ReviewTask/internal/controllers/v1"
	v1chats "ReviewTask/internal/controllers/v1/chats" // Alias for handler creation
	"ReviewTask/internal/repos/chats"
	"ReviewTask/internal/service"
	"ReviewTask/pkg/httpserver"
	"ReviewTask/pkg/logger"
	"ReviewTask/pkg/postgres"
	validate "ReviewTask/pkg/validator"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)
	gormLog := logger.NewGormAdapter(l)

	val := validate.NewValidator()

	pg, err := postgres.New(
		cfg.PG.URL,
		postgres.MaxPoolSize(cfg.PG.PoolMax),
		postgres.MaxIdleConns(cfg.PG.PoolIdle),
		postgres.GormLogger(gormLog),
	)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	chatsRepo := chats.NewRepository(l, pg)

	// Service
	chatsService := service.NewChatsService(l, chatsRepo)

	// HTTP Server
	chatsHandler := v1chats.NewChatsHandler(l, chatsService, val)

	// controllers
	mux, err := v1.NewServeMux(chatsHandler, cfg.Swagger.Enabled, l)
	if err != nil {
		l.Fatal(fmt.Errorf("failed to create v1 HTTP router: %w", err))
	}

	// Server setup
	srv := httpserver.New(
		l,
		httpserver.WithAddress(cfg.HTTP.Port),
		httpserver.WithShutdownTimeout(time.Duration(cfg.HTTP.ShutdownTime)),
	)
	srv.HTTPServer.Handler = mux
	srv.Start()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err = <-srv.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = srv.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
