package main

import (
	"io"
	"os"
	"time"

	"github.com/gvidasja/button-box-vjoy-feeder/internal/app"
	"github.com/gvidasja/button-box-vjoy-feeder/internal/appender"
	"github.com/gvidasja/button-box-vjoy-feeder/internal/buttonbox"
	"github.com/gvidasja/button-box-vjoy-feeder/internal/device"
	"github.com/gvidasja/button-box-vjoy-feeder/internal/serial"
	"github.com/gvidasja/button-box-vjoy-feeder/internal/vjoy"
	log "github.com/sirupsen/logrus"
)

func main() {
	file := appender.ForFile(`E:\dev\button-box-vjoy-feeder\button-box-vjoy-feeder.log`)

	log.SetOutput(io.MultiWriter(os.Stdout, file))
	log.SetLevel(log.DebugLevel)

	vjoyDevice := vjoy.NewDevice(1)

	buttonBoxHandler := buttonbox.NewHandler(device.New(vjoyDevice, device.DeviceConfig{
		MinimumButtonPressDuration: time.Millisecond * 20,
	}))

	app.
		New("button-box-vjoy-feeded").
		AddWorkers(vjoyDevice, serial.NewConsumer([]int{3, 15}, buttonBoxHandler)).
		Run()
}
