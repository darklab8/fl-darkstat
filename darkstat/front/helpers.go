package front

import (
	"fmt"
	"sort"

	"github.com/darklab8/fl-configs/configs/cfgtype"
	"github.com/darklab8/fl-configs/configs/configs_export"
)

func GetAmmoLimitFormatted(ammo_limit configs_export.AmmoLimit) string {
	result := "" // fmt.Sprintf("%6.0f", gun.EnergyDamagePerSec)

	if ammo_limit.AmountInCatridge != nil {
		result += fmt.Sprintf("%6d", *ammo_limit.AmountInCatridge)
	}

	if ammo_limit.MaxCatridges != nil {
		result += fmt.Sprintf("(x%d)", *ammo_limit.MaxCatridges)
	}

	if result == "" {
		return "inf."
	}

	return result
}

func SortedBases(bases_map map[cfgtype.BaseUniNick]*configs_export.GoodAtBase) []*configs_export.GoodAtBase {
	var bases []*configs_export.GoodAtBase = make([]*configs_export.GoodAtBase, 0, 10)

	for _, base := range bases_map {
		bases = append(bases, base)
	}

	sort.Slice(bases, func(i, j int) bool {
		if bases[i].BaseName != "" && bases[j].BaseName == "" {
			return true
		}
		return bases[i].BaseName < bases[j].BaseName
	})

	return bases
}

func SortedMarketGoods(goods_per_nick map[configs_export.CommodityKey]configs_export.MarketGood) []configs_export.MarketGood {
	var market_goods []configs_export.MarketGood = make([]configs_export.MarketGood, 0, 10)

	for _, good := range goods_per_nick {
		market_goods = append(market_goods, good)
	}

	sort.Slice(market_goods, func(i, j int) bool {
		if market_goods[i].Name != "" && market_goods[j].Name == "" {
			return true
		}
		return market_goods[i].Name < market_goods[j].Name
	})

	return market_goods
}

func FormattedShipClassOfCommodity(ship_class cfgtype.ShipClass) string {
	if ship_class >= 0 {
		return " (" + ship_class.ToStr() + ")"
	} else {
		return ""
	}
}

func FormattedShipClassOfCommodity2(ship_class cfgtype.ShipClass) string {
	if ship_class >= 0 {
		return fmt.Sprintf("%d,%s", ship_class, ship_class.ToStr())
	} else {
		return "nil"
	}
}
