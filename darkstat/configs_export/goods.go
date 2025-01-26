package configs_export

import (
	"math"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/initialworld/flhash"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped"
	"github.com/darklab8/go-utils/utils/ptr"
)

func NameWithSpacesOnly(word string) bool {
	for _, ch := range word {
		if ch != ' ' {
			return false
		}
	}
	return true
}

type GoodInfo struct {
	Nickname     string          `json:"nickname"`
	ShipNickname string          `json:"ship_nickname"` // market good can be ship package, if it is, then ship nickname bought by package is specified
	NicknameHash flhash.HashCode `json:"nickname_hash" format:"int64"`
	Name         string          `json:"name"`
	PriceBase    int             `json:"price_base"`
	HpType       string          `json:"hp_type"`
	Category     string          `json:"category"`
}

func (e *Exporter) GetGoodInfo(good_nickname string) GoodInfo {
	var info GoodInfo = GoodInfo{
		Nickname:     good_nickname,
		NicknameHash: flhash.HashNickname(good_nickname),
	}
	if good, found_good := e.Configs.Goods.GoodsMap[good_nickname]; found_good {
		info.PriceBase = good.Price.Get()

		info.Category = good.Category.Get()
		switch info.Category {
		default:
			if equip, ok := e.Configs.Equip.ItemsMap[good_nickname]; ok {
				info.Category = equip.Category
				info.Name = e.GetInfocardName(equip.IdsName.Get(), good_nickname)

				e.exportInfocards(InfocardKey(good_nickname), equip.IdsInfo.Get())
			}
		case "ship":
			ship := e.Configs.Goods.ShipsMap[good.Nickname.Get()]

			ship_hull := e.Configs.Goods.ShipHullsMap[ship.Hull.Get()]
			info.PriceBase = ship_hull.Price.Get()

			// Infocard data
			info.ShipNickname = ship_hull.Ship.Get()
			shiparch := e.Configs.Shiparch.ShipsMap[info.ShipNickname]

			info.Name = e.GetInfocardName(shiparch.IdsName.Get(), info.ShipNickname)

			// e.exportInfocards(InfocardKey(market_good_nickname),
			// 	shiparch.IdsInfo.Get(), shiparch.IdsInfo1.Get(), shiparch.IdsInfo2.Get(), shiparch.IdsInfo3.Get())
			e.exportInfocards(InfocardKey(good_nickname),
				shiparch.IdsInfo1.Get(), shiparch.IdsInfo.Get())
		}

		if gun, ok := e.Configs.Equip.GunMap[good_nickname]; ok {
			info.HpType, _ = gun.HPGunType.GetValue()
		}
		if shield, ok := e.Configs.Equip.ShidGenMap[good_nickname]; ok {
			info.HpType, _ = shield.HpType.GetValue()
		}
		if engine, ok := e.Configs.Equip.EnginesMap[good_nickname]; ok {
			info.HpType, _ = engine.HpType.GetValue()
		}
	}
	if NameWithSpacesOnly(info.Name) {
		info.Name = ""
	}

	return info
}

func (e *Exporter) getMarketGoods() map[cfg.BaseUniNick]map[CommodityKey]*MarketGood {

	var goods_per_base map[cfg.BaseUniNick]map[CommodityKey]*MarketGood = make(map[cfg.BaseUniNick]map[CommodityKey]*MarketGood)

	for _, base_good := range e.Configs.Market.BaseGoods {
		base_nickname := cfg.BaseUniNick(base_good.Base.Get())

		var MarketGoods map[CommodityKey]*MarketGood
		if market_goods, ok := goods_per_base[base_nickname]; ok {
			MarketGoods = market_goods
		} else {
			MarketGoods = make(map[CommodityKey]*MarketGood)
		}
		for _, market_good := range base_good.MarketGoods {

			var market_good_nickname string = market_good.Nickname.Get()

			good_to_add := &MarketGood{
				GoodInfo:      e.GetGoodInfo(market_good_nickname),
				BaseInfo:      e.GetBaseInfo(universe_mapped.BaseNickname(base_nickname)),
				LevelRequired: market_good.LevelRequired.Get(),
				RepRequired:   market_good.RepRequired.Get(),
				BaseSells:     market_good.BaseSells(),
				ShipClass:     -1,
			}
			good_to_add.PriceBaseSellsFor = int(math.Floor(float64(good_to_add.PriceBase) * market_good.PriceModifier.Get()))

			e.Hashes[market_good_nickname] = good_to_add.NicknameHash

			if good_to_add.Category == "commodity" {

				if e.Configs.Discovery != nil {
					good_to_add.PriceBaseBuysFor = ptr.Ptr(market_good.BaseSellsIPositiveAndDiscoSellPrice.Get())
				} else {
					good_to_add.PriceBaseBuysFor = ptr.Ptr(good_to_add.PriceBaseSellsFor)
				}
				equipment := e.Configs.Equip.CommoditiesMap[market_good_nickname]

				for _, volume := range equipment.Volumes {
					good_to_add2 := good_to_add
					good_to_add2.Volume = volume.Volume.Get()
					good_to_add2.ShipClass = volume.GetShipClass()
					MarketGoods[GetCommodityKey(good_to_add2.Nickname, good_to_add2.ShipClass)] = good_to_add2
				}

			} else {
				MarketGoods[GetCommodityKey(market_good_nickname, good_to_add.ShipClass)] = good_to_add
			}
		}

		if !e.TraderExists(string(base_nickname)) {
			for good_key, good := range MarketGoods {
				if good.Category == "commodity" {
					delete(MarketGoods, good_key)
				}
			}
		}

		goods_per_base[base_nickname] = MarketGoods
	}
	return goods_per_base
}
