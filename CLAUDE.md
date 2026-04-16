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
  root.go                   → Global flags (--api-key, --output, --base-url), ExitError type, config loading
  heatmap_query.go          → Core command: wraps POST /open-api/v1/heatmap/query
  heatmap_filter_values.go  → Wraps POST /open-api/v1/heatmap/filter-values
  heatmap_describe.go       → Local-only: outputs static schema JSON (no API call)
  config_set.go / config_show.go → Persistent config management
internal/api/
  client.go                 → HTTP client with auth header, rate-limit header parsing
  types.go                  → Request/response structs, CLIResponse envelope, CLIError
  schema.go                 → Static parameter definitions (metrics, filters, funName per query type)
  errors.go                 → API error code → exit code + human hint mapping
internal/config/            → Viper: flag > PTENGINE_API_KEY env > ~/.config/ptengine-cli/config.yaml
internal/output/            → PrintSuccess (stdout), PrintError (stderr), PrintJSON formatters
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
- Auth: `x-api-key` header
- Two endpoints: `/open-api/v1/heatmap/query` (4 query types) and `/open-api/v1/heatmap/filter-values`
