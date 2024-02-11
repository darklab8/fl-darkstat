package builder

import (
	"encoding/json"
	"os"

	"github.com/darklab8/fl-darkstat/darkstat/settings"

	"github.com/darklab8/fl-configs/configs/configs_export"
	"github.com/darklab8/fl-configs/configs/configs_mapped"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
	"github.com/darklab8/go-utils/goutils/utils/utils_filepath"
)

/*
Builds freelancer game data into static files accesable by front.
*/
type Builder struct {
}

func NewBuilder() *Builder {
	b := &Builder{}
	return b
}

var PermReadWrite os.FileMode = 0666

func (f *Builder) Build() {
	configs := configs_mapped.NewMappedConfigs()
	logus.Log.Debug("scanning freelancer folder=" + settings.FreelancerFolder.ToString())
	configs.Read(settings.FreelancerFolder)
	export := configs_export.NewExporter(configs)

	bases := export.Bases()
	data, err := json.Marshal(bases)
	logus.Log.CheckFatal(err, "failed to export bases at marshaling")

	err = os.WriteFile(utils_filepath.Join(settings.ProjectFolder, "web", "static", "export", "bases.json").ToString(), data, PermReadWrite)
	logus.Log.CheckFatal(err, "failed to export bases to file")
}
