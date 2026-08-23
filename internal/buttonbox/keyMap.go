package buttonbox

import "github.com/gvidasja/button-box-vjoy-feeder/internal/device"

var keyMap = map[readingID]device.ButtonID{
	enc1Neg:    device.TC_DOWN,
	enc1Pos:    device.TC_UP,
	enc2Neg:    device.ABS_DOWN,
	enc2Pos:    device.ABS_UP,
	enc3Neg:    device.ENGINE_DOWN,
	enc3Pos:    device.ENGINE_UP,
	enc4Neg:    device.RECOVERY_DOWN,
	enc4Pos:    device.RECOVERY_UP,
	button1:    device.BUTTON_9_LIGHTS,
	button2:    device.BUTTON_10_LIGHTS_2,
	button3:    device.BUTTON_11,
	button4:    device.BUTTON_12,
	button5:    device.BUTTON_13_WIPERS,
	button6:    device.BUTTON_14,
	button7:    device.BUTTON_15_HORN,
	button8:    device.BUTTON_16,
	button9:    device.BUTTON_17_STARTED,
	button10:   device.BUTTON_18_HAZARDS,
	button11:   device.BUTTON_19,
	button12:   device.BUTTON_20,
	switch1Pos: device.IGNITION_ON,
	switch2Pos: device.SWITCH_2_ON,
	switch3Pos: device.SWITCH_3_ON,
	switch4Pos: device.SWITCH_4_ON,

	// switch1Pos: device.IGNITION_OFF,
	// switch2Pos: device.SWITCH_2_OFF,
	// switch3Pos: device.SWITCH_3_OFF,
	// switch4Pos: device.SWITCH_4_OFF,
}
