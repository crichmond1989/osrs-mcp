package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/crich/osrs-mcp/internal/prices"
)

// PricesLatestTool implements the prices_latest MCP tool.
type PricesLatestTool struct {
	client prices.Client
}

// NewPricesLatestTool constructs a PricesLatestTool with the given client.
func NewPricesLatestTool(c prices.Client) *PricesLatestTool {
	return &PricesLatestTool{client: c}
}

// Definition returns the MCP tool metadata.
func (t *PricesLatestTool) Definition() mcp.Tool {
	return mcp.NewTool("prices_latest",
		mcp.WithDescription("Get the latest instant-buy and instant-sell prices for an OSRS item by its item ID."),
		mcp.WithNumber("id",
			mcp.Required(),
			mcp.Description("The item ID from the OSRS Grand Exchange, e.g. 4151 for Abyssal whip"),
		),
	)
}

// Handler is the MCP tool handler for prices_latest.
func (t *PricesLatestTool) Handler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	idVal, err := req.RequireInt("id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	res, err := t.client.GetLatest(ctx, []int{idVal})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("prices lookup failed: %v", err)), nil
	}

	idStr := fmt.Sprintf("%d", idVal)
	p, ok := res.Data[idStr]
	if !ok {
		return mcp.NewToolResultError(fmt.Sprintf("no price data found for item ID %d", idVal)), nil
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "Item ID: %d\n", idVal)
	if p.High != nil {
		fmt.Fprintf(&sb, "Instant buy:  %d gp\n", *p.High)
	} else {
		sb.WriteString("Instant buy:  no data\n")
	}
	if p.Low != nil {
		fmt.Fprintf(&sb, "Instant sell: %d gp\n", *p.Low)
	} else {
		sb.WriteString("Instant sell: no data\n")
	}
	return mcp.NewToolResultText(sb.String()), nil
}
