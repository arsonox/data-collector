package solcast

import (
	"log"
	"net/http"

	"github.com/recadra/data-collector/storage"
)

//Solcast is a connection to the Solcast API
type Solcast struct {
	addr   string
	apiKey string
	c      uint64
}

//NewSolcast creates a new Solcast to connect to the Solcast API
func NewSolcast(addr, apiKey string) *Solcast {
	return &Solcast{
		addr:   addr,
		apiKey: apiKey,
		c:      3 * 60,
	}
}

//Run gets the information from the Solcast API, parses it and
//submits it to InfluxDB
func (s *Solcast) Run(influx *storage.Influx) {
	if s.c++; s.c < (3 * 60) {
		//We run every 3 hours. Run is called every second by the
		//main loop so count until 3 hours passed.
		return
	}

	req, err := http.NewRequest("GET", s.addr, nil)
	if err != nil {
		log.Printf("Solcast Request Error: %s\n", err.Error())
		return
	}
	req.Header.Add("Authorization", "Bearer "+s.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Solcast Get Error: %s\n", err.Error())
		return
	}
	data, err := NewData(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Printf("Solcast NewData Error: %v\n", err.Error())
		log.Printf("Got status: %v\n", resp.Status)
		return
	}
	err = influx.Send(data.ToInflux())
	if err != nil {
		log.Printf("Solcast ToInflux Error: %v\n", err.Error())
		return
	}
	//No error occurred, reset counter
	s.c = 0
	log.Printf("Received Solcast data\n")
}
