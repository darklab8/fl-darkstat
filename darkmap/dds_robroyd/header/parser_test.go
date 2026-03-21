package header

import (
	"bytes"
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	assert.EqualValues(t, sizeDDTF, binary.Size(deserializer{}))
}

func TestDeserializer_Read(t *testing.T) {
	var data = [132]byte{'D', 'D', 'S', ' '}
	for i := byte(1); i < 32; i++ {
		data[i*4] = i
	}

	// expected values during verification
	data[1*4] = 124
	data[2*4] = 7
	data[2*4+1] = 16
	data[19*4] = 32
	data[21*4] = '@'

	expected := &Header{
		DDSHeader: DDSHeader{
			TextureFlags:      Flags[DDSf]{4103},
			Height:            3,
			Width:             4,
			PitchOrLinearSize: 5,
			Depth:             6,
			MipMapCount:       7,
		},
		DDPFHeader: DDPFHeader{
			PixelFlags:  Flags[DDPFf]{20},
			FourCC:      64,
			RgbBitCount: 22,
			RBitMask:    23,
			GBitMask:    24,
			BBitMask:    25,
			ABitMask:    26,
		},
		CapsHeader: CapsHeader{
			Caps1: Flags[DDSCf]{27},
			Caps2: 28,
			Caps3: 29,
			Caps4: 30,
		},
		FourCCString: "@\x00\x00\x00",
	}

	rd := bytes.NewReader(data[:])
	h, err := Read(rd)
	assert.NoError(t, err)
	assert.Equal(t, expected, h)
}
