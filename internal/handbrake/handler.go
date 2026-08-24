package handbrake

import (
	"log/slog"
	"math"
	"strconv"

	"github.com/gvidasja/button-box-vjoy-feeder/internal/device"
	"github.com/gvidasja/button-box-vjoy-feeder/internal/serial"
)

const (
	handbrakeMin = (100)
	handbrakeMax = (1024)
	vjoyMin      = (0)
	vjoyMax      = (math.MaxInt16)

	axisID = 0x32
)

var previousState = int64(vjoyMin)

func NewHandler(device device.Device) serial.Handler {
	return serial.HandlerFunc(func(data string) {
		state, _ := strconv.ParseFloat(data, 64)

		slog.Debug("handbrake", "state", state)

		scaledState := int64(vjoyMin + (vjoyMax-vjoyMin)*(state-handbrakeMin)/(handbrakeMax-handbrakeMin))

		if scaledState < vjoyMin {
			scaledState = vjoyMin
		} else if scaledState > vjoyMax {
			scaledState = vjoyMax
		}

		if math.Abs(float64(previousState-scaledState)/float64(vjoyMax-vjoyMin)) > 0.2 {
			slog.Debug("skipping handbrake", "from", previousState, "to", scaledState)
			return
		}

		previousState = scaledState

		slog.Debug("sending handbrake", "state", state, "scaled", scaledState)

		device.SetAxis(axisID, int32(scaledState))
	})
}
