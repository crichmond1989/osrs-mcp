package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/crich/osrs-mcp/internal/prices"
)

// PricesTimeSeriesTool implements the prices_timeseries MCP tool.
type PricesTimeSeriesTool struct {
	client prices.Client
}

// NewPricesTimeSeriesTool constructs a PricesTimeSeriesTool with the given client.
func NewPricesTimeSeriesTool(c prices.Client) *PricesTimeSeriesTool {
	return &PricesTimeSeriesTool{client: c}
}

// Definition returns the MCP tool metadata.
func (t *PricesTimeSeriesTool) Definition() mcp.Tool {
	return mcp.NewTool("prices_timeseries",
		mcp.WithDescription("Get averaged OSRS Grand Exchange price and volume data over a time window (5m, 1h, or 24h)."),
		mcp.WithNumber("id",
			mcp.Required(),
			mcp.Description("The item ID, e.g. 4151 for Abyssal whip"),
		),
		mcp.WithString("window",
			mcp.Required(),
			mcp.Description("Time window: '5m', '1h', or '24h'"),
			mcp.Enum("5m", "1h", "24h"),
		),
	)
}

// Handler is the MCP tool handler for prices_timeseries.
func (t *PricesTimeSeriesTool) Handler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	idVal, err := req.RequireInt("id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	window, err := req.RequireString("window")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	res, err := t.client.GetTimeSeries(ctx, window)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("timeseries lookup failed: %v", err)), nil
	}

	idStr := fmt.Sprintf("%d", idVal)
	p, ok := res.Data[idStr]
	if !ok {
		return mcp.NewToolResultError(fmt.Sprintf("no timeseries data found for item ID %d in window %s", idVal, window)), nil
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "Item ID: %d | Window: %s\n", idVal, window)
	if p.AvgHighPrice != nil {
		fmt.Fprintf(&sb, "Avg buy price:   %d gp (volume: %d)\n", *p.AvgHighPrice, p.HighPriceVolume)
	} else {
		sb.WriteString("Avg buy price:   no trades\n")
	}
	if p.AvgLowPrice != nil {
		fmt.Fprintf(&sb, "Avg sell price:  %d gp (volume: %d)\n", *p.AvgLowPrice, p.LowPriceVolume)
	} else {
		sb.WriteString("Avg sell price:  no trades\n")
	}
	return mcp.NewToolResultText(sb.String()), nil
}
