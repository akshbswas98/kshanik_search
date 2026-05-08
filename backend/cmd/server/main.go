// Kshanik Search Engine — main entry point.
//
// This server bootstraps the meta-search engine by:
// 1. Loading configuration from environment / .env
// 2. Initializing all search providers (DuckDuckGo, Wikipedia, GitHub)
// 3. Creating the MetaSearchService (concurrent aggregator)
// 4. Setting up the HTTP API with middleware
// 5. Starting the server with graceful shutdown support
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/kshanik/search-engine/internal/api"
	"github.com/kshanik/search-engine/internal/provider"
	"github.com/kshanik/search-engine/internal/service"
	"github.com/rs/zerolog"
)

func main() {
	// --- Logger setup ---
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).
		With().
		Timestamp().
		Str("app", "kshanik-search").
		Logger()

	// --- Configuration from environment ---
	port := getEnv("PORT", "8080")
	providerTimeout := getEnvDuration("PROVIDER_TIMEOUT_MS", 5000)
	searchTimeout := getEnvDuration("SEARCH_TIMEOUT_MS", 10000)
	githubToken := os.Getenv("GITHUB_TOKEN") // Optional: increases GitHub API rate limit.

	// Memory limits for 1GB RAM optimization.
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		for range ticker.C {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			logger.Info().
				Uint64("alloc_mb", m.Alloc/1024/1024).
				Uint64("sys_mb", m.Sys/1024/1024).
				Uint32("num_gc", m.NumGC).
				Msg("memory stats")
		}
	}()

	logger.Info().
		Str("port", port).
		Dur("provider_timeout", providerTimeout).
		Dur("search_timeout", searchTimeout).
		Bool("github_auth", githubToken != "").
		Str("go_version", runtime.Version()).
		Msg("configuration loaded")

	// --- Initialize search providers ---
	providers := []provider.SearchProvider{
		provider.NewDuckDuckGoProvider(providerTimeout),
		provider.NewWikipediaProvider(providerTimeout),
		provider.NewGitHubProvider(providerTimeout, githubToken),
		provider.NewRedditProvider(providerTimeout),
		provider.NewStackOverflowProvider(providerTimeout),
	}

	logger.Info().Int("provider_count", len(providers)).Msg("search providers initialized")

	// --- Create the meta-search service ---
	searchService := service.NewMetaSearchService(providers, searchTimeout, logger)

	// --- Set up HTTP server ---
	mux := http.NewServeMux()

	handler := api.NewHandler(searchService, logger)
	handler.RegisterRoutes(mux)

	// Apply middleware chain: Recovery → Logging → Handler.
	var httpHandler http.Handler = mux
	httpHandler = api.LoggingMiddleware(logger)(httpHandler)
	httpHandler = api.RecoveryMiddleware(logger)(httpHandler)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      httpHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// --- Graceful shutdown ---
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Info().Str("addr", server.Addr).Msg("starting HTTP server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("server failed")
		}
	}()

	<-done
	logger.Info().Msg("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("server forced to shutdown")
	}

	logger.Info().Msg("server exited gracefully")
}

// getEnv reads an environment variable with a fallback default.
func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

// getEnvDuration reads an env var as milliseconds and returns a time.Duration.
func getEnvDuration(key string, defaultMs int) time.Duration {
	if val := os.Getenv(key); val != "" {
		if ms, err := strconv.Atoi(val); err == nil {
			return time.Duration(ms) * time.Millisecond
		}
	}
	return time.Duration(defaultMs) * time.Millisecond
}

func init() {
	// Load .env file if it exists (simple implementation).
	loadDotEnv(".env")
}

// loadDotEnv reads a .env file and sets environment variables.
// This is a minimal implementation — for production use consider godotenv.
func loadDotEnv(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return // .env is optional.
	}

	lines := splitLines(string(data))
	for _, line := range lines {
		line = trimSpace(line)
		if line == "" || line[0] == '#' {
			continue
		}
		idx := indexOf(line, '=')
		if idx < 0 {
			continue
		}
		key := trimSpace(line[:idx])
		value := trimSpace(line[idx+1:])

		// Strip surrounding quotes.
		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			line := s[start:i]
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
			lines = append(lines, line)
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func trimSpace(s string) string {
	i := 0
	for i < len(s) && (s[i] == ' ' || s[i] == '\t') {
		i++
	}
	j := len(s)
	for j > i && (s[j-1] == ' ' || s[j-1] == '\t') {
		j--
	}
	return s[i:j]
}

func indexOf(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}

func init() {
	// Placeholder: print banner.
	fmt.Println(`
╔═══════════════════════════════════════════════╗
║       🔍 Kshanik Search Engine v1.0          ║
║       Meta Search → Classical Search         ║
╚═══════════════════════════════════════════════╝`)
}
