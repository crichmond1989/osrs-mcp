package wiki

import "encoding/json"

// OpenSearchResponse holds results from action=opensearch.
// The API returns a positional JSON array: [query, [titles], [descriptions], [urls]].
type OpenSearchResponse struct {
	Query        string
	Titles       []string
	Descriptions []string
	URLs         []string
}

// UnmarshalJSON decodes the positional array format returned by the opensearch action.
func (o *OpenSearchResponse) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if len(raw) != 4 {
		return &UnexpectedFormatError{msg: "opensearch response must have 4 elements"}
	}
	if err := json.Unmarshal(raw[0], &o.Query); err != nil {
		return err
	}
	if err := json.Unmarshal(raw[1], &o.Titles); err != nil {
		return err
	}
	if err := json.Unmarshal(raw[2], &o.Descriptions); err != nil {
		return err
	}
	return json.Unmarshal(raw[3], &o.URLs)
}

// UnexpectedFormatError is returned when an API response has an unexpected structure.
type UnexpectedFormatError struct {
	msg string
}

func (e *UnexpectedFormatError) Error() string {
	return e.msg
}

// SearchResult is one item from action=query&list=search.
type SearchResult struct {
	NS        int    `json:"ns"`
	Title     string `json:"title"`
	PageID    int    `json:"pageid"`
	Size      int    `json:"size"`
	WordCount int    `json:"wordcount"`
	Snippet   string `json:"snippet"`
	Timestamp string `json:"timestamp"`
}

// SearchResponse wraps the query.search array from action=query&list=search.
type SearchResponse struct {
	Query struct {
		SearchInfo struct {
			TotalHits int `json:"totalhits"`
		} `json:"searchinfo"`
		Search []SearchResult `json:"search"`
	} `json:"query"`
}

// PageRevision holds one revision's wikitext content.
type PageRevision struct {
	Slots struct {
		Main struct {
			Content string `json:"*"`
		} `json:"main"`
	} `json:"slots"`
}

// Page is one entry in query.pages (keyed by page id as string).
type Page struct {
	PageID    int            `json:"pageid"`
	NS        int            `json:"ns"`
	Title     string         `json:"title"`
	Revisions []PageRevision `json:"revisions"`
}

// PageQueryResponse is the outer envelope for action=query&prop=revisions.
type PageQueryResponse struct {
	Query struct {
		Pages map[string]Page `json:"pages"`
	} `json:"query"`
}
