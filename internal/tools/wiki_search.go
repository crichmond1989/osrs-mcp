package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/crich/osrs-mcp/internal/wiki"
)

// WikiSearchTool implements the wiki_search MCP tool using action=opensearch.
type WikiSearchTool struct {
	client wiki.Client
}

// NewWikiSearchTool constructs a WikiSearchTool with the given client.
func NewWikiSearchTool(c wiki.Client) *WikiSearchTool {
	return &WikiSearchTool{client: c}
}

// Definition returns the MCP tool metadata.
func (t *WikiSearchTool) Definition() mcp.Tool {
	return mcp.NewTool("wiki_search",
		mcp.WithDescription("Search the OSRS Wiki by keyword. Returns matching page titles and URLs."),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Search keywords, e.g. 'abyssal whip' or 'dragon slayer quest'"),
		),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of results (1-10, default 5)"),
		),
	)
}

// Handler is the MCP tool handler for wiki_search.
func (t *WikiSearchTool) Handler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, err := req.RequireString("query")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	limit := req.GetInt("limit", 5)
	if limit < 1 || limit > 10 {
		limit = 5
	}

	res, err := t.client.OpenSearch(ctx, query, limit)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("wiki search failed: %v", err)), nil
	}

	if len(res.Titles) == 0 {
		return mcp.NewToolResultText("No results found."), nil
	}

	var sb strings.Builder
	for i, title := range res.Titles {
		fmt.Fprintf(&sb, "%d. %s\n   %s\n", i+1, title, res.URLs[i])
	}
	return mcp.NewToolResultText(sb.String()), nil
}
