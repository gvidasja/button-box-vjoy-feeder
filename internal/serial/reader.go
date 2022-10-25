package serial

import (
	"bufio"
	"errors"
	"fmt"
	"time"

	"github.com/gvidasja/button-box-vjoy-feeder/internal/windowsservice"
	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

type Handler interface {
	Handle(reading string)
}

type Consumer struct {
	done    chan bool
	ports   []int
	handler Handler
}

var _ windowsservice.Service = (*Consumer)(nil)

func NewConsumer(ports []int, handler Handler) *Consumer {
	return &Consumer{ports: ports, handler: handler}
}

func (c *Consumer) Start() error {
	readings := make(chan string)

	go func() {
		for {
			scanner, err := c.getScanner()

			if err != nil {
				log.Error(err)
				time.Sleep(time.Second)
				continue
			}

			startTime := time.Now()

			for {
				if !scanner.Scan() {
					if err := scanner.Err(); err != nil {
						log.Error(err)
						time.Sleep(time.Second)
						break
					}
				}

				reading := scanner.Text()

				if time.Now().Before(startTime.Add(time.Second)) {
					log.Debugf("skipping %v", reading)
					continue
				}

				readings <- reading
			}
		}
	}()

	go (func() {
		for {
			select {
			case <-c.done:
				return
			case reading := <-readings:
				c.handler.Handle(reading)
			}
		}
	})()

	return nil
}

func (c *Consumer) Stop() {
	c.done <- true
}

func (c *Consumer) getScanner() (*bufio.Scanner, error) {
	for _, port := range c.ports {
		portName := fmt.Sprintf("COM%d", port)

		c := &serial.Config{Name: portName, Baud: 9600}

		port, err := serial.OpenPort(c)

		if err != nil {
			log.Errorf("could not connect to port %v: %w", portName, err)
		} else {
			return bufio.NewScanner(port), nil
		}
	}

	return nil, errors.New("no port found")
}
