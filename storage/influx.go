package storage

import (
	"io"
	"net/http"
)

type Influx struct {
	ipPort string
	db     string
}

func NewDefaultInflux() *Influx {
	return NewInflux("localhost:8086", "power")
}

func NewInflux(ipPort, db string) *Influx {
	return &Influx{
		ipPort: ipPort,
		db:     db,
	}
}

func (i *Influx) Send(buf io.Reader) (err error) {
	_, err = http.Post("http://"+i.ipPort+"/write?db="+i.db, "", buf)
	return
}
