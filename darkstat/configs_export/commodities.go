package configs_export

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped"
	"github.com/darklab8/fl-darkstat/darkcore/settings/traces"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export/infocarder"
	"github.com/darklab8/go-utils/utils/ptr"
)

type MarketGood struct {
	GoodInfo

	LevelRequired        int            `json:"level_required" validate:"required"`
	RepRequired          float64        `json:"rep_required" validate:"required"`
	PriceBaseBuysFor     *int           `json:"price_base_buys_for"`
	PriceBaseSellsFor    int            `json:"price_base_sells_for" validate:"required"`
	Volume               float64        `json:"volume" validate:"required"`
	ShipClass            *cfg.ShipClass `json:"ship_class"` // Discovery specific value. Volume can be different based on ship class. Duplicating market goods with different volumes for that
	BaseSells            bool           `json:"base_sells" validate:"required"`
	IsServerSideOverride bool           `json:"is_server_override" validate:"required"`

	NotBuyable             bool `json:"-" swaggerignore:"true"`
	IsTransportUnreachable bool `json:"-" swaggerignore:"true"`

	DiscoveryFactoryName *string `json:"-" swaggerignore:"true"`

	PoBGood *ShopItem
	PoB     *PoB

	BaseInfo
}

func (g MarketGood) GetPriceBaseBuysFor() int {
	if g.PriceBaseBuysFor == nil {
		return 0
	}
	return *g.PriceBaseBuysFor
}

type Commodity struct {
	Nickname              string                          `json:"nickname"  validate:"required"`
	PriceBase             int                             `json:"price_base"  validate:"required"`
	Name                  string                          `json:"name"  validate:"required"`
	Combinable            bool                            `json:"combinable"  validate:"required"`
	Volume                float64                         `json:"volume"  validate:"required"`
	ShipClass             *cfg.ShipClass                  `json:"ship_class"`
	NameID                int                             `json:"name_id"  validate:"required"`
	InfocardID            int                             `json:"infocard_id"  validate:"required"`
	Infocard              infocarder.InfocardKey          `json:"infocard_key"  validate:"required"`
	Bases                 map[cfg.BaseUniNick]*MarketGood `json:"-" swaggerignore:"true"`
	PriceBestBaseBuysFor  int                             `json:"price_best_base_buys_for"  validate:"required"`
	PriceBestBaseSellsFor int                             `json:"price_best_base_sells_for"  validate:"required"`
	ProffitMargin         int                             `json:"proffit_margin"  validate:"required"`
	Mass                  float64                         `json:"mass"  validate:"required"`
}

func (b Commodity) GetNickname() string { return string(b.Nickname) }

func (b Commodity) GetBases() map[cfg.BaseUniNick]*MarketGood { return b.Bases }

func GetPricePerVoume(price int, volume float64) float64 {
	if volume == 0 {
		return -1
	}
	return float64(price) / float64(volume)
}

