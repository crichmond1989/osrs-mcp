# 09 — quest_info

**Date:** 2026-03-07
**Status:** Implemented
**Package:** `internal/tools`

## Purpose

Fetches OSRS quest information from the Wiki. Returns the full wikitext of a quest
page, which contains skill requirements, item requirements, quest point rewards,
and XP rewards in structured sections.

## MCP Tool Definition

- **Name:** `quest_info`
- **Parameters:**
  - `quest` (string, required): Exact quest name as it appears on the OSRS Wiki, e.g. `Dragon Slayer I`

## API Endpoint

Reuses the existing `wiki.Client.GetPage` method:

```
GET https://oldschool.runescape.wiki/api.php
  ?action=query
  &prop=revisions
  &rvprop=content
  &rvslots=main
  &titles={quest}
  &format=json
```

## Response Format

Raw wikitext of the quest page. Look for:
- `==Requirements==` — skill and quest requirements
- `==Rewards==` — XP rewards, quest points, item rewards

Returns an error result if the quest name does not match any wiki page.

## Implementation

- `internal/tools/quest_info.go` — tool struct + handler (reuses `wiki.Client`)
- `internal/wiki/client.go` — `GetPage` method (pre-existing)
- `internal/tools/quest_info_test.go` — 100% coverage
