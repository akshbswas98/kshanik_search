// Package index is a placeholder for the future indexing subsystem.
//
// Evolution path:
//
//	Phase 1 (current) — No local index; results come from external APIs.
//	Phase 2           — Build an inverted index from crawled documents:
//	  - Tokenization & normalization pipeline
//	  - Term frequency / document frequency tracking
//	  - Positional index for phrase queries
//	Phase 3           — Elasticsearch / Bleve integration for
//	                     full-text search with faceting.
//	Phase 4           — Vector embeddings index for semantic search.
//
// This module compiles but all operations are no-ops.
package index

import (
	"context"

	"github.com/rs/zerolog"
)

// Document represents a single indexed document.
type Document struct {
	URL     string
	Title   string
	Content string
	Links   []string
}

// Index is the future inverted index / search index.
type Index struct {
	logger zerolog.Logger
}

// NewIndex creates a new Index instance.
func NewIndex(logger zerolog.Logger) *Index {
	return &Index{
		logger: logger.With().Str("component", "index").Logger(),
	}
}

// AddDocument indexes a single document. Placeholder — currently a no-op.
func (idx *Index) AddDocument(ctx context.Context, doc Document) error {
	idx.logger.Debug().Str("url", doc.URL).Msg("index: AddDocument called (no-op)")
	return nil
}

// Search queries the local index. Placeholder — currently returns nil.
func (idx *Index) Search(ctx context.Context, query string) ([]Document, error) {
	idx.logger.Debug().Str("query", query).Msg("index: Search called (no-op)")
	return nil, nil
}

// Flush persists any in-memory index data to storage. Placeholder.
func (idx *Index) Flush() error {
	idx.logger.Debug().Msg("index: Flush called (no-op)")
	return nil
}
