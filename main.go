package main

import (
	"embed"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/gvidasja/button-box-vjoy-feeder/internal/buttonbox"
	"github.com/gvidasja/button-box-vjoy-feeder/internal/device"
	appEvents "github.com/gvidasja/button-box-vjoy-feeder/internal/events"
	"github.com/gvidasja/button-box-vjoy-feeder/internal/handbrake"
	"github.com/gvidasja/button-box-vjoy-feeder/internal/serial"
	"github.com/gvidasja/button-box-vjoy-feeder/internal/vjoy"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"golang.org/x/sys/windows/registry"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	addToStartup("button-box-vjoy-feeder", os.Args[0])

	logFile, _ := os.OpenFile(`E:\dev\button-box-vjoy-feeder\button-box-vjoy-feeder.log`, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	w := io.MultiWriter(logFile, os.Stdout)
	slog.SetDefault(slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{Level: slog.LevelInfo})))

	slog.Info("working dir", "dir", getWorkingDir())

	appEventProducer := appEvents.NewAppEventProducer(time.Second / 60)

	buttonBoxVJoyDevice := vjoy.NewDevice(1)
	handbrakeVJoyDevice := vjoy.NewDevice(2)

	buttonBoxHandler := buttonbox.NewHandler(device.New(buttonBoxVJoyDevice, device.DeviceConfig{
		MinimumButtonPressDuration: time.Millisecond * 20,
	}), appEventProducer)

	handbrakeHadler := handbrake.NewHandler(device.New(handbrakeVJoyDevice, device.DeviceConfig{
		MinimumButtonPressDuration: time.Millisecond * 20,
	}), appEventProducer)

	buttonBoxSerialConsumer := serial.NewConsumer(3, buttonBoxHandler)
	handbrakeSerialConsumer := serial.NewConsumer(4, handbrakeHadler)

	app := application.New(application.Options{
		Name: "button-box-vjoy-feeder",
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Windows: application.WindowsOptions{
			DisableQuitOnLastWindowClosed: true,
		},
		SingleInstance: &application.SingleInstanceOptions{
			UniqueID: "com.gvidasja.button-box-vjoy-feeder",
		},
		OnShutdown: func() {
			slog.Info("Shutting down...")
			buttonBoxSerialConsumer.Stop()
			handbrakeSerialConsumer.Stop()
			buttonBoxVJoyDevice.Stop()
			handbrakeVJoyDevice.Stop()
			slog.Info("Shutdown complete")
		},
	})

	appEventProducer.SetApp(app)

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Button Box vJoy Feeder",
		URL:   "/",
	})

	window.RegisterHook(events.Windows.WindowClosing, func(event *application.WindowEvent) {
		slog.Info("Window closing event triggered, cancelling...")
		window.Hide()
		event.Cancel()
	})

	window.OnWindowEvent(events.Common.WindowShow, func(ctx *application.WindowEvent) {
		// window.OpenDevTools()
	})

	tray := app.SystemTray.New()

	tray.SetLabel("Button Box VJoy Feeder")

	tray.OnClick(func() {
		slog.Info("Tray icon clicked", "windowVisible", window.IsVisible())
		if window.IsVisible() {
			window.Hide()
		} else {
			window.Show()
			window.Focus()
		}
	})

	menu := application.NewMenu()

	menu.Add("Quit").OnClick(func(*application.Context) {
		slog.Info("Quit clicked")
		app.Quit()
	})

	tray.SetMenu(menu)

	buttonBoxVJoyDevice.Start()
	handbrakeVJoyDevice.Start()
	buttonBoxSerialConsumer.Start()
	handbrakeSerialConsumer.Start()

	err := app.Run()

	if err != nil {
		slog.Error("Error", "err", err)
		os.Exit(1)
	}
}

func getWorkingDir() string {
	workingDir, _ := os.Getwd()
	return workingDir
}

func addToStartup(appName, exePath string) error {
	k, _, err := registry.CreateKey(registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Run`, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer k.Close()
	return k.SetStringValue(appName, exePath)
}
