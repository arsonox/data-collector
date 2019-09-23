package solcast

import (
	"log"
	"net/http"
	"power/storage"
)

type Solcast struct {
	addr string
}

func NewSolcast(addr string) *Solcast {
	return &Solcast{addr}
}

func (s *Solcast) Run(influx *storage.Influx) {
	resp, err := http.Get(s.addr)
	if err != nil {
		log.Printf("Error: %s\n", err.Error())
		return
	}
	data, err := NewSolcastData(resp.Body)
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
