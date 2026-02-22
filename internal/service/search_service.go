package service

import (
	"context"
	"sort"
	"strings"
	"sync"

	"kshanik_search/internal/models"
	"kshanik_search/internal/provider"
)

type SearchService struct {
	providers []provider.SearchProvider
}

func NewSearchService(providers []provider.SearchProvider) *SearchService {
	return &SearchService{providers: providers}
}

func (s *SearchService) Search(ctx context.Context, query string) ([]models.SearchResult, error) {
	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		results []models.SearchResult
	)

	for _, p := range s.providers {
		p := p
		wg.Add(1)
		go func() {
			defer wg.Done()
			providerResults, err := p.Search(ctx, query)
			if err != nil {
				return
			}
			mu.Lock()
			results = append(results, providerResults...)
			mu.Unlock()
		}()
	}

	wg.Wait()
	results = dedupeByURL(results)
	results = applySimpleRanking(results)
	return results, nil
}

func dedupeByURL(input []models.SearchResult) []models.SearchResult {
	seen := make(map[string]struct{}, len(input))
	output := make([]models.SearchResult, 0, len(input))
	for _, item := range input {
		key := strings.TrimSpace(strings.ToLower(item.URL))
		if key == "" {
			continue
		}
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		output = append(output, item)
	}
	return output
}

func applySimpleRanking(results []models.SearchResult) []models.SearchResult {
	weights := map[string]float64{"duckduckgo": 1.0, "wikipedia": 0.95, "github": 0.9}
	for i := range results {
		results[i].Score = weights[results[i].Source] - float64(i)*0.01
	}
	sort.SliceStable(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	return results
}
