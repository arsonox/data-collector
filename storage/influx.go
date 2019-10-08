package storage

import (
	"io"
	"net/http"
)

//Influx is a connection to an Influx HTTP(S) server
type Influx struct {
	ipPort string
	db     string
}

//NewDefaultInflux creates a new Influx instance with default settings
//It assumes ipPort to be "localhost:8086" and db to be "power"
func NewDefaultInflux() *Influx {
	return NewInflux("localhost:8086", "power")
}

//NewInflux creates a new Influx.
//  ipPort: the server to connect to (format ip:port)
//  db: the name of the database to use as default database (default "power")
func NewInflux(ipPort, db string) *Influx {
	if db == "" {
		db = "power"
	}
	return &Influx{
		ipPort: ipPort,
		db:     db,
	}
}

//Send sends data to the default database
func (i *Influx) Send(buf io.Reader) (err error) {
	return i.SendDB(buf, i.db)
}

//SendDB sends data to the specified database (db)
func (i *Influx) SendDB(buf io.Reader, db string) (err error) {
	_, err = http.Post("http://"+i.ipPort+"/write?db="+db, "", buf)
	return
}
