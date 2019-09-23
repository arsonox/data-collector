package sunspec

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"power/storage"
	"strings"
	"time"

	"github.com/goburrow/modbus"
)

var buf bytes.Buffer

type SunSpec struct {
	handler *modbus.TCPClientHandler
	ipAddr  string
	error   bool
}

type SunSpecTable101 struct {
	Power       float64
	DCPower     float64
	Counter     float64
	Temperature float64
	Efficiency  float64
}

func (SunSpec) scale(val float64, scale int16) float64 {
	return val * math.Pow(10, float64(scale))
}

func (s *SunSpec) init() {
	s.handler = modbus.NewTCPClientHandler(s.ipAddr)
	s.handler.Timeout = 10 * time.Second
	s.handler.SlaveId = 1
}

func NewSunSpec(ipAddr string) *SunSpec {
	se := &SunSpec{
		ipAddr: ipAddr,
	}
	se.init()
	return se
}

func (s *SunSpec) Run(influx *storage.Influx) {
	client := modbus.NewClient(s.handler)
	res, err := client.ReadHoldingRegisters(40069, 52)
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		if strings.Contains(err.Error(), "broken pipe") {
			if s.error == true {
				s.init()
			} else {
				s.error = true
				s.handler.Close()
			}
		}
		return
	}
	s.error = false

	data := SunSpecTable101{
		Power:       s.scale(float64(binary.BigEndian.Uint16(res[28:])), int16(binary.BigEndian.Uint16(res[30:]))),
		DCPower:     s.scale(float64(binary.BigEndian.Uint16(res[62:])), int16(binary.BigEndian.Uint16(res[64:]))),
		Counter:     s.scale(float64(binary.BigEndian.Uint32(res[48:])), int16(binary.BigEndian.Uint16(res[52:]))),
		Temperature: s.scale(float64(binary.BigEndian.Uint16(res[68:])), int16(binary.BigEndian.Uint16(res[74:]))),
	}
	if data.DCPower > 0 {
		data.Efficiency = data.Power / data.DCPower
	}
	fmt.Fprintf(&buf, "current_watt,source=inverter value=%v\n"+
		"total_kwh,source=inverter value=%v\n"+
		"current_dcin,source=inverter value=%v\n"+
		"current_temperature,source=inverter value=%v\n"+
		"current_efficiency,source=inverter value=%f",
		data.Power, data.Counter, data.DCPower, data.Temperature, data.Efficiency)
	err = influx.Send(&buf)
	if err != nil {
		log.Printf("Error: %v\n", err.Error())
	}
	buf.Reset()
}
