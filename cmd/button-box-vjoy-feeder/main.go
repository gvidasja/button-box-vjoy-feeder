package main

import (
	"io"
	"os"

	"github.com/gvidasja/button-box-vjoy-feeder/internal/buttons"
	"github.com/gvidasja/windowsservice"
	log "github.com/sirupsen/logrus"
)

func main() {
	file, err := os.OpenFile(`C:\logs\button-box-vjoy-feeder.log`, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)

	if err != nil {
		log.Error("could not open file for logging")
	}

	log.SetOutput(io.MultiWriter(os.Stdout, file))

	service := windowsservice.New("button-box-vjoy-feeded", "button-box-vjoy-feeded", buttons.Service)
	service.Run()
}
