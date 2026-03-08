package prices

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const defaultBaseURL = "https://prices.runescape.wiki/api/v1/osrs"

// ErrInvalidWindow is returned when an unsupported time window is requested.
var ErrInvalidWindow = errors.New("prices: window must be one of: 5m, 1h, 24h")

// Client is the testable interface for OSRS prices API operations.
type Client interface {
	GetLatest(ctx context.Context, ids []int) (*LatestResponse, error)
	GetMapping(ctx context.Context) ([]MappingItem, error)
	GetTimeSeries(ctx context.Context, window string) (*TimeSeriesResponse, error)
}

// HTTPClient is the interface satisfied by *http.Client.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// RealClient implements Client using an injectable HTTP client.
type RealClient struct {
	http    HTTPClient
	baseURL string
}

// NewClient returns a production-ready Client.
func NewClient() Client {
	return &RealClient{
		http:    &http.Client{},
		baseURL: defaultBaseURL,
	}
}

// NewClientWithBase creates a RealClient with a custom base URL and HTTP client.
// Used in tests to point at an httptest.Server.
func NewClientWithBase(httpClient HTTPClient, base string) Client {
	return &RealClient{http: httpClient, baseURL: base}
}

func (c *RealClient) get(ctx context.Context, path string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("prices: build request: %w", err)
	}
	req.Header.Set("User-Agent", "osrs-mcp/1.0 (github.com/crichmond1989/osrs-mcp)")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("prices: http request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("prices: unexpected status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("prices: read body: %w", err)
	}
	return body, nil
}

// GetLatest returns the current instant-buy and instant-sell prices.
// When ids is non-empty, fetches only those items (via repeated id= params).
// When ids is empty, fetches all items.
func (c *RealClient) GetLatest(ctx context.Context, ids []int) (*LatestResponse, error) {
	path := "/latest"
	if len(ids) > 0 {
		parts := make([]string, len(ids))
		for i, id := range ids {
			parts[i] = "id=" + strconv.Itoa(id)
		}
		path += "?" + strings.Join(parts, "&")
	}

	body, err := c.get(ctx, path)
	if err != nil {
		return nil, err
	}
	var result LatestResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("prices: decode latest response: %w", err)
	}
	return &result, nil
}

// GetMapping returns metadata for all tradeable items.
func (c *RealClient) GetMapping(ctx context.Context) ([]MappingItem, error) {
	body, err := c.get(ctx, "/mapping")
	if err != nil {
		return nil, err
	}
	var result []MappingItem
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("prices: decode mapping response: %w", err)
	}
	return result, nil
}

// GetTimeSeries returns averaged price and volume data for the given window.
// window must be one of: "5m", "1h", "24h".
func (c *RealClient) GetTimeSeries(ctx context.Context, window string) (*TimeSeriesResponse, error) {
	switch window {
	case "5m", "1h", "24h":
	default:
		return nil, ErrInvalidWindow
	}

	body, err := c.get(ctx, "/"+window)
	if err != nil {
		return nil, err
	}
	var result TimeSeriesResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("prices: decode timeseries response: %w", err)
	}
	return &result, nil
}
