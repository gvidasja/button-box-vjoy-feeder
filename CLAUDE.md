# CLAUDE.md

Windows tray app (Go) that reads serial data from two Arduino peripherals (a button box and a
handbrake) and feeds it into a vJoy virtual joystick. `arduino.ino` / `arduino2.ino` are the board
firmware.

Personal, single-user project run on the dev machine: hardcoded paths/ports/IDs (COM3 button box,
COM4 handbrake, vJoy device 1, absolute log path in `main.go`) are intended, not bugs to fix.

## Build
- `.\build.ps1` — real build: icon resource (needs `windres`/mingw on PATH), `-H windowsgui`, copies
  `vJoyInterface.dll` into `bin/`. Plain `go build ./...` is fine for compile-checking only.
- No tests.

## Gotchas
- vJoy is called through a manual DLL binding (`internal/vjoy/vjoydll.go`), not cgo.
- UI is Win32 via `windigo`; git history mentions wails, but the project migrated off it.
- Single instance via a named mutex — a second launch silently exits. The app re-adds itself to
  `HKCU\...\Run` (auto-start) on every startup.
