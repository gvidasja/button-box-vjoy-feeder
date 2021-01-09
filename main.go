package main

import (
	"os"
	"strings"

	"golang.org/x/sys/windows/svc"
)

func main() {
	service := &service{
		name:        "button-box-vjoy-feeded",
		description: "button-box-vjoy-feeded",
		service:     readButtons,
	}

	err := service.runWhenInServiceMode()

	if err != nil {
		panic(err)
	}

	if len(os.Args) < 2 {
		return
	}

	cmd := strings.ToLower(os.Args[1])

	switch cmd {
	case "install", "i":
		err = service.install()
	case "uninstall", "u", "remove", "r":
		err = service.remove()
	case "start":
		err = service.start()
	case "stop":
		err = service.control(svc.Stop, svc.Stopped)
	case "pause":
		err = service.control(svc.Pause, svc.Paused)
	case "continue":
		err = service.control(svc.Continue, svc.Running)
	}

	if err != nil {
		panic(err)
	}
}
