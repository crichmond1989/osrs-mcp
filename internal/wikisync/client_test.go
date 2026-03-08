package wikisync

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

func TestGetPlayerData_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/runelite/player/") {
			t.Errorf("path = %q, want prefix /runelite/player/", r.URL.Path)
		}
		if !strings.HasSuffix(r.URL.Path, "/STANDARD") {
			t.Errorf("path = %q, want suffix /STANDARD", r.URL.Path)
		}
		if got := r.Header.Get("User-Agent"); got == "" {
			t.Error("User-Agent header not set")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"quests":{"0":2,"1":1,"2":0}}`))
	}))
	defer ts.Close()

	c := newTestClient(ts)
	res, err := c.GetPlayerData(context.Background(), "Zezima")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Quests) != 3 {
		t.Fatalf("Quests len = %d, want 3", len(res.Quests))
	}
	if res.Quests["0"] != QuestComplete {
		t.Errorf("Quests[0] = %d, want %d", res.Quests["0"], QuestComplete)
	}
}

func TestGetPlayerData_WikiSyncNotEnabled_400(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.GetPlayerData(context.Background(), "NoWikiSync")
	if err != ErrWikiSyncNotEnabled {
		t.Fatalf("expected ErrWikiSyncNotEnabled, got: %v", err)
	}
}

func TestGetPlayerData_WikiSyncNotEnabled_404(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.GetPlayerData(context.Background(), "NoWikiSync")
	if err != ErrWikiSyncNotEnabled {
		t.Fatalf("expected ErrWikiSyncNotEnabled, got: %v", err)
	}
}

func TestGetPlayerData_Non200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.GetPlayerData(context.Background(), "Zezima")
	if err == nil {
		t.Fatal("expected error for non-200 status")
	}
}

func TestGetPlayerData_BadJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{bad`))
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.GetPlayerData(context.Background(), "Zezima")
	if err == nil {
		t.Fatal("expected error for bad JSON")
	}
}

func TestGetPlayerData_CancelledContext(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"quests":{}}`))
	}))
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	c := newTestClient(ts)
	_, err := c.GetPlayerData(ctx, "Zezima")
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
}

func TestGetPlayerData_ReadBodyError(t *testing.T) {
	c := NewClientWithBase(&errBodyHTTPClient{}, "http://unused")
	_, err := c.GetPlayerData(context.Background(), "Zezima")
	if err == nil {
		t.Fatal("expected error when body read fails")
	}
}

func TestGetPlayerData_InvalidBaseURL(t *testing.T) {
	c := NewClientWithBase(&http.Client{}, "://invalid")
	_, err := c.GetPlayerData(context.Background(), "Zezima")
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
