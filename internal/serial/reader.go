package serial

import (
	"bufio"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

type Handler interface {
	Handle(reading string)
}

type HandlerFunc func(reading string)

func (f HandlerFunc) Handle(reading string) {
	f(reading)
}

type Consumer struct {
	done    chan bool
	port    int
	handler Handler
}

func NewConsumer(port int, handler Handler) *Consumer {
	return &Consumer{port: port, handler: handler}
}

func (c *Consumer) Start() error {
	c.done = make(chan bool)
	eof := false

	go func() {
		for !eof {
			port, err := c.getPort()

			if err != nil {
				log.Error(err)
				time.Sleep(time.Second)
				continue
			}

			startTime := time.Now()

			scanner := bufio.NewScanner(port)

			for !eof {
				if !scanner.Scan() {
					if err := scanner.Err(); err != nil {
						log.Errorf("scanner err: %v", err)
						time.Sleep(time.Second)
						break
					}
				}

				reading := scanner.Text()

				if time.Now().Before(startTime.Add(time.Second)) {
					log.Debugf("skipping %v", reading)
					continue
				}

				c.handler.Handle(reading)
			}
		}

		log.Info("EOF")
	}()

	go func() {
		<-c.done
		eof = true
	}()

	return nil
}

func (c *Consumer) Stop() {
	c.done <- true
}

func (c *Consumer) getPort() (*serial.Port, error) {
	portName := fmt.Sprintf("COM%d", c.port)

	cfg := &serial.Config{Name: portName, Baud: 9600}

	port, err := serial.OpenPort(cfg)

	if err != nil {
		log.Errorf("could not connect to port %v: %w", portName, err)
	}

	log.Infof("using port %s", portName)
	return port, nil
}
