package buttons

type buttonReading struct {
	buttonID buttonID
	state    bool
}

func (reading buttonReading) getButtonID() buttonID {
	if !reading.state {
		switch reading.buttonID {
		case switch1Pos:
			return switch1Neg
		case switch2Pos:
			return switch2Neg
		case switch3Pos:
			return switch3Neg
		case switch4Pos:
			return switch4Neg
		default:
			return reading.buttonID
		}
	}

	return reading.buttonID
}
