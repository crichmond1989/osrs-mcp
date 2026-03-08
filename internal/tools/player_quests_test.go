package tools

import (
	"context"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/crichmond1989/osrs-mcp/internal/wikisync"
)

func makePlayerQuestsReq(args map[string]any) mcp.CallToolRequest {
	req := mcp.CallToolRequest{}
	req.Params.Arguments = args
	return req
}

func TestPlayerQuestsTool_Definition(t *testing.T) {
	tool := NewPlayerQuestsTool(nil)
	def := tool.Definition()
	if def.Name != "player_quests" {
		t.Errorf("Name = %q, want player_quests", def.Name)
	}
}

func TestPlayerQuestsTool_Success(t *testing.T) {
	fake := &fakeWikiSyncClient{
		getPlayerDataFn: func(_ context.Context, _ string) (*wikisync.WikiSyncResponse, error) {
			return &wikisync.WikiSyncResponse{
				Quests: map[string]int{
					"Dragon Slayer I":    wikisync.QuestComplete,
					"Barbarian Training": wikisync.QuestInProgress,
					"Desert Treasure I":  wikisync.QuestNotStarted,
				},
			}, nil
		},
	}
	tool := NewPlayerQuestsTool(fake)
	req := makePlayerQuestsReq(map[string]any{"player": "Zezima"})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.IsError {
		t.Fatal("expected success result")
	}
	content := res.Content[0].(mcp.TextContent).Text
	if !strings.Contains(content, "Complete (1)") {
		t.Error("output missing complete group header")
	}
	if !strings.Contains(content, "Dragon Slayer I") {
		t.Error("output missing complete quest name")
	}
	if !strings.Contains(content, "In Progress (1)") {
		t.Error("output missing in progress group header")
	}
	if !strings.Contains(content, "Barbarian Training") {
		t.Error("output missing in-progress quest name")
	}
	if !strings.Contains(content, "Not Started (1)") {
		t.Error("output missing not started group header")
	}
	if !strings.Contains(content, "Desert Treasure I") {
		t.Error("output missing not started quest name")
	}
}

func TestPlayerQuestsTool_GuthixHer0(t *testing.T) {
	// Realistic fixture based on real WikiSync data for Guthix Her0.
	// Fight Arena is complete (2); Meat and Greet is not started (0).
	fake := &fakeWikiSyncClient{
		getPlayerDataFn: func(_ context.Context, player string) (*wikisync.WikiSyncResponse, error) {
			if player != "Guthix Her0" {
				t.Errorf("unexpected player %q", player)
			}
			return &wikisync.WikiSyncResponse{
				Quests: map[string]int{
					"Fight Arena":        wikisync.QuestComplete,
					"Demon Slayer":       wikisync.QuestComplete,
					"Dragon Slayer I":    wikisync.QuestComplete,
					"Goblin Diplomacy":   wikisync.QuestComplete,
					"Barbarian Training": wikisync.QuestInProgress,
					"Meat and Greet":     wikisync.QuestNotStarted,
					"Desert Treasure I":  wikisync.QuestNotStarted,
				},
			}, nil
		},
	}
	tool := NewPlayerQuestsTool(fake)
	req := makePlayerQuestsReq(map[string]any{"player": "Guthix Her0"})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.IsError {
		t.Fatalf("unexpected tool error: %v", res.Content)
	}
	content := res.Content[0].(mcp.TextContent).Text

	// Fight Arena should appear under Complete.
	if !strings.Contains(content, "Complete (4)") {
		t.Error("expected 4 complete quests")
	}
	if !strings.Contains(content, "  - Fight Arena") {
		t.Error("Fight Arena should be listed as complete")
	}

	// Barbarian Training should appear under In Progress.
	if !strings.Contains(content, "In Progress (1)") {
		t.Error("expected 1 in-progress quest")
	}
	if !strings.Contains(content, "  - Barbarian Training") {
		t.Error("Barbarian Training should be listed as in progress")
	}

	// Meat and Greet should appear under Not Started.
	if !strings.Contains(content, "Not Started (2)") {
		t.Error("expected 2 not-started quests")
	}
	if !strings.Contains(content, "  - Meat and Greet") {
		t.Error("Meat and Greet should be listed as not started")
	}
}

func TestPlayerQuestsTool_WikiSyncNotEnabled(t *testing.T) {
	fake := &fakeWikiSyncClient{
		getPlayerDataFn: func(_ context.Context, _ string) (*wikisync.WikiSyncResponse, error) {
			return nil, wikisync.ErrWikiSyncNotEnabled
		},
	}
	tool := NewPlayerQuestsTool(fake)
	req := makePlayerQuestsReq(map[string]any{"player": "NoWikiSync"})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected error result for WikiSync not enabled")
	}
	content := res.Content[0].(mcp.TextContent).Text
	if !strings.Contains(content, "WikiSync") {
		t.Error("error message missing WikiSync reference")
	}
}

func TestPlayerQuestsTool_ClientError(t *testing.T) {
	fake := &fakeWikiSyncClient{
		getPlayerDataFn: func(_ context.Context, _ string) (*wikisync.WikiSyncResponse, error) {
			return nil, errFixed
		},
	}
	tool := NewPlayerQuestsTool(fake)
	req := makePlayerQuestsReq(map[string]any{"player": "Zezima"})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected error result for client error")
	}
}

func TestPlayerQuestsTool_MissingPlayer(t *testing.T) {
	tool := NewPlayerQuestsTool(nil)
	req := makePlayerQuestsReq(map[string]any{})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected error result for missing player param")
	}
}
