package export_front

import (
	"context"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped"
	"github.com/darklab8/fl-darkstat/darkmap/settings"
	"github.com/darklab8/fl-darkstat/darkmap/settings/logus"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
	"github.com/darklab8/go-utils/utils/timeit"
	"github.com/darklab8/go-utils/utils/utils_logus"
)

type Export struct {
	Mapped  *configs_mapped.MappedConfigs
	Systems []*System
	Graph   SystemGraphs

	Exp *configs_export.Exporter
}

func NewExport(ctx context.Context) *Export {
	e := &Export{}

	defer timeit.NewTimer("MappedConfigs creation").Close()

	freelancer_folder := settings.Env.FreelancerFolder
	if e.Mapped == nil {
		logus.Log.Debug("scanning freelancer folder", utils_logus.FilePath(freelancer_folder))
		e.Mapped = configs_mapped.NewMappedConfigs().Read(ctx, freelancer_folder)
	}

	e.Export(ctx)

	return e
}

func (e *Export) GetInfocardName(ids_name int, nickname string) string {
	return e.Mapped.GetInfocardName(ids_name, nickname)
}

func (e *Export) Export(ctx context.Context) {
	e.Systems = ExportSystems(e.Mapped)
	e.Graph = e.GetSystemConnections(e.Systems)

	e.Exp = configs_export.NewExporter(e.Mapped)
	e.Exp.Bases = e.Exp.GetBases(ctx)
}
