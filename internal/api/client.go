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
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// NewClient creates a new API client.
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// doRequest sends a POST request and returns the parsed response, rate limit info, and any error.
func (c *Client) doRequest(path string, body interface{}) (*APIResponse, *RateLimit, error) {
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

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, nil, err
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
