package main

import (
	"time"

	"github.com/gvidasja/button-box-vjoy-feeder/vjoy"
)

type DeviceWithDelay struct {
	pressMap map[uint]time.Time
	device   *vjoy.Device
	debounce time.Duration
}

func NewDeviceWithDelay(id uint, debounceMillis int) *DeviceWithDelay {
	return &DeviceWithDelay{
		pressMap: make(map[uint]time.Time),
		device:   vjoy.NewDevice(id),
		debounce: time.Duration(debounceMillis) * time.Millisecond,
	}
}

func (d *DeviceWithDelay) Init() error {
	return d.device.Init()
}

func (d *DeviceWithDelay) Dispose() error {
	return d.Dispose()
}

func (d *DeviceWithDelay) SetButton(buttonID uint, state bool) error {
	now := time.Now()

	if state {
		d.pressMap[buttonID] = now
		return d.device.SetButton(buttonID, state)
	} else if lastPress, ok := d.pressMap[buttonID]; ok && now.Sub(lastPress) < d.debounce {
		return delay(d.debounce, func() error {
			return d.device.SetButton(buttonID, state)
		})
	} else {
		return d.device.SetButton(buttonID, state)
	}
}

func delay(duration time.Duration, f func() error) error {
	timer := time.NewTimer(duration)
	<-timer.C
	return f()
}
