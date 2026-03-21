package internal

import (
	"github.com/stretchr/testify/assert"
	"image/color"
	"testing"
)

func TestInterpolateColors(t *testing.T) {
	var tests = map[string]struct {
		in1 [2]byte
		in2 [2]byte
		out [4]color.NRGBA
	}{
		"with alpha": {
			in1: [2]byte{0b00100001, 0b00001000},
			in2: [2]byte{0b00010001, 0b10001100},
			out: [4]color.NRGBA{
				{R: 8, G: 4, B: 8, A: 255},
				{R: 136, G: 128, B: 136, A: 255},
				{R: 72, G: 66, B: 72, A: 255},
				{R: 0, G: 0, B: 0, A: 0},
			},
		},
		"without alpha": {
			in1: [2]byte{0b00010000, 0b10000100},
			in2: [2]byte{0b00100001, 0b00001000},
			out: [4]color.NRGBA{
				{R: 128, G: 128, B: 128, A: 255},
				{R: 8, G: 4, B: 8, A: 255},
				{R: 88, G: 87, B: 88, A: 255},
				{R: 48, G: 45, B: 48, A: 255},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			colors := InterpolateColors(test.in1[:], test.in2[:])
			assert.Equal(t, test.out[0], colors[0])
			assert.Equal(t, test.out[1], colors[1])
			assert.Equal(t, test.out[2], colors[2])
			assert.Equal(t, test.out[3], colors[3])
		})
	}
}
