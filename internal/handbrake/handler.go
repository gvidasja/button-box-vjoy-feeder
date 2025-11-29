package handbrake

import (
	"math"
	"strconv"

	"github.com/gvidasja/button-box-vjoy-feeder/internal/device"
	"github.com/gvidasja/button-box-vjoy-feeder/internal/events"
	"github.com/gvidasja/button-box-vjoy-feeder/internal/serial"
	log "github.com/sirupsen/logrus"
)

const (
	handbrakeMin = (100)
	handbrakeMax = (1024)
	vjoyMin      = (0)
	vjoyMax      = (math.MaxInt32)

	axisID = 0x32
)

var previousState = int64(vjoyMin)

func NewHandler(device device.Device, producer events.Producer) serial.Handler {
	return serial.HandlerFunc(func(data string) {
		state, _ := strconv.ParseFloat(data, 64)

		log.Debugf("handbrake %v", state)

		scaledState := int64(vjoyMin + (vjoyMax-vjoyMin)*(state-handbrakeMin)/(handbrakeMax-handbrakeMin))

		if scaledState < vjoyMin {
			scaledState = vjoyMin
		} else if scaledState > vjoyMax {
			scaledState = vjoyMax
		}

		if math.Abs(float64(previousState-scaledState)/float64(vjoyMax-vjoyMin)) > 0.2 {
			log.Debugf("skipping %v -> %v", previousState, scaledState)
			return
		}

		previousState = scaledState

		log.Debugf("sending %v -> %v", state, scaledState)
		device.SetAxis(axisID, int32(scaledState))
		producer.Produce("handbrake", map[string]any{
			"min":   vjoyMin,
			"max":   vjoyMax,
			"state": scaledState,
		})
	})
}
