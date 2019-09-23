package main

import (
	"github.com/recadra/data-collector/collectors/sunspec"
	"github.com/recadra/data-collector/collectors/youless"
	"github.com/recadra/data-collector/storage"
)

func main() {
	var app App
	app.influx = storage.NewDefaultInflux()
	app.ss = sunspec.NewSunSpec("192.168.192.38:1502")
	app.yl = youless.NewYouless("192.168.192.22")

	app.Run()
}
