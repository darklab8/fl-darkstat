package configs_export

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/inireader"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export/infocarder"
	"github.com/darklab8/go-utils/utils/ptr"
)

type Ammo struct {
	Name  string `json:"name" validate:"required"`
	Price int    `json:"price" validate:"required"`

	HitPts           int     `json:"hit_pts" validate:"required"`
	Volume           float64 `json:"volume" validate:"required"`
	MunitionLifetime float64 `json:"munition_lifetime" validate:"required"`

	Nickname     string  `json:"nickname" validate:"required"`
	NameID       int     `json:"name_id" validate:"required"`
	InfoID       int     `json:"info_id" validate:"required"`
	SeekerType   *string `json:"seeker_type" validate:"required"`
	SeekerRange  *int    `json:"seeker_range" validate:"required"`
	SeekerFovDeg *int    `json:"seeker_fov_deg" validate:"required"`

	Bases map[cfg.BaseUniNick]*MarketGood `json:"-" swaggerignore:"true"`

	*DiscoveryTechCompat `json:"-" swaggerignore:"true"`

	AmmoLimit AmmoLimit `json:"ammo_limit" validate:"required"`
	Mass      float64   `json:"mass" validate:"required"`
}

func (b Ammo) GetNickname() string { return string(b.Nickname) }

func (b Ammo) GetBases() map[cfg.BaseUniNick]*MarketGood { return b.Bases }

func (b Ammo) GetDiscoveryTechCompat() *DiscoveryTechCompat { return b.DiscoveryTechCompat }

// GetModelWithoutLastComments
// this function eliminates bug https://github.com/darklab8/fl-darkstat/issues/94
// when people write comments after the model ini params already for next ini sector
// we should not be printing such comments
func GetModelWithoutLastComments(item_model *semantic.Model) []*inireader.Param {
	sector := item_model.RenderModel()

	// strip from comments at the end
	comment_lines_at_the_end := 0
	for i := len(sector.Params) - 1; i > 0; i-- {

		if sector.Params[i].Key == inireader.KEY_COMMENT {
			comment_lines_at_the_end++
		} else {
			break
		}
	}

	return sector.Params[:len(sector.Params)-comment_lines_at_the_end]
}

func (e *Exporter) WriteConfigToInfocard(item_model *semantic.Model, item_nickname string) {
	// add to item name its ini config
	var infocard_addition infocarder.InfocardBuilder

	infocard_addition.WriteLineStr(string(item_model.GetOriginalType()))

	for _, param := range GetModelWithoutLastComments(item_model) {
		infocard_addition.WriteLineStr(string(param.ToString(inireader.WithComments(false))))
	}
	infocard_addition.WriteLineStr("")
	var info infocarder.InfocardBuilder
	if value, ok := e.GetInfocard2(infocarder.InfocardKey(item_nickname)); ok {
		info.Lines = value
	}
	e.PutInfocard(infocarder.InfocardKey(item_nickname), append(info.Lines, infocard_addition.Lines...))
}

func GetValuePtr[K any](some_func func() (K, bool)) *K {
	value, _ := some_func()
	return ptr.Ptr(value)
}

