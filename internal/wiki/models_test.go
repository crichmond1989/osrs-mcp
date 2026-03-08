package wiki

import (
	"encoding/json"
	"testing"
)

func TestOpenSearchResponse_UnmarshalJSON_Success(t *testing.T) {
	data := `["abyssal whip",["Abyssal whip","Abyssal whip (or)"],["",""],["https://a.example/","https://b.example/"]]`
	var r OpenSearchResponse
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Query != "abyssal whip" {
		t.Errorf("Query = %q, want %q", r.Query, "abyssal whip")
	}
	if len(r.Titles) != 2 {
		t.Errorf("Titles len = %d, want 2", len(r.Titles))
	}
	if r.Titles[0] != "Abyssal whip" {
		t.Errorf("Titles[0] = %q, want %q", r.Titles[0], "Abyssal whip")
	}
	if r.URLs[1] != "https://b.example/" {
		t.Errorf("URLs[1] = %q, want %q", r.URLs[1], "https://b.example/")
	}
}

func TestOpenSearchResponse_UnmarshalJSON_NotArray(t *testing.T) {
	var r OpenSearchResponse
	if err := json.Unmarshal([]byte(`{"key":"val"}`), &r); err == nil {
		t.Fatal("expected error for non-array input")
	}
}

func TestOpenSearchResponse_UnmarshalJSON_WrongLength(t *testing.T) {
	var r OpenSearchResponse
	if err := json.Unmarshal([]byte(`["only","three"]`), &r); err == nil {
		t.Fatal("expected error for array with wrong length")
	}
}

func TestOpenSearchResponse_UnmarshalJSON_BadQuery(t *testing.T) {
	var r OpenSearchResponse
	if err := json.Unmarshal([]byte(`[123,["t1"],["d1"],["u1"]]`), &r); err == nil {
		t.Fatal("expected error for non-string query")
	}
}

func TestOpenSearchResponse_UnmarshalJSON_BadTitles(t *testing.T) {
	var r OpenSearchResponse
	if err := json.Unmarshal([]byte(`["q","not-array",["d"],["u"]]`), &r); err == nil {
		t.Fatal("expected error for non-array titles")
	}
}

func TestOpenSearchResponse_UnmarshalJSON_BadDescriptions(t *testing.T) {
	var r OpenSearchResponse
	if err := json.Unmarshal([]byte(`["q",["t"],"not-array",["u"]]`), &r); err == nil {
		t.Fatal("expected error for non-array descriptions")
	}
}

func TestOpenSearchResponse_UnmarshalJSON_BadURLs(t *testing.T) {
	var r OpenSearchResponse
	if err := json.Unmarshal([]byte(`["q",["t"],["d"],"not-array"]`), &r); err == nil {
		t.Fatal("expected error for non-array urls")
	}
}

func TestUnexpectedFormatError(t *testing.T) {
	err := &UnexpectedFormatError{msg: "test error"}
	if err.Error() != "test error" {
		t.Errorf("Error() = %q, want %q", err.Error(), "test error")
	}
}

func TestSearchResult_JSON(t *testing.T) {
	data := `{"ns":0,"title":"Abyssal whip","pageid":123,"size":4567,"wordcount":300,"snippet":"A powerful whip","timestamp":"2024-01-01T00:00:00Z"}`
	var r SearchResult
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Title != "Abyssal whip" {
		t.Errorf("Title = %q, want %q", r.Title, "Abyssal whip")
	}
	if r.PageID != 123 {
		t.Errorf("PageID = %d, want 123", r.PageID)
	}
	if r.WordCount != 300 {
		t.Errorf("WordCount = %d, want 300", r.WordCount)
	}
}

func TestSearchResponse_JSON(t *testing.T) {
	data := `{"query":{"searchinfo":{"totalhits":1},"search":[{"ns":0,"title":"Abyssal whip","pageid":123,"size":0,"wordcount":0,"snippet":"","timestamp":""}]}}`
	var r SearchResponse
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Query.SearchInfo.TotalHits != 1 {
		t.Errorf("TotalHits = %d, want 1", r.Query.SearchInfo.TotalHits)
	}
	if len(r.Query.Search) != 1 {
		t.Fatalf("Search len = %d, want 1", len(r.Query.Search))
	}
	if r.Query.Search[0].Title != "Abyssal whip" {
		t.Errorf("Search[0].Title = %q", r.Query.Search[0].Title)
	}
}

func TestPageQueryResponse_JSON(t *testing.T) {
	data := `{"query":{"pages":{"123":{"pageid":123,"ns":0,"title":"Abyssal whip","revisions":[{"slots":{"main":{"*":"{{ItemBox}}\nsome wikitext"}}}]}}}}`
	var r PageQueryResponse
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	page, ok := r.Query.Pages["123"]
	if !ok {
		t.Fatal("expected page with id 123")
	}
	if page.Title != "Abyssal whip" {
		t.Errorf("Title = %q", page.Title)
	}
	if len(page.Revisions) != 1 {
		t.Fatalf("Revisions len = %d, want 1", len(page.Revisions))
	}
	if page.Revisions[0].Slots.Main.Content != "{{ItemBox}}\nsome wikitext" {
		t.Errorf("Content = %q", page.Revisions[0].Slots.Main.Content)
	}
}

func TestPageQueryResponse_MissingPage(t *testing.T) {
	data := `{"query":{"pages":{"-1":{"pageid":-1,"ns":0,"title":"Missing page","revisions":null}}}}`
	var r PageQueryResponse
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := r.Query.Pages["-1"]; !ok {
		t.Fatal("expected -1 key for missing page")
	}
}
