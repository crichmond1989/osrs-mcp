# 02 — wiki_page

**Date:** 2026-03-07
**Status:** Implemented
**Package:** `internal/tools`

## Purpose

Retrieves the raw wikitext content of an OSRS Wiki page by its exact title.
Useful for extracting structured item stats, quest guides, and other page data.

## MCP Tool Definition

- **Name:** `wiki_page`
- **Parameters:**
  - `title` (string, required): Exact page title, e.g. `Abyssal whip`

## API Endpoint

```
GET https://oldschool.runescape.wiki/api.php
  ?action=query
  &prop=revisions
  &rvprop=content
  &rvslots=main
  &titles={title}
  &format=json
```

## Response Format

Raw wikitext content of the page, e.g.:

```
{{Infobox Item
|name = Abyssal whip
|image = [[File:Abyssal whip.png]]
...
}}
```

Returns an error message if the page is not found or has no content.

## Special Cases

- Page ID `-1` in the API response signals "not found" — returns `ErrPageNotFound`
- Pages with no revisions return an error message

## Implementation

- `internal/tools/wiki_page.go` — tool struct + handler
- `internal/wiki/client.go` — `GetPage` method + `ErrPageNotFound` sentinel
- `internal/tools/wiki_page_test.go` — 100% coverage
