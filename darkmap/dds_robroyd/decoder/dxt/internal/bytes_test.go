package internal

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWeighted(t *testing.T) {
	var tests = map[int]struct {
		w0, w1 float64
		b0, b1 byte
		res    byte
	}{
		1: {w0: 1, b0: 1, w1: 1, b1: 1, res: 1},
		2: {w0: 1, b0: 1, w1: 1, b1: 0, res: 1},
		3: {w0: 1, b0: 0, w1: 1, b1: 10, res: 5},
		4: {w0: 10, b0: 0, w1: 1, b1: 10, res: 1},
		5: {w0: 1, b0: 0, w1: 10, b1: 10, res: 9},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			result := Weighted(test.w0, test.b0, test.w1, test.b1)
			assert.Equal(t, test.res, result)
		})
	}
}

func TestExtract(t *testing.T) {
	var in = []byte{
		0b0010_0001, 0b1000_0100, //  33, 132	|  2, 1   8, 4
		0b0011_1001, 0b1100_0110, //  57, 198	|  3, 9   8, 6
		0b1101_1110, 0b0111_1011, // 222, 123	| 13,14   7,11
		0b0000_1111, 0b1111_0000, //  15, 240	|  0,15	 15, 0
	}

	t.Run("in 4s", func(t *testing.T) {
		var out4 = []byte{
			0b0001,
			0b0010,
			0b0100,
			0b1000,
			0b1001,
			0b0011,
			0b0110,
			0b1100,
			0b1110,
			0b1101,
			0b1011,
			0b0111,
			0b1111,
			0b0000,
			0b0000,
			0b1111,
		}

		for i := byte(0); i < 16; i++ {
			t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
				assert.Equal(t, out4[i], ExtractIndex(in, i, 4))
			})
		}
	})

	t.Run("in 3s", func(t *testing.T) {
		var out3 = []byte{
			0b001,
			0b100,
			0b000,
			0b010,
			0b000,
			0b011,
			0b110,
			0b001,
			0b110,
			0b000,
			0b011,
			0b111,
			0b101,
			0b111,
			0b110,
			0b011,
			0b111,
			0b001,
			0b000,
			0b000,
			0b111,
		}

		for i := byte(0); i < 21; i++ {
			t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
				assert.Equal(t, out3[i], ExtractIndex(in, i, 3))
			})
		}
	})

	t.Run("in 6s", func(t *testing.T) {
		var out6 = []byte{
			0b100001,
			0b010000,
			0b011000,
			0b001110,
			0b000110,
			0b111011,
			0b111101,
			0b011110,
			0b001111,
			0b000000,
		}

		for i := byte(0); i < 10; i++ {
			t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
				assert.Equal(t, out6[i], ExtractIndex(in, i, 6))
			})
		}
	})
}
