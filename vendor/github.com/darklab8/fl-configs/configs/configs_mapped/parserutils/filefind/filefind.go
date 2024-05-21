/*
Package with reusable code for discovery of files and other reusable stuff like universal ini reader
*/
package filefind

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-configs/configs/settings/logus"
	"github.com/darklab8/go-typelog/typelog"

	"github.com/darklab8/go-utils/goutils/utils/utils_logus"
	"github.com/darklab8/go-utils/goutils/utils/utils_types"
)

type Filesystem struct {
	Files   []*file.File
	Hashmap map[utils_types.FilePath]*file.File
}

var FreelancerFolder Filesystem

func FindConfigs(folderpath utils_types.FilePath) *Filesystem {
	var filesystem *Filesystem = &Filesystem{}
	filesystem.Hashmap = make(map[utils_types.FilePath]*file.File)

	err := filepath.WalkDir(string(folderpath), func(path string, d fs.DirEntry, err error) error {

		if !strings.Contains(path, ".ini") &&
			!strings.Contains(path, ".txt") &&
			!strings.Contains(path, ".xml") &&
			!strings.Contains(path, ".dll") {
			return nil
		}

		logus.Log.CheckPanic(err, "unable to read file")

		file := file.NewFile(utils_types.FilePath(path))
		filesystem.Files = append(filesystem.Files, file)

		key := utils_types.FilePath(strings.ToLower(filepath.Base(path)))
		filesystem.Hashmap[key] = file

		return nil
	})

	logus.Log.CheckPanic(err, "unable to read files")
	return filesystem
}

func (file1system Filesystem) GetFile(file1names ...utils_types.FilePath) *file.File {
	for _, file1name := range file1names {
		file_, ok := file1system.Hashmap[file1name]
		if !ok {
			logus.Log.Warn("Filesystem.GetFile, failed to find find in filesystesm file trying to recover", utils_logus.FilePath(file1name))
			continue
		}
		logus.Log.Info("Filesystem.GetFile, found filepath=", utils_logus.FilePath(file_.GetFilepath()))
		result_file := file.NewFile(file_.GetFilepath())
		return result_file
	}

	logus.Log.Warn("failed to get file", typelog.Items[utils_types.FilePath]("filenames", file1names))
	return nil
}
