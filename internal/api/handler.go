package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"kshanik_search/internal/service"
)

type Handler struct {
	searchService   *service.SearchService
	requestTimeout  time.Duration
	providerTimeout time.Duration
}

func NewHandler(searchService *service.SearchService, requestTimeout, providerTimeout time.Duration) *Handler {
	return &Handler{searchService: searchService, requestTimeout: requestTimeout, providerTimeout: providerTimeout}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /search", h.search)
}

func (h *Handler) search(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		h.writeError(w, http.StatusBadRequest, "missing q query parameter")
		return
	}

	requestCtx, cancel := context.WithTimeout(r.Context(), h.requestTimeout)
	defer cancel()
	providerCtx, cancelProviders := context.WithTimeout(requestCtx, h.providerTimeout)
	defer cancelProviders()

	results, err := h.searchService.Search(providerCtx, query)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
	}
}

func (h *Handler) writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
