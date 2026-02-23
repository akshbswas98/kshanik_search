# Kshanik Search Engine — Architecture & Documentation

## 🏗️ Project Structure

```
kshanik_search/
├── backend/                          ← Go backend (meta-search engine)
│   ├── cmd/
│   │   └── server/
│   │       └── main.go               ← Entry point: config, DI, server startup
│   ├── internal/
│   │   ├── api/
│   │   │   ├── handler.go            ← HTTP handlers (GET /search, GET /health)
│   │   │   └── middleware.go          ← Request logging + panic recovery
│   │   ├── models/
│   │   │   └── search.go             ← DTOs: SearchResult, SearchRequest, SearchResponse
│   │   ├── provider/
│   │   │   ├── provider.go           ← SearchProvider interface
│   │   │   ├── duckduckgo.go         ← DuckDuckGo Instant Answer API
│   │   │   ├── wikipedia.go          ← Wikipedia OpenSearch API
│   │   │   └── github.go             ← GitHub Repository Search API
│   │   ├── service/
│   │   │   └── search_service.go     ← MetaSearchService: concurrent fan-out + merge
│   │   ├── utils/
│   │   │   └── helpers.go            ← Dedup, ranking, tokenization
│   │   ├── crawler/
│   │   │   └── crawler.go            ← 🔮 Stub: future web crawler
│   │   ├── index/
│   │   │   └── index.go              ← 🔮 Stub: future inverted index
│   │   └── ranking/
│   │       └── ranking.go            ← 🔮 Stub: future BM25/semantic ranking
│   ├── .env                          ← Local configuration
│   ├── .env.example                  ← Configuration template
│   ├── go.mod
│   └── go.sum
├── src/                              ← React frontend (existing)
│   ├── components/
│   │   ├── Results.jsx               ← Updated to call Go backend
│   │   └── ...
│   └── contexts/
│       └── ResultsContextProvider.jsx ← Updated to call Go backend
├── vite.config.js                    ← Proxy → localhost:8080
└── package.json
```

## 🔍 Architecture Overview

### Current State: Meta Search Engine (Phase 1)

```
┌─────────────────────────────────────────────────────────────────┐
│                        React Frontend                          │
│                   (Vite dev server :3000)                       │
└──────────────────────┬──────────────────────────────────────────┘
                       │ GET /api/search?q=...
                       │ (proxy strips /api prefix)
                       ▼
┌─────────────────────────────────────────────────────────────────┐
│                Go HTTP Server (:8080)                           │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │            Middleware Chain                               │   │
│  │  Recovery → Logging → Router                             │   │
│  └──────────────────────┬───────────────────────────────────┘   │
│                         │                                       │
│  ┌──────────────────────▼───────────────────────────────────┐   │
│  │            API Handler                                    │   │
│  │  GET /search?q=...  → validates query                    │   │
│  │  GET /health        → health check                       │   │
│  └──────────────────────┬───────────────────────────────────┘   │
│                         │                                       │
│  ┌──────────────────────▼───────────────────────────────────┐   │
│  │         MetaSearchService                                 │   │
│  │                                                           │   │
│  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐        │   │
│  │  │ DuckDuckGo  │ │  Wikipedia  │ │   GitHub    │        │   │
│  │  │  Provider   │ │  Provider   │ │  Provider   │        │   │
│  │  └──────┬──────┘ └──────┬──────┘ └──────┬──────┘        │   │
│  │         │               │               │                │   │
│  │         └───────────────┼───────────────┘                │   │
│  │                         │                                 │   │
│  │              ┌──────────▼──────────┐                      │   │
│  │              │  Concurrent Fan-out │                      │   │
│  │              │  (goroutines + WG)  │                      │   │
│  │              └──────────┬──────────┘                      │   │
│  │                         │                                 │   │
│  │              ┌──────────▼──────────┐                      │   │
│  │              │   Deduplication     │                      │   │
│  │              │   (URL hashing)     │                      │   │
│  │              └──────────┬──────────┘                      │   │
│  │                         │                                 │   │
│  │              ┌──────────▼──────────┐                      │   │
│  │              │     Ranking         │                      │   │
│  │              │ provider weight +   │                      │   │
│  │              │ keyword relevance   │                      │   │
│  │              └──────────┬──────────┘                      │   │
│  │                         │                                 │   │
│  └─────────────────────────┼────────────────────────────────┘   │
│                            ▼                                     │
│               JSON Response { results: [...] }                   │
└──────────────────────────────────────────────────────────────────┘
```

### Evolution Path: Meta Search → Classical Search Engine

```
╔══════════════════════════════════════════════════════════════════════════╗
║                          EVOLUTION ROADMAP                              ║
╠══════════════════════════════════════════════════════════════════════════╣
║                                                                          ║
║  PHASE 1 (CURRENT) — Meta Search Engine                                 ║
║  ├── External API aggregation (DDG + Wikipedia + GitHub)                ║
║  ├── Concurrent fan-out with goroutines                                 ║
║  ├── URL deduplication + keyword-weighted ranking                       ║
║  └── Clean provider interface for extensibility                         ║
║                                                                          ║
║  PHASE 2 — Focused Web Crawler                                          ║
║  ├── URL frontier (priority queue)                                      ║
║  ├── robots.txt compliance                                              ║
║  ├── Politeness delays + rate limiting                                  ║
║  ├── Raw page storage (disk / blob storage)                             ║
║  └── Feeds into indexing pipeline                                       ║
║                                                                          ║
║  PHASE 3 — Inverted Index + BM25 Ranking                                ║
║  ├── Tokenization & normalization pipeline                              ║
║  ├── Term frequency / document frequency tracking                       ║
║  ├── Positional index for phrase queries                                ║
║  ├── BM25 scoring (k1=1.2, b=0.75)                                     ║
║  └── Optional: Elasticsearch / Bleve integration                       ║
║                                                                          ║
║  PHASE 4 — AI Semantic Search                                           ║
║  ├── Sentence embeddings (via ONNX Runtime / external API)              ║
║  ├── Vector index (HNSW / FAISS)                                        ║
║  ├── Cosine similarity re-ranking                                       ║
║  └── Hybrid lexical + semantic scoring                                  ║
║                                                                          ║
╚══════════════════════════════════════════════════════════════════════════╝
```

