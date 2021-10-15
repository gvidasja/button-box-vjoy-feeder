package device

import (
	"time"

	"github.com/gvidasja/button-box-vjoy-feeder/internal/vjoy"
)

type deviceWithDelay struct {
	pressMap map[uint]time.Time
	device   vjoy.Device
	debounce time.Duration
}

func NewDeviceWithDelay(id uint, debounceMillis int) *deviceWithDelay {
	return &deviceWithDelay{
		pressMap: make(map[uint]time.Time),
		device:   vjoy.NewDevice(id),
		debounce: time.Duration(debounceMillis) * time.Millisecond,
	}
}

func (d *deviceWithDelay) Init() error {
	return d.device.Init()
}

func (d *deviceWithDelay) Dispose() error {
	return d.device.Dispose()
}

func (d *deviceWithDelay) SetButton(buttonID uint, state bool) error {
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
