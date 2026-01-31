package grid

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/machinebox/graphql"
	"golang.org/x/time/rate"
)

// API endpoints
const (
	CentralDataURL  = "https://api-op.grid.gg/central-data/graphql"
	SeriesStateURL  = "https://api-op.grid.gg/live-data-feed/series-state/graphql"
	FileDownloadURL = "https://api.grid.gg"
)

// Rate limits (requests per minute)
const (
	CentralDataRateLimit  = 40
	SeriesStateRateLimit  = 1200
	FileDownloadRateLimit = 20
)

// Cache interface for optional caching
type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
}

// Client is the main GRID API client
type Client struct {
	apiKey string
	cache  Cache

	// GraphQL clients
	centralDataClient *graphql.Client
	seriesStateClient *graphql.Client

	// HTTP client for file downloads
	httpClient *http.Client

	// Rate limiters
	centralDataLimiter  *rate.Limiter
	seriesStateLimiter  *rate.Limiter
	fileDownloadLimiter *rate.Limiter

	// Per-series rate limiters (75 req/min per series)
	seriesLimiters   map[string]*rate.Limiter
	seriesLimitersMu sync.RWMutex
}

// NewClient creates a new GRID API client
func NewClient(apiKey string, cache Cache) *Client {
	// Create GraphQL clients with custom HTTP client for auth
	centralDataClient := graphql.NewClient(CentralDataURL)
	seriesStateClient := graphql.NewClient(SeriesStateURL)

	return &Client{
		apiKey:              apiKey,
		cache:               cache,
		centralDataClient:   centralDataClient,
		seriesStateClient:   seriesStateClient,
		httpClient:          &http.Client{Timeout: 30 * time.Second},
		centralDataLimiter:  rate.NewLimiter(rate.Every(time.Minute/CentralDataRateLimit), 1),
		seriesStateLimiter:  rate.NewLimiter(rate.Every(time.Minute/SeriesStateRateLimit), 10),
		fileDownloadLimiter: rate.NewLimiter(rate.Every(time.Minute/FileDownloadRateLimit), 1),
		seriesLimiters:      make(map[string]*rate.Limiter),
	}
}

// getSeriesLimiter returns a rate limiter for a specific series
func (c *Client) getSeriesLimiter(seriesID string) *rate.Limiter {
	c.seriesLimitersMu.RLock()
	limiter, exists := c.seriesLimiters[seriesID]
	c.seriesLimitersMu.RUnlock()

	if exists {
		return limiter
	}

	c.seriesLimitersMu.Lock()
	defer c.seriesLimitersMu.Unlock()

	// Double-check after acquiring write lock
	if limiter, exists = c.seriesLimiters[seriesID]; exists {
		return limiter
	}

	// 75 requests per minute per series
	limiter = rate.NewLimiter(rate.Every(time.Minute/75), 1)
	c.seriesLimiters[seriesID] = limiter
	return limiter
}

// runCentralDataQuery executes a GraphQL query against Central Data API
func (c *Client) runCentralDataQuery(ctx context.Context, req *graphql.Request, resp interface{}) error {
	// Wait for rate limiter
	if err := c.centralDataLimiter.Wait(ctx); err != nil {
		return fmt.Errorf("rate limit wait: %w", err)
	}

	// Set auth header
	req.Header.Set("x-api-key", c.apiKey)

	return c.centralDataClient.Run(ctx, req, resp)
}

// runSeriesStateQuery executes a GraphQL query against Series State API
func (c *Client) runSeriesStateQuery(ctx context.Context, seriesID string, req *graphql.Request, resp interface{}) error {
	// Wait for overall rate limiter
	if err := c.seriesStateLimiter.Wait(ctx); err != nil {
		return fmt.Errorf("rate limit wait: %w", err)
	}

	// Wait for per-series rate limiter
	seriesLimiter := c.getSeriesLimiter(seriesID)
	if err := seriesLimiter.Wait(ctx); err != nil {
		return fmt.Errorf("series rate limit wait: %w", err)
	}

	// Set auth header
	req.Header.Set("x-api-key", c.apiKey)

	return c.seriesStateClient.Run(ctx, req, resp)
}

// doFileDownloadRequest executes an HTTP request against File Download API
func (c *Client) doFileDownloadRequest(ctx context.Context, url string) (*http.Response, error) {
	// Wait for rate limiter
	if err := c.fileDownloadLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("x-api-key", c.apiKey)

	return c.httpClient.Do(req)
}

// withRetry executes a function with exponential backoff retry
func withRetry(ctx context.Context, maxRetries int, fn func() error) error {
	var lastErr error
	backoff := time.Second

	for i := 0; i < maxRetries; i++ {
		if err := fn(); err != nil {
			lastErr = err

			// Check if context is cancelled
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			// Wait with exponential backoff
			time.Sleep(backoff)
			backoff *= 2
			continue
		}
		return nil
	}

	return fmt.Errorf("max retries exceeded: %w", lastErr)
}

// getCached retrieves data from cache if available
func (c *Client) getCached(ctx context.Context, key string, dest interface{}) bool {
	if c.cache == nil {
		return false
	}

	data, err := c.cache.Get(ctx, key)
	if err != nil || data == nil {
		return false
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return false
	}

	return true
}

// setCache stores data in cache
func (c *Client) setCache(ctx context.Context, key string, data interface{}, ttl time.Duration) {
	if c.cache == nil {
		return
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return
	}

	_ = c.cache.Set(ctx, key, bytes, ttl)
}

// downloadFile downloads a file from the given URL
func (c *Client) downloadFile(ctx context.Context, url string) ([]byte, error) {
	resp, err := c.doFileDownloadRequest(ctx, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
