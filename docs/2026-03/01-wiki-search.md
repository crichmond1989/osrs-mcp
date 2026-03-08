# 01 — wiki_search

**Date:** 2026-03-07
**Status:** Implemented
**Package:** `internal/tools`

## Purpose

Searches the OSRS Wiki using the MediaWiki OpenSearch API and returns matching
page titles with their canonical URLs.

## MCP Tool Definition

- **Name:** `wiki_search`
- **Parameters:**
  - `query` (string, required): Search keywords, e.g. `abyssal whip`
  - `limit` (number, optional): Max results, 1–10, default 5

## API Endpoint

```
GET https://oldschool.runescape.wiki/api.php
  ?action=opensearch
  &search={query}
  &limit={limit}
  &format=json
```

## Response Format

Numbered list of matching pages:

```
1. Abyssal whip
   https://oldschool.runescape.wiki/w/Abyssal_whip
2. Abyssal whip (or)
   https://oldschool.runescape.wiki/w/Abyssal_whip_(or)
```

Returns `"No results found."` when there are no matches.

## Implementation

- `internal/tools/wiki_search.go` — tool struct + handler
- `internal/wiki/client.go` — `OpenSearch` method on `Client` interface
- `internal/tools/wiki_search_test.go` — 100% coverage