## 🚀 How to Run

### Backend (Go)

```bash
cd backend

# Copy and configure environment
cp .env.example .env

# Build
go build -o bin/kshanik-search.exe ./cmd/server/

# Run
./bin/kshanik-search.exe

# Or simply:
go run ./cmd/server/
```

The server starts at `http://localhost:8080`.

### Frontend (React + Vite)

```bash
# In the project root
npm install
npm run dev
```

The frontend runs at `http://localhost:3000` and proxies `/api/*` to the Go backend.

### Test the API directly

```bash
curl "http://localhost:8080/search?q=golang"
curl "http://localhost:8080/health"
```

## ⚡ Concurrency Benefits vs Java Spring Boot

| Aspect | Go (this project) | Java Spring Boot |
|--------|-------------------|------------------|
| **Goroutines** | Extremely lightweight (~2KB stack). We spin up one per provider with zero overhead. | Threads are expensive (~1MB stack). Requires thread pools, `@Async`, or reactive `WebFlux`. |
| **Fan-out pattern** | Native: `go func()` + `sync.WaitGroup` + channels. 5 lines of code. | Requires `CompletableFuture.supplyAsync()` or `Mono.zip()` with reactor. Significant boilerplate. |
| **Memory** | The entire binary is ~10MB. A search query uses ~100KB. | JVM baseline is 200-500MB. Spring context adds another 100MB+. |
| **Startup time** | ~50ms cold start. | 5-15 seconds for Spring Boot context initialization. |
| **Context propagation** | `context.Context` is first-class — flows through every function. Timeouts and cancellation propagate automatically. | `RequestScope`, `ThreadLocal`, or reactive `Context` — error-prone across thread boundaries. |
| **Error handling** | Explicit `(result, error)` tuples. No hidden exceptions. Every failure path is visible. | Exception-based — easy to miss unchecked exceptions. `@ControllerAdvice` adds indirection. |
| **Deployment** | Single static binary. No runtime dependencies. `COPY binary → docker`. | Needs JRE/JDK, fat JAR or layered image. Docker images are 200MB+. |
| **HTTP server** | Built into stdlib (`net/http`). Zero dependencies for a production-grade server. | Embedded Tomcat/Netty with dozens of transitive dependencies. |

### Why goroutines shine for meta-search

In meta-search, every query triggers N parallel API calls. With goroutines:

```go
// This is the ENTIRE concurrency code:
for _, p := range providers {
    wg.Add(1)
    go func(p SearchProvider) {
        defer wg.Done()
        results, err := p.Search(ctx, query)
        ch <- providerResult{results, p.Name(), err}
    }(p)
}
```

- **No thread pool tuning** — the Go scheduler handles 100K+ goroutines.
- **No callback hell** — goroutines block naturally.
- **Context cancellation** — if the client disconnects, all provider calls are cancelled automatically.
- **Partial failure tolerance** — we collect whatever succeeds and serve partial results.

## 🔧 Configuration

| Environment Variable | Default | Description |
|---------------------|---------|-------------|
| `PORT` | `8080` | HTTP server port |
| `PROVIDER_TIMEOUT_MS` | `5000` | Timeout per provider API call |
| `SEARCH_TIMEOUT_MS` | `10000` | Timeout for entire search operation |
| `GITHUB_TOKEN` | _(empty)_ | GitHub API token for higher rate limits |

## 📡 API Reference

### `GET /search?q={query}`

Returns aggregated, deduplicated, ranked search results.

**Response:**
```json
{
  "results": [
    {
      "title": "Go (programming language)",
      "snippet": "Go is a statically typed, compiled...",
      "url": "https://en.wikipedia.org/wiki/Go_(programming_language)",
      "source": "wikipedia",
      "score": 0.85,
      "timestamp": "2026-02-22T19:30:00Z"
    }
  ],
  "total_count": 15,
  "query": "golang",
  "time_taken": "1.234s"
}
```

### `GET /health`

```json
{
  "status": "healthy",
  "engine": "kshanik-meta-search"
}
```

## 🧩 Adding a New Provider

1. Create a new file in `internal/provider/` (e.g., `brave.go`).
2. Implement the `SearchProvider` interface:

```go
type BraveProvider struct { ... }

func (b *BraveProvider) Name() string { return "brave" }
func (b *BraveProvider) Search(ctx context.Context, query string) ([]models.SearchResult, error) {
    // Your implementation
}
```

3. Register it in `cmd/server/main.go`:

```go
providers := []provider.SearchProvider{
    provider.NewDuckDuckGoProvider(providerTimeout),
    provider.NewWikipediaProvider(providerTimeout),
    provider.NewGitHubProvider(providerTimeout, githubToken),
    provider.NewBraveProvider(providerTimeout, braveAPIKey), // ← Add here
}
```

4. Optionally add a weight in `utils/helpers.go`:

```go
var ProviderWeights = map[string]float64{
    "brave": 0.85,
}
```

That's it. No other code changes needed.
