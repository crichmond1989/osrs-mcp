package tools

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/crichmond1989/osrs-mcp/internal/wiki"
)

func makeWikiSearchFullReq(args map[string]any) mcp.CallToolRequest {
	req := mcp.CallToolRequest{}
	req.Params.Arguments = args
	return req
}

func TestWikiSearchFullTool_Definition(t *testing.T) {
	tool := NewWikiSearchFullTool(nil)
	def := tool.Definition()
	if def.Name != "wiki_search_full" {
		t.Errorf("Name = %q, want wiki_search_full", def.Name)
	}
}

func TestWikiSearchFullTool_Success(t *testing.T) {
	fake := &fakeWikiClient{
		searchPagesFn: func(_ context.Context, _ string, _ int) (*wiki.SearchResponse, error) {
			var res wiki.SearchResponse
			res.Query.SearchInfo.TotalHits = 1
			res.Query.Search = []wiki.SearchResult{
				{Title: "Abyssal whip", PageID: 123, WordCount: 300, Snippet: "A powerful whip"},
			}
			return &res, nil
		},
	}
	tool := NewWikiSearchFullTool(fake)
	req := makeWikiSearchFullReq(map[string]any{"query": "abyssal whip"})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.IsError {
		t.Fatal("expected success result")
	}
}

func TestWikiSearchFullTool_EmptyResults(t *testing.T) {
	fake := &fakeWikiClient{
		searchPagesFn: func(_ context.Context, _ string, _ int) (*wiki.SearchResponse, error) {
			return &wiki.SearchResponse{}, nil
		},
	}
	tool := NewWikiSearchFullTool(fake)
	req := makeWikiSearchFullReq(map[string]any{"query": "xyz"})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestWikiSearchFullTool_NoSnippet(t *testing.T) {
	fake := &fakeWikiClient{
		searchPagesFn: func(_ context.Context, _ string, _ int) (*wiki.SearchResponse, error) {
			var res wiki.SearchResponse
			res.Query.SearchInfo.TotalHits = 1
			res.Query.Search = []wiki.SearchResult{
				{Title: "Dragon scimitar", PageID: 456, WordCount: 200, Snippet: ""},
			}
			return &res, nil
		},
	}
	tool := NewWikiSearchFullTool(fake)
	req := makeWikiSearchFullReq(map[string]any{"query": "dragon"})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.IsError {
		t.Fatal("expected success result")
	}
}

func TestWikiSearchFullTool_ClientError(t *testing.T) {
	fake := &fakeWikiClient{
		searchPagesFn: func(_ context.Context, _ string, _ int) (*wiki.SearchResponse, error) {
			return nil, errFixed
		},
	}
	tool := NewWikiSearchFullTool(fake)
	req := makeWikiSearchFullReq(map[string]any{"query": "q"})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected error result")
	}
}

func TestWikiSearchFullTool_MissingQuery(t *testing.T) {
	tool := NewWikiSearchFullTool(nil)
	req := makeWikiSearchFullReq(map[string]any{})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected error result for missing query")
	}
}

func TestWikiSearchFullTool_LimitClamping(t *testing.T) {
	var capturedLimit int
	fake := &fakeWikiClient{
		searchPagesFn: func(_ context.Context, _ string, limit int) (*wiki.SearchResponse, error) {
			capturedLimit = limit
			var res wiki.SearchResponse
			res.Query.Search = []wiki.SearchResult{{Title: "X"}}
			return &res, nil
		},
	}
	tool := NewWikiSearchFullTool(fake)

	req := makeWikiSearchFullReq(map[string]any{"query": "q", "limit": 0})
	_, _ = tool.Handler(context.Background(), req)
	if capturedLimit != 5 {
		t.Errorf("limit = %d, want 5", capturedLimit)
	}

	req = makeWikiSearchFullReq(map[string]any{"query": "q", "limit": 99})
	_, _ = tool.Handler(context.Background(), req)
	if capturedLimit != 5 {
		t.Errorf("limit = %d, want 5", capturedLimit)
	}
}
