package buttons

import (
	"bufio"
	"strconv"
	"time"

	"github.com/gvidasja/button-box-vjoy-feeder/internal/device"
	"github.com/gvidasja/button-box-vjoy-feeder/internal/log"
	"github.com/tarm/serial"
)

type buttonReading struct {
	buttonID uint
	state    bool
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

				readings <- buttonReading{buttonID: uint(button), state: actionNumber > 0}
			}
		}
	}()

	return readings
}

func Service(stoppedEvent chan<- bool, stopCommand <-chan bool) {
	device := device.NewDeviceWithDelay(1, 10)
	err := device.Init()

	if err != nil {
		log.Panic("could not connect to vjoy device:", err)
		stoppedEvent <- true
		return
	}

	readings := readButtonsSerialPort()

loop:
	for {
		select {
		case <-stopCommand:
			break loop
		case reading := <-readings:
			log.Debug("button pressed", reading.buttonID, reading.state)
			device.SetButton(reading.buttonID, reading.state)
		}
	}

	stoppedEvent <- true
}
