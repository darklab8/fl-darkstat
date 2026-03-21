package export_map

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/darklab8/fl-darkstat/darkmap/settings"
	"github.com/darklab8/fl-darkstat/darkmap/settings/logus"
	"github.com/darklab8/fl-darkstat/darkmap/utfextract"
	"github.com/darklab8/go-utils/typelog"
)

func findDirs(root, target string) ([]string, error) {
	var matches []string

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && d.Name() == target {
			matches = append(matches, path)
		}

		return nil
	})

	return matches, err
}

func GetImages() *utfextract.Shapes {
	newnavmaps, err := findDirs(string(settings.Env.FreelancerFolder), "NEWNAVMAP")
	if err != nil {
		panic(err)
	}
	if len(newnavmaps) != 1 {
		panic(fmt.Sprintln("expected to find NEWNAVMAP in freelancer folder", settings.Env.FreelancerFolder))
	}

	shapes := utfextract.NewShapes()
	err = utfextract.ExtractFromDir(newnavmaps[0], "", true, false, shapes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: %v\n", err)
	}
	fmt.Printf("Done. UTF files read: %d  Images written: %d\n", shapes.FilesRead, shapes.ImageWritten)

	if shapes.ImageWritten == 0 {
		logus.Log.Panic("expected finding inames in NEWNAVMAP folder. but not found by the path", typelog.Any("path", settings.Env.FreelancerFolder))
	}
	return shapes
}
