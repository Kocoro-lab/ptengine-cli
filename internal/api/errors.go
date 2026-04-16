package api

import "fmt"

// ExitCode constants for different error categories.
const (
	ExitOK         = 0
	ExitValidation = 1
	ExitAuth       = 2
	ExitParam      = 3
	ExitRateLimit  = 4
	ExitServer     = 5
	ExitNetwork    = 6
)

type apiErrorInfo struct {
	ExitCode int
	Message  string
	Hint     string
}

var errorMap = map[int]apiErrorInfo{
	4010: {ExitAuth, "Authentication failed: invalid API key", "Check your API key. Set via --api-key, PTENGINE_API_KEY env var, or 'ptengine-cli config set --api-key'."},
	4011: {ExitAuth, "Authentication failed: missing API key", "Provide an API key via --api-key, PTENGINE_API_KEY env var, or 'ptengine-cli config set --api-key'."},
	4007: {ExitParam, "Invalid device type", "Use one of: ALL, PC, MOBILE, TABLET."},
	4009: {ExitParam, "Invalid device type for this query", "Use one of: ALL, PC, MOBILE, TABLET."},
	4008: {ExitParam, "Block metrics not configured for this page", "Ensure blocks are defined for this page in the Ptengine dashboard before querying block_metrics."},
	4016: {ExitParam, "Element metrics not configured for this page", "Ensure elements are defined for this page in the Ptengine dashboard before querying element_metrics."},
	4290: {ExitRateLimit, "Rate limit exceeded (per-minute)", "Wait 60 seconds before retrying."},
	4291: {ExitRateLimit, "Rate limit exceeded (per-day)", "Daily limit reached. No more requests allowed today."},
	5000: {ExitServer, "Ptengine server error", "Retry after a short wait. If the issue persists, contact Ptengine support."},
}

// MapAPIError returns a CLIError and exit code for a given API error code.
func MapAPIError(code int, msg string) (*CLIError, int) {
	if info, ok := errorMap[code]; ok {
		return &CLIError{
			Code:    code,
			Message: info.Message,
			Hint:    info.Hint,
		}, info.ExitCode
	}
	// Unknown error code
	message := fmt.Sprintf("API error (code: %d)", code)
	if msg != "" {
		message = msg
	}
	return &CLIError{
		Code:    code,
		Message: message,
		Hint:    "Unexpected error. Check request parameters or contact Ptengine support.",
	}, ExitServer
}

// NewValidationError creates a CLIError for local validation failures.
func NewValidationError(message, hint string) *CLIError {
	return &CLIError{
		Code:    0,
		Message: message,
		Hint:    hint,
	}
}

// NewNetworkError creates a CLIError for network failures.
func NewNetworkError(err error) *CLIError {
	return &CLIError{
		Code:    0,
		Message: fmt.Sprintf("Network error: %v", err),
		Hint:    "Check your network connectivity and the base URL.",
	}
}
