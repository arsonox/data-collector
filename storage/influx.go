package storage

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

//Influx is a connection to an Influx HTTP(S) server
type Influx struct {
	ipPort string
	db     string
	dummy  bool
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

//SetDummy sets whether or not this Influx should be treated as a dummy for testing
//purposes. A dummy Influx writes to file (influx.txt)
func (i *Influx) SetDummy(dummy bool) {
	i.dummy = dummy
}

//Send sends data to the default database
func (i *Influx) Send(buf io.Reader) error {
	return i.SendDB(buf, i.db)
}

//SendDB sends data to the specified database (db)
func (i *Influx) SendDB(buf io.Reader, db string) error {
	if i.dummy {
		f, err := os.OpenFile("influx.txt",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
			return err
		}
		fmt.Fprintf(f, "db: %v\n", db)
		io.Copy(f, buf)
		f.Close()
		return nil
	}
	_, err := http.Post("http://"+i.ipPort+"/write?db="+db, "", buf)
	return err
}
