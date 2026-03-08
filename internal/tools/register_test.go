package tools

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/server"

	"github.com/crichmond1989/osrs-mcp/internal/hiscores"
	"github.com/crichmond1989/osrs-mcp/internal/prices"
	"github.com/crichmond1989/osrs-mcp/internal/wiki"
	"github.com/crichmond1989/osrs-mcp/internal/wikisync"
)

// noopWikiClient satisfies wiki.Client with no-op implementations.
type noopWikiClient struct{}

func (n *noopWikiClient) OpenSearch(_ context.Context, _ string, _ int) (*wiki.OpenSearchResponse, error) {
	return &wiki.OpenSearchResponse{}, nil
}
func (n *noopWikiClient) SearchPages(_ context.Context, _ string, _ int) (*wiki.SearchResponse, error) {
	return &wiki.SearchResponse{}, nil
}
func (n *noopWikiClient) GetPage(_ context.Context, _ string) (*wiki.PageQueryResponse, error) {
	return nil, wiki.ErrPageNotFound
}

// noopPricesClient satisfies prices.Client with no-op implementations.
type noopPricesClient struct{}

func (n *noopPricesClient) GetLatest(_ context.Context, _ []int) (*prices.LatestResponse, error) {
	return &prices.LatestResponse{}, nil
}
func (n *noopPricesClient) GetMapping(_ context.Context) ([]prices.MappingItem, error) {
	return nil, nil
}
func (n *noopPricesClient) GetTimeSeries(_ context.Context, _ string) (*prices.TimeSeriesResponse, error) {
	return &prices.TimeSeriesResponse{}, nil
}

// noopHiscoresClient satisfies hiscores.Client with a no-op implementation.
type noopHiscoresClient struct{}

func (n *noopHiscoresClient) GetStats(_ context.Context, _ string, _ string) (*hiscores.HiscoresResponse, error) {
	return &hiscores.HiscoresResponse{}, nil
}

// noopWikiSyncClient satisfies wikisync.Client with a no-op implementation.
type noopWikiSyncClient struct{}

func (n *noopWikiSyncClient) GetPlayerData(_ context.Context, _ string) (*wikisync.WikiSyncResponse, error) {
	return &wikisync.WikiSyncResponse{}, nil
}

func TestRegisterAll(t *testing.T) {
	s := server.NewMCPServer("test", "0.0.1")
	RegisterAll(s, &noopWikiClient{}, &noopPricesClient{}, &noopHiscoresClient{}, &noopWikiSyncClient{})
	// If RegisterAll panics, the test will fail. Reaching here means success.
}
