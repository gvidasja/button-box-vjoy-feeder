package device

type (
	ButtonID uint
	AxisID   uint
)

const (
	TC_DOWN            = ButtonID(1)
	TC_UP              = ButtonID(2)
	ABS_DOWN           = ButtonID(3)
	ABS_UP             = ButtonID(4)
	ENGINE_DOWN        = ButtonID(5)
	ENGINE_UP          = ButtonID(6)
	RECOVERY_DOWN      = ButtonID(7)
	RECOVERY_UP        = ButtonID(8)
	BUTTON_9_LIGHTS    = ButtonID(9)
	BUTTON_10_LIGHTS_2 = ButtonID(10)
	BUTTON_11          = ButtonID(11)
	BUTTON_12          = ButtonID(12)
	BUTTON_13_WIPERS   = ButtonID(13)
	BUTTON_14          = ButtonID(14)
	BUTTON_15_HORN     = ButtonID(15)
	BUTTON_16          = ButtonID(16)
	BUTTON_17_STARTED  = ButtonID(17)
	BUTTON_18_HAZARDS  = ButtonID(18)
	BUTTON_19          = ButtonID(19)
	BUTTON_20          = ButtonID(20)
	IGNITION_ON        = ButtonID(21)
	IGNITION_OFF       = ButtonID(22)
	SWITCH_2_ON        = ButtonID(23)
	SWITCH_2_OFF       = ButtonID(24)
	SWITCH_3_ON        = ButtonID(25)
	SWITCH_3_OFF       = ButtonID(26)
	SWITCH_4_ON        = ButtonID(27)
	SWITCH_4_OFF       = ButtonID(28)
)
