# 07 — player_stats

**Date:** 2026-03-07
**Status:** Implemented
**Package:** `internal/tools`

## Purpose

Fetches OSRS player skill levels, XP, and hiscores ranks for all 24 skills.
Supports all hiscores modes (standard, ironman, hardcore, ultimate).

## MCP Tool Definition

- **Name:** `player_stats`
- **Parameters:**
  - `player` (string, required): OSRS player name (case-insensitive)
  - `mode` (string, optional): Hiscores mode — `standard` (default), `ironman`, `hardcore`, or `ultimate`

## API Endpoint

```
GET https://secure.runescape.com/m=hiscore_oldschool/index_lite.json?player={player}
```

Mode variants replace `hiscore_oldschool` with:
- `hiscore_oldschool_ironman`
- `hiscore_oldschool_hardcore_ironman`
- `hiscore_oldschool_ultimate`

Returns HTTP 404 for unknown players.

## Response Format

Skill table with columns Skill / Level / XP / Rank. Unranked skills show "unranked" for rank.

```
Player: Zezima (standard)

Skill            Level           XP     Rank
-----            -----           --     ----
Overall           2277  4600000000        1
Attack              99   200000000     1234
...
```

## Implementation

- `internal/tools/player_stats.go` — tool struct + handler
- `internal/hiscores/client.go` — `GetStats` method on `Client` interface
- `internal/hiscores/models.go` — `SkillEntry`, `HiscoresResponse`
- `internal/tools/player_stats_test.go` — 100% coverage
- `internal/hiscores/client_test.go` — 100% coverage
