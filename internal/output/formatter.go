package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Kocoro-lab/ptengine-cli/internal/api"
)

// PrintSuccess writes a successful CLI response to stdout.
func PrintSuccess(resp *api.CLIResponse, format string) {
	switch format {
	case "json":
		data, _ := json.Marshal(resp)
		fmt.Println(string(data))
	case "json-pretty":
		data, _ := json.MarshalIndent(resp, "", "  ")
		fmt.Println(string(data))
	case "table":
		// For table format, print a simplified human-readable view
		if resp.Data != nil {
			data, _ := json.MarshalIndent(json.RawMessage(*resp.Data), "", "  ")
			fmt.Println(string(data))
		}
		if resp.RateLimit != nil {
			fmt.Fprintf(os.Stderr, "\n[Rate Limit] minute: %d/%d, day: %d/%d\n",
				resp.RateLimit.MinuteRemaining, resp.RateLimit.MinuteLimit,
				resp.RateLimit.DayRemaining, resp.RateLimit.DayLimit)
		}
	default:
		data, _ := json.Marshal(resp)
		fmt.Println(string(data))
	}
}

// PrintError writes a failure CLI response to stderr and returns the exit code.
func PrintError(cliErr *api.CLIError, rl *api.RateLimit, format string) int {
	resp := &api.CLIResponse{
		Success:   false,
		Error:     cliErr,
		RateLimit: rl,
	}

	var data []byte
	switch format {
	case "json-pretty":
		data, _ = json.MarshalIndent(resp, "", "  ")
	default:
		data, _ = json.Marshal(resp)
	}
	fmt.Fprintln(os.Stderr, string(data))

	if cliErr.Code > 0 {
		_, exitCode := api.MapAPIError(cliErr.Code, "")
		return exitCode
	}
	return api.ExitValidation
}

// PrintJSON writes any value as JSON to stdout.
func PrintJSON(v interface{}, format string) {
	var data []byte
	switch format {
	case "json-pretty", "table":
		data, _ = json.MarshalIndent(v, "", "  ")
	default:
		data, _ = json.Marshal(v)
	}
	fmt.Println(string(data))
}
