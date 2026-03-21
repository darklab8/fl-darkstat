package utfextract

import (
	"bytes"
	"image/jpeg"

	"github.com/darklab8/fl-darkstat/darkmap/tga_ftrvxmtrx"
)

func TransformToJpeg(data []byte) (*bytes.Buffer, error) {
	input := bytes.NewReader(data)
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
}
