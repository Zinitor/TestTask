// Package app configures and runs application.
package app

import (
	"ReviewTask/config"
	_ "ReviewTask/docs"
	"ReviewTask/internal/controller"
	"ReviewTask/internal/repos/chats"
	"ReviewTask/internal/repos/messages"
	"ReviewTask/internal/service"
	"ReviewTask/pkg/httpserver"
	"ReviewTask/pkg/logger"
	"ReviewTask/pkg/postgres"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/swaggo/swag"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)
	gormLog := logger.NewGormAdapter(l)

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
	messagesRepo := messages.NewRepository(l, pg)

	// Service
	chatsService := service.NewChatsService(l, chatsRepo, messagesRepo)

	// HTTP Server
	srv := httpserver.New(
		l,
		httpserver.WithAddress(cfg.HTTP.Port),
		httpserver.WithShutdownTimeout(time.Duration(cfg.HTTP.ShutdownTime)),
	)

	// controllers
	mux := http.NewServeMux()
	chatsHandler := controller.NewChatsHandler(l, chatsService)
	controller.RegisterRoutes(mux, chatsHandler)

	if cfg.Swagger.Enabled {
		// Read spec ONCE at startup (avoids per-request overhead)
		spec, err := swag.ReadDoc()
		if err != nil {
			l.Fatal(fmt.Errorf("app - Run - failed to load swagger spec: %w", err))
		}
		specBytes := []byte(spec)

		// Serve spec JSON
		mux.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if _, err := w.Write(specBytes); err != nil {
				l.Error("Failed to write swagger spec: %v", err)
			}
		})

		// Redirect /swagger → /swagger/ (UX improvement)
		mux.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
		})

		// Serve Swagger UI with optimized defaults
		mux.Handle("/swagger/", httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"), // Relative path to spec
			httpSwagger.DeepLinking(true),
			httpSwagger.DocExpansion("list"),
		))

		l.Info("Swagger UI enabled at /swagger/")
	}

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
