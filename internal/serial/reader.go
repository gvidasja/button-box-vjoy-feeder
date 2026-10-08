package serial

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"sync/atomic"
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
	stopped atomic.Bool
	port    int
	handler Handler
}

func NewConsumer(port int, handler Handler) *Consumer {
	return &Consumer{port: port, handler: handler}
}

func (c *Consumer) Start() error {
	c.stopped.Store(false)
	go c.run()
	return nil
}

func (c *Consumer) Stop() {
	c.stopped.Store(true)
}

// run keeps the port connected until Stop, reopening it after unplug/replug.
func (c *Consumer) run() {
	portName := fmt.Sprintf("COM%d", c.port)
	lastErr := ""

	for !c.stopped.Load() {
		port, err := serial.OpenPort(&serial.Config{Name: portName, Baud: 9600, ReadTimeout: time.Second})
		if err != nil {
			// log only on change, so an unplugged device doesn't spam the log every second
			if err.Error() != lastErr {
				slog.Error("could not connect to port", "port", portName, "err", err)
				lastErr = err.Error()
			}
			time.Sleep(time.Second)
			continue
		}

		slog.Info("using port", "port", portName)
		lastErr = ""

		c.consume(port)
		// must close: Windows COM ports are exclusive, a leaked handle makes every reopen fail with "Access is denied"
		port.Close()

		if !c.stopped.Load() {
			slog.Info("serial disconnected, reconnecting", "port", portName)
			time.Sleep(time.Second)
		}
	}

	slog.Info("EOF", "port", portName)
}

func (c *Consumer) consume(port io.Reader) {
	startTime := time.Now()
	scanner := bufio.NewScanner(&idleReader{r: port, stopped: &c.stopped})

	for scanner.Scan() {
		reading := scanner.Text()

		if time.Since(startTime) < time.Second {
			slog.Debug("skipping", "reading", reading)
			continue
		}

		c.handler.Handle(reading)
	}

	if err := scanner.Err(); err != nil {
		slog.Error("scanner err", "port", c.port, "err", err)
	}
}

// idleReader retries the empty reads a read timeout produces while the device is idle;
// bufio.Scanner would otherwise give up with io.ErrNoProgress after 100 of them.
// An empty read that returns well before the timeout means the port is dead, not idle.
type idleReader struct {
	r       io.Reader
	stopped *atomic.Bool
}

func (r *idleReader) Read(p []byte) (int, error) {
	for {
		start := time.Now()
		n, err := r.r.Read(p)
		if n > 0 || err != nil {
			return n, err
		}
		if r.stopped.Load() {
			return 0, io.EOF
		}
		if time.Since(start) < 500*time.Millisecond {
			return 0, io.ErrUnexpectedEOF
		}
	}
}
