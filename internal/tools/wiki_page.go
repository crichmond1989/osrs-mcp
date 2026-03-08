package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/crich/osrs-mcp/internal/wiki"
)

// WikiPageTool implements the wiki_page MCP tool.
type WikiPageTool struct {
	client wiki.Client
}

// NewWikiPageTool constructs a WikiPageTool with the given client.
func NewWikiPageTool(c wiki.Client) *WikiPageTool {
	return &WikiPageTool{client: c}
}

// Definition returns the MCP tool metadata.
func (t *WikiPageTool) Definition() mcp.Tool {
	return mcp.NewTool("wiki_page",
		mcp.WithDescription("Retrieve the wikitext content of an OSRS Wiki page by its exact title."),
		mcp.WithString("title",
			mcp.Required(),
			mcp.Description("Exact page title, e.g. 'Abyssal whip' or 'Dragon Slayer'"),
		),
	)
}

// Handler is the MCP tool handler for wiki_page.
func (t *WikiPageTool) Handler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	title, err := req.RequireString("title")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	res, err := t.client.GetPage(ctx, title)
	if err != nil {
		if errors.Is(err, wiki.ErrPageNotFound) {
			return mcp.NewToolResultError(fmt.Sprintf("page not found: %q", title)), nil
		}
		return mcp.NewToolResultError(fmt.Sprintf("wiki page fetch failed: %v", err)), nil
	}

	for _, page := range res.Query.Pages {
		if len(page.Revisions) == 0 {
			return mcp.NewToolResultError(fmt.Sprintf("page %q has no content", title)), nil
		}
		return mcp.NewToolResultText(page.Revisions[0].Slots.Main.Content), nil
	}

	return mcp.NewToolResultError(fmt.Sprintf("page not found: %q", title)), nil
}
