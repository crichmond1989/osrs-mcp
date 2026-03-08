package tools

import (
	"context"
	"errors"

	"github.com/crichmond1989/osrs-mcp/internal/hiscores"
	"github.com/crichmond1989/osrs-mcp/internal/prices"
	"github.com/crichmond1989/osrs-mcp/internal/wiki"
	"github.com/crichmond1989/osrs-mcp/internal/wikisync"
)

// errFixed is a sentinel error for tests.
var errFixed = errors.New("injected error")

// fakeWikiClient implements wiki.Client for unit tests with no HTTP calls.
type fakeWikiClient struct {
	openSearchFn  func(ctx context.Context, query string, limit int) (*wiki.OpenSearchResponse, error)
	searchPagesFn func(ctx context.Context, query string, limit int) (*wiki.SearchResponse, error)
	getPageFn     func(ctx context.Context, title string) (*wiki.PageQueryResponse, error)
}

func (f *fakeWikiClient) OpenSearch(ctx context.Context, q string, l int) (*wiki.OpenSearchResponse, error) {
	return f.openSearchFn(ctx, q, l)
}

func (f *fakeWikiClient) SearchPages(ctx context.Context, q string, l int) (*wiki.SearchResponse, error) {
	return f.searchPagesFn(ctx, q, l)
}

func (f *fakeWikiClient) GetPage(ctx context.Context, title string) (*wiki.PageQueryResponse, error) {
	return f.getPageFn(ctx, title)
}

// fakePricesClient implements prices.Client for unit tests with no HTTP calls.
type fakePricesClient struct {
	getLatestFn     func(ctx context.Context, ids []int) (*prices.LatestResponse, error)
	getMappingFn    func(ctx context.Context) ([]prices.MappingItem, error)
	getTimeSeriesFn func(ctx context.Context, window string) (*prices.TimeSeriesResponse, error)
}

func (f *fakePricesClient) GetLatest(ctx context.Context, ids []int) (*prices.LatestResponse, error) {
	return f.getLatestFn(ctx, ids)
}

func (f *fakePricesClient) GetMapping(ctx context.Context) ([]prices.MappingItem, error) {
	return f.getMappingFn(ctx)
}

func (f *fakePricesClient) GetTimeSeries(ctx context.Context, window string) (*prices.TimeSeriesResponse, error) {
	return f.getTimeSeriesFn(ctx, window)
}

// fakeHiscoresClient implements hiscores.Client for unit tests with no HTTP calls.
type fakeHiscoresClient struct {
	getStatsFn func(ctx context.Context, player string, mode string) (*hiscores.HiscoresResponse, error)
}

func (f *fakeHiscoresClient) GetStats(ctx context.Context, player, mode string) (*hiscores.HiscoresResponse, error) {
	return f.getStatsFn(ctx, player, mode)
}

// fakeWikiSyncClient implements wikisync.Client for unit tests with no HTTP calls.
type fakeWikiSyncClient struct {
	getPlayerDataFn func(ctx context.Context, player string) (*wikisync.WikiSyncResponse, error)
}

func (f *fakeWikiSyncClient) GetPlayerData(ctx context.Context, player string) (*wikisync.WikiSyncResponse, error) {
	return f.getPlayerDataFn(ctx, player)
}
