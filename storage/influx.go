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
	user   string
	pass   string
	cl     http.Client
}

//NewDefaultInflux creates a new Influx instance with default settings
//It assumes ipPort to be "localhost:8086" and db to be "power"
func NewDefaultInflux() *Influx {
	return NewInflux("localhost:8086", "power", "", "")
}

//NewInflux creates a new Influx.
//  ipPort: the server to connect to (format ip:port)
//  db: the name of the database to use as default database (default "power")
//  user: username to use (leave empty for no auth)
//  pass: password to use (leave empty for no auth)
func NewInflux(ipPort, db, user, pass string) *Influx {
	if db == "" {
		db = "power"
	}
	return &Influx{
		ipPort: ipPort,
		db:     db,
		user:   user,
		pass:   pass,
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

	rq, err := http.NewRequest("POST", "http://"+i.ipPort+"/write?db="+db, buf)
	if err != nil {
		return err
	}

	rq.Header.Set("Content-Type", "application/octet-stream")

	if i.user != "" {
		rq.SetBasicAuth(i.user, i.pass)
	}

	//_, err = i.cl.Post("http://"+i.ipPort+"/write?db="+db, "", buf)
	_, err = i.cl.Do(rq)
	return err
}
