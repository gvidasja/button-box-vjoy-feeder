package buttons

import (
	"bufio"
	"strconv"
	"time"

	"github.com/gvidasja/button-box-vjoy-feeder/internal/device"
	"github.com/gvidasja/button-box-vjoy-feeder/internal/vjoy"
	"github.com/gvidasja/windowsservice"
	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

func Service(s windowsservice.Service) {
	if s.IsInServiceMode() {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.DebugLevel)
	}

	vjoyDevice := vjoy.NewDevice(1)
	err := vjoyDevice.Init()

	if err != nil {
		log.Panic("could not connect to vjoy device:", err)
		s.StopService()
		return
	}

	device := device.New(vjoyDevice, device.DeviceConfig{
		MinimumButtonPressDuration: time.Millisecond * 20,
	})

	readings := readButtonsSerialPort()
loop:
	for {
		select {
		case <-s.StopCommandReceived():
			break loop
		case reading := <-readings:
			log.Debugf("button %v: %v", reading.buttonID, reading.state)

			buttonID := reading.getButtonID()

			if deviceButtonID, ok := keyMap[buttonID]; ok {
				log.Debugf("sending %v -> %v", buttonID, deviceButtonID)
				device.SetButton(deviceButtonID, reading.state)
			}
		}
	}

	vjoyDevice.Dispose()

	s.StopService()
}

func readButtonsSerialPort() <-chan buttonReading {
	readings := make(chan buttonReading)

	go func() {
		for {
			scanner, err := getScanner()

			if err != nil {
				log.Error("could not connect to port COM15:", err)
				time.Sleep(time.Second)
				continue
			}

			for {
				if !scanner.Scan() {
					if err := scanner.Err(); err != nil {
						log.Error(err)
						time.Sleep(time.Second)
						break
					}
				}

				serialString := scanner.Text()
				actionNumber, _ := strconv.ParseInt(serialString[0:1], 10, 64)
				button, _ := strconv.ParseInt(serialString[1:], 10, 64)

				readings <- buttonReading{buttonID: buttonID(button), state: actionNumber > 0}
			}
		}
	}()

	return readings
}

func getScanner() (*bufio.Scanner, error) {
	c := &serial.Config{Name: "COM15", Baud: 9600}

	port, err := serial.OpenPort(c)

	if err != nil {
		log.Error("could not connect to port COM15: ", err)
		return nil, err
	}

	return bufio.NewScanner(port), nil
}
