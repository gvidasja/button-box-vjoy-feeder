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

type buttonPosition struct{ id, x, y int }

type layout struct {
	buttons       []buttonPosition
	right         int // right edge of the button grid
	handbrakeY    int
	quitX, quitY  int
	width, height int
}

func computeLayout() layout {
	var l layout
	for rowIndex, row := range buttonGrid {
		x := gridLeft
		for _, group := range row {
			for _, id := range group {
				l.buttons = append(l.buttons, buttonPosition{id, x, gridTop + rowIndex*rowHeight})
				l.right = max(l.right, x+labelWidth)
				x += buttonSpacing
			}
			x += groupSpacing
		}
	}
	l.handbrakeY = gridTop + (len(buttonGrid)-1)*rowHeight + labelHeight + handbrakeGap
	l.quitX = gridLeft + (l.right-gridLeft-quitWidth)/2
	l.quitY = l.handbrakeY + handbrakeHeight + quitGap
	l.width = l.right - gridLeft + 2*margin
	l.height = l.quitY + quitHeight + margin
	return l
}

type controls struct {
	quitButton       *windigo.Button
	buttonLabels     map[int]*windigo.Static
	buttonStates     map[win.HWND]bool
	activeBrush      win.HBRUSH
	handbrakeBar     *windigo.Button
	handbrakePercent atomic.Int32
}

func createWindow(l layout, startVisible bool) *windigo.Main {
	cmdShow := co.SW_HIDE
	if startVisible {
		cmdShow = co.SW_SHOW
	}
	return windigo.NewMain(windigo.OptsMain().Title(windowTitle).ClassIconId(iconID).Size(windigo.Dpi(l.width, l.height)).Style(co.WS_CAPTION | co.WS_SYSMENU | co.WS_MINIMIZEBOX).CmdShow(cmdShow))
}

func createControls(window *windigo.Main, l layout) *controls {
	c := &controls{
		quitButton:   windigo.NewButton(window, windigo.OptsButton().Text("Quit").Position(windigo.Dpi(l.quitX, l.quitY)).Width(windigo.DpiX(quitWidth)).Height(windigo.DpiY(quitHeight))),
		buttonLabels: make(map[int]*windigo.Static),
		buttonStates: make(map[win.HWND]bool),
		activeBrush:  solidBrush(activeColor),
		handbrakeBar: windigo.NewButton(window, windigo.OptsButton().
			Position(windigo.Dpi(gridLeft, l.handbrakeY)).
			Width(windigo.DpiX(l.right-gridLeft)).
			Height(windigo.DpiY(handbrakeHeight)).
			CtrlStyle(co.BS_OWNERDRAW).
			WndStyle(co.WS_CHILD|co.WS_VISIBLE)),
	}
	for _, b := range l.buttons {
		c.buttonLabels[b.id] = windigo.NewStatic(window, windigo.OptsStatic().Text(fmt.Sprint(b.id)).Position(windigo.Dpi(b.x, b.y)).Size(windigo.Dpi(labelWidth, labelHeight)).CtrlStyle(co.SS_CENTER|co.SS_CENTERIMAGE|co.SS_NOTIFY).WndStyle(co.WS_CHILD|co.WS_VISIBLE|co.WS_BORDER))
	}
	return c
}

func (a *app) setupInputUpdates(updates *device.Updates) {
	c := a.controls
	updates.OnButton(func(event device.ButtonUpdate) {
		a.window.UiThread(func() {
			if label := c.buttonLabels[int(event.Button)]; label != nil {
				c.buttonStates[label.Hwnd()] = event.State
				_ = label.Hwnd().InvalidateRect(nil, true)
			}
		})
	})
	updates.OnAxis(func(event device.AxisUpdate) {
		if event.Axis != device.HANDBRAKE_AXIS {
			return
		}
		percent := (event.Value*100 + math.MaxInt16/2) / math.MaxInt16
		if c.handbrakePercent.Swap(percent) == percent {
			return
		}
		a.window.UiThread(func() {
			_ = c.handbrakeBar.Hwnd().InvalidateRect(nil, false)
		})
	})
}

func (a *app) setupPainting() {
	c := a.controls
	a.window.On().WmCtlColorStatic(func(color windigo.WmCtlColor) win.HBRUSH {
		if !c.buttonStates[color.HwndControl()] {
			return 0
		}
		_, _ = color.Hdc().SetBkColor(activeColor)
		_, _ = color.Hdc().SetTextColor(win.RGB(255, 255, 255))
		return c.activeBrush
	})
	a.window.On().WmDrawItem(func(p windigo.WmDrawItem) {
		if dis := p.DrawItemStruct(); dis.HwndItem == c.handbrakeBar.Hwnd() {
			drawHandbrakeBar(*dis, c.handbrakePercent.Load(), c.activeBrush)
		}
	})
}
