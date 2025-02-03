package configs_export

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/initialworld/flhash"
)

type Cloak struct {
	Name  string `json:"name"`
	Price int    `json:"price"`

	HitPts int     `json:"hit_pts"`
	Volume float64 `json:"volume"`

	Nickname string `json:"nickname"`
	NameID   int    `json:"name_id"`
	InfoID   int    `json:"info_id"`

	PowerUsage   float64
	CloakInTime  int
	CloakOutTime int

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
		NicknameHash := flhash.HashNickname(cloak.Nickname)
		e.Hashes[cloak.Nickname] = NicknameHash
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
