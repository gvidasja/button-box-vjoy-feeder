package main

import (
	"bufio"
	"log"
	"strconv"

	"github.com/tarm/serial"
)

type buttonReading struct {
	buttonID uint
	state    bool
	err      error
}

func readButtonsSerialPort(readings chan<- buttonReading) {
	c := &serial.Config{Name: "COM15", Baud: 9600}

	port, err := serial.OpenPort(c)

	if err != nil {
		readings <- buttonReading{err: err}
	}

	scanner := bufio.NewScanner(port)

	for {
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				log.Fatal(err)
				readings <- buttonReading{err: err}
				return
			}
		}
		serialString := scanner.Text()
		actionNumber, _ := strconv.ParseInt(serialString[0:1], 10, 64)
		button, _ := strconv.ParseInt(serialString[1:], 10, 64)

		readings <- buttonReading{buttonID: uint(button), state: actionNumber > 0}
	}
}

func readButtons(stoppedEvent chan<- bool, stopCommand <-chan bool) {
	device := NewDeviceWithDelay(1, 10)
	err := device.Init()

	if err != nil {
		log.Fatal(err)
		stoppedEvent <- true
		return
	}

	readings := make(chan buttonReading)

	go readButtonsSerialPort(readings)

loop:
	for {
		select {
		case <-stopCommand:
			break loop
		case reading := <-readings:
			if reading.err != nil {
				break loop
			}
			device.SetButton(reading.buttonID, reading.state)
		}
	}

	stoppedEvent <- true
}
