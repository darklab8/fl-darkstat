package linker

/*
Links data from exported fl-configs
into stuff rendered by fl-darkstat
*/

import (
	"sort"
	"time"

	"github.com/darklab8/fl-configs/configs/cfgtype"
	"github.com/darklab8/fl-configs/configs/configs_export"
	"github.com/darklab8/fl-configs/configs/configs_mapped"
	"github.com/darklab8/fl-darkcore/darkcore/builder"
	"github.com/darklab8/fl-darkcore/darkcore/core_static"
	"github.com/darklab8/fl-darkstat/darkstat/front"
	"github.com/darklab8/fl-darkstat/darkstat/front/static_front"
	"github.com/darklab8/fl-darkstat/darkstat/front/types"
	"github.com/darklab8/fl-darkstat/darkstat/front/urls"
	"github.com/darklab8/fl-darkstat/darkstat/settings"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
	"github.com/darklab8/go-utils/utils/timeit"
	"github.com/darklab8/go-utils/utils/utils_logus"
	"github.com/darklab8/go-utils/utils/utils_types"
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

	timeit.NewTimerF(func() {
		freelancer_folder := settings.Env.FreelancerFolder
		if l.configs == nil {
			l.mapped = configs_mapped.NewMappedConfigs()
			logus.Log.Debug("scanning freelancer folder", utils_logus.FilePath(freelancer_folder))
			l.mapped.Read(freelancer_folder)
			l.configs = configs_export.NewExporter(l.mapped)
		}
	}, timeit.WithMsg("MappedConfigs creation"))
	return l
}

