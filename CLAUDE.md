# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

A Windows tray app (Go) that reads serial data from two Arduino-based peripherals (a custom button box and a handbrake) and feeds it into a vJoy virtual joystick device, so the peripherals show up as a joystick to Windows/games. `arduino.ino` / `arduino2.ino` are the firmware sketches for the two boards.

This is a personal, single-user project: it's developed and run on the same machine, for the author's own use only. There's no need to worry about portability, multi-user support, or configurability beyond what the author needs — hardcoded paths/ports/IDs tied to the dev machine are expected and fine, not bugs to fix by default.

## Build & run

```
.\build.ps1
```

This compiles the Windows resource (icon), builds the exe with `-H windowsgui` (no console window), and copies `vJoyInterface.dll` next to the binary in `bin/`. Requires `windres` (mingw) on PATH and the vJoy SDK DLL present at repo root.

Plain `go build ./...` works for compile-checking but won't produce a fully resourced binary (no icon) and will pop a console window since `-H windowsgui` isn't set.

There are no automated tests in this repo currently.

Runtime requires: vJoy driver installed with device #1 configured, button box on COM3, handbrake on COM4 (hardcoded in `main.go`).

## Architecture

Data flow: **serial.Consumer** (per COM port) → device-specific **Handler** (parses the line) → **device.Device** chain → **vjoy.Device** (writes to the vJoy driver via `vjoy/vjoydll.go`, a manual DLL binding).

- `internal/serial`: generic line-oriented serial reader (`Consumer`), reconnects with backoff, drops the first second of readings after (re)connect as junk/startup noise. Takes a `Handler` interface (`Handle(reading string)`).
- `internal/buttonbox` and `internal/handbrake`: each implements `serial.Handler` for its board's line format and translates readings into `device.Device` calls (`SetButton`/`SetAxis`). `buttonbox/keyMap.go` maps physical wire IDs to `device.ButtonID`s; some button IDs are "switches" (on/off pairs) that get translated into a down+up pulse pair instead of a sustained state (see `switches` map in `buttonbox/handler.go`).
- `internal/device`: the `Device` interface (`SetButton`, `SetAxis`) plus decorators composed in `main.go`:
  - `NewPublishingDevice` — wraps a `Device`, also publishes every successful update to a `device.Updates` pub/sub bus (used by the UI to show live button/axis state).
  - `NewDebouncedDevice` — wraps a `Device`, enforces a minimum press duration per button by blocking the release call until the minimum duration has elapsed.
  - These are composed per-consumer in `main.go`: `outputDevice` (vjoy + publishing) is wrapped by a separate debouncer instance for each of the button box and handbrake consumers.
- `internal/vjoy`: `device.go` is the `Device` implementation; `vjoydll.go` is the raw `vJoyInterface.dll` binding (manual Windows DLL calls, not cgo). Handles the vJoy device acquire/relinquish lifecycle with retry/backoff for late-arriving driver state (`keepAcquired`).
- `internal/ui`: Win32 UI via `windigo` (no cgo, no wails despite git history mentioning it — see commit log, project migrated off wails). `ui.go` wires up the window, tray icon, quit button, and subscribes to `device.Updates` to paint live control state. `controls.go` defines the button/axis grid layout mirroring the physical button box (`buttonGrid` in `ui.go`).

## Notable behaviors / gotchas

- `main.go` has a hardcoded absolute log file path (`F:\dev\button-box-vjoy-feeder\button-box-vjoy-feeder.log`) — this only makes sense on the original dev machine, not portable.
- Single-instance enforcement via a named Windows mutex (`Local\button-box-vjoy-feeder`); second launch silently exits.
- App adds itself to `HKCU\Software\Microsoft\Windows\CurrentVersion\Run` on every startup (auto-start).
- COM ports (3 for button box, 4 for handbrake) and the vJoy device ID (1) are hardcoded in `main.go`, not configurable via flags/config file.
