package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

type service struct {
	name        string
	description string
	service     func(stoppedEvent chan<- bool, stopCommand <-chan bool)
}

func (s *service) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	changes <- svc.Status{State: svc.StartPending}

	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

	stopCommand := make(chan bool)
	stoppedEvent := make(chan bool)

	go s.service(stoppedEvent, stopCommand)

loop:
	for {
		select {
		case <-stoppedEvent:
			break loop
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				stopCommand <- true
			case svc.Continue:
				changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
			case svc.Pause:
				changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
			}
		}
	}

	changes <- svc.Status{State: svc.StopPending}
	return
}

func (s *service) runWhenInServiceMode() error {
	isService, err := svc.IsWindowsService()

	if err != nil {
		return errors.Wrap(err, "Could not check if in service mode")
	}

	if isService {
		err := svc.Run(s.name, s)

		if err != nil {
			return errors.Wrap(err, "Could not run service")
		}
	}

	return nil
}

func (s *service) install() error {
	program := os.Args[0]
	programPath, _ := filepath.Abs(program)

	log.Println(programPath)
	manager, err := mgr.Connect()

	if err != nil {
		return errors.Wrap(err, "Could not connect to service manager.")
	}

	defer manager.Disconnect()

	service, err := manager.CreateService(s.name, programPath, mgr.Config{DisplayName: s.description, StartType: mgr.StartAutomatic}, "is", "auto-started")

	if err != nil {
		return errors.Wrap(err, "Could not create service.")
	}

	defer service.Close()

	return nil
}

func (s *service) remove() error {
	manager, err := mgr.Connect()
	if err != nil {
		return errors.Wrap(err, "Could not connect to service manager")
	}
	defer manager.Disconnect()

	service, err := manager.OpenService(s.name)
	if err != nil {
		return errors.Wrapf(err, "service %s is not installed", s.name)
	}
	defer service.Close()

	err = service.Delete()
	if err != nil {
		return errors.Wrap(err, "Could not delete service")
	}

	return nil
}

func (s *service) start() error {
	manager, err := mgr.Connect()
	if err != nil {
		return errors.Wrap(err, "Could not connect to service manager")
	}

	defer manager.Disconnect()

	service, err := manager.OpenService(s.name)

	if err != nil {
		return errors.Wrap(err, "Could not access service")
	}

	defer service.Close()

	err = service.Start("is", "manual-started")

	if err != nil {
		return errors.Wrap(err, "Could not start service")
	}

	return nil
}

func (s *service) control(command svc.Cmd, newState svc.State) error {
	manager, err := mgr.Connect()
	if err != nil {
		return errors.Wrap(err, "Could not connect to service manager")
	}
	defer manager.Disconnect()

	service, err := manager.OpenService(s.name)
	if err != nil {
		return errors.Wrap(err, "Could not access service")
	}
	defer service.Close()

	status, err := service.Control(command)
	if err != nil {
		return errors.Wrapf(err, "Could not send control=%d", command)
	}

	timeout := time.Now().Add(10 * time.Second)

	for status.State != newState {
		if timeout.Before(time.Now()) {
			return fmt.Errorf("Timeout waiting for service to go to state=%d", newState)
		}

		time.Sleep(300 * time.Millisecond)

		status, err = service.Query()

		if err != nil {
			return errors.Wrap(err, "Could not retrieve service status")
		}
	}

	return nil
}
