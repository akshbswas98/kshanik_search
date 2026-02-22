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

type WikipediaProvider struct {
	client     *http.Client
	baseURL    string
	maxResults int
}

func NewWikipediaProvider(client *http.Client, baseURL string, maxResults int) *WikipediaProvider {
	return &WikipediaProvider{client: client, baseURL: baseURL, maxResults: maxResults}
}

func (p *WikipediaProvider) Name() string { return "wikipedia" }

func (p *WikipediaProvider) Search(ctx context.Context, query string) ([]models.SearchResult, error) {
	endpoint := fmt.Sprintf("%s/w/api.php?action=opensearch&search=%s&limit=%d&namespace=0&format=json", p.baseURL, url.QueryEscape(query), p.maxResults)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("wikipedia build request: %w", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("wikipedia request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wikipedia unexpected status: %d", resp.StatusCode)
	}

	var payload []interface{}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("wikipedia decode: %w", err)
	}
	if len(payload) < 4 {
		return nil, nil
	}

	titles, _ := payload[1].([]interface{})
	descriptions, _ := payload[2].([]interface{})
	links, _ := payload[3].([]interface{})

	results := make([]models.SearchResult, 0, p.maxResults)
	now := time.Now().UTC()
	for i := 0; i < len(links) && i < p.maxResults; i++ {
		u, _ := links[i].(string)
		if u == "" {
			continue
		}
		title, _ := titles[i].(string)
		snippet, _ := descriptions[i].(string)
		results = append(results, models.SearchResult{
			Title:     title,
			Snippet:   snippet,
			URL:       u,
			Source:    p.Name(),
			Timestamp: now,
		})
	}

	return results, nil
}
