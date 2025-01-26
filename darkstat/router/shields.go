package router

import (
	"sort"

	"github.com/darklab8/fl-darkstat/darkcore/builder"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
	"github.com/darklab8/fl-darkstat/darkstat/front"
	"github.com/darklab8/fl-darkstat/darkstat/front/tab"
	"github.com/darklab8/fl-darkstat/darkstat/front/types"
	"github.com/darklab8/fl-darkstat/darkstat/front/urls"
	"github.com/darklab8/go-utils/utils/utils_types"
)

func (l *Router) LinkShields(
	build *builder.Builder,
	data *configs_export.Exporter,
	shared *types.SharedData,
) {

	sort.Slice(data.Shields, func(i, j int) bool {
		if data.Shields[i].Name != "" && data.Shields[j].Name == "" {
			return true
		}
		return data.Shields[i].Name < data.Shields[j].Name
	})

	build.RegComps(
		builder.NewComponent(
			urls.Shields,
			front.ShieldT(data.FilterToUsefulShields(data.Shields), tab.ShowEmpty(false), shared),
		),
		builder.NewComponent(
			tab.AllItemsUrl(urls.Shields),
			front.ShieldT(data.Shields, tab.ShowEmpty(true), shared),
		),
	)
	for _, shield := range data.Shields {
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(front.ShieldDetailedUrl(shield)),
				front.GoodAtBaseInfoT(shield.Name, shield.Bases, front.ShowAsCommodity(false), shared),
			),
		)
	}
}
