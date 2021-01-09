package vjoy

import (
	"fmt"
)

type Device struct {
	ID uint
}

func NewDevice(id uint) *Device {
	return &Device{id}
}

func (d *Device) Dispose() error {
	return relinquishVJD(d.ID)
}

func (d *Device) SetButton(buttonID uint, state bool) error {
	return setButton(d.ID, buttonID, state)
}

func (d *Device) Init() error {
	load()

	err := validateJoystick(d.ID)

	if err != nil {
		return err
	}

	return acquireVJD(d.ID)
}

func validateJoystick(deviceID uint) error {
	if !vJoyEnabled() {
		return fmt.Errorf("vJoy is not enabled")
	}

	switch status := getVJDStatus(deviceID); status {
	case VJD_STAT_OWN, VJD_STAT_FREE:
		break
	case VJD_STAT_BUSY:
		return fmt.Errorf("Device %d is busy", deviceID)
	case VJD_STAT_MISS:
		return fmt.Errorf("Device %d not found", deviceID)
	case VJD_STAT_UNKN:
	default:
		return fmt.Errorf("Unknown error with device %d", deviceID)
	}

	return nil
}
