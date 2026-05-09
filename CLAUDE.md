# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run

```bash
go build -o ptengine-cli .         # Build
go vet ./...                        # Lint
go test ./...                       # Test
./ptengine-cli --help               # Run
```

Version info is injected via ldflags at build time (see `.goreleaser.yaml`):
```bash
go build -ldflags "-X main.version=1.0.0 -X main.commit=abc123 -X main.date=2026-04-16" -o ptengine-cli .
```

## Release

Tag push (`v*`) triggers GoReleaser via GitHub Actions, building for linux/darwin/windows × amd64/arm64.
```bash
git tag v0.1.0 && git push origin v0.1.0
```

## Architecture

```
main.go                     → Entry point, version ldflags, calls cmd.Execute()
cmd/                        → Cobra command definitions
  root.go                   → Global flags (--api-key, --output, --base-url), ExitError, failValidation/parseJSONObject helpers
  heatmap.go                → Parent for heatmap subcommands
  heatmap_query.go          → Wraps POST /open-api/v1/heatmap/query
  heatmap_filter_values.go  → Wraps POST /open-api/v1/heatmap/filter-values
  heatmap_describe.go       → Local-only: outputs static schema JSON (no API call)
  data_query.go             → Parent for data-query subcommands
  data_query_query.go       → Wraps POST /open-api/v1/data-query/query (sync NL query)
  data_query_stream.go      → Wraps POST /open-api/v1/data-query/stream (SSE, one JSON line per event)
  data_query_cancel.go      → Wraps POST /open-api/v1/data-query/{traceId}/cancel
  config_set.go / config_show.go → Persistent config management
internal/api/
  client.go                 → HTTP client; doRequest takes optional extra http.Header (e.g., X-Request-Id)
  types.go                  → Heatmap request/response structs, CLIResponse envelope, CLIError
  data_query.go             → DataQueryRequest type + Client.DataQueryQuery/Stream/Cancel methods
  schema.go                 → Static parameter definitions (heatmap metrics, filters, funName)
  errors.go                 → API error code → exit code + hint; MapHTTPStatus for stream pre-SSE errors
internal/config/            → Viper: flag > PTENGINE_API_KEY env > ~/.config/ptengine-cli/config.yaml
internal/output/            → PrintSuccess (stdout), PrintError (stderr), PrintEnvelope (dispatch), PrintJSON
```

## Key Design Patterns

**Agent-first output**: Default `--output json` emits compact JSON. API commands use envelope:
- stdout: `{"success":true, "data":{...}, "meta":{...}, "rate_limit":{...}}`
- stderr: `{"success":false, "error":{"code":N, "message":"...", "hint":"..."}}`
- `heatmap describe`, `config show`, `version` output raw JSON (no envelope)

**Error flow**: RunE handlers print structured JSON to stderr via `output.PrintError()`, then return `&ExitError{Code: N}`. `cmd.Execute()` propagates the exit code to `main.go` which calls `os.Exit()`. Never call `os.Exit()` directly from RunE.

**Exit codes**: 0=ok, 1=validation, 2=auth, 3=param, 4=rate-limit, 5=server, 6=network (defined in `internal/api/errors.go`).

**Config precedence**: `--api-key` flag > `PTENGINE_API_KEY` env > `~/.config/ptengine-cli/config.yaml`. The `--profile-id` flag is local to query/filter-values/config-set commands; it falls back to `profile_id` in the config file.

## API Reference

- **API Docs**: https://helps.ptengine.com/cn/developer/open-api
- **Base URL**: `https://xbackend.ptengine.com`
- Auth: `x-api-key` header (same key for all endpoints)

Heatmap endpoints (固定字段查询):
- `/open-api/v1/heatmap/query` (4 query types: page_metrics / page_insight / block_metrics / element_metrics)
- `/open-api/v1/heatmap/filter-values`

Data-query endpoints (自然语言查询 → SQL → 行数据):
- `/open-api/v1/data-query/query` — sync NL query, returns full envelope after backend completes (~3-25s)
- `/open-api/v1/data-query/stream` — SSE stream of phase events (`routed` / `handler` / `final` / `cancelled` / `error`); the `final` event payload matches the sync envelope. CLI emits one JSON line per event to stdout.
- `/open-api/v1/data-query/{traceId}/cancel` — abort an in-flight query; pass the `--request-id` you sent earlier or the `trace_id` from the `routed` event. 404 is treated as ExitOK with `success=false` (already finished).
- Both query/stream support a `scenario` mode: `user_overview` / `user_timeline` / `user_session_detail` / `user_benchmark` with `--params '{"userId":"..."}'`.
- Successful response carries `data.message`: `"OK"` (rows ready) or `"CLARIFICATION_NEEDED"` (multi-turn — caller passes the user's choice back via `--context`).
