package tools

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/crich/osrs-mcp/internal/wikisync"
)

// PlayerQuestsTool implements the player_quests MCP tool.
type PlayerQuestsTool struct {
	client wikisync.Client
}

// NewPlayerQuestsTool constructs a PlayerQuestsTool with the given client.
func NewPlayerQuestsTool(c wikisync.Client) *PlayerQuestsTool {
	return &PlayerQuestsTool{client: c}
}

// Definition returns the MCP tool metadata.
func (t *PlayerQuestsTool) Definition() mcp.Tool {
	return mcp.NewTool("player_quests",
		mcp.WithDescription("Fetch quest completion status for an OSRS player using WikiSync data. "+
			"Returns each quest grouped by completion state (complete, in progress, not started). "+
			"Requires the player to have enabled the WikiSync plugin in RuneLite."),
		mcp.WithString("player",
			mcp.Required(),
			mcp.Description("OSRS player name"),
		),
	)
}

// Handler is the MCP tool handler for player_quests.
func (t *PlayerQuestsTool) Handler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	player, err := req.RequireString("player")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	res, err := t.client.GetPlayerData(ctx, player)
	if err != nil {
		if errors.Is(err, wikisync.ErrWikiSyncNotEnabled) {
			return mcp.NewToolResultError(fmt.Sprintf(
				"player %q has not enabled WikiSync. "+
					"To use this tool the player must install the WikiSync plugin in RuneLite.",
				player,
			)), nil
		}
		return mcp.NewToolResultError(fmt.Sprintf("WikiSync lookup failed: %v", err)), nil
	}

	var complete, inProgress, notStarted []string
	for name, state := range res.Quests {
		switch state {
		case wikisync.QuestComplete:
			complete = append(complete, name)
		case wikisync.QuestInProgress:
			inProgress = append(inProgress, name)
		default:
			notStarted = append(notStarted, name)
		}
	}
	sort.Strings(complete)
	sort.Strings(inProgress)
	sort.Strings(notStarted)

	var sb strings.Builder
	fmt.Fprintf(&sb, "Player: %s (WikiSync data)\n\n", player)
	writeGroup := func(label string, quests []string) {
		fmt.Fprintf(&sb, "%s (%d):\n", label, len(quests))
		for _, q := range quests {
			fmt.Fprintf(&sb, "  - %s\n", q)
		}
		fmt.Fprintln(&sb)
	}
	writeGroup("Complete", complete)
	writeGroup("In Progress", inProgress)
	writeGroup("Not Started", notStarted)

	return mcp.NewToolResultText(sb.String()), nil
}
