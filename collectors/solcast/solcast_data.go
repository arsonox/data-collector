package solcast

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type SolcastData struct {
	Forecast []SolcastForecast `json:"forecast"`
}

type SolcastForecast struct {
	Estimate     float64       `json:"pv_estimate"`
	Estimate10   float64       `json:"pv_estimate10"`
	Estimate90   float64       `json:"pv_estimate90"`
	PeriodEnd    time.Time     `json:"period_end"`
	Period       string        `json:"period"`
	PeriodParsed time.Duration `json:"-"`
}

func NewSolcastData(data io.Reader) (*SolcastData, error) {
	var solcastData SolcastData
	decoder := json.NewDecoder(data)
	return &solcastData, decoder.Decode(&solcastData)
}

func (sd *SolcastData) ToInflux() *bytes.Buffer {
	var buf bytes.Buffer

	for i := range sd.Forecast {
		fmt.Fprintf(&buf, "current_watt,source=solcast value=%v %v\n",
			sd.Forecast[i].Estimate,
			sd.Forecast[i].PeriodEnd.UnixNano())
	}

	return &buf
}
