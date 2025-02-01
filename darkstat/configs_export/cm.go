package configs_export

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/initialworld/flhash"
	"github.com/darklab8/go-utils/utils/ptr"
)

type CounterMeasure struct {
	Name  string `json:"name"`
	Price int    `json:"price"`

	HitPts        int `json:"hit_pts"`
	AIRange       int `json:"ai_range"`
	Lifetime      int `json:"lifetime"`
	Range         int `json:"range"`
	DiversionPctg int `json:"diversion_pctg"`

	Lootable     bool            `json:"lootable"`
	Nickname     string          `json:"nickname"`
	NicknameHash flhash.HashCode `json:"nickname_hash" format:"int64"`
	NameID       int             `json:"name_id"`
	InfoID       int             `json:"indo_id"`

	Bases map[cfg.BaseUniNick]*MarketGood `json:"-" swaggerignore:"true"`

	*DiscoveryTechCompat `json:"-" swaggerignore:"true"`

	AmmoLimit AmmoLimit `json:"ammo_limit"`
	Mass      float64   `json:"mass"`
}

func (b CounterMeasure) GetNickname() string { return string(b.Nickname) }

func (b CounterMeasure) GetBases() map[cfg.BaseUniNick]*MarketGood { return b.Bases }

func (b CounterMeasure) GetDiscoveryTechCompat() *DiscoveryTechCompat { return b.DiscoveryTechCompat }

func (e *Exporter) GetCounterMeasures(ids []*Tractor) []CounterMeasure {
	var tractors []CounterMeasure

	for _, cm_info := range e.Configs.Equip.CounterMeasureDroppers {
		cm := CounterMeasure{
			Bases: make(map[cfg.BaseUniNick]*MarketGood),
		}
		cm.Mass, _ = cm_info.Mass.GetValue()

		cm.Nickname = cm_info.Nickname.Get()
		cm.NicknameHash = flhash.HashNickname(cm.Nickname)
		e.Hashes[cm.Nickname] = cm.NicknameHash
		cm.HitPts = cm_info.HitPts.Get()
		cm.AIRange = cm_info.AIRange.Get()
		cm.Lootable = cm_info.Lootable.Get()
		cm.NameID = cm_info.IdsName.Get()
		cm.InfoID = cm_info.IdsInfo.Get()

		if good_info, ok := e.Configs.Goods.GoodsMap[cm.Nickname]; ok {
			if price, ok := good_info.Price.GetValue(); ok {
				cm.Price = price
				cm.Bases = e.GetAtBasesSold(GetCommodityAtBasesInput{
					Nickname: good_info.Nickname.Get(),
					Price:    price,
				})
			}
		}

		cm.Name = e.GetInfocardName(cm.NameID, cm.Nickname)

		infocards := []int{cm.InfoID}
		if ammo_info, ok := e.Configs.Equip.CounterMeasureMap[cm_info.ProjectileArchetype.Get()]; ok {

			if value, ok := ammo_info.AmmoLimitAmountInCatridge.GetValue(); ok {
				cm.AmmoLimit.AmountInCatridge = ptr.Ptr(value)
			}
			if value, ok := ammo_info.AmmoLimitMaxCatridges.GetValue(); ok {
				cm.AmmoLimit.MaxCatridges = ptr.Ptr(value)
			}

			cm.Lifetime = ammo_info.Lifetime.Get()
			cm.Range = ammo_info.Range.Get()
			cm.DiversionPctg = ammo_info.DiversionPctg.Get()

			if id, ok := ammo_info.IdsInfo.GetValue(); ok {
				infocards = append(infocards, id)
			}
		}

		e.exportInfocards(InfocardKey(cm.Nickname), infocards...)
		cm.DiscoveryTechCompat = CalculateTechCompat(e.Configs.Discovery, ids, cm.Nickname)
		tractors = append(tractors, cm)
	}
	return tractors
}

func (e *Exporter) FilterToUsefulCounterMeasures(cms []CounterMeasure) []CounterMeasure {
	var useful_items []CounterMeasure = make([]CounterMeasure, 0, len(cms))
	for _, item := range cms {
		if !e.Buyable(item.Bases) {
			continue
		}
		useful_items = append(useful_items, item)
	}
	return useful_items
}
