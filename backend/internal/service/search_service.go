// Package service contains the business-logic layer that orchestrates
// search providers and produces ranked, deduplicated results.
package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/kshanik/search-engine/internal/models"
	"github.com/kshanik/search-engine/internal/provider"
	"github.com/kshanik/search-engine/internal/utils"
	"github.com/rs/zerolog"
)

// MetaSearchService fans out queries to all registered providers,
// merges the results, deduplicates, and ranks them.
type MetaSearchService struct {
	providers []provider.SearchProvider
	timeout   time.Duration
	logger    zerolog.Logger
}

// NewMetaSearchService creates a MetaSearchService with the given providers and timeout.
func NewMetaSearchService(providers []provider.SearchProvider, timeout time.Duration, logger zerolog.Logger) *MetaSearchService {
	return &MetaSearchService{
		providers: providers,
		timeout:   timeout,
		logger:    logger.With().Str("component", "meta_search_service").Logger(),
	}
}

// providerResult holds the output of a single provider call.
type providerResult struct {
	results []models.SearchResult
	source  string
	err     error
}

// Search fans out the query to all providers concurrently, aggregates
// results, deduplicates by URL, and applies ranking.
func (s *MetaSearchService) Search(ctx context.Context, query string) (*models.SearchResponse, error) {
	start := time.Now()

	// Create a timeout context for the entire fan-out.
	searchCtx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	resultsCh := make(chan providerResult, len(s.providers))
	var wg sync.WaitGroup

	// Fan out to all providers concurrently.
	for _, p := range s.providers {
		wg.Add(1)
		go func(p provider.SearchProvider) {
			defer wg.Done()
			s.logger.Info().Str("provider", p.Name()).Str("query", query).Msg("querying provider")

			results, err := p.Search(searchCtx, query)
			resultsCh <- providerResult{
				results: results,
				source:  p.Name(),
				err:     err,
			}
		}(p)
	}

	// Close channel when all goroutines are done.
	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	// Collect results from all providers.
	var allResults []models.SearchResult
	var errors []string

	for pr := range resultsCh {
		if pr.err != nil {
			s.logger.Warn().Err(pr.err).Str("provider", pr.source).Msg("provider returned error")
			errors = append(errors, fmt.Sprintf("%s: %v", pr.source, pr.err))
			continue
		}
		s.logger.Info().
			Str("provider", pr.source).
			Int("count", len(pr.results)).
			Msg("received results")
		allResults = append(allResults, pr.results...)
	}

	// If ALL providers failed, return an error.
	if len(allResults) == 0 && len(errors) > 0 {
		return nil, fmt.Errorf("all providers failed: %v", errors)
	}

	// Deduplicate by URL.
	deduped := utils.DeduplicateResults(allResults)

	// Rank results.
	ranked := utils.RankResults(deduped, query)

	elapsed := time.Since(start)

	response := &models.SearchResponse{
		Results:    ranked,
		TotalCount: len(ranked),
		Query:      query,
		TimeTaken:  elapsed.String(),
	}

	s.logger.Info().
		Int("total_results", len(ranked)).
		Dur("elapsed", elapsed).
		Msg("search completed")

	return response, nil
}
