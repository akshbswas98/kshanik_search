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

// GitHubProvider queries the GitHub Repository Search API.
// An optional token can be set via configuration to increase rate limits.
type GitHubProvider struct {
	client  *http.Client
	baseURL string
	token   string
}

// githubSearchResponse represents the GitHub search API response.
type githubSearchResponse struct {
	TotalCount int          `json:"total_count"`
	Items      []githubRepo `json:"items"`
}

type githubRepo struct {
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	HTMLURL     string `json:"html_url"`
	Language    string `json:"language"`
	Stars       int    `json:"stargazers_count"`
	UpdatedAt   string `json:"updated_at"`
}

// NewGitHubProvider creates a new GitHub search provider.
// Pass an empty string for token to use unauthenticated access (lower rate limit).
func NewGitHubProvider(timeout time.Duration, token string) *GitHubProvider {
	return &GitHubProvider{
		client: &http.Client{
			Timeout: timeout,
		},
		baseURL: "https://api.github.com/search/repositories",
		token:   token,
	}
}

func (g *GitHubProvider) Name() string {
	return "github"
}

func (g *GitHubProvider) Search(ctx context.Context, query string) ([]models.SearchResult, error) {
	params := url.Values{}
	params.Set("q", query)
	params.Set("sort", "stars")
	params.Set("order", "desc")
	params.Set("per_page", "10")

	reqURL := fmt.Sprintf("%s?%s", g.baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("github: failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "KshanikSearch/1.0")
	if g.token != "" {
		req.Header.Set("Authorization", "Bearer "+g.token)
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("github: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests {
		return nil, fmt.Errorf("github: rate limited (status %d)", resp.StatusCode)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github: unexpected status %d", resp.StatusCode)
	}

	var ghResp githubSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&ghResp); err != nil {
		return nil, fmt.Errorf("github: failed to decode response: %w", err)
	}

	return g.normalizeResults(ghResp), nil
}

func (g *GitHubProvider) normalizeResults(resp githubSearchResponse) []models.SearchResult {
	var results []models.SearchResult
	now := time.Now()

	for _, repo := range resp.Items {
		description := repo.Description
		if description == "" {
			description = "GitHub repository"
		}

		// Build a richer snippet with language and star count.
		snippet := description
		if repo.Language != "" {
			snippet = fmt.Sprintf("[%s] %s", repo.Language, description)
		}
		if repo.Stars > 0 {
			snippet = fmt.Sprintf("%s (⭐ %d)", snippet, repo.Stars)
		}

		// Score based on star count (log scale, capped).
		score := 0.5
		if repo.Stars > 0 {
			score = 0.5 + normalizeStars(repo.Stars)
		}

		results = append(results, models.SearchResult{
			Title:     repo.FullName,
			Snippet:   truncate(snippet, 300),
			URL:       repo.HTMLURL,
			Source:    "github",
			Score:     score,
			Timestamp: now,
		})
	}

	return results
}

// normalizeStars maps star count to a 0.0 – 0.5 range using a log scale.
func normalizeStars(stars int) float64 {
	if stars <= 0 {
		return 0
	}
	// log10(100_000) ≈ 5, so divide by 10 to cap at ~0.5.
	score := 0.0
	s := float64(stars)
	for s > 1 {
		score += 0.1
		s /= 10
	}
	if score > 0.5 {
		score = 0.5
	}
	return score
}
