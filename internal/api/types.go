package api

import "encoding/json"

// HeatmapQueryRequest is the request body for POST /open-api/v1/heatmap/query.
type HeatmapQueryRequest struct {
	QueryType       string   `json:"queryType"`
	ProfileID       string   `json:"profileId"`
	URL             string   `json:"url"`
	StartDate       string   `json:"startDate"`
	EndDate         string   `json:"endDate"`
	DeviceType      string   `json:"deviceType"`
	Metrics         []string `json:"metrics,omitempty"`
	ConversionNames []string `json:"conversionNames,omitempty"`
	Filters         []Filter `json:"filters,omitempty"`
	FunName         string   `json:"funName,omitempty"`
}

// Filter represents a single filter condition.
type Filter struct {
	Name  string   `json:"name"`
	Op    string   `json:"op"`
	Value []string `json:"value"`
}

// FilterValuesRequest is the request body for POST /open-api/v1/heatmap/filter-values.
type FilterValuesRequest struct {
	ProfileID string `json:"profileId"`
	Name      string `json:"name"`
	StartDate string `json:"startDate,omitempty"`
	EndDate   string `json:"endDate,omitempty"`
	Search    string `json:"search,omitempty"`
}

// APIResponse is the generic response wrapper from Ptengine API.
type APIResponse struct {
	Code int              `json:"code"`
	Data *json.RawMessage `json:"data,omitempty"`
	Meta *json.RawMessage `json:"meta,omitempty"`
	Msg  string           `json:"msg,omitempty"`
}

// RateLimit holds rate limit info parsed from response headers.
type RateLimit struct {
	MinuteLimit     int `json:"minute_limit,omitempty"`
	MinuteRemaining int `json:"minute_remaining,omitempty"`
	DayLimit        int `json:"day_limit,omitempty"`
	DayRemaining    int `json:"day_remaining,omitempty"`
}

// CLIResponse is the envelope returned to the user.
type CLIResponse struct {
	Success   bool             `json:"success"`
	Data      *json.RawMessage `json:"data,omitempty"`
	Meta      *json.RawMessage `json:"meta,omitempty"`
	RateLimit *RateLimit       `json:"rate_limit,omitempty"`
	Error     *CLIError        `json:"error,omitempty"`
}

// CLIError is a structured error for agent consumption.
type CLIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Hint    string `json:"hint"`
}
