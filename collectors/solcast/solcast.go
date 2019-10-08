package solcast

import (
	"log"
	"net/http"

	"github.com/recadra/data-collector/storage"
)

//Solcast is a connection to the Solcast API
type Solcast struct {
	addr string
}

//NewSolcast creates a new Solcast to connect to the Solcast API
func NewSolcast(addr string) *Solcast {
	return &Solcast{addr}
}

//Run gets the information from the Solcast API, parses it and
//submits it to InfluxDB
func (s *Solcast) Run(influx *storage.Influx) {
	resp, err := http.Get(s.addr)
	if err != nil {
		log.Printf("Error: %s\n", err.Error())
		return
	}
	data, err := NewData(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Printf("Error: %v\n", err.Error())
		return
	}
	err = influx.Send(data.ToInflux())
	if err != nil {
		log.Printf("Error: %v\n", err.Error())
	}
}
