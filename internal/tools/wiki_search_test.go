package tools

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/crichmond1989/osrs-mcp/internal/wiki"
)

func makeWikiSearchReq(args map[string]any) mcp.CallToolRequest {
	req := mcp.CallToolRequest{}
	req.Params.Arguments = args
	return req
}

func TestWikiSearchTool_Definition(t *testing.T) {
	tool := NewWikiSearchTool(nil)
	def := tool.Definition()
	if def.Name != "wiki_search" {
		t.Errorf("Name = %q, want wiki_search", def.Name)
	}
}

func TestWikiSearchTool_Success(t *testing.T) {
	fake := &fakeWikiClient{
		openSearchFn: func(_ context.Context, query string, limit int) (*wiki.OpenSearchResponse, error) {
			return &wiki.OpenSearchResponse{
				Query:  query,
				Titles: []string{"Abyssal whip"},
				URLs:   []string{"https://wiki.example/Abyssal_whip"},
			}, nil
		},
	}
	tool := NewWikiSearchTool(fake)
	req := makeWikiSearchReq(map[string]any{"query": "abyssal whip"})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.IsError {
		t.Fatal("expected success result")
	}
}

func TestWikiSearchTool_EmptyResults(t *testing.T) {
	fake := &fakeWikiClient{
		openSearchFn: func(_ context.Context, _ string, _ int) (*wiki.OpenSearchResponse, error) {
			return &wiki.OpenSearchResponse{Titles: []string{}}, nil
		},
	}
	tool := NewWikiSearchTool(fake)
	req := makeWikiSearchReq(map[string]any{"query": "xyznotfound"})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestWikiSearchTool_ClientError(t *testing.T) {
	fake := &fakeWikiClient{
		openSearchFn: func(_ context.Context, _ string, _ int) (*wiki.OpenSearchResponse, error) {
			return nil, errFixed
		},
	}
	tool := NewWikiSearchTool(fake)
	req := makeWikiSearchReq(map[string]any{"query": "q"})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected error result")
	}
}

func TestWikiSearchTool_MissingQuery(t *testing.T) {
	tool := NewWikiSearchTool(nil)
	req := makeWikiSearchReq(map[string]any{})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected error result for missing query")
	}
}

func TestWikiSearchTool_LimitClamping(t *testing.T) {
	var capturedLimit int
	fake := &fakeWikiClient{
		openSearchFn: func(_ context.Context, _ string, limit int) (*wiki.OpenSearchResponse, error) {
			capturedLimit = limit
			return &wiki.OpenSearchResponse{Titles: []string{"X"}, URLs: []string{"u"}}, nil
		},
	}
	tool := NewWikiSearchTool(fake)

	// limit = 0 should be clamped to 5
	req := makeWikiSearchReq(map[string]any{"query": "q", "limit": 0})
	_, _ = tool.Handler(context.Background(), req)
	if capturedLimit != 5 {
		t.Errorf("limit after clamp = %d, want 5", capturedLimit)
	}

	// limit = 99 should be clamped to 5
	req = makeWikiSearchReq(map[string]any{"query": "q", "limit": 99})
	_, _ = tool.Handler(context.Background(), req)
	if capturedLimit != 5 {
		t.Errorf("limit after clamp = %d, want 5", capturedLimit)
	}
}
