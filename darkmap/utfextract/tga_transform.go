package utfextract

import (
	"bytes"
	"image"
	"image/png"

	"github.com/darklab8/fl-darkstat/darkmap/dds"
	"github.com/darklab8/fl-darkstat/darkmap/tga"
	"golang.org/x/image/draw"
)

// 200px is the rendered size; use 400 if you want crisp output on HiDPI/retina screens.
const MaxDimension = 256

func downscale(src image.Image) image.Image {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()
	if w <= MaxDimension && h <= MaxDimension {
		return src // never upscale
	}
	scale := float64(MaxDimension) / float64(max(w, h))
	nw, nh := max(1, int(float64(w)*scale)), max(1, int(float64(h)*scale))

	dst := image.NewNRGBA(image.Rect(0, 0, nw, nh))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, b, draw.Over, nil)
	return dst
}

func TransformToJpeg(img *Image) (*bytes.Buffer, error) {
	input := bytes.NewReader(img.Data)

	if img.Extension == "tga" {
		img, err := tga.Decode(input)
		if err != nil {
			return nil, err
		}
		var output *bytes.Buffer = &bytes.Buffer{} // zero value is ready to use

		err = png.Encode(output, downscale(img))
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

		err = png.Encode(output, downscale(img))
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
