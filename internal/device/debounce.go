package device

import (
	"log/slog"
	"time"
)

type debouncedDevice struct {
	pressMap map[ButtonID]time.Time
	device   Device
	cfg      DebouncedDeviceConfig
}

type Device interface {
	SetButton(ButtonID, bool) error
	SetAxis(AxisID, int32) error
}

type DebouncedDeviceConfig struct {
	MinimumButtonPressDuration time.Duration
}

func NewDebouncedDevice(device Device, cfg DebouncedDeviceConfig) *debouncedDevice {
	return &debouncedDevice{
		pressMap: make(map[ButtonID]time.Time),
		device:   device,
		cfg:      cfg,
	}
}

func (d *debouncedDevice) SetButton(buttonID ButtonID, state bool) error {
	now := time.Now()

	slog.Debug("set button", "button", buttonID, "state", state, "now", now)

	if state {
		d.pressMap[buttonID] = now
	} else {
		lastPositiveStateAt, ok := d.pressMap[buttonID]

		earliestReleaseTime := lastPositiveStateAt.Add(d.cfg.MinimumButtonPressDuration)

		if ok && !now.After(earliestReleaseTime) {
			<-time.After(earliestReleaseTime.Sub(now))
		}
	}

	return d.device.SetButton(buttonID, state)
}

func (d *debouncedDevice) SetAxis(axisID AxisID, value int32) error {
	return d.device.SetAxis(axisID, value)
}
