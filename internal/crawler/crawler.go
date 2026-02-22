package crawler

import "context"

// Crawler is a future extension point for web crawling and document extraction.
type Crawler interface {
	Crawl(ctx context.Context, seeds []string) error
}
