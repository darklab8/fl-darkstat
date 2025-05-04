package bini

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
	"github.com/darklab8/go-utils/utils/utils_types"
)

func TestBini(t *testing.T) {
	// Converting from Bini to Txt
	if true {
		// Turned of because not wishing to commit those files
		return
	}
	vanilla_location := "~/windows10shared/fl-files-vanilla"
	var desired_filepath string
	err := filepath.WalkDir(vanilla_location, func(path string, d fs.DirEntry, err error) error {
		if !strings.Contains(path, "mbases.ini") {
			return nil
		}
		desired_filepath = path
		return errors.New("stop")
	})
	logus.Log.CheckError(err, "failed to walk dir")

	data := Dump(utils_types.FilePath(desired_filepath))
	err = os.WriteFile("output.txt", []byte(strings.Join(data, "\n")), 0644)
	logus.Log.CheckError(err, "failed to write file")

}
