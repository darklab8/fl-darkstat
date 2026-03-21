package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/darklab8/fl-darkstat/darkmap/utfextract"
	"github.com/darklab8/go-utils/utils/ptr"
	"github.com/darklab8/go-utils/utils/utils_os"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	preservePaths := ptr.Ptr(true)
	recursive := ptr.Ptr(true)
	currentDir := utils_os.GetCurrentFolder()
	inPath := ptr.Ptr(currentDir.Join("testdata").ToString())
	outPath := ptr.Ptr(currentDir.Join("testresult").ToString())

	os.RemoveAll(*outPath)

	if *preservePaths && !*recursive {
		fmt.Fprintln(os.Stderr, "note: -preserve-paths has no effect without -r")
	}
	shapes := utfextract.NewShapes()
	err := utfextract.ExtractFromDir(*inPath, *outPath, *recursive, *preservePaths, shapes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: %v\n", err)
	}
	fmt.Printf("Done. UTF files read: %d  Images written: %d\n", shapes.FilesRead, shapes.ImageWritten)
	fmt.Printf("Output: %s\n", absPath(*outPath))

	assert.Equal(t, shapes.FilesRead, 3, "expected to read 3 utf files")
	assert.Equal(t, shapes.ImageWritten, 173, "expected to extract 173 utf files")

	image_data := shapes.ShapesByNick["nav_addwaypoint"].Images[0]
	jpeg_result, err := utfextract.TransformToJpeg(image_data)
	if err != nil {
		panic(err)
	}
	file_output, err := os.Create("output.jpg")
	if err != nil {
		panic(err)
	}
	defer file_output.Close()
	file_output.Write(jpeg_result.Bytes())
}
