package configs_export

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
)

type Scanner struct {
	Name  string `json:"name" validate:"required"`
	Price int    `json:"price" validate:"required"`

	Range          int `json:"range" validate:"required"`
	CargoScanRange int `json:"cargo_scan_range" validate:"required"`

	Lootable bool   `json:"lootable" validate:"required"`
	Nickname string `json:"nickname" validate:"required"`
	NameID   int    `json:"name_id" validate:"required"`
	InfoID   int    `json:"info_id" validate:"required"`

	Bases map[cfg.BaseUniNick]*MarketGood `json:"-" swaggerignore:"true"`

	*DiscoveryTechCompat `json:"-" swaggerignore:"true"`
	Mass                 float64 `json:"mass" validate:"required"`
}

func (b Scanner) GetNickname() string { return string(b.Nickname) }

func (b Scanner) GetBases() map[cfg.BaseUniNick]*MarketGood { return b.Bases }

func (b Scanner) GetDiscoveryTechCompat() *DiscoveryTechCompat { return b.DiscoveryTechCompat }

func (e *Exporter) GetScanners(ids []*Tractor) []Scanner {
	var scanners []Scanner

	for _, scanner_info := range e.mapped.Equip.Scanners {
		item := Scanner{
			Bases: make(map[cfg.BaseUniNick]*MarketGood),
		}
		item.Mass, _ = scanner_info.Mass.GetValue()

		item.Nickname = scanner_info.Nickname.Get()

		item.Lootable = scanner_info.Lootable.Get()
		item.NameID = scanner_info.IdsName.Get()
		item.InfoID = scanner_info.IdsInfo.Get()
		item.Range = scanner_info.Range.Get()
		item.CargoScanRange = scanner_info.CargoScanRange.Get()

		if good_info, ok := e.mapped.Goods.GoodsMap[item.Nickname]; ok {
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
		item.DiscoveryTechCompat = CalculateTechCompat(e.mapped.Discovery, ids, item.Nickname)
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
