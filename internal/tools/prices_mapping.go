package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/crichmond1989/osrs-mcp/internal/prices"
)

// PricesMappingTool implements the prices_mapping MCP tool.
type PricesMappingTool struct {
	client prices.Client
}

// NewPricesMappingTool constructs a PricesMappingTool with the given client.
func NewPricesMappingTool(c prices.Client) *PricesMappingTool {
	return &PricesMappingTool{client: c}
}

// Definition returns the MCP tool metadata.
func (t *PricesMappingTool) Definition() mcp.Tool {
	return mcp.NewTool("prices_mapping",
		mcp.WithDescription("Look up OSRS item metadata (ID, name, alch values, GE limit) with an optional name filter."),
		mcp.WithString("query",
			mcp.Description("Optional name filter — returns items whose names contain this string (case-insensitive)"),
		),
	)
}

// Handler is the MCP tool handler for prices_mapping.
func (t *PricesMappingTool) Handler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query := req.GetString("query", "")

	items, err := t.client.GetMapping(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("mapping lookup failed: %v", err)), nil
	}

	var sb strings.Builder
	count := 0
	filter := strings.ToLower(query)
	for _, item := range items {
		if filter != "" && !strings.Contains(strings.ToLower(item.Name), filter) {
			continue
		}
		fmt.Fprintf(&sb, "ID: %d | %s", item.ID, item.Name)
		if item.Members {
			sb.WriteString(" [members]")
		}
		fmt.Fprintf(&sb, " | High alch: %d gp | Limit: %d\n", item.HighAlch, item.Limit)
		count++
		if count >= 20 {
			sb.WriteString("(results truncated to 20)\n")
			break
		}
	}

	if count == 0 {
		return mcp.NewToolResultText("No items found."), nil
	}
	return mcp.NewToolResultText(sb.String()), nil
}
