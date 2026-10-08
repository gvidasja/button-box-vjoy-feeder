package ui

import (
	"fmt"
	"sync/atomic"

	"github.com/rodrigocfd/windigo/co"
	windigo "github.com/rodrigocfd/windigo/ui"
	"github.com/rodrigocfd/windigo/win"
)

func setupHandbrakeBar(window *windigo.Main, controls controls) {
	window.On().WmDrawItem(func(p windigo.WmDrawItem) {
		dis := p.DrawItemStruct()
		if dis.HwndItem != controls.handbrakeBar.Hwnd() {
			return
		}
		drawHandbrakeBar(*dis, atomic.LoadInt32(controls.handbrakePercent))
	})
}

func drawHandbrakeBar(dis win.DRAWITEMSTRUCT, percent int32) {
	width := int(dis.RcItem.Right - dis.RcItem.Left)
	height := int(dis.RcItem.Bottom - dis.RcItem.Top)
	if width <= 0 || height <= 0 {
		return
	}

	memDC, err := dis.Hdc.CreateCompatibleDC()
	if err != nil {
		return
	}
	defer memDC.DeleteDC()
	bitmap, err := dis.Hdc.CreateCompatibleBitmap(width, height)
	if err != nil {
		return
	}
	defer bitmap.DeleteObject()
	oldBitmap, _ := memDC.SelectObjectBmp(bitmap)
	defer memDC.SelectObjectBmp(oldBitmap)

	full := win.RECT{Left: 0, Top: 0, Right: int32(width), Bottom: int32(height)}

	bgBrush, _ := win.CreateBrushIndirect(&win.LOGBRUSH{Style: co.BRS_SOLID, Color: win.RGB(255, 255, 255)})
	defer bgBrush.DeleteObject()
	_ = memDC.FillRect(&full, bgBrush)

	if percent > 0 {
		fillWidth := min(int32(width)*percent/100, int32(width))
		fillRect := win.RECT{Left: 0, Top: 0, Right: fillWidth, Bottom: int32(height)}
		fillBrush, _ := win.CreateBrushIndirect(&win.LOGBRUSH{Style: co.BRS_SOLID, Color: win.RGB(70, 170, 90)})
		defer fillBrush.DeleteObject()
		_ = memDC.FillRect(&fillRect, fillBrush)
	}

	borderBrush, _ := win.CreateBrushIndirect(&win.LOGBRUSH{Style: co.BRS_SOLID, Color: win.RGB(120, 120, 120)})
	defer borderBrush.DeleteObject()
	_ = memDC.FrameRect(&full, borderBrush)

	if font, _ := dis.HwndItem.SendMessage(co.WM_GETFONT, 0, 0); font != 0 {
		oldFont, _ := memDC.SelectObjectFont(win.HFONT(font))
		defer memDC.SelectObjectFont(oldFont)
	}
	_, _ = memDC.SetBkMode(co.BKMODE_TRANSPARENT)
	_, _ = memDC.SetTextColor(win.RGB(0, 0, 0))
	_, _ = memDC.DrawText(fmt.Sprintf("%d%%", percent), &full, co.DT_CENTER|co.DT_VCENTER|co.DT_SINGLELINE)

	_ = dis.Hdc.BitBlt(win.POINT{X: dis.RcItem.Left, Y: dis.RcItem.Top}, win.SIZE{Cx: int32(width), Cy: int32(height)}, memDC, win.POINT{}, co.ROP_SRCCOPY)
}
