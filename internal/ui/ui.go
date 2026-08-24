package ui

import (
	"runtime"

	"github.com/gvidasja/button-box-vjoy-feeder/internal/device"
	"github.com/rodrigocfd/windigo/co"
)

var buttonGrid = [][][]int{
	{{1, 2}, {9, 13, 17}, {21, 22}},
	{{3, 4}, {10, 14, 18}, {23, 24}},
	{{5, 6}, {11, 15, 19}, {25, 26}},
	{{7, 8}, {12, 16, 20}, {27, 28}},
}

type App struct {
	updates  *device.Updates
	shutdown func()
}

func New(updates *device.Updates, shutdown func()) *App {
	return &App{updates: updates, shutdown: shutdown}
}

func (app *App) Run() {
	runtime.LockOSThread()
	const trayMessage co.WM = co.WM_USER + 1
	const trayID uint32 = 1
	const quitCommand uint16 = 1001
	quitting := false
	windowVisible := true

	window := createWindow()
	controls := createControls(window)
	defer controls.activeBrush.DeleteObject()
	setupInputUpdates(window, app.updates, controls)
	setupControlPainting(window, controls)
	setupQuitButton(window, controls.quitButton, &quitting)
	setupTray(window, trayMessage, trayID, quitCommand, &quitting, &windowVisible)
	setupShutdown(window, trayID, &quitting, app.shutdown)
	window.RunAsMain()
}
