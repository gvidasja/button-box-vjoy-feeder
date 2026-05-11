package buttonbox

import (
	"log/slog"
	"time"

	"github.com/gvidasja/button-box-vjoy-feeder/internal/device"
	"github.com/gvidasja/button-box-vjoy-feeder/internal/events"
	"github.com/gvidasja/button-box-vjoy-feeder/internal/serial"
)

func NewHandler(device device.Device, producer events.Producer) serial.Handler {
	return serial.HandlerFunc(func(message string) {

		reading := parseButtonReading(message)

		slog.Debug("button", "id", reading.buttonID, "state", reading.state)

		buttonID := reading.getButtonID()

		if deviceButtonID, ok := keyMap[buttonID]; ok {
			slog.Debug("sending", "from", buttonID, "to", deviceButtonID)

			if err := device.SetButton(deviceButtonID, reading.state); err != nil {
				slog.Error("failed to set button, retrying once", "button", deviceButtonID, "err", err)
				time.Sleep(200 * time.Millisecond)
				if err2 := device.SetButton(deviceButtonID, reading.state); err2 != nil {
					slog.Error("retry failed for button", "button", deviceButtonID, "err", err2)
				}
			}

			if reading.state {
				producer.Produce("button", deviceButtonID)
			}
		}
	})
}
