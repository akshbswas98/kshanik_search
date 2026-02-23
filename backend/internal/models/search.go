// Package models defines the core data transfer objects used across the search engine.
package models

import "time"

// SearchResult represents a single search result returned by any provider.
type SearchResult struct {
	Title     string    `json:"title"`
	Snippet   string    `json:"snippet"`
	URL       string    `json:"url"`
	Source    string    `json:"source"`
	Score     float64   `json:"score"`
	Timestamp time.Time `json:"timestamp"`
}

// SearchRequest encapsulates the incoming search query with pagination and options.
type SearchRequest struct {
	Query  string `json:"query"`
	Page   int    `json:"page,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}

// SearchResponse is the top-level JSON response returned to clients.
type SearchResponse struct {
	Results    []SearchResult `json:"results"`
	TotalCount int            `json:"total_count"`
	Query      string         `json:"query"`
	TimeTaken  string         `json:"time_taken"`
}

// ErrorResponse is returned when the API encounters an error.
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Details string `json:"details,omitempty"`
}
