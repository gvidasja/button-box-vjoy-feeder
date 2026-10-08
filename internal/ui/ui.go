package ui

import (
	"runtime"

	"github.com/gvidasja/button-box-vjoy-feeder/internal/device"
	"github.com/rodrigocfd/windigo/co"
	"github.com/rodrigocfd/windigo/win"
)

const windowTitle = "Button Box vJoy Feeder"

var buttonGrid = [][][]int{
	{{1, 2}, {9, 13, 17}, {21, 22}},
	{{3, 4}, {10, 14, 18}, {23, 24}},
	{{5, 6}, {11, 15, 19}, {25, 26}},
	{{7, 8}, {12, 16, 20}, {27, 28}},
}

type App struct {
	updates      *device.Updates
	shutdown     func()
	startVisible bool
}

func New(updates *device.Updates, shutdown func(), startVisible bool) *App {
	return &App{updates: updates, shutdown: shutdown, startVisible: startVisible}
}

func (app *App) Run() {
	runtime.LockOSThread()
	const trayMessage co.WM = co.WM_USER + 1
	const trayID uint32 = 1
	const quitCommand uint16 = 1001
	quitting := false
	windowVisible := app.startVisible

	window := createWindow(app.startVisible)
	controls := createControls(window)
	defer controls.activeBrush.DeleteObject()
	setupInputUpdates(window, app.updates, controls)
	setupControlPainting(window, controls)
	setupQuitButton(window, controls.quitButton, &quitting)
	setupHandbrakeBar(window, controls)
	setupTray(window, trayMessage, trayID, quitCommand, &quitting, &windowVisible)
	setupShutdown(window, trayID, &quitting, app.shutdown)
	window.RunAsMain()
}

// ShowRunningInstance finds an already-running instance of this app by its
// window title and brings its window to the foreground. Returns true if a
// running instance was found.
func ShowRunningInstance() bool {
	hwnd, found := win.FindWindow(win.ClassNameNone(), windowTitle)
	if !found {
		return false
	}
	hwnd.ShowWindow(co.SW_SHOW)
	hwnd.SetForegroundWindow()
	return true
}
