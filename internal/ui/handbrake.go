package ui

import (
	"fmt"

	"github.com/rodrigocfd/windigo/co"
	"github.com/rodrigocfd/windigo/win"
)

func drawHandbrakeBar(dis win.DRAWITEMSTRUCT, percent int32, fillBrush win.HBRUSH) {
	width := dis.RcItem.Right - dis.RcItem.Left
	height := dis.RcItem.Bottom - dis.RcItem.Top
	if width <= 0 || height <= 0 {
		return
	}

	// draw off-screen and blit, so frequent updates don't flicker
	memDC, err := dis.Hdc.CreateCompatibleDC()
	if err != nil {
		return
	}
	defer memDC.DeleteDC()
	bitmap, err := dis.Hdc.CreateCompatibleBitmap(int(width), int(height))
	if err != nil {
		return
	}
	defer bitmap.DeleteObject()
	oldBitmap, _ := memDC.SelectObjectBmp(bitmap)
	defer memDC.SelectObjectBmp(oldBitmap)

	full := win.RECT{Right: width, Bottom: height}

	bgBrush := solidBrush(win.RGB(255, 255, 255))
	defer bgBrush.DeleteObject()
	_ = memDC.FillRect(&full, bgBrush)

	if percent > 0 {
		fill := win.RECT{Right: min(width*percent/100, width), Bottom: height}
		_ = memDC.FillRect(&fill, fillBrush)
	}

	borderBrush := solidBrush(win.RGB(120, 120, 120))
	defer borderBrush.DeleteObject()
	_ = memDC.FrameRect(&full, borderBrush)

	if font, _ := dis.HwndItem.SendMessage(co.WM_GETFONT, 0, 0); font != 0 {
		oldFont, _ := memDC.SelectObjectFont(win.HFONT(font))
		defer memDC.SelectObjectFont(oldFont)
	}
	_, _ = memDC.SetBkMode(co.BKMODE_TRANSPARENT)
	_, _ = memDC.SetTextColor(win.RGB(0, 0, 0))
	_, _ = memDC.DrawText(fmt.Sprintf("%d%%", percent), &full, co.DT_CENTER|co.DT_VCENTER|co.DT_SINGLELINE)

	_ = dis.Hdc.BitBlt(win.POINT{X: dis.RcItem.Left, Y: dis.RcItem.Top}, win.SIZE{Cx: width, Cy: height}, memDC, win.POINT{}, co.ROP_SRCCOPY)
}