func (e *Exporter) GetCommodities(ctx context.Context) []*Commodity {
	ctx, span := traces.Tracer.Start(ctx, "Exporter.GetCommodities")
	defer span.End()
	commodities := make([]*Commodity, 0, 100)

	for _, comm := range e.Mapped.Goods.Commodities {
		equipment_name := comm.Equipment.Get()
		equipment := e.Mapped.Equip().CommoditiesMap[equipment_name]

		for _, volume_info := range equipment.Volumes {
			commodity := &Commodity{
				Bases:     make(map[cfg.BaseUniNick]*MarketGood),
				PriceBase: comm.Price.Get(),
			}
			commodity.Mass, _ = equipment.Mass.GetValue()

			commodity.Nickname = comm.Nickname.Get()

			commodity.Combinable = comm.Combinable.Get()

			commodity.NameID = equipment.IdsName.Get()

			commodity.Name = e.GetInfocardName(equipment.IdsName.Get(), commodity.Nickname)
			e.exportInfocards(commodity.Infocard, equipment.IdsInfo.Get())
			commodity.InfocardID = equipment.IdsInfo.Get()

			commodity.Volume = volume_info.Volume.Get()
			commodity.ShipClass = volume_info.GetShipClass()
			commodity.Infocard = infocarder.InfocardKey(commodity.Nickname)

			base_item_price := comm.Price.Get()

			commodity.Bases = e.GetAtBasesSold(GetCommodityAtBasesInput{
				Nickname:  commodity.Nickname,
				Price:     base_item_price,
				Volume:    commodity.Volume,
				ShipClass: commodity.ShipClass,
			})

			for _, base_info := range commodity.Bases {
				if base_info.GetPriceBaseBuysFor() > commodity.PriceBestBaseBuysFor {
					commodity.PriceBestBaseBuysFor = base_info.GetPriceBaseBuysFor()
				}
				if base_info.PriceBaseSellsFor < commodity.PriceBestBaseSellsFor || commodity.PriceBestBaseSellsFor == 0 {
					if base_info.BaseSells && base_info.PriceBaseSellsFor > 0 {
						commodity.PriceBestBaseSellsFor = base_info.PriceBaseSellsFor
					}

				}
			}

			if commodity.PriceBestBaseBuysFor > 0 && commodity.PriceBestBaseSellsFor > 0 {
				commodity.ProffitMargin = commodity.PriceBestBaseBuysFor - commodity.PriceBestBaseSellsFor
			}

			var infocard_addition infocarder.InfocardBuilder
			if e.Mapped.Discovery != nil {
				if player_bonuses, ok := e.Mapped.Discovery.Minecontrol.PlayerBonusByOreNickname[commodity.Nickname]; ok {
					infocard_addition.WriteLineStr(`MINING BONUSES (darkstat):`)
					for _, player_bonus := range player_bonuses {
						id_nickname := player_bonus.IDNickname.Get()
						id_name := id_nickname
						if tractor, ok := e.Mapped.Equip().TractorsMap[id_nickname]; ok {
							if name_id, ok := tractor.IdsName.GetValue(); ok {
								id_name = e.GetInfocardName(name_id, string(id_nickname))
							}
						}
						infocard_addition.WriteLineStr(id_name, "= ", strconv.FormatFloat(player_bonus.Bonus.Get(), 'f', 2, 64))
					}
					infocard_addition.WriteLineStr("")
				}
			}

			var info infocarder.InfocardBuilder
			if value, ok := e.GetInfocard2(infocarder.InfocardKey(commodity.Nickname)); ok {
				info.Lines = value
			}
			if !strings.Contains(info.Lines.StringsJoin(""), "MINING BONUSES") {
				e.PutInfocard(infocarder.InfocardKey(commodity.Nickname), append(info.Lines, infocard_addition.Lines...))
			}
			commodities = append(commodities, commodity)
		}

	}

	return commodities
}

type GetCommodityAtBasesInput struct {
	Nickname  string
	Price     int
	Volume    float64
	ShipClass *cfg.ShipClass
}

func (e *Exporter) ServerSideMarketGoodsOverrides(commodity GetCommodityAtBasesInput) map[cfg.BaseUniNick]*MarketGood {
	var bases_already_found map[cfg.BaseUniNick]*MarketGood = make(map[cfg.BaseUniNick]*MarketGood)

	for _, base_market := range e.Mapped.Discovery.Prices.BasesPerGood[commodity.Nickname] {
		base_nickname := cfg.BaseUniNick(base_market.BaseNickname.Get())

		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in f", r)
				fmt.Println("recovered base_nickname", base_nickname)
				fmt.Println("recovered commodity nickname", commodity.Nickname)
				panic(r)
			}
		}()

		var base_info *MarketGood = &MarketGood{
			GoodInfo:             e.GetGoodInfo(commodity.Nickname),
			BaseInfo:             e.GetBaseInfo(universe_mapped.BaseNickname(base_nickname)),
			NotBuyable:           false,
			BaseSells:            base_market.BaseSells.Get(),
			PriceBaseBuysFor:     ptr.Ptr(base_market.PriceBaseBuysFor.Get()),
			PriceBaseSellsFor:    base_market.PriceBaseSellsFor.Get(),
			Volume:               commodity.Volume,
			ShipClass:            commodity.ShipClass,
			IsServerSideOverride: true,
		}

		if e.useful_bases_by_nick != nil {
			if _, ok := e.useful_bases_by_nick[base_info.BaseNickname]; !ok {
				base_info.NotBuyable = true
			}
		}

		bases_already_found[base_info.BaseNickname] = base_info
	}
	return bases_already_found
}

