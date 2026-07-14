package main

import (
	"log"
	"net/http"
	"time"

	"kshanik_search/internal/api"
	"kshanik_search/internal/config"
	"kshanik_search/internal/provider"
	"kshanik_search/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	httpClient := &http.Client{Timeout: cfg.ProviderTimeout}
	providers := []provider.SearchProvider{
		provider.NewDuckDuckGoProvider(httpClient, cfg.DuckDuckGoURL, cfg.MaxResults),
		provider.NewWikipediaProvider(httpClient, cfg.WikipediaURL, cfg.MaxResults),
		provider.NewGitHubProvider(httpClient, cfg.GitHubURL, cfg.MaxResults, cfg.GitHubToken),
	}

	searchService := service.NewSearchService(providers)
	handler := api.NewHandler(searchService, cfg.RequestTimeout, cfg.ProviderTimeout)

	mux := http.NewServeMux()
	handler.Register(mux)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("search API listening on :%s", cfg.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}
