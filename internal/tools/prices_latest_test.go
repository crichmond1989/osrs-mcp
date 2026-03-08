package tools

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/crich/osrs-mcp/internal/prices"
)

func makePricesLatestReq(args map[string]any) mcp.CallToolRequest {
	req := mcp.CallToolRequest{}
	req.Params.Arguments = args
	return req
}

func intPtr(v int) *int { return &v }

func TestPricesLatestTool_Definition(t *testing.T) {
	tool := NewPricesLatestTool(nil)
	def := tool.Definition()
	if def.Name != "prices_latest" {
		t.Errorf("Name = %q, want prices_latest", def.Name)
	}
}

func TestPricesLatestTool_Success(t *testing.T) {
	fake := &fakePricesClient{
		getLatestFn: func(_ context.Context, ids []int) (*prices.LatestResponse, error) {
			return &prices.LatestResponse{
				Data: map[string]prices.LatestPrice{
					"4151": {High: intPtr(2500000), Low: intPtr(2400000)},
				},
			}, nil
		},
	}
	tool := NewPricesLatestTool(fake)
	req := makePricesLatestReq(map[string]any{"id": 4151})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.IsError {
		t.Fatal("expected success result")
	}
}

func TestPricesLatestTool_NullPrices(t *testing.T) {
	fake := &fakePricesClient{
		getLatestFn: func(_ context.Context, _ []int) (*prices.LatestResponse, error) {
			return &prices.LatestResponse{
				Data: map[string]prices.LatestPrice{
					"4151": {High: nil, Low: nil},
				},
			}, nil
		},
	}
	tool := NewPricesLatestTool(fake)
	req := makePricesLatestReq(map[string]any{"id": 4151})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.IsError {
		t.Fatal("expected success result even with null prices")
	}
}

func TestPricesLatestTool_ItemNotInResponse(t *testing.T) {
	fake := &fakePricesClient{
		getLatestFn: func(_ context.Context, _ []int) (*prices.LatestResponse, error) {
			return &prices.LatestResponse{Data: map[string]prices.LatestPrice{}}, nil
		},
	}
	tool := NewPricesLatestTool(fake)
	req := makePricesLatestReq(map[string]any{"id": 9999})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected error result for missing item")
	}
}

func TestPricesLatestTool_ClientError(t *testing.T) {
	fake := &fakePricesClient{
		getLatestFn: func(_ context.Context, _ []int) (*prices.LatestResponse, error) {
			return nil, errFixed
		},
	}
	tool := NewPricesLatestTool(fake)
	req := makePricesLatestReq(map[string]any{"id": 4151})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected error result for client error")
	}
}

func TestPricesLatestTool_MissingID(t *testing.T) {
	tool := NewPricesLatestTool(nil)
	req := makePricesLatestReq(map[string]any{})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected error result for missing id")
	}
}
