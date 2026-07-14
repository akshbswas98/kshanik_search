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

type DuckDuckGoProvider struct {
	client     *http.Client
	baseURL    string
	maxResults int
}

func NewDuckDuckGoProvider(client *http.Client, baseURL string, maxResults int) *DuckDuckGoProvider {
	return &DuckDuckGoProvider{client: client, baseURL: baseURL, maxResults: maxResults}
}

func (p *DuckDuckGoProvider) Name() string { return "duckduckgo" }

type duckduckgoResponse struct {
	Heading       string `json:"Heading"`
	AbstractText  string `json:"AbstractText"`
	AbstractURL   string `json:"AbstractURL"`
	RelatedTopics []struct {
		Text     string `json:"Text"`
		FirstURL string `json:"FirstURL"`
		Topics   []struct {
			Text     string `json:"Text"`
			FirstURL string `json:"FirstURL"`
		} `json:"Topics"`
	} `json:"RelatedTopics"`
}

func (p *DuckDuckGoProvider) Search(ctx context.Context, query string) ([]models.SearchResult, error) {
	endpoint := fmt.Sprintf("%s/?q=%s&format=json&no_html=1&skip_disambig=1", p.baseURL, url.QueryEscape(query))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("duckduckgo build request: %w", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("duckduckgo request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("duckduckgo unexpected status: %d", resp.StatusCode)
	}

	var payload duckduckgoResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("duckduckgo decode: %w", err)
	}

	results := make([]models.SearchResult, 0, p.maxResults)
	now := time.Now().UTC()

	if payload.AbstractURL != "" {
		results = append(results, models.SearchResult{
			Title:     payload.Heading,
			Snippet:   payload.AbstractText,
			URL:       payload.AbstractURL,
			Source:    p.Name(),
			Timestamp: now,
		})
	}

	for _, topic := range payload.RelatedTopics {
		if len(results) >= p.maxResults {
			break
		}
		if topic.FirstURL != "" {
			results = append(results, models.SearchResult{Title: topic.Text, Snippet: topic.Text, URL: topic.FirstURL, Source: p.Name(), Timestamp: now})
			continue
		}
		for _, nested := range topic.Topics {
			if len(results) >= p.maxResults {
				break
			}
			if nested.FirstURL == "" {
				continue
			}
			results = append(results, models.SearchResult{Title: nested.Text, Snippet: nested.Text, URL: nested.FirstURL, Source: p.Name(), Timestamp: now})
		}
	}

	return results, nil
}
