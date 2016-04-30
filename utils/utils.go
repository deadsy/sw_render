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

func Round(x float32) int {
	if x < -0.5 {
		return int(x - 0.5)
	}
	if x > 0.5 {
		return int(x + 0.5)
	}
	return 0
}
