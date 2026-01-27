package controller

import (
	"ReviewTask/pkg/logger"
	"fmt"
	"net/http"
)

type ChatsService interface {
	Create() error
	GetMessages() ([]string, error)
	Get() error
}

type ChatsHandler struct {
	logger  logger.Interface
	service ChatsService
}

func NewChatsHandler(logger logger.Interface, service ChatsService) *ChatsHandler {
	return &ChatsHandler{
		logger:  logger,
		service: service,
	}
}

// RegisterRoutes attaches chat-related routes to the provided ServeMux.
// This keeps routing configuration encapsulated within the controller package.
func RegisterRoutes(mux *http.ServeMux, handler *ChatsHandler) {
	mux.HandleFunc("GET /chats/hello", handler.hello)
	mux.HandleFunc("GET /chats/headers", handler.headers)
}

// hello handles the health-check endpoint for the Chats API.
//
//	@Summary		Hello endpoint
//	@Description	Returns a simple greeting message confirming the Chats API is operational
//	@Tags			chats
//	@Produce		text/plain
//	@Success		200	{string}	string	"Hello from Chats API!"
//	@Router			/chats/hello [get]
func (h *ChatsHandler) hello(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("serving hello endpoint")
	fmt.Fprintln(w, "Hello from Chats API!")
}

// headers echoes all HTTP headers from the incoming request.
//
//	@Summary		Echo request headers
//	@Description	Returns all HTTP headers sent in the request (useful for debugging client requests)
//	@Tags			chats
//	@Produce		text/plain
//	@Success		200	{string}	string	"HeaderName: value\n..."
//	@Router			/chats/headers [get]
func (h *ChatsHandler) headers(w http.ResponseWriter, r *http.Request) {
	for name, values := range r.Header {
		for _, value := range values {
			fmt.Fprintf(w, "%s: %s\n", name, value)
		}
	}
}
