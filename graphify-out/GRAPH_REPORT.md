# Graph Report - kshanik_search  (2026-05-01)

## Corpus Check
- 35 files · ~14,575 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 147 nodes · 166 edges · 15 communities detected
- Extraction: 87% EXTRACTED · 13% INFERRED · 0% AMBIGUOUS · INFERRED: 21 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Community 0|Community 0]]
- [[_COMMUNITY_Community 1|Community 1]]
- [[_COMMUNITY_Community 2|Community 2]]
- [[_COMMUNITY_Community 3|Community 3]]
- [[_COMMUNITY_Community 4|Community 4]]
- [[_COMMUNITY_Community 5|Community 5]]
- [[_COMMUNITY_Community 6|Community 6]]
- [[_COMMUNITY_Community 7|Community 7]]
- [[_COMMUNITY_Community 8|Community 8]]
- [[_COMMUNITY_Community 9|Community 9]]
- [[_COMMUNITY_Community 10|Community 10]]
- [[_COMMUNITY_Community 11|Community 11]]
- [[_COMMUNITY_Community 12|Community 12]]
- [[_COMMUNITY_Community 13|Community 13]]
- [[_COMMUNITY_Community 14|Community 14]]

## God Nodes (most connected - your core abstractions)
1. `main()` - 12 edges
2. `loadDotEnv()` - 6 edges
3. `Handler` - 6 edges
4. `truncate()` - 6 edges
5. `Index` - 5 edges
6. `DuckDuckGoProvider` - 5 edges
7. `extractTitle()` - 5 edges
8. `NewMetaSearchService()` - 5 edges
9. `RankResults()` - 5 edges
10. `getEnv()` - 4 edges

## Surprising Connections (you probably didn't know these)
- `main()` --calls--> `NewDuckDuckGoProvider()`  [INFERRED]
  backend\cmd\server\main.go → backend\internal\provider\duckduckgo.go
- `main()` --calls--> `NewWikipediaProvider()`  [INFERRED]
  backend\cmd\server\main.go → backend\internal\provider\wikipedia.go
- `main()` --calls--> `NewGitHubProvider()`  [INFERRED]
  backend\cmd\server\main.go → backend\internal\provider\github.go
- `main()` --calls--> `NewRedditProvider()`  [INFERRED]
  backend\cmd\server\main.go → backend\internal\provider\reddit.go
- `main()` --calls--> `NewStackOverflowProvider()`  [INFERRED]
  backend\cmd\server\main.go → backend\internal\provider\stackoverflow.go

## Communities

### Community 0 - "Community 0"
Cohesion: 0.15
Nodes (8): NewHandler(), TestHandler_HandleHealth(), TestHandler_HandleSearch(), MockProvider, MockProvider, providerResult, NewMetaSearchService(), TestMetaSearchService_Search()

### Community 1 - "Community 1"
Cohesion: 0.24
Nodes (11): LoggingMiddleware(), RecoveryMiddleware(), statusResponseWriter, getEnv(), getEnvDuration(), indexOf(), init(), loadDotEnv() (+3 more)

### Community 2 - "Community 2"
Cohesion: 0.27
Nodes (7): extractTitle(), NewDuckDuckGoProvider(), truncate(), DuckDuckGoProvider, duckDuckGoResponse, duckDuckGoResult, duckDuckGoTopic

### Community 3 - "Community 3"
Cohesion: 0.28
Nodes (5): NewGitHubProvider(), normalizeStars(), GitHubProvider, githubRepo, githubSearchResponse

### Community 4 - "Community 4"
Cohesion: 0.28
Nodes (5): NewRedditProvider(), normalizeEngagement(), redditPost, RedditProvider, redditSearchResponse

### Community 5 - "Community 5"
Cohesion: 0.36
Nodes (7): MetaSearchService, computeKeywordRelevance(), DeduplicateResults(), normalizeURL(), RankResults(), sortByScore(), tokenize()

### Community 6 - "Community 6"
Cohesion: 0.29
Nodes (4): NewStackOverflowProvider(), stackOverflowItem, StackOverflowProvider, stackOverflowResponse

### Community 7 - "Community 7"
Cohesion: 0.32
Nodes (4): NewWikipediaProvider(), parseOpenSearch(), WikipediaProvider, wikipediaSearchResponse

### Community 8 - "Community 8"
Cohesion: 0.29
Nodes (2): Document, Index

### Community 9 - "Community 9"
Cohesion: 0.6
Nodes (1): Handler

### Community 10 - "Community 10"
Cohesion: 0.33
Nodes (2): CrawlConfig, Crawler

### Community 11 - "Community 11"
Cohesion: 0.33
Nodes (2): RankConfig, Ranker

### Community 12 - "Community 12"
Cohesion: 0.4
Nodes (4): ErrorResponse, SearchRequest, SearchResponse, SearchResult

### Community 13 - "Community 13"
Cohesion: 0.67
Nodes (2): ResultContextProvider(), useResultContext()

### Community 14 - "Community 14"
Cohesion: 1.0
Nodes (1): SearchProvider

## Knowledge Gaps
- **19 isolated node(s):** `CrawlConfig`, `Document`, `SearchResult`, `SearchRequest`, `SearchResponse` (+14 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **Thin community `Community 8`** (7 nodes): `index.go`, `Document`, `Index`, `.AddDocument()`, `.Flush()`, `NewIndex()`, `.Search()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 9`** (6 nodes): `Handler`, `.handleHealth()`, `.handleSearch()`, `.RegisterRoutes()`, `.writeError()`, `.writeJSON()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 10`** (6 nodes): `crawler.go`, `CrawlConfig`, `Crawler`, `NewCrawler()`, `.Start()`, `.Stop()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 11`** (6 nodes): `ranking.go`, `RankConfig`, `Ranker`, `.RankDocuments()`, `DefaultConfig()`, `NewRanker()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 13`** (4 nodes): `ResultContextProvider()`, `useResultContext()`, `ResultsContextProvider.js`, `ResultsContextProvider.jsx`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 14`** (2 nodes): `provider.go`, `SearchProvider`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `main()` connect `Community 1` to `Community 0`, `Community 2`, `Community 3`, `Community 4`, `Community 6`, `Community 7`?**
  _High betweenness centrality (0.320) - this node is a cross-community bridge._
- **Why does `NewMetaSearchService()` connect `Community 0` to `Community 1`?**
  _High betweenness centrality (0.147) - this node is a cross-community bridge._
- **Why does `NewDuckDuckGoProvider()` connect `Community 2` to `Community 1`?**
  _High betweenness centrality (0.105) - this node is a cross-community bridge._
- **Are the 9 inferred relationships involving `main()` (e.g. with `NewDuckDuckGoProvider()` and `NewWikipediaProvider()`) actually correct?**
  _`main()` has 9 INFERRED edges - model-reasoned connections that need verification._
- **Are the 2 inferred relationships involving `truncate()` (e.g. with `.normalizeResults()` and `.normalizeResults()`) actually correct?**
  _`truncate()` has 2 INFERRED edges - model-reasoned connections that need verification._
- **What connects `CrawlConfig`, `Document`, `SearchResult` to the rest of the system?**
  _19 weakly-connected nodes found - possible documentation gaps or missing edges._