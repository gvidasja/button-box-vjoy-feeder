package buttonbox

type readingID uint

const (
	enc1Pos  = readingID(32)
	enc1Neg  = readingID(31)
	enc2Pos  = readingID(30)
	enc2Neg  = readingID(29)
	enc3Pos  = readingID(27)
	enc3Neg  = readingID(28)
	enc4Pos  = readingID(25)
	enc4Neg  = readingID(26)
	button1  = readingID(3)
	button2  = readingID(8)
	button3  = readingID(13)
	button4  = readingID(18)
	button5  = readingID(23)
	button6  = readingID(22)
	button7  = readingID(21)
	button8  = readingID(0)
	button9  = readingID(5)
	button10 = readingID(10)
	button11 = readingID(15)
	button12 = readingID(20)
	// switch1Neg = readingID(4)
	// switch2Neg = readingID(14)
	// switch3Neg = readingID(7)
	// switch4Neg = readingID(17)

	// don't exist anymore - button box reports on/off state for single button ID in switches
	switch1Pos = readingID(2)
	switch2Pos = readingID(12)
	switch3Pos = readingID(6)
	switch4Pos = readingID(16)
)
