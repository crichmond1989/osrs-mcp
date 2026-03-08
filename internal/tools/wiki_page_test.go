package tools

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/crich/osrs-mcp/internal/wiki"
)

func makeWikiPageReq(args map[string]any) mcp.CallToolRequest {
	req := mcp.CallToolRequest{}
	req.Params.Arguments = args
	return req
}

func TestWikiPageTool_Definition(t *testing.T) {
	tool := NewWikiPageTool(nil)
	def := tool.Definition()
	if def.Name != "wiki_page" {
		t.Errorf("Name = %q, want wiki_page", def.Name)
	}
}

func TestWikiPageTool_Success(t *testing.T) {
	fake := &fakeWikiClient{
		getPageFn: func(_ context.Context, _ string) (*wiki.PageQueryResponse, error) {
			var res wiki.PageQueryResponse
			res.Query.Pages = map[string]wiki.Page{
				"123": {
					PageID: 123,
					Title:  "Abyssal whip",
					Revisions: []wiki.PageRevision{
						func() wiki.PageRevision {
							var r wiki.PageRevision
							r.Slots.Main.Content = "wikitext content"
							return r
						}(),
					},
				},
			}
			return &res, nil
		},
	}
	tool := NewWikiPageTool(fake)
	req := makeWikiPageReq(map[string]any{"title": "Abyssal whip"})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.IsError {
		t.Fatal("expected success result")
	}
}

func TestWikiPageTool_NoRevisions(t *testing.T) {
	fake := &fakeWikiClient{
		getPageFn: func(_ context.Context, _ string) (*wiki.PageQueryResponse, error) {
			var res wiki.PageQueryResponse
			res.Query.Pages = map[string]wiki.Page{
				"123": {PageID: 123, Title: "Empty page", Revisions: nil},
			}
			return &res, nil
		},
	}
	tool := NewWikiPageTool(fake)
	req := makeWikiPageReq(map[string]any{"title": "Empty page"})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected error result for page with no revisions")
	}
}

func TestWikiPageTool_EmptyPages(t *testing.T) {
	fake := &fakeWikiClient{
		getPageFn: func(_ context.Context, _ string) (*wiki.PageQueryResponse, error) {
			var res wiki.PageQueryResponse
			res.Query.Pages = map[string]wiki.Page{}
			return &res, nil
		},
	}
	tool := NewWikiPageTool(fake)
	req := makeWikiPageReq(map[string]any{"title": "Whatever"})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected error result for empty pages map")
	}
}

func TestWikiPageTool_NotFound(t *testing.T) {
	fake := &fakeWikiClient{
		getPageFn: func(_ context.Context, _ string) (*wiki.PageQueryResponse, error) {
			return nil, wiki.ErrPageNotFound
		},
	}
	tool := NewWikiPageTool(fake)
	req := makeWikiPageReq(map[string]any{"title": "Missing"})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected error result for not found page")
	}
}

func TestWikiPageTool_ClientError(t *testing.T) {
	fake := &fakeWikiClient{
		getPageFn: func(_ context.Context, _ string) (*wiki.PageQueryResponse, error) {
			return nil, errFixed
		},
	}
	tool := NewWikiPageTool(fake)
	req := makeWikiPageReq(map[string]any{"title": "Whatever"})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected error result for client error")
	}
}

func TestWikiPageTool_MissingTitle(t *testing.T) {
	tool := NewWikiPageTool(nil)
	req := makeWikiPageReq(map[string]any{})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected error result for missing title")
	}
}
