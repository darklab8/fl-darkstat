package configs_export

import (
	"math"
)

type MarketGood struct {
	Name     string
	Nickname string
	HpType   string
	Type     string

	LevelRequired int
	RepRequired   float64
	Infocard      InfocardKey

	IsBuyOnly     bool
	PriceModifier float64
	PriceBase     int
	Price         int
}

func NameWithSpacesOnly(word string) bool {
	for _, ch := range word {
		if ch != ' ' {
			return false
		}
	}
	return true
}

func (e *Exporter) getMarketGoods() map[string][]MarketGood {
	var goods_per_base map[string][]MarketGood = make(map[string][]MarketGood)

	for _, base_good := range e.configs.Market.BaseGoods {
		base_nickname := base_good.Base.Get()

		var MarketGoods []MarketGood = make([]MarketGood, 0, 200)
		for _, market_good := range base_good.MarketGoods {

			var market_good_nickname string = market_good.Nickname.Get()
			var price_base int
			var Name string
			var category string
			var hptype string
			if good, found_good := e.configs.Goods.GoodsMap[market_good_nickname]; found_good {
				price_base = good.Price.Get()

				category = good.Category.Get()
				switch category {
				default:
					if equip, ok := e.configs.Equip.ItemsMap[market_good_nickname]; ok {
						if infoname, ok := e.configs.Infocards.Infonames[equip.IdsName.Get()]; ok {
							Name = string(infoname)
							category = equip.Category
						}

						e.infocards_parser.Set(InfocardKey(market_good_nickname), equip.IdsInfo.Get())
					}
				case "ship":
					ship := e.configs.Goods.ShipsMap[good.Nickname.Get()]

					ship_hull := e.configs.Goods.ShipHullsMap[ship.Hull.Get()]
					price_base = ship_hull.Price.Get()

					// Infocard data
					ship_nickname := ship_hull.Ship.Get()
					shiparch := e.configs.Shiparch.ShipsMap[ship_nickname]

					if infoname, ok := e.configs.Infocards.Infonames[shiparch.IdsName.Get()]; ok {
						Name = string(infoname)
					}

					e.infocards_parser.Set(InfocardKey(market_good_nickname),
						shiparch.IdsInfo.Get(), shiparch.IdsInfo1.Get(), shiparch.IdsInfo2.Get(), shiparch.IdsInfo3.Get())
				}

				if gun, ok := e.configs.Equip.GunMap[market_good_nickname]; ok {
					hptype, _ = gun.HPGunType.GetValue()
				}
				if shield, ok := e.configs.Equip.ShidGenMap[market_good_nickname]; ok {
					hptype, _ = shield.HpType.GetValue()
				}
				if engine, ok := e.configs.Equip.EnginesMap[market_good_nickname]; ok {
					hptype, _ = engine.HpType.GetValue()
				}

			}

			if NameWithSpacesOnly(Name) {
				Name = ""
			}

			MarketGoods = append(MarketGoods, MarketGood{
				Name:          Name,
				Nickname:      market_good_nickname,
				HpType:        hptype,
				Type:          category,
				LevelRequired: market_good.LevelRequired.Get(),
				RepRequired:   market_good.RepRequired.Get(),
				IsBuyOnly:     market_good.IsBuyOnly.Get(),
				PriceModifier: market_good.PriceModifier.Get(),
				PriceBase:     price_base,
				Price:         int(math.Floor(float64(price_base) * market_good.PriceModifier.Get())),
				Infocard:      InfocardKey(market_good_nickname),
			})
		}

		goods_per_base[base_nickname] = append(goods_per_base[base_nickname], MarketGoods...)

		e.infocards_parser.makeReady() // clearing calculating queue to be not calculating repetions
	}
	return goods_per_base
}
