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

func setupTray(window *windigo.Main, trayMessage co.WM, trayID uint32, quitCommand uint16, quitting, visible *bool) {
	window.On().WmCreate(func(windigo.WmCreate) int {
		instance, err := win.GetModuleHandle("")
		if err != nil {
			slog.Error("failed to get module handle for tray icon", "err", err)
			return 0
		}
		icon, err := instance.LoadIcon(win.IconResId(101))
		if err != nil {
			slog.Error("failed to load tray icon", "err", err)
			return 0
		}
		data := win.NOTIFYICONDATA{}
		data.SetCbSize()
		data.HWnd = window.Hwnd()
		data.UID = trayID
		data.UFlags = co.NIF_MESSAGE | co.NIF_ICON | co.NIF_TIP
		data.UCallbackMessage = trayMessage
		data.HIcon = icon
		data.SetSzTip("Button Box vJoy Feeder")
		if err := win.Shell_NotifyIcon(co.NIM_ADD, &data); err != nil {
			slog.Error("failed to add tray icon", "err", err)
			return 0
		}
		slog.Info("tray icon added")
		return 0
	})
	window.On().Wm(trayMessage, func(message windigo.Wm) uintptr {
		switch co.WM(uint32(message.LParam)) {
		case co.WM_LBUTTONUP:
			if *visible {
				window.Hwnd().ShowWindow(co.SW_HIDE)
			} else {
				window.Hwnd().ShowWindow(co.SW_SHOW)
				window.Hwnd().SetForegroundWindow()
			}
			*visible = !*visible
		case co.WM_RBUTTONUP:
			showTrayMenu(window, quitCommand, quitting)
		}
		return 0
	})
}

func showTrayMenu(window *windigo.Main, quitCommand uint16, quitting *bool) {
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
	window.Hwnd().SetForegroundWindow()
	command, err := menu.TrackPopupMenu(co.TPM_RETURNCMD|co.TPM_NONOTIFY|co.TPM_RIGHTBUTTON, int(cursor.X), int(cursor.Y), window.Hwnd())
	window.Hwnd().PostMessage(co.WM_NULL, 0, 0)
	if err != nil {
		slog.Error("failed to show tray menu", "err", err)
		return
	}
	if uint16(command) == quitCommand {
		*quitting = true
		window.Hwnd().PostMessage(co.WM_CLOSE, 0, 0)
	}
}

func setupShutdown(window *windigo.Main, trayID uint32, quitting *bool, shutdown func()) {
	window.On().WmClose(func() {
		if !*quitting {
			window.Hwnd().ShowWindow(co.SW_HIDE)
			return
		}
		window.Hwnd().DestroyWindow()
	})
	window.On().WmDestroy(func() {
		data := win.NOTIFYICONDATA{}
		data.SetCbSize()
		data.HWnd = window.Hwnd()
		data.UID = trayID
		_ = win.Shell_NotifyIcon(co.NIM_DELETE, &data)
		if shutdown != nil {
			shutdown()
		}
		if *quitting {
			win.PostQuitMessage(0)
		}
	})
}
