package router

import (
	"sort"

	"github.com/a-h/templ"
	"github.com/darklab8/fl-configs/configs/configs_export"
	"github.com/darklab8/fl-darkcore/darkcore/builder"
	"github.com/darklab8/fl-darkstat/darkstat/front"
	"github.com/darklab8/fl-darkstat/darkstat/front/tab"
	"github.com/darklab8/fl-darkstat/darkstat/front/types"
	"github.com/darklab8/fl-darkstat/darkstat/front/urls"
	"github.com/darklab8/go-utils/utils/utils_types"
)

func (l *Router) LinkShips(
	build *builder.Builder,
	data *configs_export.Exporter,
	shared *types.SharedData,
) {
	sort.Slice(data.Ships, func(i, j int) bool {
		if data.Ships[i].Name != "" && data.Ships[j].Name == "" {
			return true
		}
		return data.Ships[i].Name < data.Ships[j].Name
	})

	router := NewTabRouter(build, data.Ships, data.FilterToUsefulShips)
	router.Register(urls.Ships,
		func(items []configs_export.Ship, show_empty tab.ShowEmpty) templ.Component {
			return front.ShipsT(items, front.ShipShowBases, show_empty, shared, data.Infocards)
		})
	router.Register(urls.ShipDetails,
		func(items []configs_export.Ship, show_empty tab.ShowEmpty) templ.Component {
			return front.ShipsT(items, front.ShipShowDetails, show_empty, shared, data.Infocards)
		})

	// deprecating stufff
	// build.RegComps(
	// 	builder.NewComponent(
	// 		urls.Ships,
	// 		front.ShipsT(useful_ships, front.ShipShowBases, tab.ShowEmpty(false), shared, data.Infocards),
	// 	),
	// 	builder.NewComponent(
	// 		tab.AllItemsUrl(urls.Ships),
	// 		front.ShipsT(data.Ships, front.ShipShowBases, tab.ShowEmpty(true), shared, data.Infocards),
	// 	),
	// 	builder.NewComponent(
	// 		urls.ShipDetails,
	// 		front.ShipsT(useful_ships, front.ShipShowDetails, tab.ShowEmpty(false), shared, data.Infocards),
	// 	),
	// 	builder.NewComponent(
	// 		tab.AllItemsUrl(urls.ShipDetails),
	// 		front.ShipsT(data.Ships, front.ShipShowDetails, tab.ShowEmpty(true), shared, data.Infocards),
	// 	),
	// )

	for _, ship := range data.Ships {
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(front.ShipDetailedUrl(ship, front.ShipShowBases)),
				front.GoodAtBaseInfoT(ship.Name, ship.Bases, front.ShowAsCommodity(false), shared),
			),
			builder.NewComponent(
				utils_types.FilePath(front.ShipDetailedUrl(ship, front.ShipShowDetails)),
				front.ShipDetails(ship),
			),
			builder.NewComponent(
				utils_types.FilePath(front.ShipPinnedUrl(ship, front.ShipShowBases)),
				front.ShipRow(ship, front.ShipShowBases, front.PinMode, shared, data.Infocards, true),
			),
			builder.NewComponent(
				utils_types.FilePath(front.ShipPinnedUrl(ship, front.ShipShowDetails)),
				front.ShipRow(ship, front.ShipShowDetails, front.PinMode, shared, data.Infocards, true),
			),
		)
	}
}
