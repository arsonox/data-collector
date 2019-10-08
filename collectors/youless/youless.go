package youless

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/recadra/data-collector/storage"
)

var buf bytes.Buffer

//Youless is a connection to a Youless device
type Youless struct {
	data Data
	addr string
}

//Data is part of the data received from the Youless device
type Data struct {
	Counter float64 `json:"net"`
	Power   int64   `json:"pwr"`
}

//Parse parses part of the data received by the Youless device and
//saves it into YoulessData
func (yd *Data) Parse(b []byte) (err error) {
	s := string(b)
	if i := strings.Index(s, `"pwr":`); i > 0 {
		sep := strings.IndexRune(s[i+6:], ',')
		val := strings.Trim(s[i+6:sep+i+6], " ")
		yd.Power, err = strconv.ParseInt(val, 10, 64)
		if err != nil {
			return fmt.Errorf("YoulessData(Parse): %v", err.Error())
		}
	} else {
		return fmt.Errorf("YoulessData(Parse): could not find pwr")
	}

	if i := strings.Index(s, `"net":`); i > 0 {
		sep := strings.IndexRune(s[i+6:], ',')
		val := strings.Trim(s[i+6:sep+i+6], " ")
		yd.Counter, err = strconv.ParseFloat(val, 64)
		if err != nil {
			return fmt.Errorf("YoulessData(Parse): %v", err.Error())
		}
	} else {
		return fmt.Errorf("YoulessData(Parse): could not find net")
	}

	return nil
}

//NewYouless creates a new Youless connection to the specified ipAddr
func NewYouless(ipAddr string) *Youless {
	return &Youless{
		addr: "http://" + ipAddr + "/e?f=j",
	}
}

//Run receives data from the Youless device, parses the information and
//submits it to influx
func (y *Youless) Run(influx *storage.Influx) {
	resp, err := http.Get(y.addr)
	if err != nil {
		log.Printf("Error: %s\n", err.Error())
		return
	}

	/*
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&y.data)
		if err != nil {
			log.Printf("Error: %s\n", err.Error())
			return
		}
	*/

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error: %v\n", err.Error())
		return
	}

	err = y.data.Parse(body)

	resp.Body.Close()
	if err == nil {
		fmt.Fprintf(&buf, "current_watt,source=youless value=%v\n"+
			"total_kwh,source=youless value=%v",
			y.data.Power, y.data.Counter)
		if err = influx.Send(&buf); err != nil {
			log.Printf("Error: %s\n", err.Error())
		}
		buf.Reset()
	}
}
