package tools

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/crichmond1989/osrs-mcp/internal/wiki"
)

func makeQuestInfoReq(args map[string]any) mcp.CallToolRequest {
	req := mcp.CallToolRequest{}
	req.Params.Arguments = args
	return req
}

func TestQuestInfoTool_Definition(t *testing.T) {
	tool := NewQuestInfoTool(nil)
	def := tool.Definition()
	if def.Name != "quest_info" {
		t.Errorf("Name = %q, want quest_info", def.Name)
	}
}

func TestQuestInfoTool_Success(t *testing.T) {
	fake := &fakeWikiClient{
		getPageFn: func(_ context.Context, _ string) (*wiki.PageQueryResponse, error) {
			var res wiki.PageQueryResponse
			res.Query.Pages = map[string]wiki.Page{
				"123": {
					PageID: 123,
					Title:  "Dragon Slayer I",
					Revisions: []wiki.PageRevision{
						func() wiki.PageRevision {
							var r wiki.PageRevision
							r.Slots.Main.Content = "==Requirements==\n* 32 {{Skill|Quest points}}\n==Rewards==\n* 18,650 Strength XP"
							return r
						}(),
					},
				},
			}
			return &res, nil
		},
	}
	tool := NewQuestInfoTool(fake)
	req := makeQuestInfoReq(map[string]any{"quest": "Dragon Slayer I"})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.IsError {
		t.Fatal("expected success result")
	}
}

func TestQuestInfoTool_NoRevisions(t *testing.T) {
	fake := &fakeWikiClient{
		getPageFn: func(_ context.Context, _ string) (*wiki.PageQueryResponse, error) {
			var res wiki.PageQueryResponse
			res.Query.Pages = map[string]wiki.Page{
				"123": {PageID: 123, Title: "Empty Quest", Revisions: nil},
			}
			return &res, nil
		},
	}
	tool := NewQuestInfoTool(fake)
	req := makeQuestInfoReq(map[string]any{"quest": "Empty Quest"})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected error result for quest with no revisions")
	}
}

func TestQuestInfoTool_EmptyPages(t *testing.T) {
	fake := &fakeWikiClient{
		getPageFn: func(_ context.Context, _ string) (*wiki.PageQueryResponse, error) {
			var res wiki.PageQueryResponse
			res.Query.Pages = map[string]wiki.Page{}
			return &res, nil
		},
	}
	tool := NewQuestInfoTool(fake)
	req := makeQuestInfoReq(map[string]any{"quest": "Whatever"})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected error result for empty pages map")
	}
}

func TestQuestInfoTool_NotFound(t *testing.T) {
	fake := &fakeWikiClient{
		getPageFn: func(_ context.Context, _ string) (*wiki.PageQueryResponse, error) {
			return nil, wiki.ErrPageNotFound
		},
	}
	tool := NewQuestInfoTool(fake)
	req := makeQuestInfoReq(map[string]any{"quest": "Missing Quest"})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected error result for not found quest")
	}
}

func TestQuestInfoTool_ClientError(t *testing.T) {
	fake := &fakeWikiClient{
		getPageFn: func(_ context.Context, _ string) (*wiki.PageQueryResponse, error) {
			return nil, errFixed
		},
	}
	tool := NewQuestInfoTool(fake)
	req := makeQuestInfoReq(map[string]any{"quest": "Whatever"})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected error result for client error")
	}
}

func TestQuestInfoTool_MissingQuest(t *testing.T) {
	tool := NewQuestInfoTool(nil)
	req := makeQuestInfoReq(map[string]any{})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected error result for missing quest param")
	}
}
