package hiscores

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// errorReader is an io.Reader that always returns an error.
type errorReader struct{}

func (errorReader) Read([]byte) (int, error) { return 0, errors.New("read error") }

// errBodyHTTPClient returns a 200 response whose body always errors on read.
type errBodyHTTPClient struct{}

func (e *errBodyHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(errorReader{}),
	}, nil
}

func newTestClient(ts *httptest.Server) Client {
	return NewClientWithBase(ts.Client(), ts.URL)
}

func TestGetStats_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "hiscore_oldschool") {
			t.Errorf("path = %q, want path containing hiscore_oldschool", r.URL.Path)
		}
		if r.URL.Query().Get("player") != "Zezima" {
			t.Errorf("player = %q, want Zezima", r.URL.Query().Get("player"))
		}
		if got := r.Header.Get("User-Agent"); got == "" {
			t.Error("User-Agent header not set")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"skills":[{"id":0,"name":"Overall","rank":1,"level":2277,"xp":4600000000},{"id":1,"name":"Attack","rank":1234,"level":99,"xp":200000000}]}`))
	}))
	defer ts.Close()

	c := newTestClient(ts)
	res, err := c.GetStats(context.Background(), "Zezima", "standard")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skills) != 2 {
		t.Fatalf("Skills len = %d, want 2", len(res.Skills))
	}
	if res.Skills[0].Name != "Overall" {
		t.Errorf("Skills[0].Name = %q, want Overall", res.Skills[0].Name)
	}
	if res.Skills[0].XP != 4600000000 {
		t.Errorf("Skills[0].XP = %d, want 4600000000", res.Skills[0].XP)
	}
}

func TestGetStats_PlayerNotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.GetStats(context.Background(), "NoSuchPlayer123", "standard")
	if err != ErrPlayerNotFound {
		t.Fatalf("expected ErrPlayerNotFound, got: %v", err)
	}
}

func TestGetStats_Non200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.GetStats(context.Background(), "Zezima", "standard")
	if err == nil {
		t.Fatal("expected error for non-200 status")
	}
}

func TestGetStats_BadJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{bad`))
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.GetStats(context.Background(), "Zezima", "standard")
	if err == nil {
		t.Fatal("expected error for bad JSON")
	}
}

func TestGetStats_InvalidMode(t *testing.T) {
	c := NewClientWithBase(&http.Client{}, "http://unused")
	_, err := c.GetStats(context.Background(), "Zezima", "leagues")
	if err != ErrInvalidMode {
		t.Fatalf("expected ErrInvalidMode, got: %v", err)
	}
}

func TestGetStats_CancelledContext(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"skills":[]}`))
	}))
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	c := newTestClient(ts)
	_, err := c.GetStats(ctx, "Zezima", "standard")
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
}

func TestGetStats_ReadBodyError(t *testing.T) {
	c := NewClientWithBase(&errBodyHTTPClient{}, "http://unused")
	_, err := c.GetStats(context.Background(), "Zezima", "standard")
	if err == nil {
		t.Fatal("expected error when body read fails")
	}
}

func TestGetStats_InvalidBaseURL(t *testing.T) {
	c := NewClientWithBase(&http.Client{}, "://invalid")
	_, err := c.GetStats(context.Background(), "Zezima", "standard")
	if err == nil {
		t.Fatal("expected error for invalid base URL")
	}
}

func TestNewClient(t *testing.T) {
	c := NewClient()
	if c == nil {
		t.Fatal("NewClient returned nil")
	}
}

func testGetStatsMode(t *testing.T, mode, wantSegment string) {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, wantSegment) {
			t.Errorf("path = %q, want segment %q", r.URL.Path, wantSegment)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"skills":[]}`))
	}))
	defer ts.Close()

	c := newTestClient(ts)
	res, err := c.GetStats(context.Background(), "player", mode)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestGetStats_Ironman(t *testing.T) { testGetStatsMode(t, "ironman", "hiscore_oldschool_ironman") }
func TestGetStats_Hardcore(t *testing.T) {
	testGetStatsMode(t, "hardcore", "hiscore_oldschool_hardcore_ironman")
}
func TestGetStats_Ultimate(t *testing.T) {
	testGetStatsMode(t, "ultimate", "hiscore_oldschool_ultimate")
}
