package linker

/*
Links data from exported fl-configs
into stuff rendered by fl-darkstat
*/

import (
	"sort"

	"github.com/darklab8/fl-configs/configs/configs_export"
	"github.com/darklab8/fl-configs/configs/configs_mapped"
	"github.com/darklab8/fl-configs/configs/conftypes"
	"github.com/darklab8/fl-darkstat/darkstat/builder"
	"github.com/darklab8/fl-darkstat/darkstat/common/types"
	"github.com/darklab8/fl-darkstat/darkstat/front"
	"github.com/darklab8/fl-darkstat/darkstat/front/fronttypes"
	"github.com/darklab8/fl-darkstat/darkstat/front/urls"
	"github.com/darklab8/fl-darkstat/darkstat/settings"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
	"github.com/darklab8/go-utils/goutils/utils/time_measure"
	"github.com/darklab8/go-utils/goutils/utils/utils_logus"
	"github.com/darklab8/go-utils/goutils/utils/utils_types"
)

type Linker struct {
	mapped  *configs_mapped.MappedConfigs
	configs *configs_export.Exporter
}

type LinkOption func(l *Linker)

func NewLinker(opts ...LinkOption) *Linker {
	l := &Linker{}
	for _, opt := range opts {
		opt(l)
	}

	time_measure.TimeMeasure(func(m *time_measure.TimeMeasurer) {

		if l.configs == nil {
			l.mapped = configs_mapped.NewMappedConfigs()
			logus.Log.Debug("scanning freelancer folder", utils_logus.FilePath(settings.FreelancerFolder))
			l.mapped.Read(settings.FreelancerFolder)
			l.configs = configs_export.NewExporter(l.mapped)
		}
	}, time_measure.WithMsg("MappedConfigs creation"))
	return l
}

