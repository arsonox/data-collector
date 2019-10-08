package main

import (
	"log"

	"github.com/recadra/data-collector/collectors/solcast"
	"github.com/recadra/data-collector/collectors/sunspec"
	"github.com/recadra/data-collector/collectors/youless"
	"github.com/recadra/data-collector/config"
	"github.com/recadra/data-collector/storage"
)

func main() {
	log.Print("Data-collector starting...")
	cfg := config.GetConfig()

	var app App

	if cfg.InfluxIP == "" {
		app.influx = storage.NewDefaultInflux()
	} else {
		app.influx = storage.NewInflux(cfg.InfluxIP, cfg.InfluxDB)
	}

	if cfg.InfluxDummy {
		app.influx.SetDummy(true)
	}

	app.runners = make([]Runner, 0)

	if cfg.UseSunspec {
		log.Print("Initializing Sunspec...")
		app.runners = append(app.runners, sunspec.NewSunSpec(cfg.SunspecIP))
	}

	if cfg.UseYouless {
		log.Print("Initializing Youless...")
		app.runners = append(app.runners, youless.NewYouless(cfg.YoulessIP))
	}

	if cfg.UseSolcast {
		log.Print("Initializing Solcast...")
		app.runners = append(app.runners, solcast.NewSolcast(cfg.SolcastURL))
	}

	if len(app.runners) == 0 {
		log.Fatal("No runners to run, exiting")
		return
	}

	log.Print("Starting main collection loop")
	app.Run()
}
