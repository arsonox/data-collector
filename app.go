package main

import (
	"time"

	"github.com/recadra/data-collector/collectors/sunspec"
	"github.com/recadra/data-collector/collectors/youless"
	"github.com/recadra/data-collector/storage"
)

type App struct {
	influx *storage.Influx

	ss *sunspec.SunSpec
	yl *youless.Youless
}

func (a *App) Run() {
	for {
		now := time.Now()
		a.yl.Run(a.influx)
		a.ss.Run(a.influx)
		time.Sleep(time.Second - time.Now().Sub(now))
	}
}