func (l *Linker) Link() *builder.Builder {
	var build *builder.Builder
	time_measure.TimeMeasure(func(m *time_measure.TimeMeasurer) {
		time_measure.TimeMeasure(func(m *time_measure.TimeMeasurer) {
			build = builder.NewBuilder()
		}, time_measure.WithMsg("building creation"))

		var data *configs_export.Exporter
		time_measure.TimeMeasure(func(m *time_measure.TimeMeasurer) {
			data = l.configs.Export()
		}, time_measure.WithMsg("exporting data"))

		time_measure.TimeMeasure(func(m *time_measure.TimeMeasurer) {
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

			sort.Slice(data.Scanners, func(i, j int) bool {
				if data.Scanners[i].Name != "" && data.Scanners[j].Name == "" {
					return true
				}
				return data.Scanners[i].Name < data.Scanners[j].Name
			})

			for _, base_info := range data.Scanners {
				sort.Slice(base_info.Bases, func(i, j int) bool {
					if base_info.Bases[i].BaseName != "" && base_info.Bases[j].BaseName == "" {
						return true
					}
					return base_info.Bases[i].BaseName < base_info.Bases[j].BaseName
				})
			}

		}, time_measure.WithMsg("sorting completed"))

		var useful_factions []configs_export.Faction
		var useful_ships []configs_export.Ship
		var useful_guns []configs_export.Gun
		var useful_missiles []configs_export.Gun
		tractor_id := conftypes.TractorID("")

		var disco_ids fronttypes.DiscoveryIDs

		time_measure.TimeMeasure(func(m *time_measure.TimeMeasurer) {
			useful_factions = configs_export.FilterToUsefulFactions(data.Factions)
			useful_ships = configs_export.FilterToUsefulShips(data.Ships)
			useful_guns = configs_export.FilterToUsefulGun(data.Guns)
			useful_missiles = configs_export.FilterToUsefulGun(data.Missiles)

			if l.mapped.Discovery != nil {
				disco_ids = fronttypes.DiscoveryIDs{
					Show:   true,
					Ids:    l.configs.Tractors,
					Config: l.mapped.Discovery.Techcompat,
				}
			}
		}, time_measure.WithMsg("filtering to useful stuff"))

		time_measure.TimeMeasure(func(m *time_measure.TimeMeasurer) {

			build.RegComps(
				builder.NewComponent(
					urls.Ships+utils_types.FilePath(tractor_id),
					front.ShipsT(useful_ships, front.ShipShowBases, front.ShowEmpty(false), disco_ids),
				),
				builder.NewComponent(
					front.AllItemsUrl(urls.Ships)+utils_types.FilePath(tractor_id),
					front.ShipsT(data.Ships, front.ShipShowBases, front.ShowEmpty(true), disco_ids),
				),
				builder.NewComponent(
					urls.ShipDetails+utils_types.FilePath(tractor_id),
					front.ShipsT(useful_ships, front.ShipShowDetails, front.ShowEmpty(false), disco_ids),
				),
				builder.NewComponent(
					front.AllItemsUrl(urls.ShipDetails)+utils_types.FilePath(tractor_id),
					front.ShipsT(data.Ships, front.ShipShowDetails, front.ShowEmpty(true), disco_ids),
				),
			)

			build.RegComps(

				builder.NewComponent(
					urls.Index,
					front.Index(types.ThemeLight),
				),
				builder.NewComponent(
					"dark.html",
					front.Index(types.ThemeDark),
				),
				builder.NewComponent(
					urls.Bases,
					front.BasesT(configs_export.FilterToUserfulBases(data.Bases), front.ShowEmpty(false)),
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
					front.CommoditiesT(configs_export.FilterToUsefulCommodities(data.Commodities), front.ShowEmpty(false)),
				),
				builder.NewComponent(
					front.AllItemsUrl(urls.Commodities),
					front.CommoditiesT(data.Commodities, front.ShowEmpty(true)),
				),
				builder.NewComponent(
					urls.Guns,
					front.GunsT(useful_guns, front.GunsShowBases, front.ShowEmpty(false), disco_ids),
				),
				builder.NewComponent(
					front.AllItemsUrl(urls.Guns),
					front.GunsT(data.Guns, front.GunsShowBases, front.ShowEmpty(true), disco_ids),
				),
				builder.NewComponent(
					urls.GunModifiers,
					front.GunsT(useful_guns, front.GunsShowDamageBonuses, front.ShowEmpty(false), disco_ids),
				),
				builder.NewComponent(
					front.AllItemsUrl(urls.GunModifiers),
					front.GunsT(data.Guns, front.GunsShowDamageBonuses, front.ShowEmpty(true), disco_ids),
				),
				builder.NewComponent(
					urls.Ammo,
					front.AmmoT(configs_export.FilterToUsefulAmmo(data.Ammos), front.ShowEmpty(false), disco_ids),
				),
				builder.NewComponent(
					front.AllItemsUrl(urls.Ammo),
					front.AmmoT(data.Ammos, front.ShowEmpty(true), disco_ids),
				),
				builder.NewComponent(
					urls.Missiles,
					front.GunsT(useful_missiles, front.GunsMissiles, front.ShowEmpty(false), disco_ids),
				),
				builder.NewComponent(
					front.AllItemsUrl(urls.Missiles),
					front.GunsT(data.Missiles, front.GunsMissiles, front.ShowEmpty(true), disco_ids),
				),
				builder.NewComponent(
					urls.Mines,
					front.MinesT(configs_export.FilterToUsefulMines(data.Mines), front.ShowEmpty(false), disco_ids),
				),
				builder.NewComponent(
					front.AllItemsUrl(urls.Mines),
					front.MinesT(data.Mines, front.ShowEmpty(true), disco_ids),
				),
				builder.NewComponent(
					urls.Shields,
					front.ShieldT(configs_export.FilterToUsefulShields(data.Shields), front.ShowEmpty(false), disco_ids),
				),
				builder.NewComponent(
					front.AllItemsUrl(urls.Shields),
					front.ShieldT(data.Shields, front.ShowEmpty(true), disco_ids),
				),
				builder.NewComponent(
					urls.Thrusters,
					front.ThrusterT(configs_export.FilterToUsefulThrusters(data.Thrusters), front.ShowEmpty(false), disco_ids),
				),
				builder.NewComponent(
					front.AllItemsUrl(urls.Thrusters),
					front.ThrusterT(data.Thrusters, front.ShowEmpty(true), disco_ids),
				),
				builder.NewComponent(
					urls.Tractors,
					front.TractorsT(configs_export.FilterToUsefulTractors(data.Tractors), front.ShowEmpty(false)),
				),
				builder.NewComponent(
					front.AllItemsUrl(urls.Tractors),
					front.TractorsT(data.Tractors, front.ShowEmpty(true)),
				),
				builder.NewComponent(
					urls.Engines,
					front.Engines(configs_export.FilterToUsefulEngines(data.Engines), front.ShowEmpty(false), disco_ids),
				),
				builder.NewComponent(
					front.AllItemsUrl(urls.Engines),
					front.Engines(data.Engines, front.ShowEmpty(true), disco_ids),
				),
				builder.NewComponent(
					urls.CounterMeasures,
					front.CounterMeasureT(configs_export.FilterToUsefulCounterMeasures(data.CMs), front.ShowEmpty(false), disco_ids),
				),
				builder.NewComponent(
					front.AllItemsUrl(urls.CounterMeasures),
					front.CounterMeasureT(data.CMs, front.ShowEmpty(true), disco_ids),
				),
				builder.NewComponent(
					urls.Scanners,
					front.ScannersT(configs_export.FilterToUserfulScanners(data.Scanners), front.ShowEmpty(false), disco_ids),
				),
				builder.NewComponent(
					front.AllItemsUrl(urls.Scanners),
					front.ScannersT(data.Scanners, front.ShowEmpty(true), disco_ids),
				),
			)

			for _, base := range data.Bases {
				build.RegComps(
					builder.NewComponent(
						utils_types.FilePath(front.BaseMarketGoodUrl(base)),
						front.BaseMarketGoods(base.Name, base.MarketGoods),
					),
				)
			}
		}, time_measure.WithMsg("linking main stuff"))

		time_measure.TimeMeasure(func(m *time_measure.TimeMeasurer) {
			for _, faction := range data.Factions {
				build.RegComps(
					builder.NewComponent(
						utils_types.FilePath(front.FactionRepUrl(faction, front.FactionShowBases)),
						front.FactionReps(faction, faction.Reputations),
					),
					builder.NewComponent(
						utils_types.FilePath(front.FactionRepUrl(faction, front.FactionShowRephacks)),
						front.RephackBottom(faction, faction.Rephacks),
					),
				)
			}
		}, time_measure.WithMsg("linking faction stuff"))

		time_measure.TimeMeasure(func(m *time_measure.TimeMeasurer) {
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
						front.GoodAtBaseInfoT(base_info.Name, base_info.Bases, front.ShowPricePerVolume(true)),
					),
				)
			}

			for _, gun := range data.Guns {
				build.RegComps(
					builder.NewComponent(
						utils_types.FilePath(front.GunDetailedUrl(gun, front.GunsShowBases)),
						front.GoodAtBaseInfoT(gun.Name, gun.Bases, front.ShowPricePerVolume(false)),
					),
					builder.NewComponent(
						utils_types.FilePath(front.GunDetailedUrl(gun, front.GunsShowDamageBonuses)),
						front.GunShowModifiers(gun),
					),

					builder.NewComponent(
						utils_types.FilePath(front.GunPinnedRowUrl(gun, front.GunsShowBases)),
						front.GunRow(gun, front.GunsShowBases, front.PinMode, disco_ids, true),
					),
					builder.NewComponent(
						utils_types.FilePath(front.GunPinnedRowUrl(gun, front.GunsShowDamageBonuses)),
						front.GunRow(gun, front.GunsShowDamageBonuses, front.PinMode, disco_ids, true),
					),
				)
			}
			for _, ammo := range data.Ammos {
				build.RegComps(
					builder.NewComponent(
						utils_types.FilePath(front.AmmoDetailedUrl(ammo)),
						front.GoodAtBaseInfoT(ammo.Name, ammo.Bases, front.ShowPricePerVolume(false)),
					),
				)
			}

			for _, missile := range data.Missiles {
				build.RegComps(
					builder.NewComponent(
						utils_types.FilePath(front.GunDetailedUrl(missile, front.GunsMissiles)),
						front.GoodAtBaseInfoT(missile.Name, missile.Bases, front.ShowPricePerVolume(false)),
					),
					builder.NewComponent(
						utils_types.FilePath(front.GunPinnedRowUrl(missile, front.GunsMissiles)),
						front.GunRow(missile, front.GunsMissiles, front.PinMode, disco_ids, true),
					),
				)
			}

			for _, mine := range data.Mines {
				build.RegComps(
					builder.NewComponent(
						utils_types.FilePath(front.MineDetailedUrl(mine)),
						front.GoodAtBaseInfoT(mine.Name, mine.Bases, front.ShowPricePerVolume(false)),
					),
				)
			}

			for _, shield := range data.Shields {
				build.RegComps(
					builder.NewComponent(
						utils_types.FilePath(front.ShieldDetailedUrl(shield)),
						front.GoodAtBaseInfoT(shield.Name, shield.Bases, front.ShowPricePerVolume(false)),
					),
				)
			}

			for _, thruster := range data.Thrusters {
				build.RegComps(
					builder.NewComponent(
						utils_types.FilePath(front.ThrusterDetailedUrl(thruster)),
						front.GoodAtBaseInfoT(thruster.Name, thruster.Bases, front.ShowPricePerVolume(false)),
					),
				)
			}

			for _, ship := range data.Ships {
				build.RegComps(
					builder.NewComponent(
						utils_types.FilePath(front.ShipDetailedUrl(ship, front.ShipShowBases)),
						front.GoodAtBaseInfoT(ship.Name, ship.Bases, front.ShowPricePerVolume(false)),
					),
					builder.NewComponent(
						utils_types.FilePath(front.ShipDetailedUrl(ship, front.ShipShowDetails)),
						front.ShipDetails(ship),
					),
					builder.NewComponent(
						utils_types.FilePath(front.ShipPinnedUrl(ship, front.ShipShowBases)),
						front.ShipRow(ship, front.ShipShowBases, front.PinMode, disco_ids, true),
					),
					builder.NewComponent(
						utils_types.FilePath(front.ShipPinnedUrl(ship, front.ShipShowDetails)),
						front.ShipRow(ship, front.ShipShowDetails, front.PinMode, disco_ids, true),
					),
				)

				//  id IDTractor

			}

			for _, tractor := range data.Tractors {
				build.RegComps(
					builder.NewComponent(
						utils_types.FilePath(front.TractorDetailedUrl(tractor)),
						front.GoodAtBaseInfoT(tractor.Name, tractor.Bases, front.ShowPricePerVolume(false)),
					),
				)
			}

			for _, engine := range data.Engines {
				build.RegComps(
					builder.NewComponent(
						utils_types.FilePath(front.EngineDetailedUrl(engine)),
						front.GoodAtBaseInfoT(engine.Name, engine.Bases, front.ShowPricePerVolume(false)),
					),
				)
			}

			for _, cm := range data.CMs {
				build.RegComps(
					builder.NewComponent(
						utils_types.FilePath(front.CounterMeasreDetailedUrl(cm)),
						front.GoodAtBaseInfoT(cm.Name, cm.Bases, front.ShowPricePerVolume(false)),
					),
				)
			}

			for _, item := range data.Scanners {
				build.RegComps(
					builder.NewComponent(
						utils_types.FilePath(front.ScannerDetailedUrl(item)),
						front.GoodAtBaseInfoT(item.Name, item.Bases, front.ShowPricePerVolume(false)),
					),
				)
			}
		}, time_measure.WithMsg("linking most of stuff"))

	}, time_measure.WithMsg("link, internal measure"))

	return build
}
