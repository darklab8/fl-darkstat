package utfextract

import (
	"bytes"
	"image"
	"image/png"

	"github.com/darklab8/fl-darkstat/darkmap/dds"
	"github.com/darklab8/fl-darkstat/darkmap/tga"
)

func TransformToJpeg(img *Image) (*bytes.Buffer, error) {
	input := bytes.NewReader(img.Data)

	if img.Extension == "tga" {
		img, err := tga.Decode(input)
		if err != nil {
			return nil, err
		}
		var output *bytes.Buffer = &bytes.Buffer{} // zero value is ready to use

		err = png.Encode(output, img)
		if err != nil {
			return nil, err
		}
		return output, nil

	} else if img.Extension == "dds" {
		// var input *bytes.Buffer = bytes.NewBuffer(img.Data)
		// img, _, err := image.Decode(input)
		img, err := dds.Decode(input, true)
		if err != nil {
			return nil, err
		}

		nrgba, ok := img.Image.(*image.NRGBA)
		if ok {
			img.Image = CompositeOverWhite(nrgba)
		}

		var output *bytes.Buffer = &bytes.Buffer{} // zero value is ready to use

		err = png.Encode(output, img)
		if err != nil {
			return nil, err
		}
		return output, nil
	}

	panic("not supported extension to transform to jpeg")
}

func CompositeOverWhite(img *image.NRGBA) *image.NRGBA {
	out := image.NewNRGBA(img.Bounds())
	for i := 0; i < len(img.Pix); i += 4 {
		a := uint32(img.Pix[i+3])
		invA := 255 - a
		// Blend: (pixel*alpha + 255*invAlpha + 127) / 255
		out.Pix[i+0] = uint8((uint32(img.Pix[i+0])*a + 255*invA + 127) / 255)
		out.Pix[i+1] = uint8((uint32(img.Pix[i+1])*a + 255*invA + 127) / 255)
		out.Pix[i+2] = uint8((uint32(img.Pix[i+2])*a + 255*invA + 127) / 255)
		out.Pix[i+3] = 255
	}
	return out
}
