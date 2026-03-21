// Package dds provides a decoder for the DirectDraw surface format, which is compatible with the image package.
package dds_robroyd

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"

	"github.com/darklab8/fl-darkstat/darkmap/dds_robroyd/decoder"
	"github.com/darklab8/fl-darkstat/darkmap/dds_robroyd/header"
)

// init registers the decoder for the dds image format
func init() {
	image.RegisterFormat("dds", "DDS ", Decode, DecodeConfig)
}

var ErrUnsupported = errors.New("unsupported texture format")

func DecodeConfig(r io.Reader) (image.Config, error) {
	h, err := header.Read(r)
	if err != nil {
		return image.Config{}, err
	}

	// set width and height
	c := image.Config{
		Width:  int(h.Width),
		Height: int(h.Height),
	}

	switch pf, s := h.PixelFlags, h.RgbBitCount; {
	case pf.Is(header.DDPFFourCC):
		switch h.FourCCString {
		case header.FourCCDX10:
			err = ErrUnsupported
		case "DXT1", "DXT3", "DXT5":
			c.ColorModel = color.NRGBAModel
		default:
			err = fmt.Errorf("%w; is %s", ErrUnsupported, h.FourCCString)
		}

	case pf.Has(header.DDPFRGB): // because alpha is implicit
		err = ErrUnsupported

		if s <= 32 {
			c.ColorModel = color.NRGBAModel
		} else {
			c.ColorModel = color.NRGBA64Model
		}
	case pf.Is(header.DDPFYUV):
		err = ErrUnsupported

		c.ColorModel = color.NYCbCrAModel
	case pf.Is(header.DDPFLuminance):
		err = ErrUnsupported

		if s <= 8 {
			c.ColorModel = color.GrayModel
		} else {
			c.ColorModel = color.Gray16Model
		}
	case pf.Is(header.DDPFAlpha):
		err = ErrUnsupported

		if s <= 8 {
			c.ColorModel = color.AlphaModel
		} else {
			c.ColorModel = color.Alpha16Model
		}
	case pf.Is(header.DDPFLuminance | header.DDPFAlphaPixels):
		err = ErrUnsupported

		if s <= 32 {
			c.ColorModel = color.NRGBAModel // R__A
		} else {
			c.ColorModel = color.NRGBA64Model // R__A
		}
	default:
		err = fmt.Errorf("unrecognized image format: pf.flags: %x", pf)
	}

	return c, err
}

func Decode(r io.Reader) (image.Image, error) {
	h, err := header.Read(r)
	if err != nil {
		return nil, err
	}

	d, err := decoder.Find(h)
	if err != nil {
		return nil, err
	}

	return d.Decode(r)
}
