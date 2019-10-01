package main

import (
	"github.com/recadra/data-collector/collectors/solcast"
	"github.com/recadra/data-collector/collectors/sunspec"
	"github.com/recadra/data-collector/collectors/youless"
	"github.com/recadra/data-collector/config"
	"github.com/recadra/data-collector/storage"
)

func main() {
	cfg := config.GetConfig()

	var app App

	app.influx = storage.NewDefaultInflux()
	app.runners = make([]Runner, 0)

	if cfg.UseSunspec {
		app.runners = append(app.runners, sunspec.NewSunSpec(cfg.SunspecIP))
	}

	if cfg.UseYouless {
		app.runners = append(app.runners, youless.NewYouless(cfg.YoulessIP))
	}

	if cfg.UseSolcast {
		app.runners = append(app.runners, solcast.NewSolcast(cfg.SolcastURL))
	}

	app.Run()
}
