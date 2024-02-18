package builder

import (
	"encoding/json"
	"os"

	"github.com/darklab8/fl-darkstat/darkstat/settings"

	"github.com/darklab8/fl-configs/configs/configs_export"
	"github.com/darklab8/fl-configs/configs/configs_mapped"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
	"github.com/darklab8/go-utils/goutils/utils/utils_filepath"
	"github.com/darklab8/go-utils/goutils/utils/utils_logus"
	"github.com/darklab8/go-utils/goutils/utils/utils_types"
)

/*
Builds freelancer game data into static files accesable by front.
*/
type Filesystem struct {
	files map[utils_types.FilePath]string
}

func NewFileystem() *Filesystem {
	b := &Filesystem{
		files: make(map[utils_types.FilePath]string),
	}
	return b
}

var PermReadWrite os.FileMode = 0666

func (f *Filesystem) ScanToMem() {
	configs := configs_mapped.NewMappedConfigs()

	logus.Log.Debug("scanning freelancer folder", utils_logus.FilePath(settings.FreelancerFolder))
	configs.Read(settings.FreelancerFolder)
	export := configs_export.NewExporter(configs)

	bases := export.Bases()
	data, err := json.Marshal(bases)
	logus.Log.CheckFatal(err, "failed to export bases at marshaling")

	f.files[utils_filepath.Join(settings.ProjectFolder, "web", "static", "export", "bases.json")] = string(data)
}

func (f *Filesystem) RenderToLocal() {
	for filepath, content := range f.files {
		err := os.WriteFile(filepath.ToString(), []byte(content), PermReadWrite)
		logus.Log.CheckFatal(err, "failed to export bases to file")
	}
}
