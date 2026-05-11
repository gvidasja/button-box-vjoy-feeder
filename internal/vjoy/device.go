package vjoy

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gvidasja/button-box-vjoy-feeder/internal/device"
)

type vjoyDevice struct {
	id   uint
	done chan struct{}
}

var _ device.Device = (*vjoyDevice)(nil)

func NewDevice(id uint) *vjoyDevice {
	return &vjoyDevice{id: id, done: make(chan struct{})}
}

func (d *vjoyDevice) Start() error {
	if err := loadVJoyDLL(); err != nil {
		slog.Warn("vJoy DLL not available yet, will retry in background", "err", err)
	}
	go d.keepAcquired()
	return nil
}

func (d *vjoyDevice) keepAcquired() {
	backoff := time.Second
	for {
		if err := validateJoystick(d.id); err != nil {
			slog.Warn("vjoy validate failed, retrying", "device", d.id, "err", err, "backoff", backoff)
			select {
			case <-d.done:
				return
			case <-time.After(backoff):
				backoff = min(backoff*2, 30*time.Second)
				continue
			}
		}

		if err := acquireVJD(d.id); err != nil {
			slog.Warn("vjoy acquire failed, retrying", "device", d.id, "err", err, "backoff", backoff)
			select {
			case <-d.done:
				return
			case <-time.After(backoff):
				backoff = min(backoff*2, 30*time.Second)
				continue
			}
		}

		slog.Info("vjoy device acquired", "device", d.id)
		return
	}
}

func (d *vjoyDevice) Stop() {
	select {
	case <-d.done:
	default:
		close(d.done)
	}

	err := relinquishVJD(d.id)
	if err != nil {
		slog.Error("could not relinquish VJD", "device", d.id, "err", err)
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
		return nil
	case VJD_STAT_BUSY:
		return fmt.Errorf("device %d is busy", deviceID)
	case VJD_STAT_MISS:
		return fmt.Errorf("device %d not found", deviceID)
	case VJD_STAT_UNKN:
		fallthrough
	default:
		return fmt.Errorf("unknown error with device %d", deviceID)
	}
}
