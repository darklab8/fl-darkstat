package uncompressed

import (
	"image"
	"io"

	"github.com/darklab8/fl-darkstat/darkmap/dds_robroyd/header"
)

type Decoder struct {
	flags  header.Flags[header.DDPFf]
	bounds image.Point
}

func New(header *header.Header) *Decoder {
	return &Decoder{
		flags:  header.PixelFlags,
		bounds: image.Pt(int(header.Width), int(header.Height)),
	}
}

func (d *Decoder) Decode(r io.Reader) (image.Image, error) {
	rgba := image.NewNRGBA(image.Rectangle{Max: d.bounds})
	if rgba.Rect.Empty() {
		return rgba, nil
	}

	switch d.flags.F {
	case header.DDPFAlphaPixels | header.DDPFRGB:
		for y := 0; y < d.bounds.Y; y++ {
			p := rgba.Pix[y*rgba.Stride : y*rgba.Stride+d.bounds.X*4]
			if _, err := io.ReadFull(r, p); err != nil {
				return nil, err
			}
			// BGRA to RGBA re-order.
			for i := 0; i < len(p); i += 4 {
				p[i+0], p[i+2] = p[i+2], p[i+0]
			}
		}
	case header.DDPFRGB:
		b := make([]byte, 3*d.bounds.X)
		for y := 0; y < d.bounds.Y; y++ {
			if _, err := io.ReadFull(r, b); err != nil {
				return nil, err
			}
			p := rgba.Pix[y*rgba.Stride : y*rgba.Stride+d.bounds.X*4]
			// BGRA to RGBA re-order.
			for i, j := 0, 0; i < len(p); i, j = i+4, j+3 {
				p[i+0] = b[j+2]
				p[i+1] = b[j+1]
				p[i+2] = b[j+0]
				p[i+3] = 0xFF
			}
		}
	}

	return rgba, nil
}
