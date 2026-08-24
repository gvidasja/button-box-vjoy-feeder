package buttonbox

import (
	"log/slog"

	"github.com/gvidasja/button-box-vjoy-feeder/internal/device"
	"github.com/gvidasja/button-box-vjoy-feeder/internal/serial"
)

func NewHandler(d device.Device) serial.Handler {
	send := func(buttonID device.ButtonID, state bool) {
		_ = d.SetButton(buttonID, state)
	}

	return serial.HandlerFunc(func(message string) {
		reading := parseButtonReading(message)

		slog.Debug("reading", "id", reading.id, "state", reading.state)

		if deviceButtonID, ok := keyMap[reading.id]; ok {
			if sw, ok := switches[deviceButtonID]; ok {
				send(sw[reading.state], true)
				send(sw[reading.state], false)
			} else {
				send(deviceButtonID, reading.state)
			}
		}
	})
}

var switches = map[device.ButtonID]map[bool]device.ButtonID{
	device.IGNITION_ON: {false: device.IGNITION_ON, true: device.IGNITION_OFF},
	device.SWITCH_2_ON: {false: device.SWITCH_2_ON, true: device.SWITCH_2_OFF},
	device.SWITCH_3_ON: {true: device.SWITCH_3_ON, false: device.SWITCH_3_OFF},
	device.SWITCH_4_ON: {true: device.SWITCH_4_OFF, false: device.SWITCH_4_ON},
}
