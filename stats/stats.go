package stats

var (
	searchCounter int
	fixupCounter  int
	extraCounter  int
	rotateCounter int
)

func Reset() {
	searchCounter = 0
	fixupCounter = 0
	extraCounter = 0
	rotateCounter = 0
}

func AddSearchCounter(d int) {
	searchCounter += d
}

func AddFixupCounter(d int) {
	fixupCounter += d
}

func AddExtraCounter(d int) {
	extraCounter += d
}

func AddRotateCounter(d int) {
	rotateCounter += d
}

func GetSearchCounter() int {
	return searchCounter
}

func GetFixupCounter() int {
	return fixupCounter
}

func GetExtraCounter() int {
	return extraCounter
}

func GetRotateCounter() int {
	return rotateCounter
}
