package ui

import (
	"fmt"
	"math"
	"sync/atomic"

	"github.com/gvidasja/button-box-vjoy-feeder/internal/device"
	"github.com/rodrigocfd/windigo/co"
	windigo "github.com/rodrigocfd/windigo/ui"
	"github.com/rodrigocfd/windigo/win"
)

type controls struct {
	quitButton       *windigo.Button
	buttonLabels     map[int]*windigo.Static
	buttonStates     map[win.HWND]bool
	activeBrush      win.HBRUSH
	handbrakeBar     *windigo.Button
	handbrakePercent *int32
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
	handbrakeHeight = 22
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

func createWindow(startVisible bool) *windigo.Main {
	contentWidth := gridContentWidth() - gridLeft + 2*margin
	gridBottom := gridTop + (len(buttonGrid)-1)*rowHeight + labelHeight
	handbrakeBottom := gridBottom + handbrakeGap + handbrakeHeight
	contentHeight := handbrakeBottom + quitGap + quitHeight + margin
	cmdShow := co.SW_HIDE
	if startVisible {
		cmdShow = co.SW_SHOW
	}
	return windigo.NewMain(windigo.OptsMain().Title(windowTitle).ClassIconId(101).Size(windigo.Dpi(contentWidth, contentHeight)).Style(co.WS_CAPTION | co.WS_SYSMENU | co.WS_MINIMIZEBOX).CmdShow(cmdShow))
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
		quitButton:       windigo.NewButton(window, windigo.OptsButton().Text("Quit").Position(windigo.Dpi(quitX, quitY)).Width(windigo.DpiX(quitWidth)).Height(windigo.DpiY(quitHeight))),
		buttonLabels:     make(map[int]*windigo.Static),
		buttonStates:     make(map[win.HWND]bool),
		activeBrush:      brush,
		handbrakePercent: new(int32),
		handbrakeBar: windigo.NewButton(window, windigo.OptsButton().
			Position(windigo.Dpi(gridLeft, handbrakeY)).
			Width(windigo.DpiX(contentRight-gridLeft)).
			Height(windigo.DpiY(handbrakeHeight)).
			CtrlStyle(co.BS_OWNERDRAW).
			WndStyle(co.WS_CHILD|co.WS_VISIBLE)),
	}
	for rowIndex, row := range buttonGrid {
		currentX := gridLeft
		for _, group := range row {
			for _, buttonID := range group {
				label := windigo.NewStatic(window, windigo.OptsStatic().Text(fmt.Sprintf("%d", buttonID)).Position(windigo.Dpi(currentX, gridTop+rowIndex*rowHeight)).Size(windigo.Dpi(labelWidth, labelHeight)).CtrlStyle(co.SS_CENTER|co.SS_CENTERIMAGE|co.SS_NOTIFY).WndStyle(co.WS_CHILD|co.WS_VISIBLE|co.WS_BORDER))
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
				label.Hwnd().SetWindowText(text)
				_ = label.Hwnd().InvalidateRect(nil, true)
			}
		})
	})
	updates.OnAxis(func(event device.AxisUpdate) {
		if event.Axis != handbrakeAxisID {
			return
		}
		percent := (event.Value*100 + math.MaxInt16/2) / math.MaxInt16
		if atomic.SwapInt32(controls.handbrakePercent, percent) == percent {
			return
		}
		window.UiThread(func() {
			_ = controls.handbrakeBar.Hwnd().InvalidateRect(nil, false)
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
