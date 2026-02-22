# Go Backend Architecture

## Structure

- `cmd/server/main.go`: composition root and HTTP server bootstrap.
- `internal/api`: HTTP handlers and request/response plumbing.
- `internal/service`: orchestration for meta-search aggregation + ranking hooks.
- `internal/provider`: provider adapters (DuckDuckGo, Wikipedia, GitHub).
- `internal/models`: shared domain models.
- `internal/crawler`: crawler extension-point stubs.
- `internal/index`: indexing extension-point stubs.
- `internal/config`: environment-based configuration loader.

## Runtime flow

1. `GET /search?q=...` hits API handler.
2. Handler validates input and attaches request/provider timeouts through context.
3. `SearchService` fan-outs provider searches concurrently using goroutines + `sync.WaitGroup`.
4. Aggregated results are deduplicated by URL, scored with a simple ranking heuristic, then returned as JSON.
5. Future ranking, crawler, and index layers can be injected without rewriting API/provider layers.

## Environment variables

- `PORT` (default `8080`)
- `REQUEST_TIMEOUT` (default `8s`)
- `PROVIDER_TIMEOUT` (default `4s`)
- `MAX_RESULTS_PER_PROVIDER` (default `5`)
- `DUCKDUCKGO_BASE_URL` (default `https://api.duckduckgo.com`)
- `WIKIPEDIA_BASE_URL` (default `https://en.wikipedia.org`)
- `GITHUB_BASE_URL` (default `https://api.github.com`)
- `GITHUB_TOKEN` (optional; improves GitHub API limits)
