package windowsservice

import (
	"fmt"
	"os"
	"os/signal"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/windows/svc"
)

type Service interface {
	Start() error
	Stop()
}

type cli struct {
	name        string
	description string
	services    []Service
}

func New(name, description string) *cli {
	return &cli{
		name:        name,
		description: description,
	}
}

func (c *cli) AddService(services ...Service) *cli {
	c.services = append(c.services, services...)

	return c
}

func (s *cli) Run() {
	log.Info("running")
	isService, err := s.runWhenInServiceMode()

	if err != nil {
		log.Panic("could not run in service mode:", err)
	}

	if isService {
		return
	}

	if len(os.Args) < 2 {
		cmdDoc := func(cmd command, desc string) string {
			return fmt.Sprintf("  * %s - %s", cmd, desc)
		}

		commandDescriptions := []string{
			cmdDoc(commandInstall, "install service"),
			cmdDoc(commandUninstall, "uninstall service"),
			cmdDoc(commandStart, "start service"),
			cmdDoc(commandStop, "stop service"),
			cmdDoc(commandRun, "run as a simple executable"),
		}

		fmt.Printf("\nplease use one of these commands:\n\n%s\n\n", strings.Join(commandDescriptions, "\n"))

		return
	}

	cmd := command(strings.ToLower(os.Args[1]))

	switch cmd {
	case commandInstall:
		log.Info("installing...")
		err = s.install()
	case commandUninstall:
		log.Info("removing...")
		err = s.usingManager(s.removeService)
	case commandStart:
		log.Info("starting...")
		err = s.usingManager(s.startService)
	case commandStop:
		log.Info("stopping...")
		err = s.usingManager(s.controlService(svc.Stop, svc.Stopped))
	case commandRun:
		log.Info("running...")
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)

		err := s.start()

		if err != nil {
			log.Errorf("could not run service: %w", err)
			break
		}

		<-signals
		s.stop()
	default:
		err = fmt.Errorf("unkown command: %s", cmd)
	}

	if err != nil {
		log.Panic("could not perform control command:", cmd, ". ", err)
	}
}

func (c *cli) start() error {
	for _, service := range c.services {
		err := service.Start()

		if err != nil {
			return fmt.Errorf("could not run service: %w", err)
		}
	}

	return nil
}

func (c *cli) stop() {
	for i := len(c.services) - 1; i >= 0; i-- {
		c.services[i].Stop()
	}
}

type command string

const (
	commandInstall   = command("install")
	commandUninstall = command("uninstall")
	commandStart     = command("start")
	commandStop      = command("stop")
	commandRun       = command("run")
)

func (s *cli) runWhenInServiceMode() (isService bool, err error) {
	isService, err = svc.IsWindowsService()

	if err != nil {
		return isService, fmt.Errorf("could not check if in service mode: %w", err)
	}

	if isService {
		log.Info("RUNNING SERVICE")
		err := svc.Run(s.name, &handler{service: s})

		if err != nil {
			return isService, fmt.Errorf("could not run service: %w", err)
		}
	}

	return isService, nil
}
