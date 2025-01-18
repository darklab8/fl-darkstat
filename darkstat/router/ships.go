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

	var useful_ships []configs_export.Ship = data.FilterToUsefulShips(data.Ships)

	build.RegComps(
		builder.NewComponent(
			urls.Ships,
			front.ShipsT(useful_ships, front.ShipShowBases, tab.ShowEmpty(false), shared, data.Infocards),
		),
		builder.NewComponent(
			tab.AllItemsUrl(urls.Ships),
			front.ShipsT(data.Ships, front.ShipShowBases, tab.ShowEmpty(true), shared, data.Infocards),
		),
		builder.NewComponent(
			urls.ShipDetails,
			front.ShipsT(useful_ships, front.ShipShowDetails, tab.ShowEmpty(false), shared, data.Infocards),
		),
		builder.NewComponent(
			tab.AllItemsUrl(urls.ShipDetails),
			front.ShipsT(data.Ships, front.ShipShowDetails, tab.ShowEmpty(true), shared, data.Infocards),
		),
	)

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
				front.ShipRow(ship, front.ShipShowBases, tab.PinMode, shared, data.Infocards, true),
			),
			builder.NewComponent(
				utils_types.FilePath(front.ShipPinnedUrl(ship, front.ShipShowDetails)),
				front.ShipRow(ship, front.ShipShowDetails, tab.PinMode, shared, data.Infocards, true),
			),
		)
	}

	for _, missile := range data.Missiles {
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(front.GunDetailedUrl(missile, front.GunsMissiles)),
				front.GoodAtBaseInfoT(missile.Name, missile.Bases, front.ShowAsCommodity(false), shared),
			),
			builder.NewComponent(
				utils_types.FilePath(front.GunPinnedRowUrl(missile, front.GunsMissiles)),
				front.GunRow(missile, front.GunsMissiles, tab.PinMode, shared, data.Infocards, true),
			),
		)
	}
}
