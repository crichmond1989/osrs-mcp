package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/crichmond1989/osrs-mcp/internal/wiki"
)

// WikiSearchFullTool implements the wiki_search_full MCP tool using action=query&list=search.
// It returns richer metadata than wiki_search (page ID, word count, snippet).
type WikiSearchFullTool struct {
	client wiki.Client
}

// NewWikiSearchFullTool constructs a WikiSearchFullTool with the given client.
func NewWikiSearchFullTool(c wiki.Client) *WikiSearchFullTool {
	return &WikiSearchFullTool{client: c}
}

// Definition returns the MCP tool metadata.
func (t *WikiSearchFullTool) Definition() mcp.Tool {
	return mcp.NewTool("wiki_search_full",
		mcp.WithDescription("Search the OSRS Wiki with full metadata including page IDs, word counts, and text snippets."),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Search keywords, e.g. 'abyssal whip' or 'dragon slayer quest'"),
		),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of results (1-10, default 5)"),
		),
	)
}

// Handler is the MCP tool handler for wiki_search_full.
func (t *WikiSearchFullTool) Handler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, err := req.RequireString("query")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	limit := req.GetInt("limit", 5)
	if limit < 1 || limit > 10 {
		limit = 5
	}

	res, err := t.client.SearchPages(ctx, query, limit)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("wiki search failed: %v", err)), nil
	}

	if len(res.Query.Search) == 0 {
		return mcp.NewToolResultText("No results found."), nil
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "Total hits: %d\n\n", res.Query.SearchInfo.TotalHits)
	for i, r := range res.Query.Search {
		fmt.Fprintf(&sb, "%d. %s (ID: %d, %d words)\n", i+1, r.Title, r.PageID, r.WordCount)
		if r.Snippet != "" {
			fmt.Fprintf(&sb, "   %s\n", r.Snippet)
		}
	}
	return mcp.NewToolResultText(sb.String()), nil
}
