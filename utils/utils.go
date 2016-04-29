package utils

func Sgn(x int) int {
	if x > 0 {
		return 1
	}
	if x < 0 {
		return -1
	}
	return 0
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
