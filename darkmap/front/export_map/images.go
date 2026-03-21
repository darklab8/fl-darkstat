package export_map

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

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

func GetImages(folder_name string) *utfextract.Shapes {
	folders, err := findDirs(string(settings.Env.FreelancerFolder), filepath.Base(folder_name))

	var filtered_folders []string
	for _, folder := range folders {
		if strings.Contains(folder, folder_name) {
			filtered_folders = append(filtered_folders, folder)
		}
	}
	folders = filtered_folders

	if err != nil {
		panic(err)
	}
	if len(folders) != 1 {
		panic(fmt.Sprintln("expected to find ", folder_name, " in freelancer folder", settings.Env.FreelancerFolder))
	}

	shapes := utfextract.NewShapes()
	err = utfextract.ExtractFromDir(folders[0], "", true, true, shapes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: %v\n", err)
	}
	fmt.Printf("Done. UTF files read: %d  Images written: %d\n", shapes.FilesRead, shapes.ImageWritten)

	if shapes.ImageWritten == 0 {
		logus.Log.Panic(fmt.Sprintln("expected finding inames in ", folder_name, " folder. but not found by the path"), typelog.Any("path", settings.Env.FreelancerFolder))
	}
	return shapes
}
