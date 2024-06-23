package configs_export

import (
	"fmt"

	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped"
	"github.com/darklab8/fl-configs/configs/conftypes"
)

type GoodAtBase struct {
	BaseNickname      string
	BaseName          string
	BaseSells         bool
	PriceBaseBuysFor  int
	PriceBaseSellsFor int
	Volume            float64
	LevelRequired     int
	RepRequired       float64
	SystemName        string
	Faction           string
	BasePos           conftypes.Vector
}

type Commodity struct {
	Nickname              string
	Name                  string
	Combinable            bool
	Volume                float64
	NameID                int
	InfocardID            int
	Infocard              InfocardKey
	Bases                 []*GoodAtBase
	PriceBestBaseBuysFor  int
	PriceBestBaseSellsFor int
	ProffitMargin         int
	Trades
}

func GetPricePerVoume(price int, volume float64) float64 {
	if volume == 0 {
		return -1
	}
	return float64(price) / float64(volume)
}

func (e *Exporter) GetCommodities() []*Commodity {
	commodities := make([]*Commodity, 0, 100)

	for _, comm := range e.configs.Goods.Commodities {
		commodity := &Commodity{}
		commodity.Nickname = comm.Nickname.Get()
		commodity.Combinable = comm.Combinable.Get()

		equipment_name := comm.Equipment.Get()
		equipment := e.configs.Equip.CommoditiesMap[equipment_name]

		commodity.NameID = equipment.IdsName.Get()

		commodity.Name = e.GetInfocardName(equipment.IdsName.Get(), commodity.Nickname)
		commodity.Infocard = InfocardKey(commodity.Nickname)
		e.exportInfocards(commodity.Infocard, equipment.IdsInfo.Get())
		commodity.InfocardID = equipment.IdsInfo.Get()

		volume := equipment.Volume.Get()
		commodity.Volume = volume
		base_item_price := comm.Price.Get()

		commodity.Bases = e.GetAtBasesSold(GetAtBasesInput{
			Nickname: commodity.Nickname,
			Price:    base_item_price,
			Volume:   commodity.Volume,
		})

		for _, base_info := range commodity.Bases {
			if base_info.PriceBaseBuysFor > commodity.PriceBestBaseBuysFor {
				commodity.PriceBestBaseBuysFor = base_info.PriceBaseBuysFor
			}
			if base_info.PriceBaseSellsFor < commodity.PriceBestBaseSellsFor || commodity.PriceBestBaseSellsFor == 0 {
				if base_info.BaseSells {
					commodity.PriceBestBaseSellsFor = base_info.PriceBaseSellsFor
				}

			}
		}

		if commodity.PriceBestBaseBuysFor > 0 && commodity.PriceBestBaseSellsFor > 0 {
			commodity.ProffitMargin = commodity.PriceBestBaseBuysFor - commodity.PriceBestBaseSellsFor
		}

		commodities = append(commodities, commodity)
	}

	return commodities
}

type GetAtBasesInput struct {
	Nickname string
	Price    int
	Volume   float64
}

func (e *Exporter) GetAtBasesSold(commodity GetAtBasesInput) []*GoodAtBase {
	var bases_list []*GoodAtBase
	var bases_already_found map[string]bool = make(map[string]bool)

	if e.configs.Discovery != nil {
		for _, base_market := range e.configs.Discovery.Prices.BasesPerGood[commodity.Nickname] {
			var base_info *GoodAtBase
			base_nickname := base_market.BaseNickname.Get()

			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Recovered in f", r)
					fmt.Println("recovered base_nickname", base_nickname)
					fmt.Println("recovered commodity nickname", commodity.Nickname)
					panic(r)
				}
			}()

			base_info = &GoodAtBase{
				BaseNickname:      base_nickname,
				BaseSells:         !base_market.SellOnly.Get(),
				PriceBaseBuysFor:  base_market.PriceBaseBuysFor.Get(),
				PriceBaseSellsFor: base_market.PriceBaseSellsFor.Get(),
			}

			more_info := e.GetBaseInfo(universe_mapped.BaseNickname(base_info.BaseNickname))
			base_info.BaseName = more_info.BaseName
			base_info.SystemName = more_info.SystemName
			base_info.Faction = more_info.Faction
			base_info.Volume = commodity.Volume
			base_info.BasePos = more_info.Pos

			if e.useful_bases_by_nick != nil {
				if _, ok := e.useful_bases_by_nick[base_info.BaseNickname]; !ok {
					continue
				}
			}

			bases_list = append(bases_list, base_info)
			bases_already_found[base_info.BaseNickname] = true
		}
	}

	for _, base_market := range e.configs.Market.BasesPerGood[commodity.Nickname] {
		base_nickname := base_market.Base

		// skip read from disco already
		if e.configs.Discovery != nil {
			if _, already_found := bases_already_found[base_nickname]; already_found {
				continue
			}
		}

		market_good := base_market.MarketGood
		base_info := &GoodAtBase{}
		base_info.Volume = commodity.Volume
		base_info.BaseSells = market_good.BaseSells()

		base_info.BaseNickname = base_nickname

		base_info.PriceBaseSellsFor = int(market_good.PriceModifier.Get() * float64(commodity.Price))

		if e.configs.Discovery != nil {
			base_info.PriceBaseBuysFor = market_good.BaseSellsIPositiveAndDiscoSellPrice.Get()
		} else {
			base_info.PriceBaseBuysFor = base_info.PriceBaseSellsFor
		}

		base_info.LevelRequired = market_good.LevelRequired.Get()
		base_info.RepRequired = market_good.RepRequired.Get()

		more_info := e.GetBaseInfo(universe_mapped.BaseNickname(base_info.BaseNickname))
		base_info.BaseName = more_info.BaseName
		base_info.SystemName = more_info.SystemName
		base_info.Faction = more_info.Faction
		base_info.BasePos = more_info.Pos

		if e.useful_bases_by_nick != nil {
			if _, ok := e.useful_bases_by_nick[base_info.BaseNickname]; !ok {
				continue
			}
		}

		bases_list = append(bases_list, base_info)
	}
	return bases_list
}

type BaseInfo struct {
	BaseName   string
	SystemName string
	Faction    string
	Pos        conftypes.Vector
}

func (e *Exporter) GetBaseInfo(base_nickname universe_mapped.BaseNickname) BaseInfo {
	var result BaseInfo
	universe_base, found_universe_base := e.configs.Universe_config.BasesMap[universe_mapped.BaseNickname(base_nickname)]

	if !found_universe_base {
		return result
	}

	result.BaseName = e.GetInfocardName(universe_base.StridName.Get(), string(base_nickname))
	system_nickname := universe_base.System.Get()

	if system, ok := e.configs.Universe_config.SystemMap[universe_mapped.SystemNickname(system_nickname)]; ok {
		result.SystemName = e.GetInfocardName(system.Strid_name.Get(), system_nickname)
	}

	var reputation_nickname string
	if system, ok := e.configs.Systems.SystemsMap[universe_base.System.Get()]; ok {
		for _, system_base := range system.Bases {
			if system_base.IdsName.Get() != universe_base.StridName.Get() {
				continue
			}

			reputation_nickname = system_base.RepNickname.Get()
			result.Pos = system_base.Pos.Get()
		}

	}

	var factionName string
	if group, exists := e.configs.InitialWorld.GroupsMap[reputation_nickname]; exists {
		factionName = e.GetInfocardName(group.IdsName.Get(), reputation_nickname)
	}

	result.Faction = factionName

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
