package api

import (
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	defaultMaxQueryLen   = 200
	defaultRatePerMinute = 60
)

// SecurityConfig holds runtime security settings from environment variables.
type SecurityConfig struct {
	AllowedOrigins   []string
	MaxQueryLength   int
	RatePerMinute    int
	ExposeErrorDetail bool
}

// LoadSecurityConfig reads security-related environment variables.
func LoadSecurityConfig() SecurityConfig {
	cfg := SecurityConfig{
		MaxQueryLength:    defaultMaxQueryLen,
		RatePerMinute:     defaultRatePerMinute,
		ExposeErrorDetail: !isProductionEnv(),
	}

	if raw := strings.TrimSpace(os.Getenv("ALLOWED_ORIGINS")); raw != "" {
		for _, o := range strings.Split(raw, ",") {
			if trimmed := strings.TrimSpace(o); trimmed != "" {
				cfg.AllowedOrigins = append(cfg.AllowedOrigins, trimmed)
			}
		}
	}

	if v := os.Getenv("MAX_QUERY_LENGTH"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.MaxQueryLength = n
		}
	}

	if v := os.Getenv("RATE_LIMIT_PER_MINUTE"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.RatePerMinute = n
		}
	}

	if os.Getenv("EXPOSE_ERROR_DETAIL") == "true" {
		cfg.ExposeErrorDetail = true
	}
	if os.Getenv("EXPOSE_ERROR_DETAIL") == "false" {
		cfg.ExposeErrorDetail = false
	}

	return cfg
}

func isProductionEnv() bool {
	return strings.EqualFold(os.Getenv("ENV"), "production")
}

// SecurityHeadersMiddleware adds baseline HTTP security headers for every response.
func SecurityHeadersMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			w.Header().Set("Permissions-Policy", "interest-cohort=()")
			w.Header().Set("X-Robots-Tag", "noindex, nofollow")
			next.ServeHTTP(w, r)
		})
	}
}

// CORSMiddleware restricts cross-origin access to configured frontend origins.
// In non-production, if ALLOWED_ORIGINS is unset, reflects any Origin (local dev).
func CORSMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	allowSet := make(map[string]struct{}, len(allowedOrigins))
	for _, o := range allowedOrigins {
		allowSet[o] = struct{}{}
	}

	devPermissive := len(allowSet) == 0 && !isProductionEnv()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" {
				if _, ok := allowSet[origin]; ok || devPermissive {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Vary", "Origin")
					w.Header().Add("Access-Control-Allow-Methods", "GET, OPTIONS")
					w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
				}
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// rateLimiter provides a simple per-IP fixed-window counter.
type rateLimiter struct {
	mu      sync.Mutex
	limit   int
	window  time.Duration
	buckets map[string][]time.Time
}

func newRateLimiter(perMinute int) *rateLimiter {
	return &rateLimiter{
		limit:   perMinute,
		window:  time.Minute,
		buckets: make(map[string][]time.Time),
	}
}

func (rl *rateLimiter) allow(key string) bool {
	now := time.Now()
	cutoff := now.Add(-rl.window)

	rl.mu.Lock()
	defer rl.mu.Unlock()

	times := rl.buckets[key]
	filtered := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}

	if len(filtered) >= rl.limit {
		rl.buckets[key] = filtered
		return false
	}

	rl.buckets[key] = append(filtered, now)
	return true
}

// RateLimitMiddleware limits requests per client IP to reduce abuse.
func RateLimitMiddleware(perMinute int) func(http.Handler) http.Handler {
	limiter := newRateLimiter(perMinute)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/health" {
				next.ServeHTTP(w, r)
				return
			}

			ip := clientIP(r)
			if !limiter.allow(ip) {
				w.Header().Set("Retry-After", "60")
				http.Error(w, `{"error":"rate limit exceeded","code":429}`, http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
