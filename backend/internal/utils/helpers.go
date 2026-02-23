// Package utils provides shared helper functions for deduplication, scoring,
// and other common operations across the search engine.
package utils

import (
	"math"
	"strings"

	"github.com/kshanik/search-engine/internal/models"
)

// ProviderWeights defines how much we trust each source.
// Higher weight = higher priority in final ranking.
var ProviderWeights = map[string]float64{
	"duckduckgo":    1.0,
	"wikipedia":     0.9,
	"stackoverflow": 0.85,
	"reddit":        0.75,
	"github":        0.7,
}

// DeduplicateResults removes duplicate results based on normalized URL.
// When duplicates are found, the one with the higher score is kept.
func DeduplicateResults(results []models.SearchResult) []models.SearchResult {
	seen := make(map[string]int) // normalized URL → index in deduped slice
	var deduped []models.SearchResult

	for _, r := range results {
		key := normalizeURL(r.URL)
		if idx, exists := seen[key]; exists {
			// Keep the higher-scored entry.
			if r.Score > deduped[idx].Score {
				deduped[idx] = r
			}
		} else {
			seen[key] = len(deduped)
			deduped = append(deduped, r)
		}
	}

	return deduped
}

// normalizeURL strips protocol schems, trailing slashes, and lowercases for dedup.
func normalizeURL(rawURL string) string {
	u := strings.ToLower(rawURL)
	u = strings.TrimPrefix(u, "https://")
	u = strings.TrimPrefix(u, "http://")
	u = strings.TrimPrefix(u, "www.")
	u = strings.TrimRight(u, "/")
	return u
}

// RankResults applies a combined ranking score to each result:
//   - Provider weight (source trust)
//   - Keyword match relevance (title + snippet vs query)
//
// Results are sorted in descending order by computed score.
func RankResults(results []models.SearchResult, query string) []models.SearchResult {
	queryTerms := tokenize(query)

	for i := range results {
		providerWeight := ProviderWeights[results[i].Source]
		if providerWeight == 0 {
			providerWeight = 0.5 // Default for unknown providers.
		}

		keywordScore := computeKeywordRelevance(results[i], queryTerms)

		// Combined score: 40% provider weight, 30% keyword relevance, 30% original score.
		results[i].Score = 0.4*providerWeight + 0.3*keywordScore + 0.3*results[i].Score
	}

	// Sort descending by score (simple insertion sort — result sets are small).
	sortByScore(results)

	return results
}

// computeKeywordRelevance scores how well a result matches the query terms.
func computeKeywordRelevance(result models.SearchResult, queryTerms []string) float64 {
	if len(queryTerms) == 0 {
		return 0
	}

	titleLower := strings.ToLower(result.Title)
	snippetLower := strings.ToLower(result.Snippet)

	matches := 0
	for _, term := range queryTerms {
		if strings.Contains(titleLower, term) {
			matches += 2 // Title matches are worth more.
		}
		if strings.Contains(snippetLower, term) {
			matches++
		}
	}

	// Normalize: max possible = 3 * len(queryTerms)
	maxScore := float64(3 * len(queryTerms))
	return math.Min(float64(matches)/maxScore, 1.0)
}

// tokenize splits a query into lowercase terms, filtering out stop words.
func tokenize(query string) []string {
	stopWords := map[string]bool{
		"a": true, "an": true, "the": true, "is": true, "are": true,
		"was": true, "were": true, "be": true, "been": true, "being": true,
		"have": true, "has": true, "had": true, "do": true, "does": true,
		"did": true, "will": true, "would": true, "could": true, "should": true,
		"may": true, "might": true, "shall": true, "can": true,
		"for": true, "and": true, "nor": true, "but": true, "or": true,
		"yet": true, "so": true, "in": true, "on": true, "at": true,
		"to": true, "of": true, "with": true, "by": true, "from": true,
		"it": true, "its": true, "this": true, "that": true, "what": true,
		"how": true, "who": true, "where": true, "when": true, "why": true,
	}

	words := strings.Fields(strings.ToLower(query))
	var terms []string
	for _, w := range words {
		w = strings.Trim(w, ".,!?;:'\"()[]{}") // Strip punctuation.
		if w != "" && !stopWords[w] {
			terms = append(terms, w)
		}
	}
	return terms
}

// sortByScore sorts results in descending order by Score.
func sortByScore(results []models.SearchResult) {
	// Using insertion sort — result sets are typically < 50 items.
	for i := 1; i < len(results); i++ {
		key := results[i]
		j := i - 1
		for j >= 0 && results[j].Score < key.Score {
			results[j+1] = results[j]
			j--
		}
		results[j+1] = key
	}
}
