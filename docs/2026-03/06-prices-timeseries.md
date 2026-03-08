# 06 — prices_timeseries

**Date:** 2026-03-07
**Status:** Implemented
**Package:** `internal/tools`

## Purpose

Returns averaged OSRS Grand Exchange price and trade volume data for a specific
item over a given time window (5 minutes, 1 hour, or 24 hours).

## MCP Tool Definition

- **Name:** `prices_timeseries`
- **Parameters:**
  - `id` (number, required): Item ID, e.g. `4151` for Abyssal whip
  - `window` (string, required): Time window — one of `"5m"`, `"1h"`, `"24h"`

## API Endpoints

| Window | Endpoint |
|--------|---------|
| `5m` | `GET https://prices.runescape.wiki/api/v1/osrs/5m` |
| `1h` | `GET https://prices.runescape.wiki/api/v1/osrs/1h` |
| `24h` | `GET https://prices.runescape.wiki/api/v1/osrs/24h` |

## Response Format

```
Item ID: 4151 | Window: 1h
Avg buy price:   2500000 gp (volume: 150)
Avg sell price:  2400000 gp (volume: 200)
```

When no trades occurred in the window, shows `no trades` for that price.

## Notes

- `AvgHighPrice` / `AvgLowPrice` are `null` in the API when no trades occurred
  (modeled as `*int` in `TimeSeriesPrice`)
- Useful for gauging recent market activity and price trends

## Implementation

- `internal/tools/prices_timeseries.go` — tool struct + handler
- `internal/prices/client.go` — `GetTimeSeries` method + `ErrInvalidWindow`
- `internal/prices/models.go` — `TimeSeriesResponse`, `TimeSeriesPrice`
- `internal/tools/prices_timeseries_test.go` — 100% coverage
