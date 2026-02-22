package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port            string
	RequestTimeout  time.Duration
	ProviderTimeout time.Duration
	MaxResults      int
	DuckDuckGoURL   string
	WikipediaURL    string
	GitHubURL       string
	GitHubToken     string
}

func Load() (Config, error) {
	requestTimeout, err := durationFromEnv("REQUEST_TIMEOUT", 8*time.Second)
	if err != nil {
		return Config{}, err
	}
	providerTimeout, err := durationFromEnv("PROVIDER_TIMEOUT", 4*time.Second)
	if err != nil {
		return Config{}, err
	}
	maxResults, err := intFromEnv("MAX_RESULTS_PER_PROVIDER", 5)
	if err != nil {
		return Config{}, err
	}

	return Config{
		Port:            stringFromEnv("PORT", "8080"),
		RequestTimeout:  requestTimeout,
		ProviderTimeout: providerTimeout,
		MaxResults:      maxResults,
		DuckDuckGoURL:   stringFromEnv("DUCKDUCKGO_BASE_URL", "https://api.duckduckgo.com"),
		WikipediaURL:    stringFromEnv("WIKIPEDIA_BASE_URL", "https://en.wikipedia.org"),
		GitHubURL:       stringFromEnv("GITHUB_BASE_URL", "https://api.github.com"),
		GitHubToken:     os.Getenv("GITHUB_TOKEN"),
	}, nil
}

func stringFromEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func intFromEnv(key string, fallback int) (int, error) {
	if v := os.Getenv(key); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return 0, fmt.Errorf("invalid %s: %w", key, err)
		}
		return parsed, nil
	}
	return fallback, nil
}

func durationFromEnv(key string, fallback time.Duration) (time.Duration, error) {
	if v := os.Getenv(key); v != "" {
		parsed, err := time.ParseDuration(v)
		if err != nil {
			return 0, fmt.Errorf("invalid %s: %w", key, err)
		}
		return parsed, nil
	}
	return fallback, nil
}
