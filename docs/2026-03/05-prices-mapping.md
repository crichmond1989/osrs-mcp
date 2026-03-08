# 05 — prices_mapping

**Date:** 2026-03-07
**Status:** Implemented
**Package:** `internal/tools`

## Purpose

Returns OSRS item metadata (ID, name, alch values, GE buy limit) with an optional
name filter. Useful for discovering item IDs before using other price tools.

## MCP Tool Definition

- **Name:** `prices_mapping`
- **Parameters:**
  - `query` (string, optional): Case-insensitive name filter, e.g. `whip`

## API Endpoint

```
GET https://prices.runescape.wiki/api/v1/osrs/mapping
```

## Response Format

```
ID: 4151 | Abyssal whip [members] | High alch: 108000 gp | Limit: 70
ID: 4153 | Abyssal whip (or) [members] | High alch: 108000 gp | Limit: 70
```

Results are truncated to 20 items. When no items match the filter, returns
`"No items found."`.

## Notes

- When no `query` is given, returns the first 20 tradeable items
- Members-only items are marked with `[members]`

## Implementation

- `internal/tools/prices_mapping.go` — tool struct + handler
- `internal/prices/client.go` — `GetMapping` method on `Client` interface
- `internal/prices/models.go` — `MappingItem`
- `internal/tools/prices_mapping_test.go` — 100% coverage
