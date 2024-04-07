package linker

/*
Links data from exported fl-configs
into stuff rendered by fl-darkstat
*/

import (
	"sort"
	"strings"

	"github.com/darklab8/fl-configs/configs/configs_export"
	"github.com/darklab8/fl-configs/configs/configs_mapped"
	"github.com/darklab8/fl-darkstat/darkstat/builder"
	"github.com/darklab8/fl-darkstat/darkstat/front"
	"github.com/darklab8/fl-darkstat/darkstat/front/urls"
	"github.com/darklab8/fl-darkstat/darkstat/settings"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
	"github.com/darklab8/go-utils/goutils/utils/time_measure"
	"github.com/darklab8/go-utils/goutils/utils/utils_logus"
	"github.com/darklab8/go-utils/goutils/utils/utils_types"
)

type Linker struct {
	configs *configs_export.Exporter
}

type LinkOption func(l *Linker)

func NewLinker(opts ...LinkOption) *Linker {
	l := &Linker{}
	for _, opt := range opts {
		opt(l)
	}

	if l.configs == nil {
		configs := configs_mapped.NewMappedConfigs()
		logus.Log.Debug("scanning freelancer folder", utils_logus.FilePath(settings.FreelancerFolder))
		configs.Read(settings.FreelancerFolder)
		l.configs = configs_export.NewExporter(configs)
	}

	return l
}

