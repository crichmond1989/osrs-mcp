# OSRS MCP — Feature Catalog

Master list of all implemented MCP tools. Each entry links to its detailed feature doc.

## MCP Tools

| # | Tool Name | Description | Doc |
|---|-----------|-------------|-----|
| 1 | `wiki_search` | Search OSRS Wiki by keyword (opensearch) | [docs/2026-03/01-wiki-search.md](2026-03/01-wiki-search.md) |
| 2 | `wiki_page` | Fetch wikitext of a page by exact title | [docs/2026-03/02-wiki-page.md](2026-03/02-wiki-page.md) |
| 3 | `wiki_search_full` | Search OSRS Wiki with full metadata (page ID, word count, snippet) | [docs/2026-03/03-wiki-search-full.md](2026-03/03-wiki-search-full.md) |
| 4 | `prices_latest` | Get current instant-buy/sell prices by item ID | [docs/2026-03/04-prices-latest.md](2026-03/04-prices-latest.md) |
| 5 | `prices_mapping` | Look up item metadata with optional name filter | [docs/2026-03/05-prices-mapping.md](2026-03/05-prices-mapping.md) |
| 6 | `prices_timeseries` | Get averaged price/volume data over a time window | [docs/2026-03/06-prices-timeseries.md](2026-03/06-prices-timeseries.md) |
| 7 | `player_stats` | Fetch OSRS player skill levels, XP, and hiscores ranks | [docs/2026-03/07-player-stats.md](2026-03/07-player-stats.md) |
| 8 | `player_quests` | Fetch quest completion status via WikiSync (requires RuneLite plugin) | [docs/2026-03/08-player-quests.md](2026-03/08-player-quests.md) |
| 9 | `quest_info` | Fetch quest requirements and rewards from the OSRS Wiki | [docs/2026-03/09-quest-info.md](2026-03/09-quest-info.md) |

## Adding New Features

Follow the instructions in [CLAUDE.md](../CLAUDE.md):
1. Create a feature doc at `/docs/YYYY-MM/##-feature-title.md`
2. Update this catalog with a new row
3. Run `make check` to validate
