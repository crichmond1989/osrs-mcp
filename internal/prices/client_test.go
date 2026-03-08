package prices

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

func newTestClient(ts *httptest.Server) Client {
	return NewClientWithBase(ts.Client(), ts.URL)
}

// ---- GetLatest ----

func TestGetLatest_Success_AllItems(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/latest" {
			t.Errorf("path = %q, want /latest", r.URL.Path)
		}
		if got := r.Header.Get("User-Agent"); got == "" {
			t.Error("User-Agent header not set")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"4151":{"high":2500000,"highTime":1700000000,"low":2400000,"lowTime":1699999000}}}`))
	}))
	defer ts.Close()

	c := newTestClient(ts)
	res, err := c.GetLatest(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Data) != 1 {
		t.Errorf("Data len = %d, want 1", len(res.Data))
	}
}

func TestGetLatest_Success_SpecificIDs(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.RawQuery
		if q == "" {
			t.Error("expected id query params")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{}}`))
	}))
	defer ts.Close()

	c := newTestClient(ts)
	res, err := c.GetLatest(context.Background(), []int{4151, 4153})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestGetLatest_Non200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.GetLatest(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for non-200 status")
	}
}

func TestGetLatest_BadJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{bad`))
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.GetLatest(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for bad JSON")
	}
}

func TestGetLatest_CancelledContext(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{}`))
	}))
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	c := newTestClient(ts)
	_, err := c.GetLatest(ctx, nil)
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
}

// ---- GetMapping ----

func TestGetMapping_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/mapping" {
			t.Errorf("path = %q, want /mapping", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"id":4151,"name":"Abyssal whip","examine":"A weapon.","members":true,"lowalch":72000,"highalch":108000,"limit":70,"value":120001,"icon":"icon.png"}]`))
	}))
	defer ts.Close()

	c := newTestClient(ts)
	items, err := c.GetMapping(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("items len = %d, want 1", len(items))
	}
	if items[0].Name != "Abyssal whip" {
		t.Errorf("Name = %q", items[0].Name)
	}
}

func TestGetMapping_Non200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.GetMapping(context.Background())
	if err == nil {
		t.Fatal("expected error for non-200 status")
	}
}

func TestGetMapping_BadJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`not json`))
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.GetMapping(context.Background())
	if err == nil {
		t.Fatal("expected error for bad JSON")
	}
}

func TestGetMapping_CancelledContext(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	c := newTestClient(ts)
	_, err := c.GetMapping(ctx)
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
}

// ---- GetTimeSeries ----

func testTimeSeriesSuccess(t *testing.T, window string) {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/"+window {
			t.Errorf("path = %q, want /%s", r.URL.Path, window)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"4151":{"avgHighPrice":2500000,"highPriceVolume":50,"avgLowPrice":null,"lowPriceVolume":0}}}`))
	}))
	defer ts.Close()

	c := newTestClient(ts)
	res, err := c.GetTimeSeries(context.Background(), window)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestGetTimeSeries_5m(t *testing.T)  { testTimeSeriesSuccess(t, "5m") }
func TestGetTimeSeries_1h(t *testing.T)  { testTimeSeriesSuccess(t, "1h") }
func TestGetTimeSeries_24h(t *testing.T) { testTimeSeriesSuccess(t, "24h") }

func TestGetTimeSeries_InvalidWindow(t *testing.T) {
	c := NewClientWithBase(&http.Client{}, "http://unused")
	_, err := c.GetTimeSeries(context.Background(), "7d")
	if err != ErrInvalidWindow {
		t.Fatalf("expected ErrInvalidWindow, got: %v", err)
	}
}

func TestGetTimeSeries_Non200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.GetTimeSeries(context.Background(), "1h")
	if err == nil {
		t.Fatal("expected error for non-200 status")
	}
}

func TestGetTimeSeries_BadJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{bad`))
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.GetTimeSeries(context.Background(), "5m")
	if err == nil {
		t.Fatal("expected error for bad JSON")
	}
}

func TestGetTimeSeries_CancelledContext(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{}`))
	}))
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	c := newTestClient(ts)
	_, err := c.GetTimeSeries(ctx, "1h")
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
	_, err := c.GetLatest(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for invalid base URL")
	}
}

func TestGetLatest_ReadBodyError(t *testing.T) {
	c := NewClientWithBase(&errBodyHTTPClient{}, "http://unused")
	_, err := c.GetLatest(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error when body read fails")
	}
}

func TestGetMapping_ReadBodyError(t *testing.T) {
	c := NewClientWithBase(&errBodyHTTPClient{}, "http://unused")
	_, err := c.GetMapping(context.Background())
	if err == nil {
		t.Fatal("expected error when body read fails")
	}
}

func TestGetTimeSeries_ReadBodyError(t *testing.T) {
	c := NewClientWithBase(&errBodyHTTPClient{}, "http://unused")
	_, err := c.GetTimeSeries(context.Background(), "1h")
	if err == nil {
		t.Fatal("expected error when body read fails")
	}
}
