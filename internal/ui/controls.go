package ui

import (
	"fmt"
	"math"

	"github.com/gvidasja/button-box-vjoy-feeder/internal/device"
	"github.com/rodrigocfd/windigo/co"
	windigo "github.com/rodrigocfd/windigo/ui"
	"github.com/rodrigocfd/windigo/win"
)

type controls struct {
	quitButton     *windigo.Button
	buttonLabels   map[int]*windigo.Static
	buttonStates   map[win.HWND]bool
	activeBrush    win.HBRUSH
	handbrakeLabel *windigo.Static
}

const (
	gridLeft        = 10
	gridTop         = 10
	rowHeight       = 25
	labelWidth      = 22
	labelHeight     = 20
	buttonSpacing   = 24
	groupSpacing    = 8
	handbrakeGap    = 8
	handbrakeHeight = 20
	quitGap         = 15
	quitWidth       = 70
	quitHeight      = 24
	margin          = 10
)

func gridContentWidth() int {
	maxRight := 0
	for _, row := range buttonGrid {
		currentX := gridLeft
		for _, group := range row {
			for range group {
				currentX += buttonSpacing
			}
			currentX += groupSpacing
		}
		right := currentX - groupSpacing - buttonSpacing + labelWidth
		if right > maxRight {
			maxRight = right
		}
	}
	return maxRight
}

func createWindow() *windigo.Main {
	contentWidth := gridContentWidth() - gridLeft + 2*margin
	gridBottom := gridTop + (len(buttonGrid)-1)*rowHeight + labelHeight
	handbrakeBottom := gridBottom + handbrakeGap + handbrakeHeight
	contentHeight := handbrakeBottom + quitGap + quitHeight + margin
	return windigo.NewMain(windigo.OptsMain().Title("Button Box vJoy Feeder").ClassIconId(101).Size(windigo.Dpi(contentWidth+20, contentHeight+130)).Style(co.WS_CAPTION | co.WS_SYSMENU | co.WS_MINIMIZEBOX).CmdShow(co.SW_HIDE))
}

func createControls(window *windigo.Main) controls {
	brush, err := win.CreateBrushIndirect(&win.LOGBRUSH{Style: co.BRS_SOLID, Color: win.RGB(70, 170, 90)})
	if err != nil {
		panic(err)
	}
	contentRight := gridContentWidth()
	gridBottom := gridTop + (len(buttonGrid)-1)*rowHeight + labelHeight
	handbrakeY := gridBottom + handbrakeGap
	quitX := gridLeft + (contentRight-gridLeft-quitWidth)/2
	quitY := handbrakeY + handbrakeHeight + quitGap
	result := controls{
		quitButton:     windigo.NewButton(window, windigo.OptsButton().Text("Quit").Position(windigo.Dpi(quitX, quitY)).Width(windigo.DpiX(quitWidth)).Height(windigo.DpiY(quitHeight))),
		buttonLabels:   make(map[int]*windigo.Static),
		buttonStates:   make(map[win.HWND]bool),
		activeBrush:    brush,
		handbrakeLabel: windigo.NewStatic(window, windigo.OptsStatic().Text("Handbrake: --").Position(windigo.Dpi(gridLeft, handbrakeY)).Size(windigo.Dpi(contentRight-gridLeft, handbrakeHeight))),
	}
	for rowIndex, row := range buttonGrid {
		currentX := gridLeft
		for _, group := range row {
			for _, buttonID := range group {
				label := windigo.NewStatic(window, windigo.OptsStatic().Text(fmt.Sprintf("%d", buttonID)).Position(windigo.Dpi(currentX, gridTop+rowIndex*rowHeight)).Size(windigo.Dpi(labelWidth, labelHeight)))
				result.buttonLabels[buttonID] = label
				result.buttonStates[label.Hwnd()] = false
				currentX += buttonSpacing
			}
			currentX += groupSpacing
		}
	}
	return result
}

// handbrakeAxisID must match the axis ID the handbrake handler publishes under (internal/handbrake/handler.go).
const handbrakeAxisID = device.AxisID(0x32)

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
	updates.OnAxis(func(event device.AxisUpdate) {
		if event.Axis != handbrakeAxisID {
			return
		}
		window.UiThread(func() {
			percent := event.Value * 100 / math.MaxInt16
			controls.handbrakeLabel.Hwnd().SetWindowText(fmt.Sprintf("Handbrake: %d%%", percent))
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
