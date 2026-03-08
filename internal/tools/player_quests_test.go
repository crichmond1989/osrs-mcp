package tools

import (
	"context"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/crich/osrs-mcp/internal/wikisync"
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
					"0": wikisync.QuestComplete,
					"1": wikisync.QuestInProgress,
					"2": wikisync.QuestNotStarted,
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
		t.Fatal("expected success result")
	}
	content := res.Content[0].(mcp.TextContent).Text
	if !strings.Contains(content, "1 complete") {
		t.Error("output missing complete count")
	}
	if !strings.Contains(content, "1 in progress") {
		t.Error("output missing in progress count")
	}
	if !strings.Contains(content, "1 not started") {
		t.Error("output missing not started count")
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