func (e *Exporter) GetAmmo(ids []*Tractor) []Ammo {
	var ammos []Ammo

	for _, item_info := range e.Mapped.Equip().Munitions {
		item := Ammo{
			Bases: make(map[cfg.BaseUniNick]*MarketGood),
		}
		item.Mass, _ = item_info.Mass.GetValue()

		item.Nickname = item_info.Nickname.Get()
		item.NameID, _ = item_info.IdsName.GetValue()
		item.InfoID, _ = item_info.IdsInfo.GetValue()

		item.HitPts, _ = item_info.HitPts.GetValue()

		if value, ok := item_info.AmmoLimitAmountInCatridge.GetValue(); ok {
			item.AmmoLimit.AmountInCatridge = ptr.Ptr(value)
		}
		if value, ok := item_info.AmmoLimitMaxCatridges.GetValue(); ok {
			item.AmmoLimit.MaxCatridges = ptr.Ptr(value)
		}

		item.Volume, _ = item_info.Volume.GetValue()
		item.SeekerRange = GetValuePtr(item_info.SeekerRange.GetValue)
		item.SeekerType = GetValuePtr(item_info.SeekerType.GetValue)

		item.MunitionLifetime, _ = item_info.LifeTime.GetValue()

		item.SeekerFovDeg = GetValuePtr(item_info.SeekerFovDeg.GetValue)

		if ammo_ids_name, ok := item_info.IdsName.GetValue(); ok {
			item.Name = e.GetInfocardName(ammo_ids_name, item.Nickname)
		}

		item.Price = -1
		if good_info, ok := e.Mapped.Goods.GoodsMap[item_info.Nickname.Get()]; ok {
			if price, ok := good_info.Price.GetValue(); ok {
				item.Price = price
				item.Bases = e.GetAtBasesSold(GetCommodityAtBasesInput{
					Nickname: good_info.Nickname.Get(),
					Price:    price,
				})
			}
		}

		e.ExportInfocards(infocarder.InfocardKey(item.Nickname), item.InfoID)
		item.DiscoveryTechCompat = CalculateTechCompat(e.Mapped.Discovery, ids, item.Nickname)
		ammos = append(ammos, item)

		e.WriteConfigToInfocard(&item_info.Model, item.Nickname)

	}

	for _, item_info := range e.Mapped.Equip().Mines {
		item := Ammo{
			Bases: make(map[cfg.BaseUniNick]*MarketGood),
		}
		item.Mass, _ = item_info.Mass.GetValue()

		item.Nickname = item_info.Nickname.Get()
		item.NameID, _ = item_info.IdsName.GetValue()
		item.InfoID, _ = item_info.IdsInfo.GetValue()

		item.HitPts, _ = item_info.HitPts.GetValue()

		if value, ok := item_info.AmmoLimitAmountInCatridge.GetValue(); ok {
			item.AmmoLimit.AmountInCatridge = ptr.Ptr(value)
		}
		if value, ok := item_info.AmmoLimitMaxCatridges.GetValue(); ok {
			item.AmmoLimit.MaxCatridges = ptr.Ptr(value)
		}

		item.Volume, _ = item_info.Volume.GetValue()
		item.MunitionLifetime, _ = item_info.LifeTime.GetValue()
		item.SeekerRange = GetValuePtr(item_info.SeekDist.GetValue)

		if ammo_ids_name, ok := item_info.IdsName.GetValue(); ok {
			item.Name = e.GetInfocardName(ammo_ids_name, item.Nickname)
		}

		item.Price = -1
		if good_info, ok := e.Mapped.Goods.GoodsMap[item_info.Nickname.Get()]; ok {
			if price, ok := good_info.Price.GetValue(); ok {
				item.Price = price
				item.Bases = e.GetAtBasesSold(GetCommodityAtBasesInput{
					Nickname: good_info.Nickname.Get(),
					Price:    price,
				})
			}
		}

		e.ExportInfocards(infocarder.InfocardKey(item.Nickname), item.InfoID)
		item.DiscoveryTechCompat = CalculateTechCompat(e.Mapped.Discovery, ids, item.Nickname)
		ammos = append(ammos, item)

		e.WriteConfigToInfocard(&item_info.Model, item.Nickname)

	}

	for _, item_info := range e.Mapped.Equip().CounterMeasure {
		item := Ammo{
			Bases: make(map[cfg.BaseUniNick]*MarketGood),
		}
		item.Mass, _ = item_info.Mass.GetValue()

		item.Nickname = item_info.Nickname.Get()
		item.NameID, _ = item_info.IdsName.GetValue()
		item.InfoID, _ = item_info.IdsInfo.GetValue()

		if value, ok := item_info.AmmoLimitAmountInCatridge.GetValue(); ok {
			item.AmmoLimit.AmountInCatridge = ptr.Ptr(value)
		}
		if value, ok := item_info.AmmoLimitMaxCatridges.GetValue(); ok {
			item.AmmoLimit.MaxCatridges = ptr.Ptr(value)
		}

		item.HitPts, _ = item_info.HitPts.GetValue()
		item.Volume, _ = item_info.Volume.GetValue()
		item.MunitionLifetime, _ = item_info.LifeTime.GetValue()

		if ammo_ids_name, ok := item_info.IdsName.GetValue(); ok {
			item.Name = e.GetInfocardName(ammo_ids_name, item.Nickname)
		}

		item.Price = -1
		if good_info, ok := e.Mapped.Goods.GoodsMap[item_info.Nickname.Get()]; ok {
			if price, ok := good_info.Price.GetValue(); ok {
				item.Price = price
				item.Bases = e.GetAtBasesSold(GetCommodityAtBasesInput{
					Nickname: good_info.Nickname.Get(),
					Price:    price,
				})
			}
		}

		e.ExportInfocards(infocarder.InfocardKey(item.Nickname), item.InfoID)
		item.DiscoveryTechCompat = CalculateTechCompat(e.Mapped.Discovery, ids, item.Nickname)
		ammos = append(ammos, item)

		e.WriteConfigToInfocard(&item_info.Model, item.Nickname)

	}
	return ammos
}

func (e *Exporter) FilterToUsefulAmmo(cms []Ammo) []Ammo {
	var useful_items []Ammo = make([]Ammo, 0, len(cms))
	for _, item := range cms {
		if !e.Buyable(item.Bases) {
			continue
		}
		useful_items = append(useful_items, item)
	}
	return useful_items
}
