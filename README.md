# ptengine-cli

A CLI tool for querying [Ptengine](https://ptengine.com) heatmap data via Open API. Designed for AI agents and human users.

## Installation

```bash
curl -sSL https://raw.githubusercontent.com/Kocoro-lab/ptengine-cli/main/scripts/install.sh | sh
```

Or build from source:

```bash
go install github.com/Kocoro-lab/ptengine-cli@latest
```

## Quick Start

```bash
# Set your API key
ptengine-cli config set --api-key pt-your-api-key --profile-id your-profile-id

# Discover available parameters
ptengine-cli heatmap describe
ptengine-cli heatmap describe --query-type page_metrics

# Query page metrics
ptengine-cli heatmap query \
  --query-type page_metrics \
  --profile-id 566d12f9 \
  --url "https://example.com" \
  --start-date 2026-03-01 \
  --end-date 2026-03-31 \
  --device-type ALL \
  --metrics pv,uv,bounceRate

# Query page insight (grouped by source)
ptengine-cli heatmap query \
  --query-type page_insight \
  --profile-id 566d12f9 \
  --url "https://example.com" \
  --start-date 2026-03-01 \
  --end-date 2026-03-31 \
  --device-type ALL \
  --fun-name sourceType \
  --metrics pv,uv

# Query with filters
ptengine-cli heatmap query \
  --query-type page_metrics \
  --profile-id 566d12f9 \
  --url "https://example.com" \
  --start-date 2026-03-01 \
  --end-date 2026-03-31 \
  --device-type ALL \
  --filter 'country include Japan,China' \
  --filter 'browser exclude Safari'

# Get available filter values
ptengine-cli heatmap filter-values \
  --profile-id 566d12f9 \
  --name country \
  --search "Ja"
```

## Configuration

Configuration is resolved in this order (highest priority first):

1. Command-line flags (`--api-key`, `--profile-id`)
2. Environment variable (`PTENGINE_API_KEY`)
3. Config file (`~/.config/ptengine-cli/config.yaml`)

```bash
# Save config
ptengine-cli config set --api-key pt-xxxxx --profile-id 566d12f9

# View current config
ptengine-cli config show
```

## Commands

| Command | Description |
|---------|-------------|
| `heatmap query` | Query heatmap data (page_metrics, page_insight, block_metrics, element_metrics) |
| `heatmap filter-values` | Get available values for a filter type |
| `heatmap describe` | Show available query types, metrics, filters, and parameters |
| `config set` | Save API key, profile ID, or base URL |
| `config show` | Show current effective configuration |
| `version` | Print version information |

## Output

Default output is JSON (optimized for AI agent parsing). Use `--output` flag to change:

```bash
--output json          # Compact JSON (default)
--output json-pretty   # Pretty-printed JSON
--output table         # Human-readable table
```

## Filter Syntax

```bash
--filter 'name include|exclude val1,val2,...'
```

Examples:
```bash
--filter 'country include Japan,China'
--filter 'browser exclude Safari'
--filter 'utmSource include google,facebook'
```

For complex cases, use raw JSON:
```bash
--filter-json '[{"name":"country","op":"include","value":["Japan"]}]'
```

## API Documentation

See [Ptengine Open API](https://helps.ptengine.com/cn/developer/open-api) for full API reference.

## License

MIT
