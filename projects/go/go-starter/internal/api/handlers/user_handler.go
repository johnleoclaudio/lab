package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/yourusername/go-starter/internal/api/middleware"
	"github.com/yourusername/go-starter/internal/models"
	"github.com/yourusername/go-starter/internal/service"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService service.UserService
	logger      *slog.Logger
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService service.UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

// GetUser handles GET /api/v1/users/{id} requests
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := middleware.GetRequestID(ctx)

	// Parse user ID from URL parameter
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		h.logger.WarnContext(ctx, "missing user id parameter")
		respondError(w, reqID, http.StatusBadRequest, "INVALID_ID", "User ID is required")
		return
	}

	// Parse and validate UUID
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.WarnContext(ctx, "invalid user id format",
			slog.String("id", idStr),
			slog.String("error", err.Error()),
		)
		respondError(w, reqID, http.StatusBadRequest, "INVALID_ID", "Invalid user ID format")
		return
	}

	// Get user from service
	user, err := h.userService.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			h.logger.InfoContext(ctx, "user not found",
				slog.String("id", id.String()),
			)
			respondError(w, reqID, http.StatusNotFound, "NOT_FOUND", "User not found")
			return
		}

		// Internal server error
		h.logger.ErrorContext(ctx, "failed to get user",
			slog.String("id", id.String()),
			slog.String("error", err.Error()),
		)
		respondError(w, reqID, http.StatusInternalServerError, "INTERNAL_ERROR", "An unexpected error occurred")
		return
	}

	// Convert to JSON:API format and respond
	response := JSONAPIResponse{
		Data: ToJSONAPIData(user),
	}

	h.logger.InfoContext(ctx, "user retrieved successfully",
		slog.String("id", id.String()),
	)

	respondJSON(w, http.StatusOK, response)
}
