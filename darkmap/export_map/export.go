package export_map

import (
	"context"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped"
	"github.com/darklab8/fl-darkstat/darkmap/search_bar"
	"github.com/darklab8/fl-darkstat/darkmap/settings"
	"github.com/darklab8/fl-darkstat/darkmap/settings/logus"
	"github.com/darklab8/fl-darkstat/darkmap/utfextract"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
	"github.com/darklab8/go-utils/utils/timeit"
	"github.com/darklab8/go-utils/utils/utils_logus"
)

type Export struct {
	Mapped  *configs_mapped.MappedConfigs
	Systems []*System
	Graph   SystemGraphs
	Shapes  *utfextract.Shapes

	PobsBySystemNick   map[string][]*configs_export.PoB
	MiningBySystemNick map[string][]*configs_export.Base
	MiningUsefulByNick map[string]bool
	SearchEntries      map[string]search_bar.Entry

	Exp *configs_export.Exporter
}

func NewExport(ctx context.Context) *Export {
	e := &Export{
		PobsBySystemNick:   make(map[string][]*configs_export.PoB),
		MiningBySystemNick: make(map[string][]*configs_export.Base),
		MiningUsefulByNick: make(map[string]bool),
		SearchEntries:      make(map[string]search_bar.Entry),
	}

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
	e.Shapes = GetImages("NEWNAVMAP")
	more_shapes := GetImages("DATA/SOLAR/PLANETS")
	for key, shape := range more_shapes.ShapesByNick {
		e.Shapes.ShapesByNick[key] = shape
	}
	e.Shapes.FilesRead += more_shapes.FilesRead
	e.Shapes.ImageWritten += more_shapes.ImageWritten

	e.Exp = configs_export.NewExporter(e.Mapped)

	for _, pob := range e.Exp.GetPoBs() {
		if pob.SystemNick == nil {
			continue
		}
		e.PobsBySystemNick[*pob.SystemNick] = append(e.PobsBySystemNick[*pob.SystemNick], pob)
	}

	MiningOperations := e.Exp.GetOres(ctx, []*configs_export.Commodity{}, false)
	for _, mine := range MiningOperations {
		e.MiningBySystemNick[mine.SystemNickname] = append(e.MiningBySystemNick[mine.SystemNickname], mine)
	}
	useful_mining_operations := configs_export.FitlerToUsefulOres(MiningOperations)
	for _, mine := range useful_mining_operations {
		e.MiningUsefulByNick[string(mine.Nickname)] = true
	}

	e.Systems = e.ExportSystems(e.Mapped)
	e.Graph = e.GetSystemConnections(e.Systems)
	e.Exp.Bases = e.Exp.GetBases(ctx)

}
