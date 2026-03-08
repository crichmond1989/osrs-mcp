package tools

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/crichmond1989/osrs-mcp/internal/prices"
)

func makePricesMappingReq(args map[string]any) mcp.CallToolRequest {
	req := mcp.CallToolRequest{}
	req.Params.Arguments = args
	return req
}

func TestPricesMappingTool_Definition(t *testing.T) {
	tool := NewPricesMappingTool(nil)
	def := tool.Definition()
	if def.Name != "prices_mapping" {
		t.Errorf("Name = %q, want prices_mapping", def.Name)
	}
}

func TestPricesMappingTool_NoFilter(t *testing.T) {
	fake := &fakePricesClient{
		getMappingFn: func(_ context.Context) ([]prices.MappingItem, error) {
			return []prices.MappingItem{
				{ID: 4151, Name: "Abyssal whip", Members: true, HighAlch: 108000, Limit: 70},
				{ID: 1, Name: "Coins", Members: false, HighAlch: 0, Limit: 0},
			}, nil
		},
	}
	tool := NewPricesMappingTool(fake)
	req := makePricesMappingReq(map[string]any{})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.IsError {
		t.Fatal("expected success result")
	}
}

func TestPricesMappingTool_WithFilter(t *testing.T) {
	fake := &fakePricesClient{
		getMappingFn: func(_ context.Context) ([]prices.MappingItem, error) {
			return []prices.MappingItem{
				{ID: 4151, Name: "Abyssal whip", Members: true, HighAlch: 108000, Limit: 70},
				{ID: 1, Name: "Coins", Members: false, HighAlch: 0, Limit: 0},
			}, nil
		},
	}
	tool := NewPricesMappingTool(fake)
	req := makePricesMappingReq(map[string]any{"query": "abyssal"})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.IsError {
		t.Fatal("expected success result")
	}
}

func TestPricesMappingTool_NoMatchingItems(t *testing.T) {
	fake := &fakePricesClient{
		getMappingFn: func(_ context.Context) ([]prices.MappingItem, error) {
			return []prices.MappingItem{
				{ID: 4151, Name: "Abyssal whip", Members: true},
			}, nil
		},
	}
	tool := NewPricesMappingTool(fake)
	req := makePricesMappingReq(map[string]any{"query": "xyznotexist"})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.IsError {
		t.Fatal("expected success result (no items found message)")
	}
}

func TestPricesMappingTool_Truncation(t *testing.T) {
	items := make([]prices.MappingItem, 25)
	for i := range items {
		items[i] = prices.MappingItem{ID: i, Name: "Item"}
	}
	fake := &fakePricesClient{
		getMappingFn: func(_ context.Context) ([]prices.MappingItem, error) {
			return items, nil
		},
	}
	tool := NewPricesMappingTool(fake)
	req := makePricesMappingReq(map[string]any{})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.IsError {
		t.Fatal("expected success result")
	}
}

func TestPricesMappingTool_ClientError(t *testing.T) {
	fake := &fakePricesClient{
		getMappingFn: func(_ context.Context) ([]prices.MappingItem, error) {
			return nil, errFixed
		},
	}
	tool := NewPricesMappingTool(fake)
	req := makePricesMappingReq(map[string]any{})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected error result for client error")
	}
}