func (l *Linker) Link() *builder.Builder {
	var data *configs_export.Exporter
	time_measure.TimeMeasure(func(m *time_measure.TimeMeasurer) {
		data = l.configs.Export()
	}, time_measure.WithMsg("exporting data"))

	sort.Slice(data.Bases, func(i, j int) bool {
		return data.Bases[i].Name < data.Bases[j].Name
	})

	for _, base := range data.Bases {
		sort.Slice(base.MarketGoods, func(i, j int) bool {
			if base.MarketGoods[i].Name != "" && base.MarketGoods[j].Name == "" {
				return true
			}
			return base.MarketGoods[i].Name < base.MarketGoods[j].Name
		})
	}

	sort.Slice(data.Factions, func(i, j int) bool {
		if data.Factions[i].Name != "" && data.Factions[j].Name == "" {
			return true
		}
		return data.Factions[i].Name < data.Factions[j].Name
	})

	for _, faction := range data.Factions {
		sort.Slice(faction.Reputations, func(i, j int) bool {
			if faction.Reputations[i].Name != "" && faction.Reputations[j].Name == "" {
				return true
			}
			return faction.Reputations[i].Name < faction.Reputations[j].Name
		})
	}

	sort.Slice(data.Commodities, func(i, j int) bool {
		if data.Commodities[i].Name != "" && data.Commodities[j].Name == "" {
			return true
		}
		return data.Commodities[i].Name < data.Commodities[j].Name
	})

	for _, base_info := range data.Commodities {
		sort.Slice(base_info.Bases, func(i, j int) bool {
			if base_info.Bases[i].BaseName != "" && base_info.Bases[j].BaseName == "" {
				return true
			}
			return base_info.Bases[i].BaseName < base_info.Bases[j].BaseName
		})
	}

	sort.Slice(data.Guns, func(i, j int) bool {
		if data.Guns[i].Name != "" && data.Guns[j].Name == "" {
			return true
		}
		return data.Guns[i].Name < data.Guns[j].Name
	})

	for _, base_info := range data.Guns {
		sort.Slice(base_info.Bases, func(i, j int) bool {
			if base_info.Bases[i].BaseName != "" && base_info.Bases[j].BaseName == "" {
				return true
			}
			return base_info.Bases[i].BaseName < base_info.Bases[j].BaseName
		})
	}

	for _, base_info := range data.Mines {
		sort.Slice(base_info.Bases, func(i, j int) bool {
			if base_info.Bases[i].BaseName != "" && base_info.Bases[j].BaseName == "" {
				return true
			}
			return base_info.Bases[i].BaseName < base_info.Bases[j].BaseName
		})
	}

	sort.Slice(data.Shields, func(i, j int) bool {
		if data.Shields[i].Name != "" && data.Shields[j].Name == "" {
			return true
		}
		return data.Shields[i].Name < data.Shields[j].Name
	})

	for _, base_info := range data.Shields {
		sort.Slice(base_info.Bases, func(i, j int) bool {
			if base_info.Bases[i].BaseName != "" && base_info.Bases[j].BaseName == "" {
				return true
			}
			return base_info.Bases[i].BaseName < base_info.Bases[j].BaseName
		})
	}

	sort.Slice(data.Thrusters, func(i, j int) bool {
		if data.Thrusters[i].Name != "" && data.Thrusters[j].Name == "" {
			return true
		}
		return data.Thrusters[i].Name < data.Thrusters[j].Name
	})

	for _, base_info := range data.Thrusters {
		sort.Slice(base_info.Bases, func(i, j int) bool {
			if base_info.Bases[i].BaseName != "" && base_info.Bases[j].BaseName == "" {
				return true
			}
			return base_info.Bases[i].BaseName < base_info.Bases[j].BaseName
		})
	}

	sort.Slice(data.Ships, func(i, j int) bool {
		if data.Ships[i].Name != "" && data.Ships[j].Name == "" {
			return true
		}
		return data.Ships[i].Name < data.Ships[j].Name
	})

	for _, base_info := range data.Ships {
		sort.Slice(base_info.Bases, func(i, j int) bool {
			if base_info.Bases[i].BaseName != "" && base_info.Bases[j].BaseName == "" {
				return true
			}
			return base_info.Bases[i].BaseName < base_info.Bases[j].BaseName
		})
	}

	sort.Slice(data.Tractors, func(i, j int) bool {
		if data.Tractors[i].Name != "" && data.Tractors[j].Name == "" {
			return true
		}
		return data.Tractors[i].Name < data.Tractors[j].Name
	})

	for _, base_info := range data.Tractors {
		sort.Slice(base_info.Bases, func(i, j int) bool {
			if base_info.Bases[i].BaseName != "" && base_info.Bases[j].BaseName == "" {
				return true
			}
			return base_info.Bases[i].BaseName < base_info.Bases[j].BaseName
		})
	}

	sort.Slice(data.Engines, func(i, j int) bool {
		if data.Engines[i].Name != "" && data.Engines[j].Name == "" {
			return true
		}
		return data.Engines[i].Name < data.Engines[j].Name
	})

	for _, base_info := range data.Engines {
		sort.Slice(base_info.Bases, func(i, j int) bool {
			if base_info.Bases[i].BaseName != "" && base_info.Bases[j].BaseName == "" {
				return true
			}
			return base_info.Bases[i].BaseName < base_info.Bases[j].BaseName
		})
	}

	build := builder.NewBuilder()

	var useful_bases []configs_export.Base = make([]configs_export.Base, 0, len(data.Bases))
	for _, item := range data.Bases {
		if item.Name == "" {
			continue
		}
		useful_bases = append(useful_bases, item)
	}

	var useful_factions []configs_export.Faction = make([]configs_export.Faction, 0, len(data.Factions))
	for _, item := range data.Factions {
		if front.Empty(item.Name) || strings.Contains(item.Name, "_grp") {
			continue
		}
		useful_factions = append(useful_factions, item)
	}

	var buyable_tractors []configs_export.Tractor = make([]configs_export.Tractor, 0, len(data.Tractors))
	for _, item := range data.Tractors {
		if !front.Buyable(item.Bases) {
			continue
		}
		buyable_tractors = append(buyable_tractors, item)
	}

	var buyable_engines []configs_export.Engine = make([]configs_export.Engine, 0, len(data.Engines))
	for _, engine := range data.Engines {
		if !front.Buyable(engine.Bases) {
			continue
		}
		buyable_engines = append(buyable_engines, engine)
	}
	build.RegComps(
		builder.NewComponent(
			urls.Index,
			front.Index(),
		),
		builder.NewComponent(
			urls.Bases,
			front.BasesT(useful_bases, front.ShowEmpty(false)),
		),
		builder.NewComponent(
			front.AllItemsUrl(urls.Bases),
			front.BasesT(data.Bases, front.ShowEmpty(true)),
		),
		builder.NewComponent(
			urls.Factions,
			front.FactionsT(useful_factions, front.FactionShowBases, front.ShowEmpty(false)),
		),
		builder.NewComponent(
			front.AllItemsUrl(urls.Factions),
			front.FactionsT(data.Factions, front.FactionShowBases, front.ShowEmpty(true)),
		),
		builder.NewComponent(
			urls.Rephacks,
			front.FactionsT(useful_factions, front.FactionShowRephacks, front.ShowEmpty(false)),
		),
		builder.NewComponent(
			front.AllItemsUrl(urls.Rephacks),
			front.FactionsT(data.Factions, front.FactionShowRephacks, front.ShowEmpty(true)),
		),
		builder.NewComponent(
			urls.Commodities,
			front.CommoditiesT(data.Commodities, front.ShowEmpty(false)),
		),
		builder.NewComponent(
			front.AllItemsUrl(urls.Commodities),
			front.CommoditiesT(data.Commodities, front.ShowEmpty(true)),
		),
		builder.NewComponent(
			urls.Guns,
			front.GunsT(data.Guns, front.GunsShowBases, front.ShowEmpty(false)),
		),
		builder.NewComponent(
			front.AllItemsUrl(urls.Guns),
			front.GunsT(data.Guns, front.GunsShowBases, front.ShowEmpty(true)),
		),
		builder.NewComponent(
			urls.GunModifiers,
			front.GunsT(data.Guns, front.GunsShowDamageBonuses, front.ShowEmpty(false)),
		),
		builder.NewComponent(
			front.AllItemsUrl(urls.GunModifiers),
			front.GunsT(data.Guns, front.GunsShowDamageBonuses, front.ShowEmpty(true)),
		),
		builder.NewComponent(
			urls.Missiles,
			front.GunsT(data.Missiles, front.GunsMissiles, front.ShowEmpty(false)),
		),
		builder.NewComponent(
			front.AllItemsUrl(urls.Missiles),
			front.GunsT(data.Missiles, front.GunsMissiles, front.ShowEmpty(true)),
		),
		builder.NewComponent(
			urls.Mines,
			front.MinesT(data.Mines, front.ShowEmpty(false)),
		),
		builder.NewComponent(
			front.AllItemsUrl(urls.Mines),
			front.MinesT(data.Mines, front.ShowEmpty(true)),
		),
		builder.NewComponent(
			urls.Shields,
			front.ShieldT(data.Shields, front.ShowEmpty(false)),
		),
		builder.NewComponent(
			front.AllItemsUrl(urls.Shields),
			front.ShieldT(data.Shields, front.ShowEmpty(true)),
		),
		builder.NewComponent(
			urls.Thrusters,
			front.ThrusterT(data.Thrusters, front.ShowEmpty(false)),
		),
		builder.NewComponent(
			front.AllItemsUrl(urls.Thrusters),
			front.ThrusterT(data.Thrusters, front.ShowEmpty(true)),
		),
		builder.NewComponent(
			urls.Ships,
			front.ShipsT(data.Ships, front.ShipShowBases, front.ShowEmpty(false)),
		),
		builder.NewComponent(
			front.AllItemsUrl(urls.Ships),
			front.ShipsT(data.Ships, front.ShipShowBases, front.ShowEmpty(true)),
		),
		builder.NewComponent(
			urls.ShipDetails,
			front.ShipsT(data.Ships, front.ShipShowDetails, front.ShowEmpty(false)),
		),
		builder.NewComponent(
			front.AllItemsUrl(urls.ShipDetails),
			front.ShipsT(data.Ships, front.ShipShowDetails, front.ShowEmpty(true)),
		),
		builder.NewComponent(
			urls.Tractors,
			front.TractorsT(buyable_tractors, front.ShowEmpty(false)),
		),
		builder.NewComponent(
			front.AllItemsUrl(urls.Tractors),
			front.TractorsT(data.Tractors, front.ShowEmpty(true)),
		),
		builder.NewComponent(
			urls.Engines,
			front.Engines(buyable_engines, front.ShowEmpty(false)),
		),
		builder.NewComponent(
			front.AllItemsUrl(urls.Engines),
			front.Engines(data.Engines, front.ShowEmpty(true)),
		),
	)

	for _, base := range data.Bases {
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(front.BaseMarketGoodUrl(base)),
				front.BaseMarketGoods(base.MarketGoods),
			),
		)
	}

	for _, faction := range data.Factions {
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(front.FactionRepUrl(faction, front.FactionShowBases)),
				front.FactionReps(faction.Reputations),
			),
			builder.NewComponent(
				utils_types.FilePath(front.FactionRepUrl(faction, front.FactionShowRephacks)),
				front.RephackBottom(faction.Rephacks),
			),
		)
	}

	for nickname, infocard := range data.Infocards {
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(front.InfocardURL(nickname)),
				front.Infocard(infocard),
			),
		)
	}

	for _, base_info := range data.Commodities {
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(front.GoodAtBaseInfoTUrl(base_info)),
				front.GoodAtBaseInfoT(base_info.Bases, front.ShowPricePerVolume(true)),
			),
		)
	}

	for _, gun := range data.Guns {
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(front.GunDetailedUrl(gun, front.GunsShowBases)),
				front.GoodAtBaseInfoT(gun.Bases, front.ShowPricePerVolume(false)),
			),
			builder.NewComponent(
				utils_types.FilePath(front.GunDetailedUrl(gun, front.GunsShowDamageBonuses)),
				front.GunShowModifiers(gun),
			),
		)
	}

	for _, missile := range data.Missiles {
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(front.GunDetailedUrl(missile, front.GunsMissiles)),
				front.GoodAtBaseInfoT(missile.Bases, front.ShowPricePerVolume(false)),
			),
		)
	}

	for _, mine := range data.Mines {
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(front.MineDetailedUrl(mine)),
				front.GoodAtBaseInfoT(mine.Bases, front.ShowPricePerVolume(false)),
			),
		)
	}

	for _, shield := range data.Shields {
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(front.ShieldDetailedUrl(shield)),
				front.GoodAtBaseInfoT(shield.Bases, front.ShowPricePerVolume(false)),
			),
		)
	}

	for _, thruster := range data.Thrusters {
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(front.ThrusterDetailedUrl(thruster)),
				front.GoodAtBaseInfoT(thruster.Bases, front.ShowPricePerVolume(false)),
			),
		)
	}

	for _, ship := range data.Ships {
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(front.ShipDetailedUrl(ship, front.ShipShowBases)),
				front.GoodAtBaseInfoT(ship.Bases, front.ShowPricePerVolume(false)),
			),
			builder.NewComponent(
				utils_types.FilePath(front.ShipDetailedUrl(ship, front.ShipShowDetails)),
				front.ShipDetails(ship),
			),
		)
	}

	for _, tractor := range data.Tractors {
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(front.TractorDetailedUrl(tractor)),
				front.GoodAtBaseInfoT(tractor.Bases, front.ShowPricePerVolume(false)),
			),
		)
	}

	for _, engine := range data.Engines {
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(front.EngineDetailedUrl(engine)),
				front.GoodAtBaseInfoT(engine.Bases, front.ShowPricePerVolume(false)),
			),
		)
	}

	return build
}
