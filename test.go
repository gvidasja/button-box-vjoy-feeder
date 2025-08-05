package main

import (
	"log/slog"

	"github.com/gvidasja/button-box-vjoy-feeder/internal/app"
	"github.com/gvidasja/button-box-vjoy-feeder/internal/serial"
)

func main() {
	reader := serial.NewConsumer(4, serial.HandlerFunc(func(reading string) {
		slog.Info("read", "msg", reading)
	}))

	app.New("com port reader").
		AddWorkers(reader).
		Run()
}
