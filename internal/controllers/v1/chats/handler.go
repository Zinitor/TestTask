package chats

import (
	"ReviewTask/internal"
	"ReviewTask/internal/service"
	"ReviewTask/pkg/logger"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
)

type chatsServiceInterface interface {
	Create(ctx context.Context, title string) (*service.Chat, error)
	Get(ctx context.Context, chatID int, limit int, lastID int) (*service.Chat, error)
	AddMessage(ctx context.Context, chatID int, text string) (*service.Message, error)
	Delete(ctx context.Context, chatID int) error
}

type Handler struct {
	logger    logger.Interface
	service   chatsServiceInterface
	validator *validator.Validate
}

func NewChatsHandler(
	logger logger.Interface,
	service chatsServiceInterface,
	validator *validator.Validate,
) *Handler {
	return &Handler{
		logger:    logger,
		service:   service,
		validator: validator,
	}
}

func RegisterRoutes(mux *http.ServeMux, handler *Handler) {
	mux.HandleFunc("GET /health", handler.health)
	mux.HandleFunc("POST /chats/", handler.add)
	mux.HandleFunc("POST /chats/{id}/messages/", handler.addMessage)
	mux.HandleFunc("GET /chats/{id}", handler.get)
	mux.HandleFunc("DELETE /chats/{id}", handler.delete)
}

func (h *Handler) handleServiceError(w http.ResponseWriter, err error) {
	statusCode := http.StatusInternalServerError
	errorMsg := "Internal server error"

	switch {
	case errors.Is(err, internal.ErrChatNotFound):
		statusCode = http.StatusNotFound
		errorMsg = "Chat not found"
	case errors.Is(err, internal.ErrInvalidInput),
		errors.Is(err, internal.ErrMessageEmpty),
		errors.Is(err, internal.ErrMessageTooLong):
		statusCode = http.StatusBadRequest
		if err != nil && err.Error() != "" {
			errorMsg = err.Error()
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": errorMsg}); encodeErr != nil {
		h.logger.Error("failed to encode error response: %v", encodeErr)
	}
}

func writeJSONResponse(w http.ResponseWriter, log logger.Interface, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Error("failed to encode JSON response [status=%d]: %v", statusCode, err)
	}
}

// Health handles the health-check endpoint for the Chats API.
//
//	@Summary		Health endpoint
//	@Description	Returns a simple greeting message confirming the Chats API is operational
//	@Produce		text/plain
//	@Success		200	{string}	string	"Hello from Chats API!"
//	@Router			/health [get]
func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	h.logger.Info("serving hello endpoint")
	fmt.Fprintln(w, "Hello from Chats API!")
}

// @Summary		Create a new chat
// @Description	Creates a chat resource using the title provided in the request body
// @Tags			chats
// @Accept			json
// @Produce		json
// @Param			chat	body		CreateChatRequest	true	"Chat creation payload"	SchemaExample({ "title": "Team Discussion" })
// @Success		201		{object}	Chat				"Chat created successfully"
// @Failure		400		{object}	map[string]string	"Invalid request body or validation error"
// @Failure		500		{object}	map[string]string	"Internal server error"
// @Router			/chats/ [post]
func (h *Handler) add(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer func() {
		if closeErr := r.Body.Close(); closeErr != nil {
			h.logger.Error("failed to close request body: %v", closeErr)
		}
	}()

	var req CreateChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode create chat request: %v", err)
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	//validation
	var maxTitleLen int = 200
	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	trimmedTitle := strings.TrimSpace(req.Title)
	if utf8.RuneCountInString(trimmedTitle) > maxTitleLen {
		http.Error(w, fmt.Sprintf("Title should be no more than %d symbols", maxTitleLen), http.StatusBadRequest)
		return
	}

	chat, err := h.service.Create(ctx, trimmedTitle)
	if err != nil {
		h.logger.Error("failed to create chat '%s': %v", trimmedTitle, err)
		h.handleServiceError(w, err)
		return
	}

	response := TransformToSchema(chat)
	writeJSONResponse(w, h.logger, http.StatusCreated, response)
}

