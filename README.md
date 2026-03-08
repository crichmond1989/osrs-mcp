# osrs-mcp

A [Model Context Protocol (MCP)](https://modelcontextprotocol.io) server that gives AI assistants real-time access to Old School RuneScape data: wiki pages, item prices, player stats, and quest information.

## Tools

| Tool | Description |
|------|-------------|
| `wiki_search` | Search the OSRS Wiki by keyword |
| `wiki_page` | Fetch full wikitext of a page by exact title |
| `wiki_search_full` | Search with rich metadata (page ID, word count, snippet) |
| `prices_latest` | Current instant buy/sell prices by item ID |
| `prices_mapping` | Item metadata with optional name filter |
| `prices_timeseries` | Averaged price and volume data over 5m/1h/24h windows |
| `player_stats` | Skill levels, XP, and hiscores ranks for any player |
| `player_quests` | Quest completion status via WikiSync (requires RuneLite plugin) |
| `quest_info` | Quest requirements and rewards from the OSRS Wiki |

See [docs/feature-catalog.md](docs/feature-catalog.md) for the full catalog with links to per-feature docs.

## Requirements

- Go 1.23+
- No API keys required — all data sources are public

## Build

```sh
go build -o osrs-mcp ./cmd/osrs-mcp
```

Or install directly:

```sh
go install github.com/crichmond1989/osrs-mcp/cmd/osrs-mcp@latest
```

## Integration

The server supports two transports:

- **stdio** (default) — for local desktop clients
- **HTTP** (`--addr :8080`) — for remote/mobile clients and hosted deployments

### Claude Desktop (stdio)

Add an entry to your Claude Desktop config file:

**macOS:** `~/Library/Application Support/Claude/claude_desktop_config.json`
**Windows:** `%APPDATA%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "osrs": {
      "command": "/absolute/path/to/osrs-mcp"
    }
  }
}
```

Replace `/absolute/path/to/osrs-mcp` with the actual path to the built binary (e.g. `$GOPATH/bin/osrs-mcp` if installed via `go install`).

Restart Claude Desktop after editing the config. The OSRS tools will appear in the tools panel.

### Claude Code (CLI, stdio)

```sh
claude mcp add osrs /absolute/path/to/osrs-mcp
```

### Cursor (stdio)

Open **Settings → MCP** and add a new server:

```json
{
  "osrs": {
    "command": "/absolute/path/to/osrs-mcp"
  }
}
```

### Claude Mobile / Remote Clients (HTTP)

Mobile apps and remote clients connect over HTTP. Run the server with `--addr`:

```sh
osrs-mcp --addr :8080
```

The MCP endpoint is served at `/mcp`. Point your client at `https://your-host/mcp`.

In Claude's mobile app, add a remote connector URL of `https://your-host/mcp`.

You can also set the `ADDR` environment variable instead of the flag:

```sh
ADDR=:8080 osrs-mcp
```

### Docker

```sh
make docker-build
docker run -p 8080:8080 osrs-mcp
```

The Docker image defaults to HTTP mode on port 8080.

### Azure Container Apps (recommended for hosting)

```sh
# Build and push to a registry
docker tag osrs-mcp <registry>.azurecr.io/osrs-mcp:latest
docker push <registry>.azurecr.io/osrs-mcp:latest

# Deploy (scale-to-zero, ~$0/month for typical usage)
az containerapp up \
  --name osrs-mcp \
  --resource-group <rg> \
  --image <registry>.azurecr.io/osrs-mcp:latest \
  --target-port 8080 \
  --ingress external
```

The deployed URL will be `https://osrs-mcp.<region>.azurecontainerapps.io`. Use `/mcp` as the connector path.

### Other MCP Clients

- **stdio**: pass the binary path as the command. No flags needed.
- **HTTP**: run with `--addr :PORT` and point the client at `https://your-host/mcp`.

## Development

```sh
make lint    # run golangci-lint
make build   # compile
make test    # run tests with race detector and coverage
make check   # lint + build + test (run this before every commit)
make cover   # show per-function coverage (must be 100%)
```

### Adding a New Tool

1. Add the client method to the relevant interface in `internal/wiki/client.go` or `internal/prices/client.go`.
2. Add the real implementation.
3. Create `internal/tools/<tool_name>.go` and `internal/tools/<tool_name>_test.go`.
4. Register the tool in `internal/tools/register.go`.
5. Add a row to `docs/feature-catalog.md` and create a feature doc at `docs/YYYY-MM/##-feature-title.md`.
6. Run `make check` and `make cover` — both must be clean.

### Keeping This README Up to Date

Update this file whenever:

- A new tool is added or removed — update the **Tools** table
- The minimum Go version changes — update **Requirements**
- The module path or binary name changes — update all build/install commands
- A new integration method is tested and confirmed — add it to **Integration**
- The HTTP endpoint path or transport changes — update the **HTTP** and **Docker** sections

The [feature catalog](docs/feature-catalog.md) is the source of truth for tool details; this README is the entry point summary.

## Data Sources

| Source | URL |
|--------|-----|
| OSRS Wiki (MediaWiki API) | https://oldschool.runescape.wiki/api.php |
| OSRS Wiki Prices API | https://prices.runescape.wiki/api/v1 |
| OSRS Hiscores API | https://secure.runescape.com/m=hiscore_oldschool |
| WikiSync API | https://sync.runescape.wiki/runelite/player |

## License

MIT
