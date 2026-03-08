package wiki

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
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

// newTestClient creates a RealClient pointed at the given httptest server.
func newTestClient(ts *httptest.Server) Client {
	return NewClientWithBase(ts.Client(), ts.URL)
}

// ---- OpenSearch ----

func TestOpenSearch_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("action") != "opensearch" {
			t.Errorf("action = %q, want opensearch", r.URL.Query().Get("action"))
		}
		if got := r.Header.Get("User-Agent"); got == "" {
			t.Error("User-Agent header not set")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`["abyssal whip",["Abyssal whip"],[""],["https://wiki.example/"]]`))
	}))
	defer ts.Close()

	c := newTestClient(ts)
	res, err := c.OpenSearch(context.Background(), "abyssal whip", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Query != "abyssal whip" {
		t.Errorf("Query = %q", res.Query)
	}
	if len(res.Titles) != 1 || res.Titles[0] != "Abyssal whip" {
		t.Errorf("Titles = %v", res.Titles)
	}
}

func TestOpenSearch_Non200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.OpenSearch(context.Background(), "q", 5)
	if err == nil {
		t.Fatal("expected error for non-200 status")
	}
}

func TestOpenSearch_BadJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{bad json`))
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.OpenSearch(context.Background(), "q", 5)
	if err == nil {
		t.Fatal("expected error for bad JSON")
	}
}

func TestOpenSearch_CancelledContext(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	c := newTestClient(ts)
	_, err := c.OpenSearch(ctx, "q", 5)
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
}

// ---- SearchPages ----

func TestSearchPages_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("list") != "search" {
			t.Errorf("list param = %q, want search", r.URL.Query().Get("list"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"query":{"searchinfo":{"totalhits":1},"search":[{"ns":0,"title":"Abyssal whip","pageid":123,"size":0,"wordcount":0,"snippet":"A whip","timestamp":""}]}}`))
	}))
	defer ts.Close()

	c := newTestClient(ts)
	res, err := c.SearchPages(context.Background(), "abyssal whip", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Query.SearchInfo.TotalHits != 1 {
		t.Errorf("TotalHits = %d", res.Query.SearchInfo.TotalHits)
	}
	if len(res.Query.Search) != 1 {
		t.Fatalf("Search len = %d", len(res.Query.Search))
	}
}

func TestSearchPages_Non200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.SearchPages(context.Background(), "q", 5)
	if err == nil {
		t.Fatal("expected error for non-200 status")
	}
}

func TestSearchPages_BadJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`not json`))
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.SearchPages(context.Background(), "q", 5)
	if err == nil {
		t.Fatal("expected error for bad JSON")
	}
}

func TestSearchPages_CancelledContext(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{}`))
	}))
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	c := newTestClient(ts)
	_, err := c.SearchPages(ctx, "q", 5)
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
}

// ---- GetPage ----

func TestGetPage_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("titles") != "Abyssal whip" {
			t.Errorf("titles param = %q", r.URL.Query().Get("titles"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"query":{"pages":{"123":{"pageid":123,"ns":0,"title":"Abyssal whip","revisions":[{"slots":{"main":{"*":"wikitext here"}}}]}}}}`))
	}))
	defer ts.Close()

	c := newTestClient(ts)
	res, err := c.GetPage(context.Background(), "Abyssal whip")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	page, ok := res.Query.Pages["123"]
	if !ok {
		t.Fatal("expected page 123")
	}
	if page.Revisions[0].Slots.Main.Content != "wikitext here" {
		t.Errorf("Content = %q", page.Revisions[0].Slots.Main.Content)
	}
}

func TestGetPage_NotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"query":{"pages":{"-1":{"pageid":-1,"ns":0,"title":"Missing","revisions":null}}}}`))
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.GetPage(context.Background(), "Missing")
	if err != ErrPageNotFound {
		t.Fatalf("expected ErrPageNotFound, got: %v", err)
	}
}

func TestGetPage_Non200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.GetPage(context.Background(), "Whatever")
	if err == nil {
		t.Fatal("expected error for non-200 status")
	}
}

func TestGetPage_BadJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{bad`))
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.GetPage(context.Background(), "Whatever")
	if err == nil {
		t.Fatal("expected error for bad JSON")
	}
}

func TestGetPage_CancelledContext(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{}`))
	}))
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	c := newTestClient(ts)
	_, err := c.GetPage(ctx, "Whatever")
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
}

func TestNewClient(t *testing.T) {
	c := NewClient()
	if c == nil {
		t.Fatal("NewClient returned nil")
	}
}

func TestGet_InvalidBaseURL(t *testing.T) {
	c := NewClientWithBase(&http.Client{}, "://invalid")
	_, err := c.OpenSearch(context.Background(), "q", 5)
	if err == nil {
		t.Fatal("expected error for invalid base URL")
	}
}

func TestGet_ReadBodyError(t *testing.T) {
	c := NewClientWithBase(&errBodyHTTPClient{}, "http://unused")
	_, err := c.OpenSearch(context.Background(), "q", 5)
	if err == nil {
		t.Fatal("expected error when body read fails")
	}
}

func TestSearchPages_ReadBodyError(t *testing.T) {
	c := NewClientWithBase(&errBodyHTTPClient{}, "http://unused")
	_, err := c.SearchPages(context.Background(), "q", 5)
	if err == nil {
		t.Fatal("expected error when body read fails")
	}
}

func TestGetPage_ReadBodyError(t *testing.T) {
	c := NewClientWithBase(&errBodyHTTPClient{}, "http://unused")
	_, err := c.GetPage(context.Background(), "Whatever")
	if err == nil {
		t.Fatal("expected error when body read fails")
	}
}
