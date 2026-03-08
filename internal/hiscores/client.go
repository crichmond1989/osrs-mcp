package hiscores

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const defaultBaseURL = "https://secure.runescape.com"

// ErrPlayerNotFound is returned when the hiscores API reports no player with the given name.
var ErrPlayerNotFound = errors.New("hiscores: player not found")

// ErrInvalidMode is returned when an unsupported game mode is requested.
var ErrInvalidMode = errors.New("hiscores: mode must be one of: standard, ironman, hardcore, ultimate")

// modeToEndpoint maps the mode parameter to the hiscores URL path segment.
var modeToEndpoint = map[string]string{
	"standard": "hiscore_oldschool",
	"ironman":  "hiscore_oldschool_ironman",
	"hardcore": "hiscore_oldschool_hardcore_ironman",
	"ultimate": "hiscore_oldschool_ultimate",
}

// Client is the testable interface for OSRS hiscores operations.
type Client interface {
	GetStats(ctx context.Context, player string, mode string) (*HiscoresResponse, error)
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

// GetStats fetches skill levels, XP, and ranks for the given player.
// mode must be one of: "standard", "ironman", "hardcore", "ultimate".
// Returns ErrPlayerNotFound if the player does not appear on the hiscores.
// Returns ErrInvalidMode if mode is not a recognised value.
func (c *RealClient) GetStats(ctx context.Context, player, mode string) (*HiscoresResponse, error) {
	endpoint, ok := modeToEndpoint[mode]
	if !ok {
		return nil, ErrInvalidMode
	}

	rawURL := c.baseURL + "/m=" + endpoint + "/index_lite.json?player=" + url.QueryEscape(player)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, fmt.Errorf("hiscores: build request: %w", err)
	}
	req.Header.Set("User-Agent", "osrs-mcp/1.0 (github.com/crich/osrs-mcp)")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("hiscores: http request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrPlayerNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("hiscores: unexpected status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("hiscores: read body: %w", err)
	}

	var result HiscoresResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("hiscores: decode response: %w", err)
	}
	return &result, nil
}
