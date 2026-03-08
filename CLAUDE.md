# CLAUDE.md — Agent Instructions for osrs-mcp

## Git

**Never create git commits automatically.** All commits must be made manually by
the user, or in response to an explicit prompt like "commit this" or "create a
commit". Do not stage files or commit as part of any other task.

## Validation Loop

Every change MUST pass all three steps in order before being considered done:

```
make lint    # golangci-lint run ./...
make build   # go build ./...
make test    # go test -race -coverprofile=coverage.out -covermode=atomic ./...
```

Run them as `make check` for convenience. A change that passes lint and build but
not test is NOT done. A change that passes test but has < 100% coverage is NOT done.

## Coverage Requirement

**100% code coverage is mandatory.** Every new file must have a corresponding
`_test.go` file. HTTP-dependent code MUST be tested via the interface/fake
pattern (see `internal/wiki/client.go` for the canonical example). No real
HTTP calls in tests — use `net/http/httptest` servers or interface fakes.

Check coverage:
```
make cover
```

Any function below 100% coverage is a failure. `cmd/osrs-mcp/main.go` is the
only file excluded from this requirement (it wires dependencies and calls
`server.ServeStdio` which cannot be unit tested).

## Documentation Structure

### Master Catalog
All features must be listed in `/docs/feature-catalog.md`.

### Per-Feature Docs
New feature docs live at:
```
/docs/YYYY-MM/##-feature-title.md
```
Example: `/docs/2026-03/01-wiki-search.md`

The `##` is a zero-padded sequence number within the month.
Update `feature-catalog.md` whenever you add a new feature doc.

## Adding a New MCP Tool

1. Add the client method to the relevant interface in `internal/wiki/client.go`
   or `internal/prices/client.go`.
2. Add the real implementation and update `NewClientWithBase` if needed.
3. Create `internal/tools/<tool_name>.go` and `internal/tools/<tool_name>_test.go`.
4. Register the tool in `internal/tools/register.go`.
5. Add the feature to `docs/feature-catalog.md` and create the feature doc.
6. Run `make check` and `make cover` — both must be clean.

## Dependency Injection Pattern

All tool handlers receive their client dependencies via constructor injection,
never via global state:

```go
type WikiSearchTool struct {
    client wiki.Client
}

func NewWikiSearchTool(c wiki.Client) *WikiSearchTool { ... }
func (t *WikiSearchTool) Definition() mcp.Tool { ... }
func (t *WikiSearchTool) Handler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) { ... }
```

## HTTP Testability Pattern

Use the two-seam pattern established in `internal/wiki/client.go`:

- **Seam 1** — `HTTPClient` interface with `Do(*http.Request) (*http.Response, error)`:
  Allows injecting `httptest.NewServer(...).Client()` in tests.
- **Seam 2** — `wiki.Client` / `prices.Client` interfaces:
  Allows injecting inline fake structs in tool handler tests.

Constructor for tests: `NewClientWithBase(httpClient HTTPClient, baseURL string)`

## SDK

This project uses `github.com/mark3labs/mcp-go`. Do NOT switch SDKs without
updating this file and all tool registrations.

## Go Version

Requires Go 1.23+. Do not use features from later versions without updating `go.mod`.