func (l *Linker) Link() *builder.Builder {
	var build *builder.Builder
	defer timeit.NewTimer("link, internal measure").Close()
	timer_building_creation := timeit.NewTimer("building creation")
	tractor_tab_name := settings.Env.TractorTabName
	if l.mapped.Discovery != nil {
		tractor_tab_name = "IDs"
	}
	staticPrefix := "static/"
	siteRoot := settings.Env.SiteRoot
	params := &types.GlobalParams{
		Buildpath: "",
		Theme:     types.ThemeLight,
		Themes: []string{
			siteRoot + urls.Index.ToString(),
			siteRoot + urls.DarkIndex.ToString(),
			siteRoot + urls.VanillaIndex.ToString(),
		},
		TractorTabName: tractor_tab_name,
		SiteRoot:       siteRoot,
		StaticRoot:     siteRoot + staticPrefix,
		Heading:        settings.Env.AppHeading,
		Timestamp:      time.Now().UTC(),
	}

	static_files := []builder.StaticFile{
		builder.NewStaticFileFromCore(core_static.HtmxJS),
		builder.NewStaticFileFromCore(core_static.HtmxPreloadJS),
		builder.NewStaticFileFromCore(core_static.SortableJS),
		builder.NewStaticFileFromCore(core_static.ResetCSS),
		builder.NewStaticFileFromCore(core_static.FaviconIco),

		builder.NewStaticFileFromCore(static_front.CommonCSS),
		builder.NewStaticFileFromCore(static_front.CustomCSS),
		builder.NewStaticFileFromCore(static_front.CustomJS),
		builder.NewStaticFileFromCore(static_front.CustomJSResizer),
		builder.NewStaticFileFromCore(static_front.CustomJSFiltering),
		builder.NewStaticFileFromCore(static_front.CustomJSFilteringRoutes),
	}

	build = builder.NewBuilder(params, static_files)

	timer_building_creation.Close()

	var data *configs_export.Exporter
	timeit.NewTimerMF("exporting data", func() { data = l.configs.Export() })

	timeit.NewTimerMF("sorting completed", func() {
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

			sort.Slice(base.TradeRoutes, func(i, j int) bool {
				return base.TradeRoutes[i].Transport.GetProffitPerTime() > base.TradeRoutes[j].Transport.GetProffitPerTime()
			})
		}

		sort.Slice(data.Factions, func(i, j int) bool {
			if data.Factions[i].Name != "" && data.Factions[j].Name == "" {
				return true
			}
			return data.Factions[i].Name < data.Factions[j].Name
		})

		for fac_index, faction := range data.Factions {
			var reps []configs_export.Reputation = make([]configs_export.Reputation, 0, len(faction.Reputations))
			for _, rep := range faction.Reputations {
				if rep.Name != "" {
					reps = append(reps, rep)
				}
			}
			sort.Slice(reps, func(i, j int) bool {
				return reps[i].Rep > reps[j].Rep
			})
			data.Factions[fac_index].Reputations = reps
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

	})

	var useful_factions []configs_export.Faction
	var useful_ships []configs_export.Ship
	var useful_guns []configs_export.Gun
	var useful_missiles []configs_export.Gun
	tractor_id := cfgtype.TractorID("")

	var disco_ids types.DiscoveryIDs

	timeit.NewTimerMF("filtering to useful stuff", func() {
		useful_factions = configs_export.FilterToUsefulFactions(data.Factions)
		useful_ships = data.FilterToUsefulShips(data.Ships)
		useful_guns = data.FilterToUsefulGun(data.Guns)
		useful_missiles = data.FilterToUsefulGun(data.Missiles)

		if l.mapped.Discovery != nil {
			disco_ids = types.DiscoveryIDs{
				Show:         true,
				Ids:          l.configs.Tractors,
				TractorsByID: l.configs.TractorsByID,
				Config:       l.mapped.Discovery.Techcompat,
				LatestPatch:  l.mapped.Discovery.LatestPatch,
			}
		}
		disco_ids.Infocards = l.configs.Infocards
	})

	timeit.NewTimerMF("linking main stuff", func() {

		build.RegComps(
			builder.NewComponent(
				urls.Ships+utils_types.FilePath(tractor_id),
				front.ShipsT(useful_ships, front.ShipShowBases, front.ShowEmpty(false), disco_ids, data.Infocards),
			),
			builder.NewComponent(
				front.AllItemsUrl(urls.Ships)+utils_types.FilePath(tractor_id),
				front.ShipsT(data.Ships, front.ShipShowBases, front.ShowEmpty(true), disco_ids, data.Infocards),
			),
			builder.NewComponent(
				urls.ShipDetails+utils_types.FilePath(tractor_id),
				front.ShipsT(useful_ships, front.ShipShowDetails, front.ShowEmpty(false), disco_ids, data.Infocards),
			),
			builder.NewComponent(
				front.AllItemsUrl(urls.ShipDetails)+utils_types.FilePath(tractor_id),
				front.ShipsT(data.Ships, front.ShipShowDetails, front.ShowEmpty(true), disco_ids, data.Infocards),
			),
		)

		build.RegComps(

			builder.NewComponent(
				urls.Index,
				front.Index(types.ThemeLight),
			),
			builder.NewComponent(
				urls.DarkIndex,
				front.Index(types.ThemeDark),
			),
			builder.NewComponent(
				urls.VanillaIndex,
				front.Index(types.ThemeVanilla),
			),
			builder.NewComponent(
				urls.Bases,
				front.BasesT(configs_export.FilterToUserfulBases(data.Bases), front.BaseShowShops, front.ShowEmpty(false), disco_ids),
			),
			builder.NewComponent(
				front.AllItemsUrl(urls.Bases),
				front.BasesT(data.Bases, front.BaseShowShops, front.ShowEmpty(true), disco_ids),
			),
			builder.NewComponent(
				urls.Missions,
				front.BasesT(configs_export.FilterToUserfulBases(data.Bases), front.BaseShowMissions, front.ShowEmpty(false), disco_ids),
			),
			builder.NewComponent(
				front.AllItemsUrl(urls.Missions),
				front.BasesT(data.Bases, front.BaseShowMissions, front.ShowEmpty(true), disco_ids),
			),
			builder.NewComponent(
				urls.Factions,
				front.FactionsT(useful_factions, front.FactionShowBases, front.ShowEmpty(false), disco_ids),
			),
			builder.NewComponent(
				front.AllItemsUrl(urls.Factions),
				front.FactionsT(data.Factions, front.FactionShowBases, front.ShowEmpty(true), disco_ids),
			),
			builder.NewComponent(
				urls.Rephacks,
				front.FactionsT(useful_factions, front.FactionShowRephacks, front.ShowEmpty(false), disco_ids),
			),
			builder.NewComponent(
				front.AllItemsUrl(urls.Rephacks),
				front.FactionsT(data.Factions, front.FactionShowRephacks, front.ShowEmpty(true), disco_ids),
			),
			builder.NewComponent(
				urls.Commodities,
				front.CommoditiesT(data.FilterToUsefulCommodities(data.Commodities), front.ShowEmpty(false), disco_ids),
			),
			builder.NewComponent(
				front.AllItemsUrl(urls.Commodities),
				front.CommoditiesT(data.Commodities, front.ShowEmpty(true), disco_ids),
			),
			builder.NewComponent(
				urls.Guns,
				front.GunsT(useful_guns, front.GunsShowBases, front.ShowEmpty(false), disco_ids, data.Infocards),
			),
			builder.NewComponent(
				front.AllItemsUrl(urls.Guns),
				front.GunsT(data.Guns, front.GunsShowBases, front.ShowEmpty(true), disco_ids, data.Infocards),
			),
			builder.NewComponent(
				urls.GunModifiers,
				front.GunsT(useful_guns, front.GunsShowDamageBonuses, front.ShowEmpty(false), disco_ids, data.Infocards),
			),
			builder.NewComponent(
				front.AllItemsUrl(urls.GunModifiers),
				front.GunsT(data.Guns, front.GunsShowDamageBonuses, front.ShowEmpty(true), disco_ids, data.Infocards),
			),
			builder.NewComponent(
				urls.Ammo,
				front.AmmoT(data.FilterToUsefulAmmo(data.Ammos), front.ShowEmpty(false), disco_ids),
			),
			builder.NewComponent(
				front.AllItemsUrl(urls.Ammo),
				front.AmmoT(data.Ammos, front.ShowEmpty(true), disco_ids),
			),
			builder.NewComponent(
				urls.Missiles,
				front.GunsT(useful_missiles, front.GunsMissiles, front.ShowEmpty(false), disco_ids, data.Infocards),
			),
			builder.NewComponent(
				front.AllItemsUrl(urls.Missiles),
				front.GunsT(data.Missiles, front.GunsMissiles, front.ShowEmpty(true), disco_ids, data.Infocards),
			),
			builder.NewComponent(
				urls.Mines,
				front.MinesT(data.FilterToUsefulMines(data.Mines), front.ShowEmpty(false), disco_ids),
			),
			builder.NewComponent(
				front.AllItemsUrl(urls.Mines),
				front.MinesT(data.Mines, front.ShowEmpty(true), disco_ids),
			),
			builder.NewComponent(
				urls.Shields,
				front.ShieldT(data.FilterToUsefulShields(data.Shields), front.ShowEmpty(false), disco_ids),
			),
			builder.NewComponent(
				front.AllItemsUrl(urls.Shields),
				front.ShieldT(data.Shields, front.ShowEmpty(true), disco_ids),
			),
			builder.NewComponent(
				urls.Thrusters,
				front.ThrusterT(data.FilterToUsefulThrusters(data.Thrusters), front.ShowEmpty(false), disco_ids),
			),
			builder.NewComponent(
				front.AllItemsUrl(urls.Thrusters),
				front.ThrusterT(data.Thrusters, front.ShowEmpty(true), disco_ids),
			),
			builder.NewComponent(
				urls.Tractors,
				front.TractorsT(data.FilterToUsefulTractors(data.Tractors), front.ShowEmpty(false), front.TractorModShop, disco_ids),
			),
			builder.NewComponent(
				front.AllItemsUrl(urls.Tractors),
				front.TractorsT(data.Tractors, front.ShowEmpty(true), front.TractorModShop, disco_ids),
			),
			builder.NewComponent(
				urls.IDRephacks,
				front.TractorsT(data.FilterToUsefulTractors(data.Tractors), front.ShowEmpty(false), front.TractorIDRephacks, disco_ids),
			),
			builder.NewComponent(
				front.AllItemsUrl(urls.IDRephacks),
				front.TractorsT(data.Tractors, front.ShowEmpty(true), front.TractorIDRephacks, disco_ids),
			),
			builder.NewComponent(
				urls.Engines,
				front.Engines(data.FilterToUsefulEngines(data.Engines), front.ShowEmpty(false), disco_ids),
			),
			builder.NewComponent(
				front.AllItemsUrl(urls.Engines),
				front.Engines(data.Engines, front.ShowEmpty(true), disco_ids),
			),
			builder.NewComponent(
				urls.CounterMeasures,
				front.CounterMeasureT(data.FilterToUsefulCounterMeasures(data.CMs), front.ShowEmpty(false), disco_ids),
			),
			builder.NewComponent(
				front.AllItemsUrl(urls.CounterMeasures),
				front.CounterMeasureT(data.CMs, front.ShowEmpty(true), disco_ids),
			),
			builder.NewComponent(
				urls.Scanners,
				front.ScannersT(data.FilterToUserfulScanners(data.Scanners), front.ShowEmpty(false), disco_ids),
			),
			builder.NewComponent(
				front.AllItemsUrl(urls.Scanners),
				front.ScannersT(data.Scanners, front.ShowEmpty(true), disco_ids),
			),
		)

		for _, base := range data.Bases {
			build.RegComps(
				builder.NewComponent(
					utils_types.FilePath(front.BaseDetailedUrl(base, front.BaseShowShops)),
					front.BaseMarketGoods(base.Name, base.MarketGoods, front.BaseShowShops),
				),
				builder.NewComponent(
					utils_types.FilePath(front.BaseDetailedUrl(base, front.BaseShowMissions)),
					front.BaseMissions(base.Name, base.Missions, front.BaseShowMissions),
				),
				builder.NewComponent(
					utils_types.FilePath(front.BaseDetailedUrl(base, front.BaseTabTrades)),
					front.BaseTrades(base.Name, base.Trades, front.BaseTabTrades, disco_ids),
				),
			)

			for _, combo_route := range base.TradeRoutes {

				build.RegComps(
					builder.NewComponent(
						utils_types.FilePath(front.TradeRouteUrl(combo_route)),
						front.TradeRouteInfo(combo_route, disco_ids),
					),
				)
			}
		}

		for _, base := range data.MiningOperations {
			build.RegComps(
				builder.NewComponent(
					utils_types.FilePath(front.BaseDetailedUrl(base, front.BaseTabOres)),
					front.BaseTrades(base.Name, base.Trades, front.BaseTabOres, disco_ids),
				),
			)

			for _, combo_route := range base.TradeRoutes {

				build.RegComps(
					builder.NewComponent(
						utils_types.FilePath(front.TradeRouteUrl(combo_route)),
						front.TradeRouteInfo(combo_route, disco_ids),
					),
				)
			}

		}
	})

	timeit.NewTimerMF("linking faction stuff", func() {
		for _, faction := range data.Factions {
			build.RegComps(
				builder.NewComponent(
					utils_types.FilePath(front.FactionRepUrl(faction, front.FactionShowBases)),
					front.FactionReps(faction, faction.Reputations),
				),
				builder.NewComponent(
					utils_types.FilePath(front.FactionRepUrl(faction, front.FactionShowRephacks)),
					front.RephackBottom(faction, faction.Bribes),
				),
			)
		}
	})

	timeit.NewTimerMF("linking most of stuff", func() {
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
					front.GoodAtBaseInfoT(base_info.Name, base_info.Bases, front.ShowAsCommodity(true)),
				),
			)
		}

		for _, gun := range data.Guns {
			build.RegComps(
				builder.NewComponent(
					utils_types.FilePath(front.GunDetailedUrl(gun, front.GunsShowBases)),
					front.GoodAtBaseInfoT(gun.Name, gun.Bases, front.ShowAsCommodity(false)),
				),
				builder.NewComponent(
					utils_types.FilePath(front.GunDetailedUrl(gun, front.GunsShowDamageBonuses)),
					front.GunShowModifiers(gun),
				),

				builder.NewComponent(
					utils_types.FilePath(front.GunPinnedRowUrl(gun, front.GunsShowBases)),
					front.GunRow(gun, front.GunsShowBases, front.PinMode, disco_ids, data.Infocards, true),
				),
				builder.NewComponent(
					utils_types.FilePath(front.GunPinnedRowUrl(gun, front.GunsShowDamageBonuses)),
					front.GunRow(gun, front.GunsShowDamageBonuses, front.PinMode, disco_ids, data.Infocards, true),
				),
			)
		}
		for _, ammo := range data.Ammos {
			build.RegComps(
				builder.NewComponent(
					utils_types.FilePath(front.AmmoDetailedUrl(ammo)),
					front.GoodAtBaseInfoT(ammo.Name, ammo.Bases, front.ShowAsCommodity(false)),
				),
			)
		}

		for _, missile := range data.Missiles {
			build.RegComps(
				builder.NewComponent(
					utils_types.FilePath(front.GunDetailedUrl(missile, front.GunsMissiles)),
					front.GoodAtBaseInfoT(missile.Name, missile.Bases, front.ShowAsCommodity(false)),
				),
				builder.NewComponent(
					utils_types.FilePath(front.GunPinnedRowUrl(missile, front.GunsMissiles)),
					front.GunRow(missile, front.GunsMissiles, front.PinMode, disco_ids, data.Infocards, true),
				),
			)
		}

		for _, mine := range data.Mines {
			build.RegComps(
				builder.NewComponent(
					utils_types.FilePath(front.MineDetailedUrl(mine)),
					front.GoodAtBaseInfoT(mine.Name, mine.Bases, front.ShowAsCommodity(false)),
				),
			)
		}

		for _, shield := range data.Shields {
			build.RegComps(
				builder.NewComponent(
					utils_types.FilePath(front.ShieldDetailedUrl(shield)),
					front.GoodAtBaseInfoT(shield.Name, shield.Bases, front.ShowAsCommodity(false)),
				),
			)
		}

		for _, thruster := range data.Thrusters {
			build.RegComps(
				builder.NewComponent(
					utils_types.FilePath(front.ThrusterDetailedUrl(thruster)),
					front.GoodAtBaseInfoT(thruster.Name, thruster.Bases, front.ShowAsCommodity(false)),
				),
			)
		}

		for _, ship := range data.Ships {
			build.RegComps(
				builder.NewComponent(
					utils_types.FilePath(front.ShipDetailedUrl(ship, front.ShipShowBases)),
					front.GoodAtBaseInfoT(ship.Name, ship.Bases, front.ShowAsCommodity(false)),
				),
				builder.NewComponent(
					utils_types.FilePath(front.ShipDetailedUrl(ship, front.ShipShowDetails)),
					front.ShipDetails(ship),
				),
				builder.NewComponent(
					utils_types.FilePath(front.ShipPinnedUrl(ship, front.ShipShowBases)),
					front.ShipRow(ship, front.ShipShowBases, front.PinMode, disco_ids, data.Infocards, true),
				),
				builder.NewComponent(
					utils_types.FilePath(front.ShipPinnedUrl(ship, front.ShipShowDetails)),
					front.ShipRow(ship, front.ShipShowDetails, front.PinMode, disco_ids, data.Infocards, true),
				),
			)
		}

		for _, tractor := range data.Tractors {
			build.RegComps(
				builder.NewComponent(
					utils_types.FilePath(front.TractorDetailedUrl(tractor, front.TractorModShop)),
					front.GoodAtBaseInfoT(tractor.Name, tractor.Bases, front.ShowAsCommodity(false)),
				),
				builder.NewComponent(
					utils_types.FilePath(front.TractorDetailedUrl(tractor, front.TractorIDRephacks)),
					front.IDRephacksT(tractor),
				),
			)
		}

		for _, engine := range data.Engines {
			build.RegComps(
				builder.NewComponent(
					utils_types.FilePath(front.EngineDetailedUrl(engine)),
					front.GoodAtBaseInfoT(engine.Name, engine.Bases, front.ShowAsCommodity(false)),
				),
			)
		}

		for _, cm := range data.CMs {
			build.RegComps(
				builder.NewComponent(
					utils_types.FilePath(front.CounterMeasreDetailedUrl(cm)),
					front.GoodAtBaseInfoT(cm.Name, cm.Bases, front.ShowAsCommodity(false)),
				),
			)
		}

		for _, item := range data.Scanners {
			build.RegComps(
				builder.NewComponent(
					utils_types.FilePath(front.ScannerDetailedUrl(item)),
					front.GoodAtBaseInfoT(item.Name, item.Bases, front.ShowAsCommodity(false)),
				),
			)
		}

		sort.Slice(data.Bases, func(i, j int) bool {
			if data.Bases[j].BestTransportRoute == nil {
				return true
			}
			if data.Bases[i].BestTransportRoute == nil {
				return false
			}
			return data.Bases[i].BestTransportRoute.GetProffitPerTime() > data.Bases[j].BestTransportRoute.GetProffitPerTime()
		})

		build.RegComps(
			builder.NewComponent(
				urls.Trades,
				front.BasesT(configs_export.FilterToUserfulBases(data.Bases), front.BaseTabTrades, front.ShowEmpty(false), disco_ids),
			),
			builder.NewComponent(
				front.AllItemsUrl(urls.Trades),
				front.BasesT(data.Bases, front.BaseTabTrades, front.ShowEmpty(true), disco_ids),
			),
			builder.NewComponent(
				urls.Asteroids,
				front.BasesT(configs_export.FitlerToUsefulOres(data.MiningOperations), front.BaseTabOres, front.ShowEmpty(false), disco_ids),
			),
			builder.NewComponent(
				front.AllItemsUrl(urls.Asteroids),
				front.BasesT(data.MiningOperations, front.BaseTabOres, front.ShowEmpty(true), disco_ids),
			),
		)
	})

	return build
}
