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

// RedditProvider queries the Reddit public JSON search API.
// No authentication required — uses the public .json endpoints.
type RedditProvider struct {
	client  *http.Client
	baseURL string
}

// redditSearchResponse represents the Reddit search API response.
type redditSearchResponse struct {
	Data struct {
		Children []struct {
			Data redditPost `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

type redditPost struct {
	Title       string  `json:"title"`
	Selftext    string  `json:"selftext"`
	Permalink   string  `json:"permalink"`
	URL         string  `json:"url"`
	Subreddit   string  `json:"subreddit"`
	Score       int     `json:"score"`
	NumComments int     `json:"num_comments"`
	CreatedUTC  float64 `json:"created_utc"`
}

// NewRedditProvider creates a new Reddit search provider with the given timeout.
func NewRedditProvider(timeout time.Duration) *RedditProvider {
	return &RedditProvider{
		client: &http.Client{
			Timeout: timeout,
		},
		baseURL: "https://www.reddit.com/search.json",
	}
}

func (r *RedditProvider) Name() string {
	return "reddit"
}

func (r *RedditProvider) Search(ctx context.Context, query string) ([]models.SearchResult, error) {
	params := url.Values{}
	params.Set("q", query)
	params.Set("sort", "relevance")
	params.Set("limit", "10")
	params.Set("type", "link")

	reqURL := fmt.Sprintf("%s?%s", r.baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("reddit: failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "KshanikSearch/1.0 (search engine project)")
	req.Header.Set("Accept", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("reddit: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, fmt.Errorf("reddit: rate limited (status 429)")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("reddit: unexpected status %d", resp.StatusCode)
	}

	var redditResp redditSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&redditResp); err != nil {
		return nil, fmt.Errorf("reddit: failed to decode response: %w", err)
	}

	return r.normalizeResults(redditResp), nil
}

func (r *RedditProvider) normalizeResults(resp redditSearchResponse) []models.SearchResult {
	var results []models.SearchResult
	now := time.Now()

	for _, child := range resp.Data.Children {
		post := child.Data

		// Build snippet from self text or subreddit info.
		snippet := post.Selftext
		if snippet == "" {
			snippet = fmt.Sprintf("Posted in r/%s", post.Subreddit)
		}
		snippet = truncate(snippet, 300)

		// Add engagement info.
		if post.Score > 0 || post.NumComments > 0 {
			snippet = fmt.Sprintf("%s (↑%d • 💬%d comments)", snippet, post.Score, post.NumComments)
		}

		// Full Reddit URL.
		resultURL := fmt.Sprintf("https://www.reddit.com%s", post.Permalink)

		// Score based on Reddit upvotes (log scale).
		score := 0.5
		if post.Score > 0 {
			score = 0.5 + normalizeEngagement(post.Score)
		}

		results = append(results, models.SearchResult{
			Title:     post.Title,
			Snippet:   snippet,
			URL:       resultURL,
			Source:    "reddit",
			Score:     score,
			Timestamp: now,
		})
	}

	return results
}

// normalizeEngagement maps engagement count to a 0.0 – 0.4 range.
func normalizeEngagement(count int) float64 {
	if count <= 0 {
		return 0
	}
	score := 0.0
	c := float64(count)
	for c > 1 {
		score += 0.08
		c /= 10
	}
	if score > 0.4 {
		score = 0.4
	}
	return score
}
