# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**ptengine-cli** — Go CLI tool for querying Ptengine heatmap data via Open API. Primarily designed for AI agent consumption, also usable by humans.

- **Repository**: Kocoro-lab/ptengine-cli
- **API Docs**: https://helps.ptengine.com/cn/developer/open-api
- **Base URL**: `https://xbackend.ptengine.com/`

## Build & Run

```bash
go build -o ptengine-cli .         # Build
go vet ./...                        # Lint
go test ./...                       # Test
./ptengine-cli --help               # Run
```

## Architecture

```
main.go              → Entry point, version ldflags
cmd/                  → Cobra commands (root, version, heatmap/*, config/*)
internal/api/         → HTTP client, request/response types, error mapping, schema
internal/config/      → Viper-based config (flag > env > ~/.config/ptengine-cli/config.yaml)
internal/output/      → JSON/table output formatting
```

Key design decisions:
- **Output**: JSON by default (agent-first). `--output json-pretty|table` for humans.
- **Config precedence**: `--api-key` flag > `PTENGINE_API_KEY` env > config file
- **Error handling**: Structured JSON on stderr with distinct exit codes (2=auth, 3=param, 4=rate-limit, 5=server, 6=network)
- **Schema discovery**: `heatmap describe` command outputs all valid query types, metrics, filters, funName values as JSON

## API Endpoints Wrapped

1. `POST /open-api/v1/heatmap/query` → `ptengine-cli heatmap query`
2. `POST /open-api/v1/heatmap/filter-values` → `ptengine-cli heatmap filter-values`

## Release

- CI: GitHub Actions on push/PR (`go vet`, `go test`, `go build`)
- Release: GoReleaser on tag push (`v*`), builds for linux/darwin/windows × amd64/arm64
- Install: `curl -sSL https://raw.githubusercontent.com/Kocoro-lab/ptengine-cli/main/scripts/install.sh | sh`
