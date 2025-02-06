package configs_export

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
)

type Cloak struct {
	Name  string `json:"name" validate:"required"`
	Price int    `json:"price" validate:"required"`

	HitPts int     `json:"hit_pts" validate:"required"`
	Volume float64 `json:"volume" validate:"required"`

	Nickname string `json:"nickname" validate:"required"`
	NameID   int    `json:"name_id" validate:"required"`
	InfoID   int    `json:"info_id" validate:"required"`

	PowerUsage   float64 `json:"power_usage" validate:"required"`
	CloakInTime  int     `json:"cloakintime" validate:"required"`
	CloakOutTime int     `json:"cloakouttime" validate:"required"`

	Bases                map[cfg.BaseUniNick]*MarketGood `json:"-" swaggerignore:"true"`
	*DiscoveryTechCompat `json:"-" swaggerignore:"true"`
}

func (b Cloak) GetNickname() string { return string(b.Nickname) }

func (b Cloak) GetBases() map[cfg.BaseUniNick]*MarketGood { return b.Bases }

func (b Cloak) GetDiscoveryTechCompat() *DiscoveryTechCompat { return b.DiscoveryTechCompat }

func (e *Exporter) GetCloaks(ids []*Tractor) []Cloak {
	var items []Cloak

	for _, cloak_info := range e.Configs.Equip.Cloaks {
		cloak := Cloak{
			Bases: make(map[cfg.BaseUniNick]*MarketGood),
		}
		cloak.PowerUsage, _ = cloak_info.PowerUsage.GetValue()
		cloak.CloakInTime, _ = cloak_info.CloakInTime.GetValue()
		cloak.CloakOutTime, _ = cloak_info.CloakOutTime.GetValue()
		cloak.Volume, _ = cloak_info.Volume.GetValue()

		cloak.Nickname = cloak_info.Nickname.Get()
		cloak.NameID, _ = cloak_info.IdsName.GetValue()
		cloak.InfoID, _ = cloak_info.IdsInfo.GetValue()
		cloak.HitPts, _ = cloak_info.HitPts.GetValue()

		if item_name, ok := cloak_info.IdsName.GetValue(); ok {
			cloak.Name = e.GetInfocardName(item_name, cloak.Nickname)
		}

		cloak.Price = -1
		if good_info, ok := e.Configs.Goods.GoodsMap[cloak_info.Nickname.Get()]; ok {
			if price, ok := good_info.Price.GetValue(); ok {
				cloak.Price = price
				cloak.Bases = e.GetAtBasesSold(GetCommodityAtBasesInput{
					Nickname: good_info.Nickname.Get(),
					Price:    price,
				})
			}
		}

		e.exportInfocards(InfocardKey(cloak.Nickname), cloak.InfoID)
		cloak.DiscoveryTechCompat = CalculateTechCompat(e.Configs.Discovery, ids, cloak.Nickname)
		items = append(items, cloak)
	}
	return items
}