func (e *Exporter) GetAtBasesSold(commodity GetCommodityAtBasesInput) map[cfg.BaseUniNick]*MarketGood {
	var goods_per_base map[cfg.BaseUniNick]*MarketGood = make(map[cfg.BaseUniNick]*MarketGood)

	for _, base_market := range e.Mapped.Market().BasesPerGood[commodity.Nickname] {
		base_nickname := base_market.Base

		market_good := base_market.MarketGood
		base_info := &MarketGood{
			GoodInfo:   e.GetGoodInfo(commodity.Nickname),
			BaseInfo:   e.GetBaseInfo(universe_mapped.BaseNickname(base_nickname)),
			NotBuyable: false,
			Volume:     commodity.Volume,
			ShipClass:  commodity.ShipClass,
		}
		base_info.BaseSells = market_good.BaseSells()
		base_info.BaseNickname = base_nickname

		base_info.PriceBaseSellsFor = int(market_good.PriceModifier.Get() * float64(commodity.Price))

		if e.Mapped.Discovery != nil {
			base_info.PriceBaseBuysFor = ptr.Ptr(market_good.BaseSellsIPositiveAndDiscoSellPrice.Get())
		} else {
			base_info.PriceBaseBuysFor = ptr.Ptr(base_info.PriceBaseSellsFor)
		}

		base_info.LevelRequired = market_good.LevelRequired.Get()
		base_info.RepRequired = market_good.RepRequired.Get()

		base_info.BaseInfo = e.GetBaseInfo(universe_mapped.BaseNickname(base_info.BaseNickname))

		if e.useful_bases_by_nick != nil {
			if _, ok := e.useful_bases_by_nick[base_info.BaseNickname]; !ok {
				base_info.NotBuyable = true
			}
		}

		goods_per_base[base_info.BaseNickname] = base_info
	}

	if e.Mapped.Discovery != nil {
		serverside_overrides := e.ServerSideMarketGoodsOverrides(commodity)
		for _, item := range serverside_overrides {
			goods_per_base[item.BaseNickname] = item
		}

	}
	if e.Mapped.Discovery != nil || e.Mapped.FLSR != nil {
		pob_produced := e.pob_produced()
		if _, ok := pob_produced[commodity.Nickname]; ok {
			good_to_add := &MarketGood{
				GoodInfo:             e.GetGoodInfo(commodity.Nickname),
				BaseSells:            true,
				IsServerSideOverride: true,
				Volume:               commodity.Volume,
				ShipClass:            commodity.ShipClass,
				BaseInfo: BaseInfo{
					BaseNickname: pob_crafts_nickname,
					BaseName:     e.Mapped.CraftableBaseName(),
					SystemName:   "Neverwhere",
					Region:       "NEVERWHERE",
					FactionName:  "Neverwhere",
				},
			}
			goods_per_base[pob_crafts_nickname] = good_to_add

		}
	}

	loot_findable := e.findable_in_loot()
	if _, ok := loot_findable[commodity.Nickname]; ok {
		good_to_add := &MarketGood{
			GoodInfo:             e.GetGoodInfo(commodity.Nickname),
			BaseSells:            true,
			IsServerSideOverride: false,
			Volume:               commodity.Volume,
			ShipClass:            commodity.ShipClass,
			BaseInfo: BaseInfo{
				BaseNickname: BaseLootableNickname,
				BaseName:     BaseLootableName,
				SystemName:   "Neverwhere",
				Region:       "NEVERWHERE",
				FactionName:  BaseLootableFaction,
			},
		}
		goods_per_base[BaseLootableNickname] = good_to_add
	}

	// this thing propagates data to commodities. hmm. needs to be refreshing.
	// Also missing defined PoB good.
	if e.Mapped.Discovery != nil {
		pob_buyable := e.get_pob_buyable()
		if goods, ok := pob_buyable[commodity.Nickname]; ok {
			for _, good := range goods {
				good_to_add := &MarketGood{
					GoodInfo:             e.GetGoodInfo(commodity.Nickname),
					BaseSells:            good.Quantity > good.MinStock,
					IsServerSideOverride: true,
					PriceBaseBuysFor:     ptr.Ptr(good.PriceBaseBuysFor),
					PriceBaseSellsFor:    good.PriceBaseSellsFor,
					Volume:               commodity.Volume,
					ShipClass:            commodity.ShipClass,
					BaseInfo: BaseInfo{
						BaseNickname: cfg.BaseUniNick(good.PobNickname),
						BaseName:     "(PoB) " + good.PoBName,
						SystemName:   good.SystemName,
						FactionName:  good.FactionName,
					},
				}

				if good.System != nil {
					good_to_add.BaseInfo.Region = e.GetRegionName(good.System)
				}
				if good.BasePos != nil && good.System != nil {
					good_to_add.BasePos = *good.BasePos
					good_to_add.SectorCoord = VectorToSectorCoord(good.System, *good.BasePos)
				}
				goods_per_base[cfg.BaseUniNick(good.PobNickname)] = good_to_add
			}
		}
	}

	for _, good := range goods_per_base {
		good.GoodInfo = e.GetGoodInfo(commodity.Nickname)
	}

	for base_nickname, good := range goods_per_base {
		if !e.TraderExists(string(base_nickname)) {
			if good.Category == "commodity" {
				delete(goods_per_base, base_nickname)
			}
		}
	}

	return goods_per_base
}

