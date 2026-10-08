package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"slices"
	"time"

	"github.com/gvidasja/button-box-vjoy-feeder/internal/buttonbox"
	"github.com/gvidasja/button-box-vjoy-feeder/internal/device"
	"github.com/gvidasja/button-box-vjoy-feeder/internal/handbrake"
	"github.com/gvidasja/button-box-vjoy-feeder/internal/serial"
	"github.com/gvidasja/button-box-vjoy-feeder/internal/ui"
	"github.com/gvidasja/button-box-vjoy-feeder/internal/vjoy"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

func main() {
	startMinimized := slices.Contains(os.Args[1:], "--minimized")

	mutex, alreadyRunning, err := acquireSingleInstance()
	if err != nil || alreadyRunning {
		if alreadyRunning {
			ui.ShowRunningInstance()
		}
		return
	}
	defer windows.CloseHandle(mutex)

	_ = addToStartup("button-box-vjoy-feeder", os.Args[0])

	logFile, _ := os.OpenFile(`F:\dev\button-box-vjoy-feeder\button-box-vjoy-feeder.log`, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	defer logFile.Close()
	slog.SetDefault(slog.New(slog.NewTextHandler(io.MultiWriter(logFile, os.Stdout), &slog.HandlerOptions{Level: slog.LevelInfo})))

	vjoyDevice := vjoy.NewDevice(1)
	updates := device.NewUpdates()
	outputDevice := device.NewPublishingDevice(vjoyDevice, updates)
	// each consumer gets its own debouncer: they run on separate goroutines
	debounced := func() device.Device {
		return device.NewDebouncedDevice(outputDevice, device.DebouncedDeviceConfig{MinimumButtonPressDuration: 20 * time.Millisecond})
	}
	buttonBoxConsumer := serial.NewConsumer(3, buttonbox.NewHandler(debounced()))
	handbrakeConsumer := serial.NewConsumer(4, handbrake.NewHandler(debounced()))

	if err := vjoyDevice.Start(); err != nil {
		slog.Warn("vjoy start", "err", err)
	}
	_ = buttonBoxConsumer.Start()
	_ = handbrakeConsumer.Start()

	ui.Run(updates, func() {
		buttonBoxConsumer.Stop()
		handbrakeConsumer.Stop()
		vjoyDevice.Stop()
	}, !startMinimized)
}

func acquireSingleInstance() (windows.Handle, bool, error) {
	name, err := windows.UTF16PtrFromString("Local\\button-box-vjoy-feeder")
	if err != nil {
		return 0, false, err
	}
	handle, createErr := windows.CreateMutex(nil, false, name)
	if handle == 0 {
		return 0, false, createErr
	}
	if createErr == windows.ERROR_ALREADY_EXISTS {
		_ = windows.CloseHandle(handle)
		return 0, true, nil
	}
	return handle, false, createErr
}

func addToStartup(appName, exePath string) error {
	key, _, err := registry.CreateKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer key.Close()
	return key.SetStringValue(appName, fmt.Sprintf(`"%s" --minimized`, exePath))
}
