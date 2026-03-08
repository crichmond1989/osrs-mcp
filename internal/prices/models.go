package prices

// MappingItem represents one tradeable item from the /mapping endpoint.
type MappingItem struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Examine  string `json:"examine"`
	Members  bool   `json:"members"`
	LowAlch  int    `json:"lowalch"`
	HighAlch int    `json:"highalch"`
	Limit    int    `json:"limit"`
	Value    int    `json:"value"`
	Icon     string `json:"icon"`
}

// LatestPrice holds the current instant-buy and instant-sell prices for one item.
type LatestPrice struct {
	High     *int  `json:"high"`
	HighTime int64 `json:"highTime"`
	Low      *int  `json:"low"`
	LowTime  int64 `json:"lowTime"`
}

// LatestResponse is the envelope returned by the /latest endpoint.
type LatestResponse struct {
	Data map[string]LatestPrice `json:"data"`
}

// TimeSeriesPrice holds averaged price and volume data over a time window.
// AvgHighPrice and AvgLowPrice are pointers because the API returns null when
// no trades occurred in the window.
type TimeSeriesPrice struct {
	AvgHighPrice    *int `json:"avgHighPrice"`
	HighPriceVolume int  `json:"highPriceVolume"`
	AvgLowPrice     *int `json:"avgLowPrice"`
	LowPriceVolume  int  `json:"lowPriceVolume"`
}

// TimeSeriesResponse is the shared envelope for /5m, /1h, and /24h endpoints.
type TimeSeriesResponse struct {
	Data map[string]TimeSeriesPrice `json:"data"`
}