type BaseInfo struct {
	BaseNickname cfg.BaseUniNick `json:"base_nickname" validate:"required"`
	BaseName     string          `json:"base_name" validate:"required"`
	SystemName   string          `json:"system_name" validate:"required"`
	Region       string          `json:"region_name" validate:"required"`
	FactionName  string          `json:"faction_name" validate:"required"`
	BasePos      cfg.Vector      `json:"base_pos" validate:"required"`
	SectorCoord  string          `json:"sector_coord" validate:"required"`
}

func (e *Exporter) GetRegionName(system *universe_mapped.System) string {
	return e.Mapped.GetRegionName(system)
}
func (e *ExporterRelay) GetRegionName(system *universe_mapped.System) string {
	return e.Mapped.GetRegionName(system)
}

func (e *Exporter) GetBaseInfo(base_nickname universe_mapped.BaseNickname) BaseInfo {
	var result BaseInfo = BaseInfo{
		BaseNickname: cfg.BaseUniNick(base_nickname),
	}
	universe_base, found_universe_base := e.Mapped.Universe.BasesMap[universe_mapped.BaseNickname(base_nickname)]

	if !found_universe_base {
		return result
	}

	result.BaseName = e.GetInfocardName(universe_base.StridName.Get(), string(base_nickname))
	system_nickname := universe_base.System.Get()

	system, system_ok := e.Mapped.Universe.SystemMap[universe_mapped.SystemNickname(system_nickname)]
	if system_ok {
		result.SystemName = e.GetInfocardName(system.StridName.Get(), system_nickname)
		result.Region = e.GetRegionName(system)
	}

	var reputation_nickname string
	if system, ok := e.Mapped.Systems.SystemsMap[universe_base.System.Get()]; ok {
		for _, system_base := range system.Bases {
			if system_base.IdsName.Get() != universe_base.StridName.Get() {
				continue
			}

			reputation_nickname = system_base.RepNickname.Get()
			result.BasePos = system_base.Pos.Get()
		}

	}

	result.SectorCoord = VectorToSectorCoord(system, result.BasePos)

	var factionName string
	if group, exists := e.Mapped.InitialWorld.GroupsMap[reputation_nickname]; exists {
		factionName = e.GetInfocardName(group.IdsName.Get(), reputation_nickname)
	}

	result.FactionName = factionName

	return result
}

func (e *Exporter) FilterToUsefulCommodities(commodities []*Commodity) []*Commodity {
	var items []*Commodity = make([]*Commodity, 0, len(commodities))
	for _, item := range commodities {
		if !e.Buyable(item.Bases) {
			continue
		}
		items = append(items, item)
	}
	return items
}
