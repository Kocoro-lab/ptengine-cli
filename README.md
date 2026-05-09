# ptengine-cli

A CLI tool for querying [Ptengine](https://ptengine.com) analytics data via Open API. Two command groups:

- `ptengine-cli heatmap *` — fixed-field heatmap metrics queries (page / block / element)
- `ptengine-cli data-query *` — natural-language data queries that return rows + the SQL the backend ran

Designed for AI agents and human users.

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

# Natural-language data query (sync)
ptengine-cli data-query query --question "昨天的 UV"

# Natural-language data query with cancel-able request id
ptengine-cli data-query query \
  --question "中国大陆 vs 海外的支付完成率" \
  --request-id $(uuidgen)

# User deep-dive scenario (no LLM cost — runs canned SQL templates)
ptengine-cli data-query query \
  --scenario user_overview \
  --params '{"userId":"c78c1a57-a89b-4086-9fcd-5733c04dcdef"}'

# Stream phase events (one JSON line per SSE frame)
ptengine-cli data-query stream \
  --question "本周 UV vs 上周" \
  --request-id $(uuidgen)

# Cancel an in-flight query by trace id
ptengine-cli data-query cancel <trace-id>
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
| `data-query query` | Submit a natural-language analytics question (sync, returns full envelope) |
| `data-query stream` | Submit a question and stream phase events (SSE, one JSON line per event) |
| `data-query cancel` | Abort an in-flight `data-query` by trace id |
| `config set` | Save API key, profile ID, or base URL |
| `config show` | Show current effective configuration |
| `version` | Print version information |

### `data-query` flags

| Flag | Required | Notes |
|---|---|---|
| `--question` | yes (free mode) | Natural language question (zh / en / ja). **Pass the user's wording verbatim** — the backend's multilingual parser normalizes it. |
| `--profile-id` | yes | Falls back to config file if omitted. Must equal the API key's bound profile. |
| `--scenario` | yes (scenario mode) | One of `user_overview` / `user_timeline` / `user_session_detail` / `user_benchmark`. |
| `--params` | scenario-only | JSON object literal, e.g. `'{"userId":"abc..."}'`. |
| `--context` | optional | Multi-turn clarification reply or time-anchor hints; ≤2000 chars. |
| `--request-id` | recommended | UUID; sent as `X-Request-Id`, becomes the `traceId`. Without it, `cancel` cannot target the query mid-flight. |

### Response envelope

`data-query` wraps an inner open-api envelope — reach `data.data.results[]` for actual rows, and branch on `data.message`:

```json
{
  "success": true,
  "data": {
    "code": 200,
    "message": "OK",
    "data": {
      "intent": { "type": "sql", "matchedBy": "skill:data_query", "confidence": "exact" },
      "results": [
        { "title": "[atomic:metric_total]", "columns": ["uv"], "rows": [{"uv": 1600}],
          "sql": "SELECT COUNT(DISTINCT user_id) AS uv FROM behavior_fact ...", "rowCount": 1 }
      ],
      "traceId": "tq_xxx"
    }
  },
  "rate_limit": {
    "minute_limit": 30, "minute_remaining": 28,
    "day_limit": 3000, "day_remaining": 2997
  }
}
```

| Field | Meaning |
|---|---|
| `success` | CLI-side dispatch flag. `true` → stdout / OK; `false` → stderr (or stdout for soft-fail like cancel `not_found`). |
| `data.message` | `"OK"` (rows in `data.data.results[]`) or `"CLARIFICATION_NEEDED"` (multi-turn — see below). |
| `data.data.traceId` | Echoes `--request-id` if you set one; otherwise server-generated. Use it for `data-query cancel`. |
| `rate_limit.minute_remaining` | Requests left this minute. Plan tiers: Free 3 RPM / Free Trial 10 RPM / Growth 30 RPM. |
| `rate_limit.day_remaining` | Requests left today. Plans: Free 100 RPD / Free Trial 1,000 RPD / Growth 3,000 RPD. |

### Multi-turn clarification

When the question is ambiguous (e.g., `"支付完成率"` could map to several events), the backend returns `data.message === "CLARIFICATION_NEEDED"` instead of rows:

```bash
ptengine-cli data-query query --question "中国大陆 vs 海外的支付完成率"
```

```json
{
  "success": true,
  "data": {
    "code": 200,
    "message": "CLARIFICATION_NEEDED",
    "data": {
      "type": "event_disambiguation",
      "question": "「支付完成」可能对应以下事件，请选择：",
      "options": [
        { "label": "Payment Success",  "value": "payment_success"  },
        { "label": "Payment Failed",   "value": "payment_failed"   },
        { "label": "Payment Canceled", "value": "payment_canceled" },
        { "label": "Payment Refund",   "value": "payment_refund"   }
      ],
      "allowCustom": true,
      "traceId": "..."
    }
  }
}
```

Resubmit with the user's choice via `--context`:

```bash
ptengine-cli data-query query \
  --question "中国大陆 vs 海外的支付完成率" \
  --context "用户选择: payment_success (对应「支付完成」)"
```

The CLI does **not** loop automatically — your agent (or a human) inspects `options[]`, picks one, and recalls. Other clarification `type` values: `missing_identifier`, `ambiguous_scope`, `custom`.

## Output

Default output is JSON (optimized for AI agent parsing). Use `--output` flag to change:

```bash
--output json          # Compact JSON (default)
--output json-pretty   # Pretty-printed JSON
--output table         # Human-readable table
```

Success envelope is printed to **stdout**; error envelope is printed to **stderr** — except for "soft failures" like `data-query cancel` returning `not_found`, which goes to stdout (the envelope still carries useful fields).

## Exit Codes

Branch on the process exit code first, then on `data.code` from the JSON for finer category:

| Code | Class | When | Recommended client behaviour |
|---|---|---|---|
| 0 | OK | Success **or** clarification (check `data.message === "CLARIFICATION_NEEDED"`) **or** cancel `not_found` (`success=false` but exit 0) | Process the response |
| 1 | validation | Caller-side: missing flag, bad value, malformed `--params` JSON | Don't retry — fix the invocation |
| 2 | auth | API codes 4010 / 4011 / 4030 (key invalid, missing, or `profileId` mismatch) | Don't retry — verify config |
| 3 | param | Other 4xx including 4220 (SQL safety validation) and 4990 (cancelled) | **4220 means never retry the same `--question`** — ask user to rephrase |
| 4 | rate-limit | 4290 (per-minute) / 4291 (per-day) | Wait then retry. Respect `rate_limit.minute_remaining` |
| 5 | server | 5xxx (5000 internal, 5020 upstream LLM/DB, 5040 upstream timeout) | Backoff: 1s → 3s → 9s, max 3 attempts |
| 6 | network | DNS failure, connection reset, no response | Retry once; surface to user if still failing |

`data-query cancel` 404 (already-finished or unknown trace-id) is intentionally exit 0 — the original query already completed, nothing to abort.

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
