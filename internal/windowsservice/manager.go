package windowsservice

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

func (s *cli) install() error {
	program := os.Args[0]
	programPath, _ := filepath.Abs(program)

	log.Infof("installation path: %s", programPath)
	manager, err := mgr.Connect()

	if err != nil {
		return fmt.Errorf("could not connect to service manager.: %w", err)
	}

	defer manager.Disconnect()

	service, err := manager.CreateService(s.name, programPath, mgr.Config{DisplayName: s.description, StartType: mgr.StartAutomatic}, "is", "auto-started")

	if err != nil {
		return fmt.Errorf("could not create service.: %w", err)
	}

	defer service.Close()

	return nil
}

func (*cli) removeService(service *mgr.Service) error {
	err := service.Delete()
	if err != nil {
		return fmt.Errorf("could not delete service: %w", err)
	}

	return nil
}

func (*cli) startService(service *mgr.Service) error {
	err := service.Start()

	if err != nil {
		return fmt.Errorf("could not start service: %w", err)
	}

	return nil
}

func (*cli) controlService(command svc.Cmd, expectedState svc.State) func(*mgr.Service) error {
	return func(service *mgr.Service) error {
		status, err := service.Control(command)
		if err != nil {
			return fmt.Errorf("could not send control=%d: %w", command, err)
		}

		timeout := time.Now().Add(10 * time.Second)

		for status.State != expectedState {
			if timeout.Before(time.Now()) {
				return fmt.Errorf("timeout waiting for service to go to state=%d", expectedState)
			}

			time.Sleep(300 * time.Millisecond)

			status, err = service.Query()

			if err != nil {
				return fmt.Errorf("could not retrieve service status: %w", err)
			}
		}

		return nil
	}
}

func (c *cli) usingManager(action func(*mgr.Service) error) error {
	manager, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("could not connect to service manager: %w", err)
	}

	defer manager.Disconnect()

	service, err := manager.OpenService(c.name)

	if err != nil {
		return fmt.Errorf("could not access service: %w", err)
	}

	defer service.Close()

	return action(service)
}
