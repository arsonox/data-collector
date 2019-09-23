package main

import (
	"power/collectors/sunspec"
	"power/collectors/youless"
	"power/storage"
	"time"
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
