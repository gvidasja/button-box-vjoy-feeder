package device

import "sync"

type ButtonUpdate struct {
	Button ButtonID
	State  bool
}

type AxisUpdate struct {
	Axis  AxisID
	Value int32
}

type Updates struct {
	mu             sync.RWMutex
	buttonHandlers []func(ButtonUpdate)
	axisHandlers   []func(AxisUpdate)
}

type publishingDevice struct {
	device  Device
	updates *Updates
}

var _ Device = (*publishingDevice)(nil)

func NewPublishingDevice(device Device, updates *Updates) Device {
	return &publishingDevice{device: device, updates: updates}
}

func (d *publishingDevice) SetButton(button ButtonID, state bool) error {
	err := d.device.SetButton(button, state)
	if err == nil {
		d.updates.PublishButton(ButtonUpdate{Button: button, State: state})
	}
	return err
}

func (d *publishingDevice) SetAxis(axis AxisID, value int32) error {
	err := d.device.SetAxis(axis, value)
	if err == nil {
		d.updates.PublishAxis(AxisUpdate{Axis: axis, Value: value})
	}
	return err
}

func NewUpdates() *Updates {
	return &Updates{}
}

func (u *Updates) OnButton(handler func(ButtonUpdate)) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.buttonHandlers = append(u.buttonHandlers, handler)
}

func (u *Updates) OnAxis(handler func(AxisUpdate)) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.axisHandlers = append(u.axisHandlers, handler)
}

func (u *Updates) PublishButton(update ButtonUpdate) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	for _, handler := range u.buttonHandlers {
		handler(update)
	}
}

func (u *Updates) PublishAxis(update AxisUpdate) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	for _, handler := range u.axisHandlers {
		handler(update)
	}
}
