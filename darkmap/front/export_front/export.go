package export_front

import (
	"github.com/darklab8/fl-darkstat/configs/configs_mapped"
	"github.com/darklab8/fl-darkstat/darkmap/settings"
	"github.com/darklab8/fl-darkstat/darkmap/settings/logus"
	"github.com/darklab8/go-utils/utils/timeit"
	"github.com/darklab8/go-utils/utils/utils_logus"
)

type Export struct {
	Mapped  *configs_mapped.MappedConfigs
	Systems []System
}

func NewExport() *Export {
	e := &Export{}

	defer timeit.NewTimer("MappedConfigs creation").Close()

	freelancer_folder := settings.Env.FreelancerFolder
	if e.Mapped == nil {
		logus.Log.Debug("scanning freelancer folder", utils_logus.FilePath(freelancer_folder))
		e.Mapped = configs_mapped.NewMappedConfigs().Read(freelancer_folder)
	}

	e.export()

	return e
}

func (e *Export) GetInfocardName(ids_name int, nickname string) string {
	return e.Mapped.GetInfocardName(ids_name, nickname)
}

func (e *Export) export() {
	e.Systems = exportSystems(e.Mapped)
}
