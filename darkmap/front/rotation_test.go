package front

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRotation1(t *testing.T) {
	assert.Equal(t, -89.0, ObjRotAngle(0, 89, 0))
}

func TestRotation2(t *testing.T) {
	Printdetails(-180, 27, 180)
	assert.Equal(t, -153.0, ObjRotAngle(-180, 27, 180))

	fmt.Println(ObjRotAngle(-180, 27.1, 180))
}

func TestRotation_zone_bw08_corridor_08(t *testing.T) {
	x := -180.0
	y := -58.1
	z := 180.0

	// -58.1 degres or its opposite is ok. 360-58.1

	Printdetails(x, y, z)
	// assert.Equal(t, 153.0, ProjectRotationAngle(-180, -58.1, 180))
	// fmt.Println(ProjectRotationAngle(x, x, z))
}

func TestRotation_zone_bw08_path_img6_1(t *testing.T) {
	x := 90.0
	y := -71.6
	z := 0.0

	// -58.1 degres or its opposite is ok. 360-58.1

	Printdetails(x, y, z)
	// assert.Equal(t, 153.0, ProjectRotationAngle(-180, -58.1, 180))
	// fmt.Println(ProjectRotationAngle(x, x, z))
}

func TestRotation_zone_bw08_corridor_10(t *testing.T) {
	x := 0.0
	y := 23.4
	z := 0.0

	Printdetails(x, y, z)

	// should be -23.4
	// or 360-24.4
}

func Printdetails(rx float64, ry float64, rz float64) {
	angle, projLen := ProjectToNavMap(rx, ry, rz)

	fmt.Printf("input:  X=%.1f°  Y=%.1f°  Z=%.1f°\n", rx, ry, rz)
	fmt.Printf("Projected angle onto XY: %.2f°\n", angle)
	fmt.Printf("Projection length:      %.4f\n", projLen)
}
