package configs_export

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/initialworld/flhash"
	"github.com/darklab8/go-utils/utils/ptr"
)

type Ammo struct {
	Name  string `json:"name"`
	Price int    `json:"price"`

	HitPts           int     `json:"hit_pts"`
	Volume           float64 `json:"volume"`
	MunitionLifetime float64 `json:"munition_lifetime"`

	Nickname     string          `json:"nickname"`
	NicknameHash flhash.HashCode `json:"nickname_hash" format:"int64"`
	NameID       int             `json:"name_id"`
	InfoID       int             `json:"info_id"`
	SeekerType   string          `json:"seeker_type"`
	SeekerRange  int             `json:"seeker_range"`
	SeekerFovDeg int             `json:"seeker_fov_deg"`

	Bases map[cfg.BaseUniNick]*MarketGood `json:"-" swaggerignore:"true"`

	*DiscoveryTechCompat `json:"-" swaggerignore:"true"`

	AmmoLimit AmmoLimit `json:"ammo_limit"`
	Mass      float64   `json:"mass"`
}

func (b Ammo) GetNickname() string { return string(b.Nickname) }

func (b Ammo) GetBases() map[cfg.BaseUniNick]*MarketGood { return b.Bases }

func (b Ammo) GetDiscoveryTechCompat() *DiscoveryTechCompat { return b.DiscoveryTechCompat }

func (e *Exporter) GetAmmo(ids []*Tractor) []Ammo {
	var tractors []Ammo

	for _, munition_info := range e.Configs.Equip.Munitions {
		munition := Ammo{
			Bases: make(map[cfg.BaseUniNick]*MarketGood),
		}
		munition.Mass, _ = munition_info.Mass.GetValue()

		munition.Nickname = munition_info.Nickname.Get()
		munition.NicknameHash = flhash.HashNickname(munition.Nickname)
		e.Hashes[munition.Nickname] = munition.NicknameHash
		munition.NameID, _ = munition_info.IdsName.GetValue()
		munition.InfoID, _ = munition_info.IdsInfo.GetValue()

		munition.HitPts, _ = munition_info.HitPts.GetValue()

		if value, ok := munition_info.AmmoLimitAmountInCatridge.GetValue(); ok {
			munition.AmmoLimit.AmountInCatridge = ptr.Ptr(value)
		}
		if value, ok := munition_info.AmmoLimitMaxCatridges.GetValue(); ok {
			munition.AmmoLimit.MaxCatridges = ptr.Ptr(value)
		}

		munition.Volume, _ = munition_info.Volume.GetValue()
		munition.SeekerRange, _ = munition_info.SeekerRange.GetValue()
		munition.SeekerType, _ = munition_info.SeekerType.GetValue()

		munition.MunitionLifetime, _ = munition_info.LifeTime.GetValue()

		munition.SeekerFovDeg, _ = munition_info.SeekerFovDeg.GetValue()

		if ammo_ids_name, ok := munition_info.IdsName.GetValue(); ok {
			munition.Name = e.GetInfocardName(ammo_ids_name, munition.Nickname)
		}

		munition.Price = -1
		if good_info, ok := e.Configs.Goods.GoodsMap[munition_info.Nickname.Get()]; ok {
			if price, ok := good_info.Price.GetValue(); ok {
				munition.Price = price
				munition.Bases = e.GetAtBasesSold(GetCommodityAtBasesInput{
					Nickname: good_info.Nickname.Get(),
					Price:    price,
				})
			}
		}

		if !e.Buyable(munition.Bases) {
			continue
		}

		e.exportInfocards(InfocardKey(munition.Nickname), munition.InfoID)
		munition.DiscoveryTechCompat = CalculateTechCompat(e.Configs.Discovery, ids, munition.Nickname)
		tractors = append(tractors, munition)
	}
	return tractors
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
