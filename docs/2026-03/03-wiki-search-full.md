# 03 — wiki_search_full

**Date:** 2026-03-07
**Status:** Implemented
**Package:** `internal/tools`

## Purpose

Searches the OSRS Wiki using the full MediaWiki search API (`action=query&list=search`),
returning richer metadata than `wiki_search`: page IDs, word counts, and text snippets.

## MCP Tool Definition

- **Name:** `wiki_search_full`
- **Parameters:**
  - `query` (string, required): Search keywords, e.g. `dragon slayer quest`
  - `limit` (number, optional): Max results, 1–10, default 5

## API Endpoint

```
GET https://oldschool.runescape.wiki/api.php
  ?action=query
  &list=search
  &srsearch={query}
  &srlimit={limit}
  &format=json
```

## Response Format

```
Total hits: 42

1. Dragon Slayer I (ID: 12345, 1200 words)
   Complete the quest to wear rune platebody...
2. Dragon Slayer II (ID: 67890, 2000 words)
   ...
```

Returns `"No results found."` when there are no matches.

## Difference from wiki_search

| Feature | wiki_search | wiki_search_full |
|---------|------------|-----------------|
| API | opensearch | query&list=search |
| Page ID | No | Yes |
| Word count | No | Yes |
| Text snippet | No | Yes |
| Total hit count | No | Yes |

## Implementation

- `internal/tools/wiki_search_full.go` — tool struct + handler
- `internal/wiki/client.go` — `SearchPages` method on `Client` interface
- `internal/tools/wiki_search_full_test.go` — 100% coverage
