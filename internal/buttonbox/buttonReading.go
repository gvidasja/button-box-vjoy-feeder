package buttonbox

import "strconv"

type buttonReading struct {
	id    readingID
	state bool
}

func parseButtonReading(serialString string) buttonReading {
	actionNumber, _ := strconv.ParseInt(serialString[0:1], 10, 64)
	button, _ := strconv.ParseInt(serialString[1:], 10, 64)
	return buttonReading{id: readingID(button), state: actionNumber > 0}
}
