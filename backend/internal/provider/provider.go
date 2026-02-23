// Package provider defines the SearchProvider interface and concrete
// implementations that fetch results from external APIs.
package provider

import (
	"context"

	"github.com/kshanik/search-engine/internal/models"
)

// SearchProvider is the interface that all search backends must implement.
// Implementations are expected to be safe for concurrent use.
type SearchProvider interface {
	// Search executes a query and returns normalized results.
	// The context should be used for timeout / cancellation propagation.
	Search(ctx context.Context, query string) ([]models.SearchResult, error)

	// Name returns a human-readable identifier for the provider (e.g. "duckduckgo").
	Name() string
}
