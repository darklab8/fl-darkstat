package builder

import (
	"encoding/json"
	"fldarkstat/fldarkstat/settings"
	"os"

	"github.com/darklab8/darklab_flconfigs/flconfigs/configs_export"
	"github.com/darklab8/darklab_flconfigs/flconfigs/configs_mapped"
	"github.com/darklab8/darklab_flconfigs/flconfigs/settings/logus"
	"github.com/darklab8/darklab_goutils/goutils/utils/utils_filepath"
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
