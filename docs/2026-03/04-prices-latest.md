# 04 — prices_latest

**Date:** 2026-03-07
**Status:** Implemented
**Package:** `internal/tools`

## Purpose

Returns the current instant-buy (high) and instant-sell (low) prices for an OSRS
Grand Exchange item, identified by its item ID.

## MCP Tool Definition

- **Name:** `prices_latest`
- **Parameters:**
  - `id` (number, required): Item ID, e.g. `4151` for Abyssal whip

## API Endpoint

```
GET https://prices.runescape.wiki/api/v1/osrs/latest?id={id}
```

## Response Format

```
Item ID: 4151
Instant buy:  2500000 gp
Instant sell: 2400000 gp
```

When no recent trade exists for buy or sell, shows `no data`.

## Notes

- Item IDs can be found via the `prices_mapping` tool
- Prices update in real-time as trades occur on the Grand Exchange

## Implementation

- `internal/tools/prices_latest.go` — tool struct + handler
- `internal/prices/client.go` — `GetLatest` method on `Client` interface
- `internal/prices/models.go` — `LatestResponse`, `LatestPrice` (uses `*int` for nullable prices)
- `internal/tools/prices_latest_test.go` — 100% coverage
