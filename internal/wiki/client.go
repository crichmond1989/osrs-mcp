package wiki

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

const defaultBaseURL = "https://oldschool.runescape.wiki/api.php"

// ErrPageNotFound is returned when the wiki reports no page exists for a title.
var ErrPageNotFound = errors.New("wiki: page not found")

// Client is the testable interface for all OSRS wiki operations.
type Client interface {
	OpenSearch(ctx context.Context, query string, limit int) (*OpenSearchResponse, error)
	SearchPages(ctx context.Context, query string, limit int) (*SearchResponse, error)
	GetPage(ctx context.Context, title string) (*PageQueryResponse, error)
}

// HTTPClient is the interface satisfied by *http.Client.
// It is the seam that allows httptest.Server injection in tests.
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

func (c *RealClient) get(ctx context.Context, params url.Values) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("wiki: build request: %w", err)
	}
	req.Header.Set("User-Agent", "osrs-mcp/1.0 (github.com/crichmond1989/osrs-mcp)")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("wiki: http request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wiki: unexpected status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("wiki: read body: %w", err)
	}
	return body, nil
}

// OpenSearch calls action=opensearch and returns matching titles and URLs.
func (c *RealClient) OpenSearch(ctx context.Context, query string, limit int) (*OpenSearchResponse, error) {
	params := url.Values{
		"action": {"opensearch"},
		"search": {query},
		"limit":  {strconv.Itoa(limit)},
		"format": {"json"},
	}
	body, err := c.get(ctx, params)
	if err != nil {
		return nil, err
	}
	var result OpenSearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("wiki: decode opensearch response: %w", err)
	}
	return &result, nil
}

// SearchPages calls action=query&list=search and returns richer search results.
func (c *RealClient) SearchPages(ctx context.Context, query string, limit int) (*SearchResponse, error) {
	params := url.Values{
		"action":   {"query"},
		"list":     {"search"},
		"srsearch": {query},
		"srlimit":  {strconv.Itoa(limit)},
		"format":   {"json"},
	}
	body, err := c.get(ctx, params)
	if err != nil {
		return nil, err
	}
	var result SearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("wiki: decode search response: %w", err)
	}
	return &result, nil
}

// GetPage calls action=query&prop=revisions to retrieve the wikitext of a page.
// Returns ErrPageNotFound if the wiki reports no page for the given title.
func (c *RealClient) GetPage(ctx context.Context, title string) (*PageQueryResponse, error) {
	params := url.Values{
		"action":  {"query"},
		"prop":    {"revisions"},
		"rvprop":  {"content"},
		"rvslots": {"main"},
		"titles":  {title},
		"format":  {"json"},
	}
	body, err := c.get(ctx, params)
	if err != nil {
		return nil, err
	}
	var result PageQueryResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("wiki: decode page response: %w", err)
	}
	if _, missing := result.Query.Pages["-1"]; missing {
		return nil, ErrPageNotFound
	}
	return &result, nil
}
