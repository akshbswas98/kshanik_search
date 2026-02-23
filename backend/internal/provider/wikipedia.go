package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/kshanik/search-engine/internal/models"
)

// WikipediaProvider queries the Wikipedia Search API (MediaWiki opensearch).
type WikipediaProvider struct {
	client  *http.Client
	baseURL string
}

// wikipediaSearchResponse represents the opensearch JSON response:
// [ query, [titles...], [descriptions...], [urls...] ]
type wikipediaSearchResponse struct {
	Query        string
	Titles       []string
	Descriptions []string
	URLs         []string
}

// NewWikipediaProvider creates a new Wikipedia search provider with the given timeout.
func NewWikipediaProvider(timeout time.Duration) *WikipediaProvider {
	return &WikipediaProvider{
		client: &http.Client{
			Timeout: timeout,
		},
		baseURL: "https://en.wikipedia.org/w/api.php",
	}
}

func (w *WikipediaProvider) Name() string {
	return "wikipedia"
}

func (w *WikipediaProvider) Search(ctx context.Context, query string) ([]models.SearchResult, error) {
	params := url.Values{}
	params.Set("action", "opensearch")
	params.Set("search", query)
	params.Set("limit", "10")
	params.Set("namespace", "0")
	params.Set("format", "json")

	reqURL := fmt.Sprintf("%s?%s", w.baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("wikipedia: failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "KshanikSearch/1.0 (search engine project)")

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("wikipedia: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wikipedia: unexpected status %d", resp.StatusCode)
	}

	// Opensearch returns a heterogeneous JSON array:
	// [query_string, [titles], [descriptions], [urls]]
	var raw []json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("wikipedia: failed to decode response: %w", err)
	}

	if len(raw) < 4 {
		return nil, nil // No results or unexpected format.
	}

	parsed, err := parseOpenSearch(raw)
	if err != nil {
		return nil, fmt.Errorf("wikipedia: failed to parse opensearch: %w", err)
	}

	return w.normalizeResults(parsed), nil
}

func parseOpenSearch(raw []json.RawMessage) (*wikipediaSearchResponse, error) {
	var result wikipediaSearchResponse

	if err := json.Unmarshal(raw[0], &result.Query); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(raw[1], &result.Titles); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(raw[2], &result.Descriptions); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(raw[3], &result.URLs); err != nil {
		return nil, err
	}

	return &result, nil
}

func (w *WikipediaProvider) normalizeResults(parsed *wikipediaSearchResponse) []models.SearchResult {
	var results []models.SearchResult
	now := time.Now()

	count := len(parsed.Titles)
	if count > len(parsed.Descriptions) {
		count = len(parsed.Descriptions)
	}
	if count > len(parsed.URLs) {
		count = len(parsed.URLs)
	}

	for i := 0; i < count; i++ {
		snippet := parsed.Descriptions[i]
		if snippet == "" {
			snippet = "Wikipedia article"
		}

		results = append(results, models.SearchResult{
			Title:     parsed.Titles[i],
			Snippet:   snippet,
			URL:       parsed.URLs[i],
			Source:    "wikipedia",
			Score:     0.8 - float64(i)*0.02, // Slight rank decay.
			Timestamp: now,
		})
	}

	return results
}
