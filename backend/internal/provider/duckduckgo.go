package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kshanik/search-engine/internal/models"
)

// DuckDuckGoProvider queries the DuckDuckGo Instant Answer API.
// This is a free, public API that does not require authentication.
type DuckDuckGoProvider struct {
	client  *http.Client
	baseURL string
}

// duckDuckGoResponse maps the relevant fields from the DDG API response.
type duckDuckGoResponse struct {
	Abstract       string              `json:"Abstract"`
	AbstractURL    string              `json:"AbstractURL"`
	AbstractSource string              `json:"AbstractSource"`
	Heading        string              `json:"Heading"`
	RelatedTopics  []duckDuckGoTopic   `json:"RelatedTopics"`
	Results        []duckDuckGoResult  `json:"Results"`
}

type duckDuckGoTopic struct {
	FirstURL string `json:"FirstURL"`
	Text     string `json:"Text"`
	Result   string `json:"Result"`
	// Nested topics (category groups)
	Topics []duckDuckGoTopic `json:"Topics"`
}

type duckDuckGoResult struct {
	FirstURL string `json:"FirstURL"`
	Text     string `json:"Text"`
}

// NewDuckDuckGoProvider creates a new DuckDuckGo provider with the given timeout.
func NewDuckDuckGoProvider(timeout time.Duration) *DuckDuckGoProvider {
	return &DuckDuckGoProvider{
		client: &http.Client{
			Timeout: timeout,
		},
		baseURL: "https://api.duckduckgo.com/",
	}
}

func (d *DuckDuckGoProvider) Name() string {
	return "duckduckgo"
}

func (d *DuckDuckGoProvider) Search(ctx context.Context, query string) ([]models.SearchResult, error) {
	params := url.Values{}
	params.Set("q", query)
	params.Set("format", "json")
	params.Set("no_redirect", "1")
	params.Set("no_html", "1")
	params.Set("skip_disambig", "1")

	reqURL := fmt.Sprintf("%s?%s", d.baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("duckduckgo: failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "KshanikSearch/1.0")

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("duckduckgo: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("duckduckgo: unexpected status %d", resp.StatusCode)
	}

	var ddgResp duckDuckGoResponse
	if err := json.NewDecoder(resp.Body).Decode(&ddgResp); err != nil {
		return nil, fmt.Errorf("duckduckgo: failed to decode response: %w", err)
	}

	return d.normalizeResults(ddgResp), nil
}

func (d *DuckDuckGoProvider) normalizeResults(resp duckDuckGoResponse) []models.SearchResult {
	var results []models.SearchResult
	now := time.Now()

	// Add the main abstract result if present.
	if resp.Abstract != "" && resp.AbstractURL != "" {
		heading := resp.Heading
		if heading == "" {
			heading = resp.AbstractSource
		}
		results = append(results, models.SearchResult{
			Title:     heading,
			Snippet:   truncate(resp.Abstract, 300),
			URL:       resp.AbstractURL,
			Source:    "duckduckgo",
			Score:     1.0, // Abstract answers score highest from DDG.
			Timestamp: now,
		})
	}

	// Add related topics.
	for _, topic := range resp.RelatedTopics {
		results = append(results, d.flattenTopics(topic, now)...)
	}

	// Add direct results.
	for _, r := range resp.Results {
		if r.FirstURL != "" {
			results = append(results, models.SearchResult{
				Title:     extractTitle(r.Text),
				Snippet:   truncate(r.Text, 300),
				URL:       r.FirstURL,
				Source:    "duckduckgo",
				Score:     0.7,
				Timestamp: now,
			})
		}
	}

	return results
}

// flattenTopics recursively extracts results from nested topic groups.
func (d *DuckDuckGoProvider) flattenTopics(topic duckDuckGoTopic, now time.Time) []models.SearchResult {
	var results []models.SearchResult

	if topic.FirstURL != "" {
		results = append(results, models.SearchResult{
			Title:     extractTitle(topic.Text),
			Snippet:   truncate(topic.Text, 300),
			URL:       topic.FirstURL,
			Source:    "duckduckgo",
			Score:     0.6,
			Timestamp: now,
		})
	}

	for _, nested := range topic.Topics {
		results = append(results, d.flattenTopics(nested, now)...)
	}

	return results
}

// extractTitle pulls a reasonable title from a DDG text blob.
func extractTitle(text string) string {
	if idx := strings.Index(text, " - "); idx > 0 && idx < 120 {
		return text[:idx]
	}
	return truncate(text, 80)
}

// truncate shortens s to at most maxLen characters.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
