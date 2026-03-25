// Command utfextract extracts TGA/DDS textures from Freelancer UTF files.
//
// Usage:
//
//	utfextract -in <file_or_dir> -out <output_dir> [-r] [-preserve-paths]
//
// Examples:
//
//	# Extract images from a single file
//	utfextract -in DATA/INTERFACE/NEURONET/NAVMAP/NEWNAVMAP/navmap.txm -out ./images
//
//	# Extract from all UTF files in a directory (non-recursive)
//	utfextract -in DATA/INTERFACE/NEURONET/NAVMAP/NEWNAVMAP -out ./images
//
//	# Extract recursively, mirroring the source directory tree in the output
//	# (useful for validating results against a reference Perl extraction)
//	utfextract -in DATA -out ./images -r -preserve-paths
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/darklab8/fl-darkstat/darkmap/utfextract"
)

/*
Expected:
Total UTF Files Read: 67
Total UTF Texture Libraries: 61
Total UTF Files Modified: 0
Total Images Exported: 278
Total Time: 1 seconds

Received:
Done. UTF files read: 67  Images written: 278
Output: /home/naa/repos/pet_projects/fl-darkstat/darkmap/utfextract/script/received

diff:
nav_navmap_right.cmp => 30.tgaframetexture.tga backgroundpattern.dds edgecolor.tgaagain4.tga
*/

// go run . -preserve-paths -r -in /home/naa/apps/freelancer_related/wine_prefix_freelancer_online2/drive_c/Discovery/DATA/INTERFACE/NEURONET/NAVMAP/NEWNAVMAP/ -out ./received

// go run . -in nav_addwaypoint.cmp -out ./images

func main() {
	inPath := flag.String("in", "", "Input file or directory (required)")
	outPath := flag.String("out", ".", "Output directory (default: current directory)")
	recursive := flag.Bool("r", false, "Recurse into sub-directories")
	preservePaths := flag.Bool("preserve-paths", false, "Mirror source sub-directory structure under output dir (useful for diff-based testing)")
	flag.Parse()

	if *inPath == "" {
		fmt.Fprintln(os.Stderr, "error: -in is required")
		flag.Usage()
		os.Exit(1)
	}

	if err := os.MkdirAll(*outPath, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "error: cannot create output dir: %v\n", err)
		os.Exit(1)
	}

	info, err := os.Stat(*inPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	shapes := utfextract.NewShapes()

	if info.IsDir() {
		if *preservePaths && !*recursive {
			fmt.Fprintln(os.Stderr, "note: -preserve-paths has no effect without -r")
		}
		err := utfextract.ExtractFromDir(*inPath, *outPath, *recursive, *preservePaths, shapes)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: %v\n", err)
		}
		fmt.Printf("Done. UTF files read: %d  Images written: %d\n", shapes.FilesRead, shapes.ImageWritten)
		fmt.Printf("Output: %s\n", absPath(*outPath))
	} else {

		err := utfextract.ExtractFromFile(*inPath, *outPath, shapes)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Done. Images written: %d\n", shapes.ImageWritten)
		fmt.Printf("Output: %s%c%s\n", absPath(*outPath), filepath.Separator, filepath.Base(*inPath))
	}
}

func absPath(p string) string {
	a, err := filepath.Abs(p)
	if err != nil {
		return p
	}
	return a
}
