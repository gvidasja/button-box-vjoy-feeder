package vjoy

import (
	"fmt"

	"github.com/gvidasja/button-box-vjoy-feeder/internal/device"
	log "github.com/sirupsen/logrus"
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

	if err := validateJoystick(d.id); err != nil {
		return fmt.Errorf("invalid Joystick: %w", err)
	}

	if err := acquireVJD(d.id); err != nil {
		return fmt.Errorf("cannot acquire VJD %w", err)
	}

	return nil
}

func (d *vjoyDevice) Stop() {
	err := relinquishVJD(d.id)

	if err != nil {
		log.Errorf("could not relinquish VJD: %v", err)
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
