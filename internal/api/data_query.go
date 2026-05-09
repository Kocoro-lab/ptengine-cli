package api

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ValidScenarios enumerates the predefined scenario modes for data-query.
// Used by the CLI for client-side validation before hitting the server.
var ValidScenarios = []string{
	"user_overview",
	"user_timeline",
	"user_session_detail",
	"user_benchmark",
}

// DataQueryRequest is the request body for POST /open-api/v1/data-query/{query,stream}.
type DataQueryRequest struct {
	Question  string                 `json:"question,omitempty"`
	ProfileID string                 `json:"profileId"`
	Context   string                 `json:"context,omitempty"`
	Scenario  string                 `json:"scenario,omitempty"`
	Params    map[string]interface{} `json:"params,omitempty"`
}

// DataQueryQuery calls POST /open-api/v1/data-query/query (sync).
//
// The optional requestID is sent as X-Request-Id; the server echoes it as
// `data.data.traceId` in the response and accepts it as the path arg for
// /<traceId>/cancel. Note: SSE event payloads (stream) use `trace_id` snake_case
// instead — that's a server-side inconsistency, not a typo.
func (c *Client) DataQueryQuery(req *DataQueryRequest, requestID string) (*CLIResponse, int) {
	apiResp, rl, err := c.doRequest(
		"/open-api/v1/data-query/query", req, requestHeaders(requestID),
	)
	if err != nil {
		return &CLIResponse{Success: false, Error: NewNetworkError(err)}, ExitNetwork
	}
	if apiResp.Code != 200 {
		cliErr, exitCode := MapAPIError(apiResp.Code, apiResp.Msg)
		return &CLIResponse{Success: false, Error: cliErr, RateLimit: rl}, exitCode
	}
	return &CLIResponse{Success: true, Data: apiResp.Data, Meta: apiResp.Meta, RateLimit: rl}, ExitOK
}

// DataQueryStream calls POST /open-api/v1/data-query/stream and writes each
// parsed SSE event as one JSON line ({"event":"...","data":{...}}) to `out`.
//
// Returns ExitOK on a clean stream end, or a transport-level exit code on
// error. Logical query failures arrive as a `final` or `error` event in the
// stream — callers should inspect the line stream rather than rely on the
// exit code alone.
func (c *Client) DataQueryStream(req *DataQueryRequest, requestID string, out io.Writer) int {
	body, err := json.Marshal(req)
	if err != nil {
		return ExitValidation
	}
	httpReq, err := http.NewRequest("POST", c.BaseURL+"/open-api/v1/data-query/stream", bytes.NewReader(body))
	if err != nil {
		return ExitNetwork
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")
	httpReq.Header.Set("x-api-key", c.APIKey)
	httpReq.Header.Set("User-Agent", "ptengine-cli")
	if requestID != "" {
		httpReq.Header.Set("X-Request-Id", requestID)
	}

	resp, err := c.StreamClient.Do(httpReq)
	if err != nil {
		return ExitNetwork
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		// Error before stream — body is JSON envelope, not SSE
		respBody, _ := io.ReadAll(resp.Body)
		fmt.Fprintln(out, string(respBody))
		return MapHTTPStatus(resp.StatusCode)
	}

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	var event, data string
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "event: "):
			event = strings.TrimPrefix(line, "event: ")
		case strings.HasPrefix(line, "data: "):
			data = strings.TrimPrefix(line, "data: ")
		case line == "":
			if event != "" || data != "" {
				writeSSEFrame(out, event, data)
				event, data = "", ""
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return ExitNetwork
	}
	return ExitOK
}

// writeSSEFrame emits one parsed SSE event as a JSON line: {"event":"...","data":...}.
// `data` is treated as a JSON value (object/null/scalar); empty becomes JSON null.
func writeSSEFrame(out io.Writer, event, data string) {
	rawData := json.RawMessage(data)
	if data == "" {
		rawData = json.RawMessage("null")
	}
	frame := struct {
		Event string          `json:"event"`
		Data  json.RawMessage `json:"data"`
	}{Event: event, Data: rawData}
	b, err := json.Marshal(frame)
	if err != nil {
		return
	}
	out.Write(append(b, '\n'))
}

// DataQueryCancel calls POST /open-api/v1/data-query/{traceId}/cancel.
//
// 404 (apiResp.Code == 4040) means "already finished or unknown id" — surfaced
// as success=false but treated as ExitOK since it's not a failure for the caller.
func (c *Client) DataQueryCancel(traceID string) (*CLIResponse, int) {
	path := "/open-api/v1/data-query/" + traceID + "/cancel"
	apiResp, rl, err := c.doRequest(path, struct{}{})
	if err != nil {
		return &CLIResponse{Success: false, Error: NewNetworkError(err)}, ExitNetwork
	}
	if apiResp.Code == 4040 {
		return &CLIResponse{Success: false, Data: apiResp.Data, RateLimit: rl}, ExitOK
	}
	if apiResp.Code != 200 {
		cliErr, exitCode := MapAPIError(apiResp.Code, apiResp.Msg)
		return &CLIResponse{Success: false, Error: cliErr, RateLimit: rl}, exitCode
	}
	return &CLIResponse{Success: true, Data: apiResp.Data, RateLimit: rl}, ExitOK
}

func requestHeaders(requestID string) http.Header {
	h := http.Header{}
	if requestID != "" {
		h.Set("X-Request-Id", requestID)
	}
	return h
}
