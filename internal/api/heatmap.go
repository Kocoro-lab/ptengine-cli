package api

// HeatmapQuery calls POST /open-api/v1/heatmap/query.
func (c *Client) HeatmapQuery(req *HeatmapQueryRequest) (*CLIResponse, int) {
	apiResp, rl, err := c.doRequest("/open-api/v1/heatmap/query", req)
	if err != nil {
		cliErr := NewNetworkError(err)
		return &CLIResponse{Success: false, Error: cliErr}, ExitNetwork
	}

	if apiResp.Code != 200 {
		cliErr, exitCode := MapAPIError(apiResp.Code, apiResp.Msg)
		return &CLIResponse{Success: false, Error: cliErr, RateLimit: rl}, exitCode
	}

	return &CLIResponse{
		Success:   true,
		Data:      apiResp.Data,
		Meta:      apiResp.Meta,
		RateLimit: rl,
	}, ExitOK
}

// HeatmapFilterValues calls POST /open-api/v1/heatmap/filter-values.
func (c *Client) HeatmapFilterValues(req *FilterValuesRequest) (*CLIResponse, int) {
	apiResp, rl, err := c.doRequest("/open-api/v1/heatmap/filter-values", req)
	if err != nil {
		cliErr := NewNetworkError(err)
		return &CLIResponse{Success: false, Error: cliErr}, ExitNetwork
	}

	if apiResp.Code != 200 {
		cliErr, exitCode := MapAPIError(apiResp.Code, apiResp.Msg)
		return &CLIResponse{Success: false, Error: cliErr, RateLimit: rl}, exitCode
	}

	return &CLIResponse{
		Success:   true,
		Data:      apiResp.Data,
		RateLimit: rl,
	}, ExitOK
}
