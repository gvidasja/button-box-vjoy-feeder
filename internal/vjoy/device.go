package vjoy

import (
	"fmt"
)

type Device interface {
	Init() error
	Dispose() error
	SetButton(buttonID uint, state bool) error
}

type device struct {
	id uint
}

func NewDevice(id uint) *device {
	return &device{id}
}

func (d *device) Dispose() error {
	return relinquishVJD(d.id)
}

func (d *device) SetButton(buttonID uint, state bool) error {
	return setButton(d.id, buttonID, state)
}

func (d *device) Init() error {
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
