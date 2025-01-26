package builder

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/darklab8/fl-darkstat/darkcore/settings/logus"
	"github.com/darklab8/go-typelog/typelog"
	"github.com/darklab8/go-utils/utils/utils_filepath"
	"github.com/darklab8/go-utils/utils/utils_types"
)

/*
Filesystem allows us to write to files to memory for later reusage in web app serving static assets from memory
Optionally same filesystem supports rendering to local, for deployment of static assets
*/
type Filesystem struct {
	Files      map[utils_types.FilePath]MemFile
	mu         sync.Mutex
	build_root utils_types.FilePath
}

type MemFile interface {
	Render() []byte
}

type MemComp struct {
	comp *Component
	b    *Builder
}

func (m *MemComp) Render() []byte {
	return m.comp.Write(m.b.params).bytes
}

type MemStatic struct {
	content []byte
}

func (m *MemStatic) Render() []byte {
	return m.content
}

func NewFileystem(build_root utils_types.FilePath) *Filesystem {
	b := &Filesystem{
		Files:      make(map[utils_types.FilePath]MemFile),
		build_root: build_root,
	}
	return b
}

var PermReadWrite os.FileMode = 0666

func (f *Filesystem) GetBuildRoot() utils_types.FilePath {
	return f.build_root
}

func (f *Filesystem) WriteToMem(path utils_types.FilePath, content MemFile) {
	f.mu.Lock()
	f.Files[path] = content
	f.mu.Unlock()
}

func (f *Filesystem) WriteToFile(path utils_types.FilePath, content []byte) {

	final_path := utils_filepath.Join(f.build_root, path)
	haveParentFoldersCreated(final_path)
	// TODO add check for creating all missing folders in the path
	err := os.WriteFile(final_path.ToString(), []byte(content), PermReadWrite)
	logus.Log.CheckFatal(err, "failed to export bases to file")
}

func (f *Filesystem) CreateBuildFolder() {
	os.RemoveAll(f.build_root.ToString())
	os.MkdirAll(f.build_root.ToString(), os.ModePerm)
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
