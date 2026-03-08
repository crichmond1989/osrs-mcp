package tools

import (
	"context"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/crichmond1989/osrs-mcp/internal/hiscores"
)

func makePlayerStatsReq(args map[string]any) mcp.CallToolRequest {
	req := mcp.CallToolRequest{}
	req.Params.Arguments = args
	return req
}

func TestPlayerStatsTool_Definition(t *testing.T) {
	tool := NewPlayerStatsTool(nil)
	def := tool.Definition()
	if def.Name != "player_stats" {
		t.Errorf("Name = %q, want player_stats", def.Name)
	}
}

func TestPlayerStatsTool_Success(t *testing.T) {
	fake := &fakeHiscoresClient{
		getStatsFn: func(_ context.Context, _ string, _ string) (*hiscores.HiscoresResponse, error) {
			return &hiscores.HiscoresResponse{
				Skills: []hiscores.SkillEntry{
					{ID: 0, Name: "Overall", Rank: 1, Level: 2277, XP: 4600000000},
					{ID: 1, Name: "Attack", Rank: 1234, Level: 99, XP: 200000000},
				},
			}, nil
		},
	}
	tool := NewPlayerStatsTool(fake)
	req := makePlayerStatsReq(map[string]any{"player": "Zezima"})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.IsError {
		t.Fatal("expected success result")
	}
	content := res.Content[0].(mcp.TextContent).Text
	if !strings.Contains(content, "Zezima") {
		t.Error("output missing player name")
	}
	if !strings.Contains(content, "Overall") {
		t.Error("output missing skill name")
	}
	if !strings.Contains(content, "Attack") {
		t.Error("output missing skill name")
	}
}

func TestPlayerStatsTool_DefaultMode(t *testing.T) {
	var capturedMode string
	fake := &fakeHiscoresClient{
		getStatsFn: func(_ context.Context, _ string, mode string) (*hiscores.HiscoresResponse, error) {
			capturedMode = mode
			return &hiscores.HiscoresResponse{}, nil
		},
	}
	tool := NewPlayerStatsTool(fake)
	req := makePlayerStatsReq(map[string]any{"player": "Zezima"})
	_, _ = tool.Handler(context.Background(), req)
	if capturedMode != "standard" {
		t.Errorf("default mode = %q, want standard", capturedMode)
	}
}

func TestPlayerStatsTool_AllModes(t *testing.T) {
	modes := []string{"standard", "ironman", "hardcore", "ultimate"}
	for _, mode := range modes {
		mode := mode
		t.Run(mode, func(t *testing.T) {
			var capturedMode string
			fake := &fakeHiscoresClient{
				getStatsFn: func(_ context.Context, _ string, m string) (*hiscores.HiscoresResponse, error) {
					capturedMode = m
					return &hiscores.HiscoresResponse{}, nil
				},
			}
			tool := NewPlayerStatsTool(fake)
			req := makePlayerStatsReq(map[string]any{"player": "x", "mode": mode})
			_, _ = tool.Handler(context.Background(), req)
			if capturedMode != mode {
				t.Errorf("mode = %q, want %q", capturedMode, mode)
			}
		})
	}
}

func TestPlayerStatsTool_PlayerNotFound(t *testing.T) {
	fake := &fakeHiscoresClient{
		getStatsFn: func(_ context.Context, _ string, _ string) (*hiscores.HiscoresResponse, error) {
			return nil, hiscores.ErrPlayerNotFound
		},
	}
	tool := NewPlayerStatsTool(fake)
	req := makePlayerStatsReq(map[string]any{"player": "NoSuchPlayer"})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected error result for player not found")
	}
}

func TestPlayerStatsTool_ClientError(t *testing.T) {
	fake := &fakeHiscoresClient{
		getStatsFn: func(_ context.Context, _ string, _ string) (*hiscores.HiscoresResponse, error) {
			return nil, errFixed
		},
	}
	tool := NewPlayerStatsTool(fake)
	req := makePlayerStatsReq(map[string]any{"player": "Zezima"})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected error result for client error")
	}
}

func TestPlayerStatsTool_MissingPlayer(t *testing.T) {
	tool := NewPlayerStatsTool(nil)
	req := makePlayerStatsReq(map[string]any{})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected error result for missing player param")
	}
}

func TestPlayerStatsTool_UnrankedSkill(t *testing.T) {
	fake := &fakeHiscoresClient{
		getStatsFn: func(_ context.Context, _ string, _ string) (*hiscores.HiscoresResponse, error) {
			return &hiscores.HiscoresResponse{
				Skills: []hiscores.SkillEntry{
					{ID: 0, Name: "Overall", Rank: -1, Level: 10, XP: 1000},
				},
			}, nil
		},
	}
	tool := NewPlayerStatsTool(fake)
	req := makePlayerStatsReq(map[string]any{"player": "Zezima"})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.IsError {
		t.Fatal("expected success result")
	}
	content := res.Content[0].(mcp.TextContent).Text
	if !strings.Contains(content, "unranked") {
		t.Error("output missing 'unranked' for rank=-1 skill")
	}
}
