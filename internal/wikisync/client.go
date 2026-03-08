package wikisync

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const defaultBaseURL = "https://sync.runescape.wiki"

// ErrWikiSyncNotEnabled is returned when the player has not enabled the WikiSync plugin in RuneLite.
var ErrWikiSyncNotEnabled = errors.New("wikisync: player has not enabled WikiSync")

// Client is the testable interface for WikiSync API operations.
type Client interface {
	GetPlayerData(ctx context.Context, player string) (*WikiSyncResponse, error)
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

// GetPlayerData fetches WikiSync data for the given player.
// Returns ErrWikiSyncNotEnabled if the player has not enabled the WikiSync plugin.
func (c *RealClient) GetPlayerData(ctx context.Context, player string) (*WikiSyncResponse, error) {
	rawURL := c.baseURL + "/runelite/player/" + url.PathEscape(player) + "/STANDARD"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, fmt.Errorf("wikisync: build request: %w", err)
	}
	req.Header.Set("User-Agent", "osrs-mcp/1.0 (github.com/crichmond1989/osrs-mcp)")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("wikisync: http request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusNotFound {
		return nil, ErrWikiSyncNotEnabled
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wikisync: unexpected status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("wikisync: read body: %w", err)
	}

	var result WikiSyncResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("wikisync: decode response: %w", err)
	}
	return &result, nil
}
