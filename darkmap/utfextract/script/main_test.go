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
	fr, iw, err := utfextract.ExtractFromDir(*inPath, *outPath, *recursive, *preservePaths)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: %v\n", err)
	}
	fmt.Printf("Done. UTF files read: %d  Images written: %d\n", fr, iw)
	fmt.Printf("Output: %s\n", absPath(*outPath))

	assert.Equal(t, fr, 3, "expected to read 3 utf files")
	assert.Equal(t, iw, 173, "expected to extract 173 utf files")
}
