package provider

import (
	"context"

	"kshanik_search/internal/models"
)

// SearchProvider describes providers that can search a remote data source.
type SearchProvider interface {
	Search(ctx context.Context, query string) ([]models.SearchResult, error)
	Name() string
}
