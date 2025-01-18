package router

import (
	"sort"

	"github.com/darklab8/fl-darkcore/darkcore/builder"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
	"github.com/darklab8/fl-darkstat/darkstat/front"
	"github.com/darklab8/fl-darkstat/darkstat/front/tab"
	"github.com/darklab8/fl-darkstat/darkstat/front/types"
	"github.com/darklab8/fl-darkstat/darkstat/front/urls"
	"github.com/darklab8/go-utils/utils/utils_types"
)

func (l *Router) LinkAmmo(
	build *builder.Builder,
	data *configs_export.Exporter,
	shared *types.SharedData,
) {
	sort.Slice(data.Ammos, func(i, j int) bool {
		if data.Ammos[i].Name != "" && data.Ammos[j].Name == "" {
			return true
		}
		return data.Ammos[i].Name < data.Ammos[j].Name
	})

	build.RegComps(
		builder.NewComponent(
			urls.Ammo,
			front.AmmoT(data.FilterToUsefulAmmo(data.Ammos), tab.ShowEmpty(false), shared),
		),
		builder.NewComponent(
			tab.AllItemsUrl(urls.Ammo),
			front.AmmoT(data.Ammos, tab.ShowEmpty(true), shared),
		),
	)
	for _, ammo := range data.Ammos {
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(front.AmmoDetailedUrl(ammo)),
				front.GoodAtBaseInfoT(ammo.Name, ammo.Bases, front.ShowAsCommodity(false), shared),
			),
		)
	}

}
