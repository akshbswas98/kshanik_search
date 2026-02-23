// Package crawler is a placeholder for the future web crawling subsystem.
//
// Evolution path:
//
//	Phase 1 (current) — Meta search only, no crawling.
//	Phase 2           — Implement a focused crawler that:
//	  - Maintains a URL frontier (priority queue)
//	  - Respects robots.txt
//	  - Implements politeness delays
//	  - Stores raw pages for the indexing pipeline
//	Phase 3           — Distributed crawling with work partitioning.
//
// This module compiles but all operations are no-ops.
package crawler

import (
	"context"

	"github.com/rs/zerolog"
)

// Crawler is the future web crawler.
type Crawler struct {
	logger zerolog.Logger
}

// NewCrawler creates a new Crawler instance.
func NewCrawler(logger zerolog.Logger) *Crawler {
	return &Crawler{
		logger: logger.With().Str("component", "crawler").Logger(),
	}
}

// CrawlConfig holds configuration for the crawler (future use).
type CrawlConfig struct {
	SeedURLs        []string
	MaxDepth        int
	PolitenessDelay int // milliseconds
	MaxPages        int
	RespectRobots   bool
	UserAgent       string
}

// Start begins the crawling process with the given config.
// Currently a no-op placeholder.
func (c *Crawler) Start(ctx context.Context, config CrawlConfig) error {
	c.logger.Info().Msg("crawler module is a placeholder — not yet implemented")
	return nil
}

// Stop gracefully shuts down the crawler.
func (c *Crawler) Stop() error {
	c.logger.Info().Msg("crawler stop called (no-op)")
	return nil
}
