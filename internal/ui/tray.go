package ui

import (
	"log/slog"
	"syscall"
	"unsafe"

	"github.com/rodrigocfd/windigo/co"
	windigo "github.com/rodrigocfd/windigo/ui"
	"github.com/rodrigocfd/windigo/win"
	"golang.org/x/sys/windows"
)

func (a *app) setupTray() {
	a.window.On().WmCreate(func(windigo.WmCreate) int {
		instance, err := win.GetModuleHandle("")
		if err != nil {
			slog.Error("failed to get module handle for tray icon", "err", err)
			return 0
		}
		icon, err := instance.LoadIcon(win.IconResId(iconID))
		if err != nil {
			slog.Error("failed to load tray icon", "err", err)
			return 0
		}
		data := a.trayIconData()
		data.UFlags = co.NIF_MESSAGE | co.NIF_ICON | co.NIF_TIP
		data.UCallbackMessage = trayMessage
		data.HIcon = icon
		data.SetSzTip(windowTitle)
		if err := win.Shell_NotifyIcon(co.NIM_ADD, &data); err != nil {
			slog.Error("failed to add tray icon", "err", err)
			return 0
		}
		slog.Info("tray icon added")
		return 0
	})
	a.window.On().Wm(trayMessage, func(message windigo.Wm) uintptr {
		switch co.WM(uint32(message.LParam)) {
		case co.WM_LBUTTONUP:
			hwnd := a.window.Hwnd()
			if windows.IsWindowVisible(windows.HWND(hwnd)) {
				hwnd.ShowWindow(co.SW_HIDE)
			} else {
				show(hwnd)
			}
		case co.WM_RBUTTONUP:
			a.showTrayMenu()
		}
		return 0
	})
}

func (a *app) trayIconData() win.NOTIFYICONDATA {
	data := win.NOTIFYICONDATA{}
	data.SetCbSize()
	data.HWnd = a.window.Hwnd()
	data.UID = trayID
	return data
}

func (a *app) showTrayMenu() {
	menu, err := win.CreatePopupMenu()
	if err != nil {
		slog.Error("failed to create tray menu", "err", err)
		return
	}
	defer menu.DestroyMenu()
	text, _ := syscall.UTF16PtrFromString("Quit")
	appendMenu := windows.NewLazySystemDLL("user32.dll").NewProc("AppendMenuW")
	if result, _, callErr := appendMenu.Call(uintptr(menu), uintptr(co.MF_STRING), uintptr(quitCommand), uintptr(unsafe.Pointer(text))); result == 0 {
		slog.Error("failed to add quit menu item", "err", callErr)
		return
	}
	cursor, err := win.GetCursorPos()
	if err != nil {
		slog.Error("failed to get tray cursor position", "err", err)
		return
	}
	hwnd := a.window.Hwnd()
	hwnd.SetForegroundWindow()
	command, err := menu.TrackPopupMenu(co.TPM_RETURNCMD|co.TPM_NONOTIFY|co.TPM_RIGHTBUTTON, int(cursor.X), int(cursor.Y), hwnd)
	hwnd.PostMessage(co.WM_NULL, 0, 0)
	if err != nil {
		slog.Error("failed to show tray menu", "err", err)
		return
	}
	if uint16(command) == quitCommand {
		a.quit()
	}
}

// setupShutdown hides the window on close unless the user chose Quit.
// windigo's main window posts WM_QUIT itself once destroyed.
func (a *app) setupShutdown(shutdown func()) {
	a.window.On().WmClose(func() {
		if a.quitting {
			a.window.Hwnd().DestroyWindow()
		} else {
			a.window.Hwnd().ShowWindow(co.SW_HIDE)
		}
	})
	a.window.On().WmDestroy(func() {
		data := a.trayIconData()
		_ = win.Shell_NotifyIcon(co.NIM_DELETE, &data)
		shutdown()
	})
}
