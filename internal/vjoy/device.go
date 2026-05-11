package vjoy

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gvidasja/button-box-vjoy-feeder/internal/device"
)

type vjoyDevice struct {
	id uint
}

var _ device.Device = (*vjoyDevice)(nil)

func NewDevice(id uint) *vjoyDevice {
	return &vjoyDevice{id}
}

func (d *vjoyDevice) Start() error {
	if err := loadVJoyDLL(); err != nil {
		return fmt.Errorf("cannot load vJoy DLL: %w", err)
	}

	var lastErr error
	// try a few times to handle transient busy/missing states
	for i := 0; i < 5; i++ {
		if err := validateJoystick(d.id); err != nil {
			lastErr = fmt.Errorf("invalid Joystick: %w", err)
			time.Sleep(time.Duration(200*(i+1)) * time.Millisecond)
			continue
		}

		if err := acquireVJD(d.id); err != nil {
			lastErr = fmt.Errorf("cannot acquire VJD: %w", err)
			time.Sleep(time.Duration(200*(i+1)) * time.Millisecond)
			continue
		}

		return nil
	}

	return lastErr
}

func (d *vjoyDevice) Stop() {
	err := relinquishVJD(d.id)

	if err != nil {
		slog.Error("could not relinquish VJD", "err", err)
	}
}

func (d *vjoyDevice) SetButton(buttonID device.ButtonID, state bool) error {
	return setButton(d.id, buttonID, state)
}

func (d *vjoyDevice) SetAxis(axisID device.AxisID, value int32) error {
	return setAxis(d.id, axisID, value)
}

func validateJoystick(deviceID uint) error {
	if !vJoyEnabled() {
		return fmt.Errorf("vJoy is not enabled")
	}

	switch status := getVJDStatus(deviceID); status {
	case VJD_STAT_OWN, VJD_STAT_FREE:
		break
	case VJD_STAT_BUSY:
		return fmt.Errorf("device %d is busy", deviceID)
	case VJD_STAT_MISS:
		return fmt.Errorf("device %d not found", deviceID)
	case VJD_STAT_UNKN:
	default:
		return fmt.Errorf("unknown error with device %d", deviceID)
	}

	return nil
}
