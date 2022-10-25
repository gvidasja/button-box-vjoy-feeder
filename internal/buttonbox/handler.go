package buttonbox

import (
	"github.com/gvidasja/button-box-vjoy-feeder/internal/device"
	"github.com/gvidasja/button-box-vjoy-feeder/internal/serial"
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	device device.Device
}

var _ serial.Handler = (*Handler)(nil)

func NewHandler(d device.Device) *Handler {
	return &Handler{
		device: d,
	}
}

func (h *Handler) Handle(message string) {
	reading := parseButtonReading(message)

	log.Debugf("button %v: %v", reading.buttonID, reading.state)

	buttonID := reading.getButtonID()

	if deviceButtonID, ok := keyMap[buttonID]; ok {
		log.Debugf("sending %v -> %v", buttonID, deviceButtonID)
		h.device.SetButton(deviceButtonID, reading.state)
	}
}
