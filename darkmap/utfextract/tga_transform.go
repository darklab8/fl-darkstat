package utfextract

import (
	"bytes"
	"image/jpeg"

	"github.com/darklab8/fl-darkstat/darkmap/dds"
	"github.com/darklab8/fl-darkstat/darkmap/tga_ftrvxmtrx"
)

func TransformToJpeg(img *Image) (*bytes.Buffer, error) {
	input := bytes.NewReader(img.Data)

	if img.Extension == "tga" {
		img, err := tga_ftrvxmtrx.Decode(input)
		if err != nil {
			return nil, err
		}
		var output *bytes.Buffer = &bytes.Buffer{} // zero value is ready to use

		err = jpeg.Encode(output, img, &jpeg.Options{
			Quality: 90,
		})
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
		var output *bytes.Buffer = &bytes.Buffer{} // zero value is ready to use

		err = jpeg.Encode(output, img, &jpeg.Options{})
		if err != nil {
			return nil, err
		}
		return output, nil
	}

	panic("not supported extension to transform to jpeg")
}
