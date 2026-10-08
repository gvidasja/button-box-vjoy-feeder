package ui

import (
	"runtime"

	"github.com/gvidasja/button-box-vjoy-feeder/internal/device"
	"github.com/rodrigocfd/windigo/co"
	windigo "github.com/rodrigocfd/windigo/ui"
	"github.com/rodrigocfd/windigo/win"
)

const (
	windowTitle = "Button Box vJoy Feeder"
	iconID      = 101 // see appicon.rc
	trayMessage = co.WM_USER + 1
	trayID      = 1
	quitCommand = 1001
)

var activeColor = win.RGB(70, 170, 90)

var buttonGrid = [][][]int{
	{{1, 2}, {9, 13, 17}, {21, 22}},
	{{3, 4}, {10, 14, 18}, {23, 24}},
	{{5, 6}, {11, 15, 19}, {25, 26}},
	{{7, 8}, {12, 16, 20}, {27, 28}},
}

type app struct {
	window   *windigo.Main
	controls *controls
	quitting bool
}

// Run creates the window and tray icon and blocks until the user quits.
func Run(updates *device.Updates, shutdown func(), startVisible bool) {
	runtime.LockOSThread()
	l := computeLayout()
	a := &app{window: createWindow(l, startVisible)}
	a.controls = createControls(a.window, l)
	defer a.controls.activeBrush.DeleteObject()
	a.controls.quitButton.On().BnClicked(a.quit)
	a.setupInputUpdates(updates)
	a.setupPainting()
	a.setupTray()
	a.setupShutdown(shutdown)
	a.window.RunAsMain()
}

func (a *app) quit() {
	a.quitting = true
	a.window.Hwnd().PostMessage(co.WM_CLOSE, 0, 0)
}

// ShowRunningInstance brings an already-running instance's window to the foreground.
func ShowRunningInstance() {
	if hwnd, found := win.FindWindow(win.ClassNameNone(), windowTitle); found {
		show(hwnd)
	}
}

func show(hwnd win.HWND) {
	hwnd.ShowWindow(co.SW_SHOW)
	hwnd.SetForegroundWindow()
}

func solidBrush(color win.COLORREF) win.HBRUSH {
	brush, _ := win.CreateBrushIndirect(&win.LOGBRUSH{Style: co.BRS_SOLID, Color: color})
	return brush
}
