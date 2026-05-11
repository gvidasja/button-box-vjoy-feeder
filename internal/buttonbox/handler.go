package buttonbox

import (
	"time"

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

			if err := device.SetButton(deviceButtonID, reading.state); err != nil {
				log.Errorf("failed to set button %v: %v; retrying once", deviceButtonID, err)
				time.Sleep(200 * time.Millisecond)
				if err2 := device.SetButton(deviceButtonID, reading.state); err2 != nil {
					log.Errorf("retry failed for button %v: %v", deviceButtonID, err2)
				}
			}

			if reading.state {
				producer.Produce("button", deviceButtonID)
			}
		}
	})
}
