package configs_export

import (
	"math"

	"github.com/darklab8/fl-configs/configs/cfgtype"
	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/initialworld/flhash"
	"github.com/darklab8/go-utils/utils/ptr"
)

type MarketGood struct {
	Name         string
	Nickname     string
	NicknameHash flhash.HashCode
	HpType       string
	Type         string

	LevelRequired int
	RepRequired   float64
	Infocard      InfocardKey

	BaseSells     bool
	PriceModifier float64
	PriceBase     int
	PriceToBuy    int
	PriceToSell   *int
	Volume        float64
	ShipClass     cfgtype.ShipClass

	IsServerSideOverride bool
}

func NameWithSpacesOnly(word string) bool {
	for _, ch := range word {
		if ch != ' ' {
			return false
		}
	}
	return true
}

func (e *Exporter) getMarketGoods() map[cfgtype.BaseUniNick]map[CommodityKey]MarketGood {

	var goods_per_base map[cfgtype.BaseUniNick]map[CommodityKey]MarketGood = make(map[cfgtype.BaseUniNick]map[CommodityKey]MarketGood)

	for _, base_good := range e.Configs.Market.BaseGoods {
		base_nickname := cfgtype.BaseUniNick(base_good.Base.Get())

		var MarketGoods map[CommodityKey]MarketGood
		if market_goods, ok := goods_per_base[base_nickname]; ok {
			MarketGoods = market_goods
		} else {
			MarketGoods = make(map[CommodityKey]MarketGood)
		}
		for _, market_good := range base_good.MarketGoods {

			var market_good_nickname string = market_good.Nickname.Get()
			var price_base int
			var Name string
			var category string
			var hptype string
			if good, found_good := e.Configs.Goods.GoodsMap[market_good_nickname]; found_good {
				price_base = good.Price.Get()

				category = good.Category.Get()
				switch category {
				default:
					if equip, ok := e.Configs.Equip.ItemsMap[market_good_nickname]; ok {
						category = equip.Category
						Name = e.GetInfocardName(equip.IdsName.Get(), market_good_nickname)

						e.exportInfocards(InfocardKey(market_good_nickname), equip.IdsInfo.Get())
					}
				case "ship":
					ship := e.Configs.Goods.ShipsMap[good.Nickname.Get()]

					ship_hull := e.Configs.Goods.ShipHullsMap[ship.Hull.Get()]
					price_base = ship_hull.Price.Get()

					// Infocard data
					ship_nickname := ship_hull.Ship.Get()
					shiparch := e.Configs.Shiparch.ShipsMap[ship_nickname]

					Name = e.GetInfocardName(shiparch.IdsName.Get(), ship_nickname)

					// e.exportInfocards(InfocardKey(market_good_nickname),
					// 	shiparch.IdsInfo.Get(), shiparch.IdsInfo1.Get(), shiparch.IdsInfo2.Get(), shiparch.IdsInfo3.Get())
					e.exportInfocards(InfocardKey(market_good_nickname),
						shiparch.IdsInfo1.Get(), shiparch.IdsInfo.Get())
				}

				if gun, ok := e.Configs.Equip.GunMap[market_good_nickname]; ok {
					hptype, _ = gun.HPGunType.GetValue()
				}
				if shield, ok := e.Configs.Equip.ShidGenMap[market_good_nickname]; ok {
					hptype, _ = shield.HpType.GetValue()
				}
				if engine, ok := e.Configs.Equip.EnginesMap[market_good_nickname]; ok {
					hptype, _ = engine.HpType.GetValue()
				}

			}

			if NameWithSpacesOnly(Name) {
				Name = ""
			}

			good_to_add := MarketGood{
				Name:          Name,
				Nickname:      market_good_nickname,
				NicknameHash:  flhash.HashNickname(market_good_nickname),
				HpType:        hptype,
				Type:          category,
				LevelRequired: market_good.LevelRequired.Get(),
				RepRequired:   market_good.RepRequired.Get(),
				BaseSells:     market_good.BaseSells(),
				PriceModifier: market_good.PriceModifier.Get(),
				PriceBase:     price_base,
				PriceToBuy:    int(math.Floor(float64(price_base) * market_good.PriceModifier.Get())),
				Infocard:      InfocardKey(market_good_nickname),
				ShipClass:     -1,
			}

			e.Hashes[market_good_nickname] = good_to_add.NicknameHash

			if category == "commodity" {

				if e.Configs.Discovery != nil {
					good_to_add.PriceToSell = ptr.Ptr(market_good.BaseSellsIPositiveAndDiscoSellPrice.Get())
				} else {
					good_to_add.PriceToSell = &good_to_add.PriceToBuy
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

		goods_per_base[base_nickname] = MarketGoods
	}
	return goods_per_base
}
