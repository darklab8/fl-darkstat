package configs_export

import (
	"fmt"

	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped"
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
}

type Commodity struct {
	Nickname              string
	Name                  string
	Combinable            bool
	Volume                float64
	NameID                int
	InfocardID            int
	Infocard              InfocardKey
	Bases                 []GoodAtBase
	PriceBestBaseBuysFor  int
	PriceBestBaseSellsFor int
	ProffitMargin         int
}

func GetPricePerVoume(price int, volume float64) float64 {
	if volume == 0 {
		return -1
	}
	return float64(price) / float64(volume)
}

func (e *Exporter) GetCommodities() []Commodity {
	commodities := make([]Commodity, 0, 100)

	for _, comm := range e.configs.Goods.Commodities {

		var name string
		commodity := Commodity{}
		commodity.Nickname = comm.Nickname.Get()
		commodity.Combinable = comm.Combinable.Get()

		equipment_name := comm.Equipment.Get()
		equipment := e.configs.Equip.CommoditiesMap[equipment_name]

		commodity.NameID = equipment.IdsName.Get()
		if infoname, ok := e.configs.Infocards.Infonames[equipment.IdsName.Get()]; ok {
			name = string(infoname)
		}
		commodity.Name = name
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

func (e *Exporter) GetAtBasesSold(commodity GetAtBasesInput) []GoodAtBase {
	var bases_list []GoodAtBase
	var bases_already_found map[string]bool = make(map[string]bool)

	if e.configs.Discovery != nil {
		for _, base_market := range e.configs.Discovery.Prices.BasesPerGood[commodity.Nickname] {
			var base_info GoodAtBase
			base_nickname := base_market.BaseNickname.Get()

			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Recovered in f", r)
					fmt.Println("recovered base_nickname", base_nickname)
					fmt.Println("recovered commodity nickname", commodity.Nickname)
					panic(r)
				}
			}()

			base_info = GoodAtBase{
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
		base_info := GoodAtBase{}
		base_info.Volume = commodity.Volume
		base_info.BaseSells = market_good.BaseSellsIfAboveZero.Get() > 0

		base_info.BaseNickname = base_nickname

		base_info.PriceBaseSellsFor = int(market_good.PriceModifier.Get() * float64(commodity.Price))

		if e.configs.Discovery != nil {
			base_info.PriceBaseBuysFor = market_good.SellPrice.Get()
		} else {
			base_info.PriceBaseBuysFor = base_info.PriceBaseSellsFor
		}

		base_info.LevelRequired = market_good.LevelRequired.Get()
		base_info.RepRequired = market_good.RepRequired.Get()

		more_info := e.GetBaseInfo(universe_mapped.BaseNickname(base_info.BaseNickname))
		base_info.BaseName = more_info.BaseName
		base_info.SystemName = more_info.SystemName
		base_info.Faction = more_info.Faction

		bases_list = append(bases_list, base_info)
	}
	return bases_list
}

type BaseInfo struct {
	BaseName   string
	SystemName string
	Faction    string
}

func (e *Exporter) GetBaseInfo(base_nickname universe_mapped.BaseNickname) BaseInfo {
	var result BaseInfo
	if universe_base, ok := e.configs.Universe_config.BasesMap[universe_mapped.BaseNickname(base_nickname)]; ok {

		if infoname, ok := e.configs.Infocards.Infonames[universe_base.StridName.Get()]; ok {
			result.BaseName = string(infoname)
		}
		system_nickname := universe_base.System.Get()

		if system, ok := e.configs.Universe_config.SystemMap[universe_mapped.SystemNickname(system_nickname)]; ok {
			if infoname, ok := e.configs.Infocards.Infonames[system.Strid_name.Get()]; ok {
				result.SystemName = string(infoname)
			}
		}

		var reputation_nickname string
		if system, ok := e.configs.Systems.SystemsMap[universe_base.System.Get()]; ok {
			for _, system_base := range system.Bases {
				if system_base.IdsName.Get() == universe_base.StridName.Get() {
					reputation_nickname = system_base.RepNickname.Get()
				}
			}
		}

		var factionName string
		if group, exists := e.configs.InitialWorld.GroupsMap[reputation_nickname]; exists {
			if faction_name, exists := e.configs.Infocards.Infonames[group.IdsName.Get()]; exists {
				factionName = string(faction_name)
			}
		}

		result.Faction = factionName
	}
	return result
}

func FilterToUsefulCommodities(commodities []Commodity) []Commodity {
	var items []Commodity = make([]Commodity, 0, len(commodities))
	for _, item := range commodities {
		if !Buyable(item.Bases) {
			continue
		}
		items = append(items, item)
	}
	return items
}
