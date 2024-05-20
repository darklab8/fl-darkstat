package builder

import (
	"os"
	"path/filepath"

	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
	"github.com/darklab8/go-typelog/typelog"
	"github.com/darklab8/go-utils/goutils/utils/utils_filepath"
	"github.com/darklab8/go-utils/goutils/utils/utils_types"
)

/*
Filesystem allows us to write to files to memory for later reusage in web app serving static assets from memory
Optionally same filesystem supports rendering to local, for deployment of static assets
*/
type Filesystem struct {
	Files      map[utils_types.FilePath][]byte
	build_root utils_types.FilePath
}

func NewFileystem(build_root utils_types.FilePath) *Filesystem {
	b := &Filesystem{
		Files:      make(map[utils_types.FilePath][]byte),
		build_root: build_root,
	}
	return b
}

var PermReadWrite os.FileMode = 0666

func (f *Filesystem) WriteToMem(path utils_types.FilePath, content []byte) {
	f.Files[path] = content
}

func (f *Filesystem) RenderToLocal() {
	os.RemoveAll(f.build_root.ToString())
	os.MkdirAll(f.build_root.ToString(), os.ModePerm)

	for path, content := range f.Files {
		final_path := utils_filepath.Join(f.build_root, path)
		haveParentFoldersCreated(final_path)
		// TODO add check for creating all missing folders in the path
		err := os.WriteFile(final_path.ToString(), []byte(content), PermReadWrite)
		logus.Log.CheckFatal(err, "failed to export bases to file")
	}
}

func haveParentFoldersCreated(buildpath utils_types.FilePath) {
	path := buildpath.ToString()
	folder_path := filepath.Dir(path)
	err := os.MkdirAll(folder_path, os.ModePerm)
	logus.Log.CheckError(err,
		"haveParentFoldersCreated finished",
		typelog.String("folderpath", folder_path),
		typelog.String("path", path),
	)
}
