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

// StackOverflowProvider queries the StackExchange public API for Q&A results.
// No authentication required — uses the public API with a generous rate limit.
type StackOverflowProvider struct {
	client  *http.Client
	baseURL string
}

// stackOverflowResponse represents the StackExchange search API response.
type stackOverflowResponse struct {
	Items []stackOverflowItem `json:"items"`
}

type stackOverflowItem struct {
	Title        string   `json:"title"`
	Link         string   `json:"link"`
	Score        int      `json:"score"`
	AnswerCount  int      `json:"answer_count"`
	IsAnswered   bool     `json:"is_answered"`
	ViewCount    int      `json:"view_count"`
	Tags         []string `json:"tags"`
	CreationDate int64    `json:"creation_date"`
}

// NewStackOverflowProvider creates a new StackOverflow search provider.
func NewStackOverflowProvider(timeout time.Duration) *StackOverflowProvider {
	return &StackOverflowProvider{
		client: &http.Client{
			Timeout: timeout,
		},
		baseURL: "https://api.stackexchange.com/2.3/search/excerpts",
	}
}

func (s *StackOverflowProvider) Name() string {
	return "stackoverflow"
}

func (s *StackOverflowProvider) Search(ctx context.Context, query string) ([]models.SearchResult, error) {
	params := url.Values{}
	params.Set("order", "desc")
	params.Set("sort", "relevance")
	params.Set("q", query)
	params.Set("site", "stackoverflow")
	params.Set("pagesize", "10")
	params.Set("filter", "default")

	// Use the /search endpoint instead for better title results.
	searchURL := fmt.Sprintf("https://api.stackexchange.com/2.3/search?%s", params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("stackoverflow: failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "KshanikSearch/1.0")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("stackoverflow: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, fmt.Errorf("stackoverflow: rate limited (status 429)")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("stackoverflow: unexpected status %d", resp.StatusCode)
	}

	var soResp stackOverflowResponse
	if err := json.NewDecoder(resp.Body).Decode(&soResp); err != nil {
		return nil, fmt.Errorf("stackoverflow: failed to decode response: %w", err)
	}

	return s.normalizeResults(soResp), nil
}

func (s *StackOverflowProvider) normalizeResults(resp stackOverflowResponse) []models.SearchResult {
	var results []models.SearchResult
	now := time.Now()

	for _, item := range resp.Items {
		// Build a rich snippet with tags and answer info.
		snippet := ""
		if len(item.Tags) > 0 {
			tagStr := ""
			maxTags := 4
			if len(item.Tags) < maxTags {
				maxTags = len(item.Tags)
			}
			for i := 0; i < maxTags; i++ {
				if i > 0 {
					tagStr += ", "
				}
				tagStr += item.Tags[i]
			}
			snippet = fmt.Sprintf("[%s] ", tagStr)
		}

		answerStatus := "No answers yet"
		if item.IsAnswered {
			answerStatus = fmt.Sprintf("✅ Answered (%d answers)", item.AnswerCount)
		} else if item.AnswerCount > 0 {
			answerStatus = fmt.Sprintf("%d answers (not accepted)", item.AnswerCount)
		}
		snippet += fmt.Sprintf("%s • 👁 %d views • ↑%d", answerStatus, item.ViewCount, item.Score)

		// Score based on SO vote count and answer status.
		score := 0.5
		if item.IsAnswered {
			score += 0.15
		}
		if item.Score > 0 {
			score += normalizeEngagement(item.Score) * 0.5
		}
		if score > 1.0 {
			score = 1.0
		}

		results = append(results, models.SearchResult{
			Title:     item.Title,
			Snippet:   snippet,
			URL:       item.Link,
			Source:    "stackoverflow",
			Score:     score,
			Timestamp: now,
		})
	}

	return results
}
