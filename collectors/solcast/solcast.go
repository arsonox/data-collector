package solcast

import (
	"log"
	"net/http"

	"github.com/recadra/data-collector/storage"
)

//Solcast is a connection to the Solcast API
type Solcast struct {
	addr    string
	apiKey  string
	c       uint64
	cl      *http.Client
	runChan chan *storage.Influx
}

//NewSolcast creates a new Solcast to connect to the Solcast API
func NewSolcast(addr, apiKey string) *Solcast {
	solcast := &Solcast{
		addr:    addr,
		apiKey:  apiKey,
		c:       3 * 60 * 60,
		cl:      &http.Client{},
		runChan: make(chan *storage.Influx, 30),
	}
	go solcast.run()
	return solcast
}

//Run gets the information from the Solcast API, parses it and
//submits it to InfluxDB
func (s *Solcast) Run(influx *storage.Influx) {
	s.runChan <- influx
}

//run is a parallel runner thread for the Solcast data, because Solcast data
//takes some time to receive and parse
func (s *Solcast) run() {
	for {
		influx := <-s.runChan
		if s.c++; s.c < (3 * 60 * 60) {
			//We run every 3 hours. Run is called every second by the
			//main loop so count until 3 hours passed.
			continue
		}

		req, err := http.NewRequest("GET", s.addr, nil)
		if err != nil {
			log.Printf("Solcast Request Error: %s\n", err.Error())
			continue
		}
		req.Header.Add("Authorization", "Bearer "+s.apiKey)

		resp, err := s.cl.Do(req)
		if err != nil {
			log.Printf("Solcast Get Error: %s\n", err.Error())
			continue
		}
		data, err := NewData(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Printf("Solcast NewData Error: %v\n", err.Error())
			log.Printf("Got status: %v\n", resp.Status)
			continue
		}
		err = influx.Send(data.ToInflux())
		if err != nil {
			log.Printf("Solcast ToInflux Error: %v\n", err.Error())
			continue
		}
		//No error occurred, reset counter
		s.c = 0
		log.Printf("Received Solcast data\n")
	}
}
