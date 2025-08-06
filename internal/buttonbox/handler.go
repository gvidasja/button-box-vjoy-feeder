package buttonbox

import (
	"github.com/gvidasja/button-box-vjoy-feeder/internal/device"
	"github.com/gvidasja/button-box-vjoy-feeder/internal/events"
	"github.com/gvidasja/button-box-vjoy-feeder/internal/serial"
	log "github.com/sirupsen/logrus"
)

func NewHandler(device device.Device, producer events.Producer) serial.Handler {
	return serial.HandlerFunc(func(message string) {

		reading := parseButtonReading(message)

		log.Debugf("button %v: %v", reading.buttonID, reading.state)

		buttonID := reading.getButtonID()

		if deviceButtonID, ok := keyMap[buttonID]; ok {
			log.Debugf("sending %v -> %v", buttonID, deviceButtonID)
			device.SetButton(deviceButtonID, reading.state)

			if reading.state {
				producer.Produce("button", deviceButtonID)
			}
		}
	})
}
