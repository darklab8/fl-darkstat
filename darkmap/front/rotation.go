package front

import "math"

func rotateForward(rx, ry, rz float64) float64 {
	xFlipped := math.Abs(math.Abs(rx)-180) < 0.01
	zFlipped := math.Abs(math.Abs(rz)-180) < 0.01

	if xFlipped && zFlipped {
		return ry // double flip cancels out: -(-ry)
	}
	return -ry
}
