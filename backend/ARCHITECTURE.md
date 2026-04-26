# Kshanik Search Engine — Redesigned Hybrid Architecture

## 🏗️ Evolution: Meta-Search to Distributed Hybrid Search

The Kshanik Search Engine is evolving from a pure meta-search aggregator to a **Hybrid Search System**. This redesign leverages the high-memory ARM architecture of Oracle Cloud to support local indexing, distributed crawling, and semantic ranking.

### 🧩 System Components

1.  **Search Core (Go Monolith/Microservices):**
    *   **Meta-Aggregator:** Concurrent fan-out to external providers (DDG, GitHub, etc.).
    *   **Local-Search Indexer:** Queries a local high-performance index (Meilisearch or custom Inverted Index).
    *   **Hybrid Ranker:** Merges results using a weighted multi-factor algorithm (Source Reliability + Recency + Keyword Density).
2.  **The Crawler Pipeline (Asynchronous):**
    *   **URL Frontier:** A priority queue (Redis-backed) managing millions of URLs.
    *   **Distributed Crawler:** High-concurrency workers utilizing Go's `net/http` and `colly`.
    *   **Storage Layer:** Raw data in S3-compatible Object Storage (OCI Object Storage).
3.  **Data Processing:**
    *   **Extractor:** Transforms raw HTML into clean text and metadata.
    *   **Vector Engine (Future):** Generates embeddings for semantic search.

---

## 🔍 Detailed Architecture Diagram (Hybrid Phase)

```
┌─────────────────────────────────────────────────────────────────┐
│                        React Frontend                          │
│                   (Vite / Netlify / Vercel)                    │
└──────────────────────┬──────────────────────────────────────────┘
                       │ HTTPS (Caddy / Cloudflare)
                       ▼
┌─────────────────────────────────────────────────────────────────┐
│              Oracle Cloud ARM Instance (Go)                     │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                    API Gateway                           │   │
│  │         (Auth, Rate Limiting, Caching: Redis)            │   │
│  └──────────────────────┬───────────────────────────────────┘   │
│                         │                                       │
│  ┌──────────────────────▼───────────────────────────────────┐   │
│  │                Hybrid Search Orchestrator                 │   │
│  │                                                           │   │
│  │  ┌──────────────┐         ┌───────────────┐               │   │
│  │  │ Meta-Search  │         │  Local Index  │               │   │
│  │  │ (External)   │         │ (Meilisearch) │               │   │
│  │  └──────┬───────┘         └───────┬───────┘               │   │
│  │         │                         │                       │   │
│  │         └──────────┬──────────────┘                       │   │
│  │                    ▼                                      │   │
│  │           Hybrid Ranking Engine                           │   │
│  │    (BM25 + Cross-Source Deduplication)                    │   │
│  └────────────────────┬─────────────────────────────────────┘   │
│                       │                                         │
│  ┌────────────────────▼─────────────────────────────────────┐   │
│  │              Background Crawler System                   │   │
│  │  (Managed by Go Workers + OCI Object Storage)             │   │
│  └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

## ⚡ Recomputed Performance Metrics (ARM A1 Flex)

By moving to this hybrid architecture on OCI Ampere (ARM64), we achieve:

| Metric | Target | Computation / Reason |
|--------|--------|----------------------|
| **Search Latency** | < 300ms | 80% of top queries served from Redis cache / Local Index. |
| **Concurrency** | 10k+ Req/s | ARM instances offer better vertical scaling for Go's M:N scheduler. |
| **Cost** | $0/mo | Optimized to stay within OCI "Always Free" (4 OCPU / 24GB RAM). |
| **Memory Efficiency** | 24GB Cap | Ample RAM for high-speed Inverted Index and Crawler Frontier. |

---

## 🛠️ Internal Package Structure

```
backend/
├── cmd/server/main.go          ← Entry point
├── internal/
│   ├── api/                    ← HTTP Handlers & Middleware
│   ├── hybrid/                 ← Hybrid Search Logic (The "Brain")
│   ├── provider/               ← External Meta-Search Providers
│   ├── crawler/                ← URL Frontier & Scraping Logic
│   ├── index/                  ← Meilisearch/Bleve local indexing
│   ├── ranking/                ← Advanced Scoring (BM25, PR)
│   ├── models/                 ← Domain DTOs
│   └── cache/                  ← Redis integration
```

## 🚀 Future: Semantic Vector Search

The architecture is designed to integrate **Milvus** or **Pinecone** for vector-based semantic search, allowing users to find "stories that define history" not just by keywords, but by meaning.
