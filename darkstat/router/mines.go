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

func (l *Router) LinkMines(
	build *builder.Builder,
	data *configs_export.Exporter,
	shared *types.SharedData,
) {
	sort.Slice(data.Mines, func(i, j int) bool {
		if data.Mines[i].Name != "" && data.Mines[j].Name == "" {
			return true
		}
		return data.Mines[i].Name < data.Mines[j].Name
	})

	build.RegComps(
		builder.NewComponent(
			urls.Mines,
			front.MinesT(data.FilterToUsefulMines(data.Mines), tab.ShowEmpty(false), shared),
		),
		builder.NewComponent(
			tab.AllItemsUrl(urls.Mines),
			front.MinesT(data.Mines, tab.ShowEmpty(true), shared),
		),
	)

	for _, mine := range data.Mines {
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(front.MineDetailedUrl(mine)),
				front.GoodAtBaseInfoT(mine.Name, mine.Bases, front.ShowAsCommodity(false), shared),
			),
		)
	}
}
