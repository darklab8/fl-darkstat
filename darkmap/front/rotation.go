package front

import (
	"math"
)

func toRad(deg float64) float64 {
	return deg * math.Pi / 180
}

func rotX(v [3]float64, a float64) [3]float64 {
	c, s := math.Cos(a), math.Sin(a)
	return [3]float64{v[0], c*v[1] - s*v[2], s*v[1] + c*v[2]}
}

func rotY(v [3]float64, a float64) [3]float64 {
	c, s := math.Cos(a), math.Sin(a)
	return [3]float64{c*v[0] + s*v[2], v[1], -s*v[0] + c*v[2]}
}

func rotZ(v [3]float64, a float64) [3]float64 {
	c, s := math.Cos(a), math.Sin(a)
	return [3]float64{c*v[0] - s*v[1], s*v[0] + c*v[1], v[2]}
}

func ObjRotAngle(rx, ry, rz float64) (angleDeg float64) {
	angle, _ := ProjectToNavMap(rx, ry, rz)
	return angle
}

func ObjRotLength(rx, ry, rz float64) float64 {
	_, length := ProjectToNavMap(rx, ry, rz)
	return length
}

// Freelancer: vector forward = [0, 0, -1] (NEG_Z_AXIS)
// Projection on a navmap (X/Z plane, Y-up)
func ProjectToNavMap(rx, ry, rz float64) (angleDeg float64, projLen float64) {
	v := [3]float64{0, 0, -1} // NEG_Z_AXIS — "forward"
	v = rotX(v, toRad(rx))    // pitch
	v = rotY(v, toRad(ry))    // yaw
	v = rotZ(v, toRad(rz))    // roll

	fx, fz := v[0], v[2]
	projLen = math.Sqrt(fx*fx + fz*fz)

	// atan2(fx, -fz): angle from "north" of a map (-Z direction)
	angleDeg = math.Atan2(fx, -fz) * 180 / math.Pi

	if projLen < 1e-9 {
		return angleDeg, 0 // you could yield error here, but we do not need it :)
	}

	return angleDeg, projLen
}
