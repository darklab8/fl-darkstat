package configs_export

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/inireader"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export/infocarder"
	"github.com/darklab8/go-utils/utils/ptr"
)

type ExtraItem struct {
	Name  string `json:"name"  validate:"required"`
	Price int    `json:"price"  validate:"required"`

	Category string `json:"category"  validate:"required"`
	Lootable bool   `json:"lootable"  validate:"required"`
	Nickname string `json:"nickname"  validate:"required"`
	NameID   int    `json:"name_id"  validate:"required"`
	InfoID   int    `json:"indo_id"  validate:"required"`

	Mass   float64 `json:"mass"`
	HpType string  `json:"hp_type"`

	PowerCapacity   *int `json:"power_capacity"`
	PowerChargeRate *int `json:"power_charge_rate"`

	Bases                map[cfg.BaseUniNick]*MarketGood `json:"-" swaggerignore:"true"`
	*DiscoveryTechCompat `json:"-" swaggerignore:"true"`
}

func (b ExtraItem) GetNickname() string { return string(b.Nickname) }

func (b ExtraItem) GetBases() map[cfg.BaseUniNick]*MarketGood { return b.Bases }

func (b ExtraItem) GetDiscoveryTechCompat() *DiscoveryTechCompat { return b.DiscoveryTechCompat }

func (e *Exporter) GetExtraItems(ids []*Tractor) []ExtraItem {
	var tractors []ExtraItem

	for _, item_info := range e.Mapped.Equip().Extras {
		item := ExtraItem{
			Bases: make(map[cfg.BaseUniNick]*MarketGood),
		}
		item.Mass, _ = item_info.Mass.GetValue()

		item.Nickname = item_info.Nickname.Get()
		item.NameID = item_info.IdsName.Get()
		item.InfoID = item_info.IdsInfo.Get()
		item.Category = item_info.Category

		if item.Category == "attachedfx" ||
			item.Category == "cargopod" ||
			item.Category == "internalfx" ||
			item.Category == "light" ||
			item.Category == "lootcrate" ||
			item.Category == "motor" ||
			item.Category == "tradelane" {
			continue
		}

		item.Mass, _ = item_info.Mass.GetValue()
		item.HpType, _ = item_info.HpType.GetValue()

		if item.Category == "power" {
			generator := e.Mapped.Equip().PowersMap[item.Nickname]
			item.PowerCapacity = ptr.Ptr(generator.Capacity.Get())
			item.PowerChargeRate = ptr.Ptr(generator.ChargeRate.Get())
		}

		if good_info, ok := e.Mapped.Goods.GoodsMap[item.Nickname]; ok {
			if price, ok := good_info.Price.GetValue(); ok {
				item.Price = price
				item.Bases = e.GetAtBasesSold(GetCommodityAtBasesInput{
					Nickname: good_info.Nickname.Get(),
					Price:    price,
				})
			}
		}

		item.Name = e.GetInfocardName(item.NameID, item.Nickname)
		e.exportInfocards(infocarder.InfocardKey(item.Nickname), item.InfoID)

		// add to item name its ini config
		var infocard_addition infocarder.InfocardBuilder
		sector := item_info.Model.RenderModel()
		infocard_addition.WriteLineStr(string(sector.OriginalType))
		for _, param := range sector.Params {
			infocard_addition.WriteLineStr(string(param.ToString(inireader.WithComments(false))))
		}
		infocard_addition.WriteLineStr("")

		var info infocarder.InfocardBuilder
		if value, ok := e.GetInfocard2(infocarder.InfocardKey(item.Nickname)); ok {
			info.Lines = value
		}
		e.PutInfocard(infocarder.InfocardKey(item.Nickname), append(info.Lines, infocard_addition.Lines...))

		item.DiscoveryTechCompat = CalculateTechCompat(e.Mapped.Discovery, ids, item.Nickname)
		tractors = append(tractors, item)
	}
	return tractors
}

func (e *Exporter) FilterToUsefulItems(cms []ExtraItem) []ExtraItem {
	var useful_items []ExtraItem = make([]ExtraItem, 0, len(cms))
	for _, item := range cms {
		if !e.Buyable(item.Bases) {
			continue
		}
		useful_items = append(useful_items, item)
	}
	return useful_items
}
