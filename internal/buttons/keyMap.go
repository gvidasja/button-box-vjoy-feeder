package buttons

import "github.com/gvidasja/button-box-vjoy-feeder/internal/device"

var keyMap = map[buttonID]device.ButtonID{
	enc1Neg:    device.Enc1Neg,
	enc1Pos:    device.Enc1Pos,
	enc2Neg:    device.Enc2Neg,
	enc2Pos:    device.Enc2Pos,
	enc3Neg:    device.Enc3Neg,
	enc3Pos:    device.Enc3Pos,
	enc4Neg:    device.Enc4Neg,
	enc4Pos:    device.Enc4Pos,
	button1:    device.Button1,
	button2:    device.Button2,
	button3:    device.Button3,
	button4:    device.Button4,
	button5:    device.Button5,
	button6:    device.Button6,
	button7:    device.Button7,
	button8:    device.Button8,
	button9:    device.Button9,
	button10:   device.Button10,
	button11:   device.Button11,
	button12:   device.Button12,
	switch1Neg: device.Switch1Neg,
	switch1Pos: device.Switch1Pos,
	switch2Neg: device.Switch2Neg,
	switch2Pos: device.Switch2Pos,
	switch3Neg: device.Switch3Neg,
	switch3Pos: device.Switch3Pos,
	switch4Neg: device.Switch4Neg,
	switch4Pos: device.Switch4Pos,
}
