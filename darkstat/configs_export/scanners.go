package configs_export

import (
	"github.com/darklab8/fl-configs/configs/cfgtype"
	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/initialworld/flhash"
)

type Scanner struct {
	Name  string
	Price int

	Range          int
	CargoScanRange int

	Lootable     bool
	Nickname     string
	NicknameHash flhash.HashCode
	NameID       int
	InfoID       int

	Bases map[cfgtype.BaseUniNick]*GoodAtBase

	*DiscoveryTechCompat
	Mass float64
}

func (e *Exporter) GetScanners(ids []Tractor) []Scanner {
	var scanners []Scanner

	for _, scanner_info := range e.Configs.Equip.Scanners {
		item := Scanner{
			Bases: make(map[cfgtype.BaseUniNick]*GoodAtBase),
		}
		item.Mass, _ = scanner_info.Mass.GetValue()

		item.Nickname = scanner_info.Nickname.Get()
		item.NicknameHash = flhash.HashNickname(item.Nickname)
		e.Hashes[item.Nickname] = item.NicknameHash

		item.Lootable = scanner_info.Lootable.Get()
		item.NameID = scanner_info.IdsName.Get()
		item.InfoID = scanner_info.IdsInfo.Get()
		item.Range = scanner_info.Range.Get()
		item.CargoScanRange = scanner_info.CargoScanRange.Get()

		if good_info, ok := e.Configs.Goods.GoodsMap[item.Nickname]; ok {
			if price, ok := good_info.Price.GetValue(); ok {
				item.Price = price
				item.Bases = e.GetAtBasesSold(GetCommodityAtBasesInput{
					Nickname: good_info.Nickname.Get(),
					Price:    price,
				})
			}
		}

		item.Name = e.GetInfocardName(item.NameID, item.Nickname)

		e.exportInfocards(InfocardKey(item.Nickname), item.InfoID)
		item.DiscoveryTechCompat = CalculateTechCompat(e.Configs.Discovery, ids, item.Nickname)
		scanners = append(scanners, item)
	}
	return scanners
}

func (e *Exporter) FilterToUserfulScanners(items []Scanner) []Scanner {
	var useful_items []Scanner = make([]Scanner, 0, len(items))
	for _, item := range items {
		if !e.Buyable(item.Bases) {
			continue
		}
		useful_items = append(useful_items, item)
	}
	return useful_items
}
