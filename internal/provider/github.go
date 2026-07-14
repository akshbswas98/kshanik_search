package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"kshanik_search/internal/models"
)

type GitHubProvider struct {
	client     *http.Client
	baseURL    string
	maxResults int
	token      string
}

func NewGitHubProvider(client *http.Client, baseURL string, maxResults int, token string) *GitHubProvider {
	return &GitHubProvider{client: client, baseURL: baseURL, maxResults: maxResults, token: token}
}

func (p *GitHubProvider) Name() string { return "github" }

type githubSearchResponse struct {
	Items []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		HTMLURL     string `json:"html_url"`
	} `json:"items"`
}

func (p *GitHubProvider) Search(ctx context.Context, query string) ([]models.SearchResult, error) {
	endpoint := fmt.Sprintf("%s/search/repositories?q=%s&per_page=%d", p.baseURL, url.QueryEscape(query), p.maxResults)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("github build request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "kshanik-search-go")
	if p.token != "" {
		req.Header.Set("Authorization", "Bearer "+p.token)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("github request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github unexpected status: %d", resp.StatusCode)
	}

	var payload githubSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("github decode: %w", err)
	}

	results := make([]models.SearchResult, 0, len(payload.Items))
	now := time.Now().UTC()
	for _, item := range payload.Items {
		results = append(results, models.SearchResult{
			Title:     item.Name,
			Snippet:   item.Description,
			URL:       item.HTMLURL,
			Source:    p.Name(),
			Timestamp: now,
		})
	}

	return results, nil
}
