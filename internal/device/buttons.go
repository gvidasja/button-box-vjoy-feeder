package device

type ButtonID uint

const (
	Enc1Neg    = ButtonID(1)
	Enc1Pos    = ButtonID(2)
	Enc2Neg    = ButtonID(3)
	Enc2Pos    = ButtonID(4)
	Enc3Neg    = ButtonID(5)
	Enc3Pos    = ButtonID(6)
	Enc4Neg    = ButtonID(7)
	Enc4Pos    = ButtonID(8)
	Button1    = ButtonID(9)
	Button2    = ButtonID(10)
	Button3    = ButtonID(11)
	Button4    = ButtonID(12)
	Button5    = ButtonID(13)
	Button6    = ButtonID(14)
	Button7    = ButtonID(15)
	Button8    = ButtonID(16)
	Button9    = ButtonID(17)
	Button10   = ButtonID(18)
	Button11   = ButtonID(19)
	Button12   = ButtonID(20)
	Switch1Neg = ButtonID(21)
	Switch1Pos = ButtonID(22)
	Switch2Neg = ButtonID(23)
	Switch2Pos = ButtonID(24)
	Switch3Neg = ButtonID(25)
	Switch3Pos = ButtonID(26)
	Switch4Neg = ButtonID(27)
	Switch4Pos = ButtonID(28)
)

var switches = map[ButtonID]bool{
	Switch1Neg: true,
	Switch1Pos: true,
	Switch2Neg: true,
	Switch2Pos: true,
	Switch3Neg: true,
	Switch3Pos: true,
	Switch4Neg: true,
	Switch4Pos: true,
}
