package prices

import (
	"encoding/json"
	"testing"
)

func TestMappingItem_JSON(t *testing.T) {
	data := `{"id":4151,"name":"Abyssal whip","examine":"A weapon from the Abyss.","members":true,"lowalch":72000,"highalch":108000,"limit":70,"value":120001,"icon":"Abyssal_whip.png"}`
	var item MappingItem
	if err := json.Unmarshal([]byte(data), &item); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.ID != 4151 {
		t.Errorf("ID = %d, want 4151", item.ID)
	}
	if item.Name != "Abyssal whip" {
		t.Errorf("Name = %q", item.Name)
	}
	if !item.Members {
		t.Error("Members should be true")
	}
	if item.Limit != 70 {
		t.Errorf("Limit = %d, want 70", item.Limit)
	}
}

func TestLatestPrice_JSON_WithValues(t *testing.T) {
	high := 2500000
	low := 2400000
	data := `{"high":2500000,"highTime":1700000000,"low":2400000,"lowTime":1699999000}`
	var p LatestPrice
	if err := json.Unmarshal([]byte(data), &p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.High == nil || *p.High != high {
		t.Errorf("High = %v, want %d", p.High, high)
	}
	if p.Low == nil || *p.Low != low {
		t.Errorf("Low = %v, want %d", p.Low, low)
	}
	if p.HighTime != 1700000000 {
		t.Errorf("HighTime = %d", p.HighTime)
	}
}

func TestLatestPrice_JSON_NullPrices(t *testing.T) {
	data := `{"high":null,"highTime":0,"low":null,"lowTime":0}`
	var p LatestPrice
	if err := json.Unmarshal([]byte(data), &p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.High != nil {
		t.Errorf("High should be nil, got %v", p.High)
	}
	if p.Low != nil {
		t.Errorf("Low should be nil, got %v", p.Low)
	}
}

func TestLatestResponse_JSON(t *testing.T) {
	data := `{"data":{"4151":{"high":2500000,"highTime":1700000000,"low":2400000,"lowTime":1699999000}}}`
	var r LatestResponse
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p, ok := r.Data["4151"]
	if !ok {
		t.Fatal("expected entry for item 4151")
	}
	if p.High == nil || *p.High != 2500000 {
		t.Errorf("High = %v", p.High)
	}
}

func TestTimeSeriesPrice_JSON_WithValues(t *testing.T) {
	avgHigh := 2500000
	avgLow := 2400000
	data := `{"avgHighPrice":2500000,"highPriceVolume":150,"avgLowPrice":2400000,"lowPriceVolume":200}`
	var p TimeSeriesPrice
	if err := json.Unmarshal([]byte(data), &p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.AvgHighPrice == nil || *p.AvgHighPrice != avgHigh {
		t.Errorf("AvgHighPrice = %v, want %d", p.AvgHighPrice, avgHigh)
	}
	if p.AvgLowPrice == nil || *p.AvgLowPrice != avgLow {
		t.Errorf("AvgLowPrice = %v", p.AvgLowPrice)
	}
	if p.HighPriceVolume != 150 {
		t.Errorf("HighPriceVolume = %d", p.HighPriceVolume)
	}
}

func TestTimeSeriesPrice_JSON_NullPrices(t *testing.T) {
	data := `{"avgHighPrice":null,"highPriceVolume":0,"avgLowPrice":null,"lowPriceVolume":0}`
	var p TimeSeriesPrice
	if err := json.Unmarshal([]byte(data), &p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.AvgHighPrice != nil {
		t.Errorf("AvgHighPrice should be nil, got %v", p.AvgHighPrice)
	}
	if p.AvgLowPrice != nil {
		t.Errorf("AvgLowPrice should be nil, got %v", p.AvgLowPrice)
	}
}

func TestTimeSeriesResponse_JSON(t *testing.T) {
	avgHigh := 2500000
	data := `{"data":{"4151":{"avgHighPrice":2500000,"highPriceVolume":150,"avgLowPrice":null,"lowPriceVolume":0}}}`
	var r TimeSeriesResponse
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p, ok := r.Data["4151"]
	if !ok {
		t.Fatal("expected entry for item 4151")
	}
	if p.AvgHighPrice == nil || *p.AvgHighPrice != avgHigh {
		t.Errorf("AvgHighPrice = %v", p.AvgHighPrice)
	}
	if p.AvgLowPrice != nil {
		t.Error("AvgLowPrice should be nil")
	}
}
