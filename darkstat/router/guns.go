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

func (l *Router) LinkGuns(
	build *builder.Builder,
	data *configs_export.Exporter,
	shared *types.SharedData,
) {
	sort.Slice(data.Guns, func(i, j int) bool {
		if data.Guns[i].Name != "" && data.Guns[j].Name == "" {
			return true
		}
		return data.Guns[i].Name < data.Guns[j].Name
	})
	sort.Slice(data.Missiles, func(i, j int) bool {
		if data.Missiles[i].Name != "" && data.Missiles[j].Name == "" {
			return true
		}
		return data.Missiles[i].Name < data.Missiles[j].Name
	})
	var useful_guns []configs_export.Gun = data.FilterToUsefulGun(data.Guns)
	var useful_missiles []configs_export.Gun = data.FilterToUsefulGun(data.Missiles)

	build.RegComps(
		builder.NewComponent(
			urls.Guns,
			front.GunsT(useful_guns, front.GunsShowBases, tab.ShowEmpty(false), shared, data.Infocards),
		),
		builder.NewComponent(
			tab.AllItemsUrl(urls.Guns),
			front.GunsT(data.Guns, front.GunsShowBases, tab.ShowEmpty(true), shared, data.Infocards),
		),
		builder.NewComponent(
			urls.GunModifiers,
			front.GunsT(useful_guns, front.GunsShowDamageBonuses, tab.ShowEmpty(false), shared, data.Infocards),
		),
		builder.NewComponent(
			tab.AllItemsUrl(urls.GunModifiers),
			front.GunsT(data.Guns, front.GunsShowDamageBonuses, tab.ShowEmpty(true), shared, data.Infocards),
		),

		builder.NewComponent(
			urls.Missiles,
			front.GunsT(useful_missiles, front.GunsMissiles, tab.ShowEmpty(false), shared, data.Infocards),
		),
		builder.NewComponent(
			tab.AllItemsUrl(urls.Missiles),
			front.GunsT(data.Missiles, front.GunsMissiles, tab.ShowEmpty(true), shared, data.Infocards),
		),
	)

	for _, gun := range data.Guns {
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(front.GunDetailedUrl(gun, front.GunsShowBases)),
				front.GoodAtBaseInfoT(gun.Name, gun.Bases, front.ShowAsCommodity(false), shared),
			),
			builder.NewComponent(
				utils_types.FilePath(front.GunDetailedUrl(gun, front.GunsShowDamageBonuses)),
				front.GunShowModifiers(gun),
			),

			builder.NewComponent(
				utils_types.FilePath(front.GunPinnedRowUrl(gun, front.GunsShowBases)),
				front.GunRow(gun, front.GunsShowBases, tab.PinMode, shared, data.Infocards, true),
			),
			builder.NewComponent(
				utils_types.FilePath(front.GunPinnedRowUrl(gun, front.GunsShowDamageBonuses)),
				front.GunRow(gun, front.GunsShowDamageBonuses, tab.PinMode, shared, data.Infocards, true),
			),
		)
	}
}
