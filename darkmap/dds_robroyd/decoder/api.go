package decoder

import (
	"fmt"
	"image"
	"io"

	"github.com/darklab8/fl-darkstat/darkmap/dds_robroyd/decoder/dxt"
	"github.com/darklab8/fl-darkstat/darkmap/dds_robroyd/decoder/uncompressed"
	"github.com/darklab8/fl-darkstat/darkmap/dds_robroyd/header"
)

// Decoder is the default interface for actual decoding operations.
type Decoder interface {
	// Decode takes the header-less reader and tries to read an parse the image-data from it.
	Decode(io.Reader) (image.Image, error)
}

// Find takes a parsed header.Header and tries to find a fitting Decoder or returns an error.
func Find(h *header.Header) (d Decoder, err error) {
	if h.PixelFlags.Is(header.DDPFFourCC) {
		switch h.FourCC {
		case 0:
			if !h.PixelFlags.Has(header.DDPFRGB) {
				err = fmt.Errorf("unsupported pixel format %x", h.PixelFlags)
			} else {
				d = uncompressed.New(h)
			}

		default:
			switch h.FourCCString {
			case "DXT1", "DXT2", "DXT3", "DXT4", "DXT5":
				d, err = dxt.New(h.FourCCString, int(h.Width), int(h.Height))

			default:
				err = fmt.Errorf("texture with compression '%v' is unsupported", h.FourCC)
			}
		}
	} else {
		err = fmt.Errorf("texture without compression '%v' is unsupported", h.PixelFlags)
	}

	return
}
