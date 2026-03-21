package main

import (
	"image/jpeg"
	"os"

	"github.com/darklab8/fl-darkstat/darkmap/tga_ftrvxmtrx"
)

func main() {
	input, err := os.Open("../../discovery_navmap/nav_addwaypoint.cmp/nav_buttonback.tga")
	if err != nil {
		panic(err)
	}
	defer input.Close()

	img, err := tga_ftrvxmtrx.Decode(input)
	if err != nil {
		panic(err)
	}

	output, err := os.Create("output.jpg")
	if err != nil {
		panic(err)
	}
	defer output.Close()

	err = jpeg.Encode(output, img, &jpeg.Options{
		Quality: 90,
	})
	if err != nil {
		panic(err)
	}
}
