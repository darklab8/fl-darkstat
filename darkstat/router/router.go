package router

import (
	"github.com/darklab8/fl-darkstat/darkstat/builder"
	"github.com/darklab8/fl-darkstat/darkstat/front"
	"github.com/darklab8/fl-darkstat/darkstat/front/urls"
)

var Builder *builder.Builder

func init() {
	// write(utils_filepath.Join(settings.ProjectFolder, "web", "static", "export", "bases.json"), data)

	// configs := configs_mapped.NewMappedConfigs()

	// logus.Log.Debug("scanning freelancer folder", utils_logus.FilePath(settings.FreelancerFolder))
	// configs.Read(settings.FreelancerFolder)
	// export := configs_export.NewExporter(configs)

	// bases := export.Bases()
	// data, err := json.Marshal(bases)
	// logus.Log.CheckFatal(err, "failed to export bases at marshaling")

	Builder = builder.NewBuilder()
	Builder.RegComps(
		builder.NewComponent(
			urls.Bases,
			front.BasesT(),
		),
	)
}
