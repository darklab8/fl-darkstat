package frmt

import (
	"fmt"
	"sort"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
)

func SortedMarketGoods(goods_per_nick map[configs_export.CommodityKey]*configs_export.MarketGood) []*configs_export.MarketGood {
	var market_goods []*configs_export.MarketGood = make([]*configs_export.MarketGood, 0, 10)

	for _, good := range goods_per_nick {
		market_goods = append(market_goods, good)
	}

	sort.Slice(market_goods, func(i, j int) bool {
		if market_goods[i].Category != market_goods[j].Category {
			return market_goods[i].Category < market_goods[j].Category
		}
		return market_goods[i].Name < market_goods[j].Name
	})

	return market_goods
}

func FormattedShipClassOfCommodity(ship_class cfg.ShipClass) string {
	if ship_class >= 0 {
		return " (" + ship_class.ToStr() + ")"
	} else {
		return ""
	}
}

func FormattedShipClassOfCommodity2(ship_class cfg.ShipClass) string {
	if ship_class >= 0 {
		return fmt.Sprintf("%d,%s", ship_class, ship_class.ToStr())
	} else {
		return "nil"
	}
}
