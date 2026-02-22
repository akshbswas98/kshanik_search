package models

import "time"

// SearchResult is the normalized result shape returned by all providers.
type SearchResult struct {
	Title     string    `json:"title"`
	Snippet   string    `json:"snippet"`
	URL       string    `json:"url"`
	Source    string    `json:"source"`
	Score     float64   `json:"score"`
	Timestamp time.Time `json:"timestamp"`
}
