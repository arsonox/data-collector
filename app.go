package main

import (
	"time"

	"github.com/recadra/data-collector/storage"
)

type Runner interface {
	Run(*storage.Influx)
}

type App struct {
	influx  *storage.Influx
	runners []Runner
}

func (a *App) Run() {
	for {
		now := time.Now()

		for _, runner := range a.runners {
			runner.Run(a.influx)
		}

		time.Sleep(time.Second - time.Now().Sub(now))
	}
}
