package export_map

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/darklab8/fl-darkstat/darkmap/settings"
	"github.com/stretchr/testify/assert"
)

func TestGrabColrs(t *testing.T) {
	folder_name := "DATA/SOLAR"
	folders, _ := findDirs(string(settings.Env.FreelancerFolder), filepath.Base(folder_name))

	var filtered_folders []string
	for _, folder := range folders {
		if strings.Contains(folder, folder_name) {
			filtered_folders = append(filtered_folders, folder)
		}
	}
	folders = filtered_folders

	fmt.Println(folders)
	assert.Len(t, folders, 1)
}
