# 08 — player_quests

**Date:** 2026-03-07
**Status:** Implemented
**Package:** `internal/tools`

## Purpose

Fetches quest completion status for an OSRS player using WikiSync data.
Returns counts of quests by completion state (complete, in progress, not started).

**Requirement:** The player must have the WikiSync plugin enabled in RuneLite.
If WikiSync is not enabled, the tool returns an informative error message.

## MCP Tool Definition

- **Name:** `player_quests`
- **Parameters:**
  - `player` (string, required): OSRS player name

## API Endpoint

```
GET https://sync.runescape.wiki/runelite/player/{player}/STANDARD
```

Returns HTTP 400 or 404 if the player has not enabled WikiSync.

## Response Format

```
Player: Guthix Her0 (WikiSync data)

Quest completion: 247 complete, 3 in progress, 20 not started
```

Returns an error result if the player has not enabled WikiSync, with a message
explaining how to enable it.

## Implementation

- `internal/tools/player_quests.go` — tool struct + handler
- `internal/wikisync/client.go` — `GetPlayerData` method on `Client` interface
- `internal/wikisync/models.go` — `WikiSyncResponse`, quest state constants
- `internal/tools/player_quests_test.go` — 100% coverage
- `internal/wikisync/client_test.go` — 100% coverage
