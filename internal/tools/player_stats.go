package tools

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/crich/osrs-mcp/internal/hiscores"
)

// PlayerStatsTool implements the player_stats MCP tool.
type PlayerStatsTool struct {
	client hiscores.Client
}

// NewPlayerStatsTool constructs a PlayerStatsTool with the given client.
func NewPlayerStatsTool(c hiscores.Client) *PlayerStatsTool {
	return &PlayerStatsTool{client: c}
}

// Definition returns the MCP tool metadata.
func (t *PlayerStatsTool) Definition() mcp.Tool {
	return mcp.NewTool("player_stats",
		mcp.WithDescription("Fetch OSRS player skill levels, XP, and hiscores ranks for all 24 skills."),
		mcp.WithString("player",
			mcp.Required(),
			mcp.Description("OSRS player name (case-insensitive)"),
		),
		mcp.WithString("mode",
			mcp.Description("Hiscores mode: 'standard' (default), 'ironman', 'hardcore', or 'ultimate'"),
			mcp.Enum("standard", "ironman", "hardcore", "ultimate"),
		),
	)
}

// Handler is the MCP tool handler for player_stats.
func (t *PlayerStatsTool) Handler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	player, err := req.RequireString("player")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	mode := req.GetString("mode", "standard")

	res, err := t.client.GetStats(ctx, player, mode)
	if err != nil {
		if errors.Is(err, hiscores.ErrPlayerNotFound) {
			return mcp.NewToolResultError(fmt.Sprintf("player not found: %q", player)), nil
		}
		return mcp.NewToolResultError(fmt.Sprintf("hiscores lookup failed: %v", err)), nil
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "Player: %s (%s)\n\n", player, mode)
	fmt.Fprintf(&sb, "%-16s %5s %12s %8s\n", "Skill", "Level", "XP", "Rank")
	fmt.Fprintf(&sb, "%-16s %5s %12s %8s\n", "-----", "-----", "--", "----")

	for _, skill := range res.Skills {
		rank := fmt.Sprintf("%d", skill.Rank)
		if skill.Rank == -1 {
			rank = "unranked"
		}
		fmt.Fprintf(&sb, "%-16s %5d %12d %8s\n", skill.Name, skill.Level, skill.XP, rank)
	}

	return mcp.NewToolResultText(sb.String()), nil
}
