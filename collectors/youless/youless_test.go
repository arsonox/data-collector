package youless

import (
	"encoding/json"
	"testing"
)

var YoulessDataParseTestData = []byte(`[{"tm":1568486967,"net": 117.444,"pwr": 595,"ts0":1483228800,"cs0": 0.000,"ps0": 0,"p1": 64.802,"p2": 74.612,"n1": 9.003,"n2": 12.967,"gas": 0.000,"gts":0}]`)

func TestDataParse(t *testing.T) {
	var data Data
	err := data.Parse(YoulessDataParseTestData)
	if err != nil {
		t.Logf("Error during parse: %v\n", err.Error())
		t.FailNow()
	}
	if data.Counter != 117.444 {
		t.Logf("counter: expected 117.444, got %v", data.Counter)
		t.Fail()
	}
	if data.Power != 595 {
		t.Logf("power: expected 595, got %v", data.Power)
		t.Fail()
	}
}

func BenchmarkDataParse(b *testing.B) {
	var data Data
	b.ResetTimer()
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		data.Parse(YoulessDataParseTestData)
	}
}

func BenchmarkJSONDataParse(b *testing.B) {
	var data []Data
	b.ResetTimer()
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		json.Unmarshal(YoulessDataParseTestData, &data)
	}
}
