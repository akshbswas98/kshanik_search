package index

import "context"

// Index is a future extension point for inverted index + vector store.
type Index interface {
	AddDocument(ctx context.Context, id string, content string) error
	Search(ctx context.Context, query string, limit int) ([]string, error)
}
