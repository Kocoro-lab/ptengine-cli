package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// Client is the HTTP client for the Ptengine API.
//
// HTTPClient is for one-shot request/response calls (heatmap, data-query/query,
// data-query/cancel). StreamClient is for long-lived SSE responses (data-query/stream)
// where the wall-clock duration is bounded by the server, not the client.
type Client struct {
	BaseURL      string
	APIKey       string
	HTTPClient   *http.Client
	StreamClient *http.Client
}

// NewClient creates a new API client.
//
// HTTPClient timeout is 120s — covers heatmap (~1s) and data-query/query
// (worst-case ~25s LLM inference + network jitter) with safe margin.
// StreamClient has no timeout; the server controls SSE duration.
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL:      baseURL,
		APIKey:       apiKey,
		HTTPClient:   &http.Client{Timeout: 120 * time.Second},
		StreamClient: &http.Client{},
	}
}

// reservedHeaders are set by doRequest from Client fields and must not be
// overridden by per-call extra headers (would break auth or content negotiation).
var reservedHeaders = map[string]struct{}{
	"Content-Type": {},
	"X-Api-Key":    {},
	"User-Agent":   {},
}

// doRequest sends a POST request and returns the parsed response, rate limit info, and any error.
// Optional extra headers (e.g., X-Request-Id) can be passed; later headers override earlier ones.
func (c *Client) doRequest(path string, body interface{}, extraHeaders ...http.Header) (*APIResponse, *RateLimit, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := c.BaseURL + path
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("User-Agent", "ptengine-cli")
	for _, h := range extraHeaders {
		for k, vv := range h {
			if _, reserved := reservedHeaders[http.CanonicalHeaderKey(k)]; reserved {
				continue
			}
			for _, v := range vv {
				req.Header.Set(k, v)
			}
		}
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("POST %s: %w", path, err)
	}
	defer resp.Body.Close()

	// Parse rate limit headers
	rl := parseRateLimit(resp.Header)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, rl, fmt.Errorf("failed to read response: %w", err)
	}

	var apiResp APIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, rl, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	return &apiResp, rl, nil
}

func parseRateLimit(h http.Header) *RateLimit {
	rl := &RateLimit{}
	hasValue := false

	if v := h.Get("X-RateLimit-Limit-Minute"); v != "" {
		rl.MinuteLimit, _ = strconv.Atoi(v)
		hasValue = true
	}
	if v := h.Get("X-RateLimit-Remaining-Minute"); v != "" {
		rl.MinuteRemaining, _ = strconv.Atoi(v)
		hasValue = true
	}
	if v := h.Get("X-RateLimit-Limit-Day"); v != "" {
		rl.DayLimit, _ = strconv.Atoi(v)
		hasValue = true
	}
	if v := h.Get("X-RateLimit-Remaining-Day"); v != "" {
		rl.DayRemaining, _ = strconv.Atoi(v)
		hasValue = true
	}

	if !hasValue {
		return nil
	}
	return rl
}
