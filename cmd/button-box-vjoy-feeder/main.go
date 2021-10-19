package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/gvidasja/button-box-vjoy-feeder/internal/buttons"
	"github.com/gvidasja/button-box-vjoy-feeder/internal/log"
	"github.com/gvidasja/button-box-vjoy-feeder/internal/service"
	"github.com/pkg/errors"
	"golang.org/x/sys/windows/svc"
)

func main() {
	log.SetFile("C:\\logs\\button-box-vjoy-feeder.log")
	defer log.Close()

	service := service.New("button-box-vjoy-feeded", "button-box-vjoy-feeded", buttons.Service)

	isService, err := service.RunWhenInServiceMode()

	if err != nil {
		log.Panic("could not run in service mode:", err)
	}

	if isService {
		return
	}

	log.SetDebug(true)

	if len(os.Args) < 2 {
		cmdDoc := func(cmd command, desc string) string {
			return fmt.Sprintf("  * %s - %s", cmd, desc)
		}

		commandDescriptions := []string{
			cmdDoc(commandInstall, "install service"),
			cmdDoc(commandUninstall, "uninstall service"),
			cmdDoc(commandStart, "start service"),
			cmdDoc(commandStop, "stop service"),
			cmdDoc(commandPause, "pause service"),
			cmdDoc(commandContinue, "unpaused service"),
			cmdDoc(commandRun, "run as a simple executable"),
		}

		fmt.Printf("\nplease use one of these commands:\n\n%s\n\n", strings.Join(commandDescriptions, "\n"))

		return
	}

	cmd := command(strings.ToLower(os.Args[1]))

	switch cmd {
	case commandInstall:
		err = service.Install()
	case commandUninstall:
		err = service.Remove()
	case commandStart:
		err = service.Start()
	case commandStop:
		err = service.Control(svc.Stop, svc.Stopped)
	case commandPause:
		err = service.Control(svc.Pause, svc.Paused)
	case commandContinue:
		err = service.Control(svc.Continue, svc.Running)
	case commandRun:
		signals := make(chan os.Signal, 1)
		done := make(chan bool)
		signal.Notify(signals, os.Interrupt)

		go func() {
			buttons.Service(make(chan<- bool), done)
		}()

		<-signals
		done <- true
	default:
		err = errors.Errorf("unkown command: %s", cmd)
	}

	if err != nil {
		log.Panic("could not perform control command:", cmd, ". ", err)
	}
}

type command string

const (
	commandInstall   = command("install")
	commandUninstall = command("uninstall")
	commandStart     = command("start")
	commandStop      = command("stop")
	commandPause     = command("pause")
	commandContinue  = command("continue")
	commandRun       = command("run")
)
