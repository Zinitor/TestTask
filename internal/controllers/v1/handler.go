package v1

import (
	"ReviewTask/internal/controllers/v1/chats"
	"ReviewTask/pkg/logger"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/swaggo/swag"
)

// NewServeMux creates and configures a ServeMux with all v1 API routes.
// Registers business routes and (optionally) Swagger documentation routes.
// Returns a ready-to-use ServeMux. Safe to call once during application startup.
//
// Note: Requires _ "ReviewTask/docs" import in main package to register Swagger spec
// before this function is called.
func NewServeMux(
	chatsHandler *chats.Handler,
	enableSwagger bool,
	logger logger.Interface,
) (*http.ServeMux, error) {
	mux := http.NewServeMux()
	chats.RegisterRoutes(mux, chatsHandler)
	if enableSwagger {
		if err := registerSwaggerRoutes(mux, logger); err != nil {
			return nil, err
		}
		logger.Info("Swagger UI registered at /swagger/")
	}

	return mux, nil
}

func registerSwaggerRoutes(mux *http.ServeMux, logger logger.Interface) error {
	spec, err := swag.ReadDoc()
	if err != nil {
		return err
	}
	specBytes := []byte(spec)

	mux.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if _, err = w.Write(specBytes); err != nil {
			logger.Error("Failed to write swagger spec: %v", err)
		}
	})

	mux.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	})

	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("list"),
	))
	return nil
}
