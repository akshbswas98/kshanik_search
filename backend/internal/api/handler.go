// Package api implements the HTTP handlers for the search engine API.
package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/kshanik/search-engine/internal/models"
	"github.com/kshanik/search-engine/internal/service"
	"github.com/rs/zerolog"
)

// Handler holds the HTTP handler dependencies.
type Handler struct {
	searchService *service.MetaSearchService
	logger        zerolog.Logger
	security      SecurityConfig
}

// NewHandler creates a new Handler with the given search service.
func NewHandler(searchService *service.MetaSearchService, logger zerolog.Logger, security SecurityConfig) *Handler {
	return &Handler{
		searchService: searchService,
		logger:        logger.With().Str("component", "api_handler").Logger(),
		security:      security,
	}
}

// RegisterRoutes registers all API routes on the given mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/search", h.handleSearch)
	mux.HandleFunc("/health", h.handleHealth)
}

// handleSearch processes GET /search?q=query
func (h *Handler) handleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "only GET method is allowed")
		return
	}

	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		h.writeError(w, http.StatusBadRequest, "query parameter 'q' is required")
		return
	}
	if len(query) > h.security.MaxQueryLength {
		h.writeError(w, http.StatusBadRequest, "query is too long")
		return
	}

	h.logger.Info().Str("query", query).Str("remote", r.RemoteAddr).Msg("search request received")

	result, err := h.searchService.Search(r.Context(), query)
	if err != nil {
		h.logger.Error().Err(err).Str("query", query).Msg("search failed")
		msg := "search failed"
		if h.security.ExposeErrorDetail {
			msg = "search failed: " + err.Error()
		}
		h.writeError(w, http.StatusInternalServerError, msg)
		return
	}

	h.writeJSON(w, http.StatusOK, result)
}

// handleHealth returns the server health status.
func (h *Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	h.writeJSON(w, http.StatusOK, map[string]string{
		"status": "healthy",
		"engine": "kshanik-meta-search",
	})
}

// writeJSON serializes data to JSON and writes it to the response.
func (h *Handler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error().Err(err).Msg("failed to encode JSON response")
	}
}

// writeError sends a structured error response.
func (h *Handler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, models.ErrorResponse{
		Error: message,
		Code:  status,
	})
}
