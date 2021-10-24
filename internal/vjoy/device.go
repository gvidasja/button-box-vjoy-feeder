package vjoy

import (
	"fmt"

	"github.com/gvidasja/button-box-vjoy-feeder/internal/device"
)

type vjoyDevice struct {
	id uint
}

func NewDevice(id uint) *vjoyDevice {
	return &vjoyDevice{id}
}

func (d *vjoyDevice) Dispose() error {
	return relinquishVJD(d.id)
}

func (d *vjoyDevice) SetButton(buttonID device.ButtonID, state bool) error {
	return setButton(d.id, buttonID, state)
}

func (d *vjoyDevice) Init() error {
	load()

	err := validateJoystick(d.id)

	if err != nil {
		return err
	}

	return acquireVJD(d.id)
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
