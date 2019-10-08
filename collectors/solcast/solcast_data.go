package solcast

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

//Data contains the (parsed) data received from the Solcast API
type Data struct {
	Forecast []Forecast `json:"forecast"`
}

//Forecast is a single forecasted record
type Forecast struct {
	Estimate     float64       `json:"pv_estimate"`
	Estimate10   float64       `json:"pv_estimate10"`
	Estimate90   float64       `json:"pv_estimate90"`
	PeriodEnd    time.Time     `json:"period_end"`
	Period       string        `json:"period"`
	PeriodParsed time.Duration `json:"-"`
}

//NewData creates a SolcastData and parses the data from `data`
func NewData(data io.Reader) (*Data, error) {
	var solcastData Data
	decoder := json.NewDecoder(data)
	return &solcastData, decoder.Decode(&solcastData)
}

//ToInflux converts the data to data understood by Influx
func (sd *Data) ToInflux() *bytes.Buffer {
	var buf bytes.Buffer

	for i := range sd.Forecast {
		fmt.Fprintf(&buf, "current_watt,source=solcast value=%v %v\n",
			sd.Forecast[i].Estimate,
			sd.Forecast[i].PeriodEnd.UnixNano())
	}

	return &buf
}
