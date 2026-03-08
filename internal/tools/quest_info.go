package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/crich/osrs-mcp/internal/wiki"
)

// QuestInfoTool implements the quest_info MCP tool.
type QuestInfoTool struct {
	client wiki.Client
}

// NewQuestInfoTool constructs a QuestInfoTool with the given client.
func NewQuestInfoTool(c wiki.Client) *QuestInfoTool {
	return &QuestInfoTool{client: c}
}

// Definition returns the MCP tool metadata.
func (t *QuestInfoTool) Definition() mcp.Tool {
	return mcp.NewTool("quest_info",
		mcp.WithDescription("Fetch OSRS quest information from the Wiki. "+
			"Returns wikitext for the quest page — look for ==Requirements== and ==Rewards== sections "+
			"to find skill requirements, item requirements, quest point rewards, and XP rewards."),
		mcp.WithString("quest",
			mcp.Required(),
			mcp.Description("Exact quest name as it appears on the OSRS Wiki, e.g. 'Dragon Slayer I'"),
		),
	)
}

// Handler is the MCP tool handler for quest_info.
func (t *QuestInfoTool) Handler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	quest, err := req.RequireString("quest")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	res, err := t.client.GetPage(ctx, quest)
	if err != nil {
		if errors.Is(err, wiki.ErrPageNotFound) {
			return mcp.NewToolResultError(fmt.Sprintf("quest not found: %q", quest)), nil
		}
		return mcp.NewToolResultError(fmt.Sprintf("quest info fetch failed: %v", err)), nil
	}

	for _, page := range res.Query.Pages {
		if len(page.Revisions) == 0 {
			return mcp.NewToolResultError(fmt.Sprintf("quest page %q has no content", quest)), nil
		}
		return mcp.NewToolResultText(page.Revisions[0].Slots.Main.Content), nil
	}

	return mcp.NewToolResultError(fmt.Sprintf("quest not found: %q", quest)), nil
}
