package service

import "kshanik_search/internal/models"

// AdvancedRanker is reserved for future BM25/vector hybrid ranking logic.
type AdvancedRanker interface {
	Rank(query string, input []models.SearchResult) ([]models.SearchResult, error)
}