// @Summary		Add message to chat
// @Description	Adds a new message to an existing chat specified by ID
// @Tags			chats
// @Accept			json
// @Produce		json
// @Param			id		path		int						true	"Chat ID"			Format(int64)
// @Param			message	body		CreateMessageRequest	true	"Message payload"	example({ "text": "Hello, team!" })
// @Success		201		{object}	Message					"Message created successfully"
// @Failure		400		{object}	map[string]string		"Invalid request: missing text or invalid format"
// @Failure		404		{object}	map[string]string		"Chat not found"
// @Failure		500		{object}	map[string]string		"Internal server error"
// @Router			/chats/{id}/messages/ [post]
func (h *Handler) addMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer func() {
		if closeErr := r.Body.Close(); closeErr != nil {
			h.logger.Error("failed to close request body: %v", closeErr)
		}
	}()

	path := strings.TrimPrefix(r.URL.Path, "/chats/")
	path = strings.TrimSuffix(path, "/")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] != "messages" {
		h.logger.Warn("invalid path format: %s", r.URL.Path)
		http.Error(w, "Invalid request path", http.StatusBadRequest)
		return
	}

	chatIDStr := parts[0]
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil || chatID <= 0 {
		http.Error(w, "Invalid chat ID format", http.StatusBadRequest)
		return
	}

	var req CreateMessageRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode create message request: %v", err)
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	//validation
	var maxMessageLen int = 5000
	if req.Text == "" {
		http.Error(w, "Text is required", http.StatusBadRequest)
		return
	}
	if utf8.RuneCountInString(req.Text) > maxMessageLen {
		http.Error(w, fmt.Sprintf("Message should be no more than %d symbols", maxMessageLen), http.StatusBadRequest)
		return
	}

	message, err := h.service.AddMessage(ctx, chatID, req.Text)
	if err != nil {
		h.logger.Error("failed to create message '%s': %v", req.Text, err)
		h.handleServiceError(w, err)
		return
	}

	response := TransformMessageToSchema(message)
	writeJSONResponse(w, h.logger, http.StatusCreated, response)
}

// @Summary		Delete chat by ID
// @Description	Soft deletes a chat and its messages
// @Tags			chats
// @Param			id	path		int					true	"Chat ID"	Format(int64)
// @Success		204	{object}	string				"Chat deleted successfully"
// @Failure		400	{object}	map[string]string	"Invalid chat ID"
// @Failure		404	{object}	map[string]string	"Chat not found"
// @Failure		500	{object}	map[string]string	"Internal server error"
// @Router			/chats/{id} [delete]
func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	path := strings.TrimPrefix(r.URL.Path, "/chats/")
	path = strings.TrimSuffix(path, "/")

	if path == "" {
		http.Error(w, "Chat ID is required", http.StatusBadRequest)
		return
	}

	chatID, err := strconv.Atoi(path)
	if err != nil || chatID <= 0 {
		http.Error(w, "Invalid chat ID format", http.StatusBadRequest)
		return
	}

	if err = h.service.Delete(ctx, chatID); err != nil {
		h.logger.Error("failed to delete chat: %v | id: %d", err, chatID)
		h.handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary		Get chat by ID
// @Description	Retrieves a chat with limited messages (default 20, max 100)
// @Tags			chats
// @Produce		json
// @Param			id		path		int					true	"Chat ID"					Format(int64)
// @Param			limit	query		int					false	"Max messages to return"	default(20)	minimum(1)	maximum(100)
// @Param			last_id	query		int					false	"id of last message"
// @Success		200		{object}	Chat				"Chat retrieved successfully"
// @Failure		400		{object}	map[string]string	"Invalid chat ID or limit"
// @Failure		404		{object}	map[string]string	"Chat not found"
// @Failure		500		{object}	map[string]string	"Internal server error"
// @Router			/chats/{id} [get]
func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	path := strings.TrimPrefix(r.URL.Path, "/chats/")
	path = strings.TrimSuffix(path, "/")
	if path == "" {
		http.Error(w, "Chat ID is required", http.StatusBadRequest)
		return
	}

	chatID, err := strconv.Atoi(path)
	if err != nil || chatID <= 0 {
		http.Error(w, "Invalid chat ID format", http.StatusBadRequest)
		return
	}
	qParams := r.URL.Query()
	limit := 20 // default
	if limitStr := qParams.Get("limit"); limitStr != "" {
		var parsed int
		if parsed, err = strconv.Atoi(limitStr); err == nil {
			switch {
			case parsed < 1:
				limit = 1
			case parsed > 100:
				limit = 100
			default:
				limit = parsed
			}
		}
	}
	var lastID int
	if rawLastID := qParams.Get("last_id"); rawLastID != "" {
		lastID, err = strconv.Atoi(rawLastID)
		if err != nil {
			http.Error(w, "Invalid last_id format: must be a positive integer", http.StatusBadRequest)
			return
		}
		if lastID <= 0 {
			http.Error(w, "last_id must be a positive integer", http.StatusBadRequest)
			return
		}
	}

	chat, err := h.service.Get(ctx, chatID, limit, lastID)
	if err != nil {
		h.logger.Error("failed to get chat: %v | id: %d", err, chatID)
		h.handleServiceError(w, err)
		return
	}

	response := TransformToSchema(chat)
	writeJSONResponse(w, h.logger, http.StatusOK, response)
}
