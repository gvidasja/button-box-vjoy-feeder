package serial

import (
	"bufio"
	"fmt"
	"log/slog"
	"time"

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
				slog.Error("serial error", "err", err)
				time.Sleep(time.Second)
				continue
			}

			startTime := time.Now()

			scanner := bufio.NewScanner(port)

			for !eof {
				if !scanner.Scan() {
					if err := scanner.Err(); err != nil {
						slog.Error("scanner err", "err", err)
					}
					slog.Debug("serial disconnected, reconnecting...")
					time.Sleep(time.Second)
					break
				}

				reading := scanner.Text()

				if time.Now().Before(startTime.Add(time.Second)) {
					slog.Debug("skipping", "reading", reading)
					continue
				}

				c.handler.Handle(reading)
			}
		}

		slog.Info("EOF")
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

	cfg := &serial.Config{Name: portName, Baud: 9600, ReadTimeout: time.Second}

	port, err := serial.OpenPort(cfg)

	if err != nil {
		slog.Error("could not connect to port", "port", portName, "err", err)
		return nil, err
	}

	slog.Info("using port", "port", portName)
	return port, nil
}
