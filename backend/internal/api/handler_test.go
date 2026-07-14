package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kshanik/search-engine/internal/models"
	"github.com/kshanik/search-engine/internal/provider"
	"github.com/kshanik/search-engine/internal/service"
	"github.com/rs/zerolog"
)

// MockProvider is used to control search service behavior during testing.
type MockProvider struct {
	results []models.SearchResult
	err     error
}

func (m *MockProvider) Search(ctx context.Context, query string) ([]models.SearchResult, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.results, nil
}

func (m *MockProvider) Name() string {
	return "MockAPI"
}

func TestHandler_HandleSearch(t *testing.T) {
	logger := zerolog.Nop() // Disables logging for test cleanliness

	t.Run("Successful query", func(t *testing.T) {
		providers := []provider.SearchProvider{
			&MockProvider{
				results: []models.SearchResult{
					{Title: "API Test", URL: "http://test.loc/1"},
				},
			},
		}
		searchService := service.NewMetaSearchService(providers, time.Second, logger)
		handler := NewHandler(searchService, logger, LoadSecurityConfig())

		req, err := http.NewRequest(http.MethodGet, "/search?q=golang", nil)
		if err != nil {
			t.Fatal(err)
		}

		// We use httptest.NewRecorder() to simulate the ResponseWriter
		rr := httptest.NewRecorder()
		handler.handleSearch(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var response models.SearchResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatal("Failed to decode response JSON")
		}

		if response.Query != "golang" {
			t.Errorf("expected query 'golang', got '%s'", response.Query)
		}
		if len(response.Results) != 1 {
			t.Errorf("expected 1 result, got %d", len(response.Results))
		}
	})

	t.Run("Missing query parameter", func(t *testing.T) {
		searchService := service.NewMetaSearchService([]provider.SearchProvider{}, time.Second, logger)
		handler := NewHandler(searchService, logger, LoadSecurityConfig())

		req, err := http.NewRequest(http.MethodGet, "/search", nil) // no ?q=...
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.handleSearch(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("Method Not Allowed (POST instead of GET)", func(t *testing.T) {
		searchService := service.NewMetaSearchService([]provider.SearchProvider{}, time.Second, logger)
		handler := NewHandler(searchService, logger, LoadSecurityConfig())

		req, err := http.NewRequest(http.MethodPost, "/search?q=test", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.handleSearch(rr, req)

		if status := rr.Code; status != http.StatusMethodNotAllowed {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
		}
	})

}

func TestHandler_HandleHealth(t *testing.T) {
	logger := zerolog.Nop()
	searchService := service.NewMetaSearchService([]provider.SearchProvider{}, time.Second, logger)
	handler := NewHandler(searchService, logger, LoadSecurityConfig())

	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.handleHealth(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal("Failed to decode response JSON")
	}

	if response["status"] != "healthy" {
		t.Errorf("Expected status healthy, got %s", response["status"])
	}
}
