package tools

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/crich/osrs-mcp/internal/prices"
)

func makePricesTimeSeriesReq(args map[string]any) mcp.CallToolRequest {
	req := mcp.CallToolRequest{}
	req.Params.Arguments = args
	return req
}

func TestPricesTimeSeriesTool_Definition(t *testing.T) {
	tool := NewPricesTimeSeriesTool(nil)
	def := tool.Definition()
	if def.Name != "prices_timeseries" {
		t.Errorf("Name = %q, want prices_timeseries", def.Name)
	}
}

func TestPricesTimeSeriesTool_Success_WithPrices(t *testing.T) {
	fake := &fakePricesClient{
		getTimeSeriesFn: func(_ context.Context, window string) (*prices.TimeSeriesResponse, error) {
			return &prices.TimeSeriesResponse{
				Data: map[string]prices.TimeSeriesPrice{
					"4151": {
						AvgHighPrice:    intPtr(2500000),
						HighPriceVolume: 50,
						AvgLowPrice:     intPtr(2400000),
						LowPriceVolume:  60,
					},
				},
			}, nil
		},
	}
	tool := NewPricesTimeSeriesTool(fake)
	req := makePricesTimeSeriesReq(map[string]any{"id": 4151, "window": "1h"})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.IsError {
		t.Fatal("expected success result")
	}
}

func TestPricesTimeSeriesTool_NullPrices(t *testing.T) {
	fake := &fakePricesClient{
		getTimeSeriesFn: func(_ context.Context, _ string) (*prices.TimeSeriesResponse, error) {
			return &prices.TimeSeriesResponse{
				Data: map[string]prices.TimeSeriesPrice{
					"4151": {AvgHighPrice: nil, AvgLowPrice: nil},
				},
			}, nil
		},
	}
	tool := NewPricesTimeSeriesTool(fake)
	req := makePricesTimeSeriesReq(map[string]any{"id": 4151, "window": "5m"})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.IsError {
		t.Fatal("expected success result even with null prices")
	}
}

func TestPricesTimeSeriesTool_ItemNotInResponse(t *testing.T) {
	fake := &fakePricesClient{
		getTimeSeriesFn: func(_ context.Context, _ string) (*prices.TimeSeriesResponse, error) {
			return &prices.TimeSeriesResponse{Data: map[string]prices.TimeSeriesPrice{}}, nil
		},
	}
	tool := NewPricesTimeSeriesTool(fake)
	req := makePricesTimeSeriesReq(map[string]any{"id": 9999, "window": "24h"})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected error result for missing item")
	}
}

func TestPricesTimeSeriesTool_ClientError(t *testing.T) {
	fake := &fakePricesClient{
		getTimeSeriesFn: func(_ context.Context, _ string) (*prices.TimeSeriesResponse, error) {
			return nil, errFixed
		},
	}
	tool := NewPricesTimeSeriesTool(fake)
	req := makePricesTimeSeriesReq(map[string]any{"id": 4151, "window": "1h"})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected error result for client error")
	}
}

func TestPricesTimeSeriesTool_MissingID(t *testing.T) {
	tool := NewPricesTimeSeriesTool(nil)
	req := makePricesTimeSeriesReq(map[string]any{"window": "1h"})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected error result for missing id")
	}
}

func TestPricesTimeSeriesTool_MissingWindow(t *testing.T) {
	tool := NewPricesTimeSeriesTool(nil)
	req := makePricesTimeSeriesReq(map[string]any{"id": 4151})
	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected error result for missing window")
	}
}
