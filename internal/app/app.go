package app

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/gvidasja/button-box-vjoy-feeder/internal/windowsservice"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/windows/svc"
)

type App struct {
	name        string
	description string
	services    []Worker
}

func New(name string) *App {
	return &App{
		name:        name,
		description: name,
	}
}

func (c *App) AddWorkers(services ...Worker) *App {
	c.services = append(c.services, services...)

	return c
}

func (s *App) Run() {
	log.Info("running")

	isService, err := svc.IsWindowsService()

	if err != nil {
		log.Panicf("could not check if in service mode: %w", err)
	}

	if isService {
		s.runWindowsService()
	} else {
		s.runStandalone()
	}
}

func (s *App) runWindowsService() {
	err := svc.Run(s.name, windowsservice.NewHandler(s))

	if err != nil {
		log.Panicf("could not run service: %w", err)
	}
}

func (s *App) runStandalone() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	err := s.Start()

	if err != nil {
		log.Panicf("could not run service: %w", err)
	}

	<-signals
	s.Stop()
}

func (c *App) Start() error {
	log.Debug("starting...")
	for _, service := range c.services {
		err := service.Start()

		if err != nil {
			return fmt.Errorf("could not run service: %w", err)
		}
	}
	log.Debug("started")

	return nil
}

func (c *App) Stop() {
	log.Debug("stopping...")
	for i := len(c.services) - 1; i >= 0; i-- {
		c.services[i].Stop()
	}
	log.Debug("stopped")
}
