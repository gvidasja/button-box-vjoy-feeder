package device

import (
	"time"
)

type deviceWithDelay struct {
	pressMap map[ButtonID]time.Time
	device   Device
	cfg      DeviceConfig
}

type Device interface {
	SetButton(ButtonID, bool) error
}

type DeviceConfig struct {
	MinimumButtonPressDuration time.Duration
}

func New(device Device, cfg DeviceConfig) *deviceWithDelay {
	return &deviceWithDelay{
		pressMap: make(map[ButtonID]time.Time),
		device:   device,
		cfg:      cfg,
	}
}

func (d *deviceWithDelay) SetButton(buttonID ButtonID, state bool) error {
	now := time.Now()

	if switches[buttonID] {
		err := d.device.SetButton(buttonID, true)

		if err != nil {
			return err
		}

		<-time.NewTimer(d.cfg.MinimumButtonPressDuration).C
		return d.device.SetButton(buttonID, false)
	}

	if state {
		d.pressMap[buttonID] = now
	}

	if lastPress, ok := d.pressMap[buttonID]; !state && ok && now.Before(lastPress.Add(d.cfg.MinimumButtonPressDuration)) {
		<-time.NewTimer(d.cfg.MinimumButtonPressDuration).C
	}

	return d.device.SetButton(buttonID, state)
}
