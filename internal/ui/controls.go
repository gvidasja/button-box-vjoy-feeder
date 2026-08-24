package ui

import (
	"fmt"

	"github.com/gvidasja/button-box-vjoy-feeder/internal/device"
	"github.com/rodrigocfd/windigo/co"
	windigo "github.com/rodrigocfd/windigo/ui"
	"github.com/rodrigocfd/windigo/win"
)

type controls struct {
	quitButton   *windigo.Button
	buttonLabels map[int]*windigo.Static
	buttonStates map[win.HWND]bool
	activeBrush  win.HBRUSH
}

func createWindow() *windigo.Main {
	return windigo.NewMain(windigo.OptsMain().Title("Button Box vJoy Feeder").ClassIconId(101).Size(windigo.Dpi(380, 320)).Style(co.WS_CAPTION | co.WS_SYSMENU | co.WS_MINIMIZEBOX | co.WS_SIZEBOX | co.WS_VISIBLE))
}

func createControls(window *windigo.Main) controls {
	brush, err := win.CreateBrushIndirect(&win.LOGBRUSH{Style: co.BRS_SOLID, Color: win.RGB(70, 170, 90)})
	if err != nil {
		panic(err)
	}
	result := controls{
		quitButton:   windigo.NewButton(window, windigo.OptsButton().Text("Quit").Position(windigo.Dpi(290, 10)).Width(windigo.DpiX(70)).Height(windigo.DpiY(24))),
		buttonLabels: make(map[int]*windigo.Static),
		buttonStates: make(map[win.HWND]bool),
		activeBrush:  brush,
	}
	for rowIndex, row := range buttonGrid {
		currentX := 10
		for _, group := range row {
			for _, buttonID := range group {
				label := windigo.NewStatic(window, windigo.OptsStatic().Text(fmt.Sprintf("%d", buttonID)).Position(windigo.Dpi(currentX, 95+rowIndex*25)).Size(windigo.Dpi(22, 20)))
				result.buttonLabels[buttonID] = label
				result.buttonStates[label.Hwnd()] = false
				currentX += 24
			}
			currentX += 8
		}
	}
	return result
}

func setupInputUpdates(window *windigo.Main, updates *device.Updates, controls controls) {
	if updates == nil {
		return
	}
	updates.OnButton(func(event device.ButtonUpdate) {
		window.UiThread(func() {
			if label := controls.buttonLabels[int(event.Button)]; label != nil {
				controls.buttonStates[label.Hwnd()] = event.State
				text := fmt.Sprintf("%d", event.Button)
				if event.State {
					text = fmt.Sprintf("[%d]", event.Button)
				}
				label.Hwnd().SetWindowText(text)
				_ = label.Hwnd().InvalidateRect(nil, true)
			}
		})
	})
}

func setupQuitButton(window *windigo.Main, button *windigo.Button, quitting *bool) {
	button.On().BnClicked(func() {
		*quitting = true
		window.Hwnd().PostMessage(co.WM_CLOSE, 0, 0)
	})
}

func setupControlPainting(window *windigo.Main, controls controls) {
	window.On().WmCtlColorStatic(func(color windigo.WmCtlColor) win.HBRUSH {
		if controls.buttonStates[color.HwndControl()] {
			_, _ = color.Hdc().SetBkColor(win.RGB(70, 170, 90))
			_, _ = color.Hdc().SetTextColor(win.RGB(255, 255, 255))
			return controls.activeBrush
		}
		return 0
	})
}
