package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kshanik/search-engine/internal/models"
	"github.com/kshanik/search-engine/internal/provider"
	"github.com/rs/zerolog"
)

type MockProvider struct {
	name    string
	results []models.SearchResult
	err     error
	delay   time.Duration
}

func (m *MockProvider) Search(ctx context.Context, query string) ([]models.SearchResult, error) {
	if m.delay > 0 {
		select {
		case <-time.After(m.delay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	if m.err != nil {
		return nil, m.err
	}
	return m.results, nil
}

func (m *MockProvider) Name() string {
	return m.name
}

func TestMetaSearchService_Search(t *testing.T) {
	logger := zerolog.Nop() // Disables logging during tests

	t.Run("Success all providers", func(t *testing.T) {
		providers := []provider.SearchProvider{
			&MockProvider{
				name: "Provider1",
				results: []models.SearchResult{
					{Title: "Result 1", URL: "http://example.com/1"},
				},
			},
			&MockProvider{
				name: "Provider2",
				results: []models.SearchResult{
					{Title: "Result 2", URL: "http://example.com/2"},
				},
			},
		}

		service := NewMetaSearchService(providers, 2*time.Second, logger)
		resp, err := service.Search(context.Background(), "test query")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if resp == nil {
			t.Fatal("Expected response, got nil")
		}

		if resp.TotalCount != 2 {
			t.Errorf("Expected 2 results, got %d", resp.TotalCount)
		}
		if len(resp.Results) != 2 {
			t.Errorf("Expected 2 results, got %d", len(resp.Results))
		}
	})

	t.Run("Partial failure", func(t *testing.T) {
		providers := []provider.SearchProvider{
			&MockProvider{
				name: "ProviderSuccess",
				results: []models.SearchResult{
					{Title: "Result 1", URL: "http://example.com/1"},
				},
			},
			&MockProvider{
				name: "ProviderFail",
				err:  errors.New("API down"),
			},
		}

		service := NewMetaSearchService(providers, 2*time.Second, logger)
		resp, err := service.Search(context.Background(), "test query")

		if err != nil {
			t.Fatalf("Expected no error on partial failure, got %v", err)
		}

		if resp.TotalCount != 1 {
			t.Errorf("Expected 1 result from successful provider, got %d", resp.TotalCount)
		}
	})

	t.Run("Total failure", func(t *testing.T) {
		providers := []provider.SearchProvider{
			&MockProvider{
				name: "ProviderFail1",
				err:  errors.New("API 1 down"),
			},
			&MockProvider{
				name: "ProviderFail2",
				err:  errors.New("API 2 down"),
			},
		}

		service := NewMetaSearchService(providers, 2*time.Second, logger)
		_, err := service.Search(context.Background(), "test query")

		if err == nil {
			t.Fatal("Expected an error when all providers fail, but got nil")
		}
	})

	t.Run("Timeout handling", func(t *testing.T) {
		providers := []provider.SearchProvider{
			&MockProvider{
				name: "FastProvider",
				results: []models.SearchResult{
					{Title: "Fast Result", URL: "http://example.com/fast"},
				},
			},
			&MockProvider{
				name:  "SlowProvider",
				delay: 2 * time.Second, // Intentionally slow
				results: []models.SearchResult{
					{Title: "Slow Result", URL: "http://example.com/slow"},
				},
			},
		}

		// Set a very aggressive timeout that the FastProvider passes but SlowProvider fails
		service := NewMetaSearchService(providers, 100*time.Millisecond, logger)
		resp, err := service.Search(context.Background(), "test timeout")

		if err != nil {
			t.Fatalf("Expected no outer error, got %v", err)
		}

		// It should only have the fast provider's result
		if resp.TotalCount != 1 {
			t.Errorf("Expected 1 result due to timeout, got %d", resp.TotalCount)
		}
		if resp.Results[0].URL != "http://example.com/fast" {
			t.Errorf("Expected fast result, got %s", resp.Results[0].URL)
		}
	})
}
