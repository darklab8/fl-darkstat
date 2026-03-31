package front

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRotation1(t *testing.T) {
	assert.Equal(t, rotateForward(0, 89, 0), -89.0)
}

func TestRotation2(t *testing.T) {
	assert.Equal(t, 27.0, rotateForward(-180, 27, 180))

	fmt.Println(rotateForward(-180, 27.1, 180))
}
