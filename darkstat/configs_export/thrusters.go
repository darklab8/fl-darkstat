package configs_export

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/initialworld/flhash"
)

type Thruster struct {
	Name         string          `json:"name"`
	Price        int             `json:"price"`
	MaxForce     int             `json:"max_force"`
	PowerUsage   int             `json:"power_usage"`
	Efficiency   float64         `json:"efficiency"`
	Value        float64         `json:"value"`
	HitPts       int             `json:"hit_pts"`
	Lootable     bool            `json:"lootable"`
	Nickname     string          `json:"nickname"`
	NicknameHash flhash.HashCode `json:"nickname_hash" format:"int64"`
	NameID       int             `json:"name_id"`
	InfoID       int             `json:"info_id"`

	Bases map[cfg.BaseUniNick]*MarketGood `json:"-" swaggerignore:"true"`

	*DiscoveryTechCompat `json:"-" swaggerignore:"true"`
	Mass                 float64 `json:"mass"`
}

func (b Thruster) GetNickname() string { return string(b.Nickname) }

func (b Thruster) GetBases() map[cfg.BaseUniNick]*MarketGood { return b.Bases }

func (b Thruster) GetDiscoveryTechCompat() *DiscoveryTechCompat { return b.DiscoveryTechCompat }

func (e *Exporter) GetThrusters(ids []*Tractor) []Thruster {
	var thrusters []Thruster

	for _, thruster_info := range e.Configs.Equip.Thrusters {
		thruster := Thruster{
			Bases: make(map[cfg.BaseUniNick]*MarketGood),
		}
		thruster.Mass, _ = thruster_info.Mass.GetValue()

		thruster.Nickname = thruster_info.Nickname.Get()
		thruster.NicknameHash = flhash.HashNickname(thruster.Nickname)
		e.Hashes[thruster.Nickname] = thruster.NicknameHash

		thruster.MaxForce = thruster_info.MaxForce.Get()
		thruster.PowerUsage = thruster_info.PowerUsage.Get()
		thruster.HitPts = thruster_info.HitPts.Get()
		thruster.Lootable = thruster_info.Lootable.Get()
		thruster.NameID = thruster_info.IdsName.Get()
		thruster.InfoID = thruster_info.IdsInfo.Get()

		if good_info, ok := e.Configs.Goods.GoodsMap[thruster.Nickname]; ok {
			if price, ok := good_info.Price.GetValue(); ok {
				thruster.Price = price
				thruster.Bases = e.GetAtBasesSold(GetCommodityAtBasesInput{
					Nickname: good_info.Nickname.Get(),
					Price:    price,
				})
			}
		}

		thruster.Name = e.GetInfocardName(thruster.NameID, thruster.Nickname)

		/*
			Copy paste of Adoxa's changelog
			* Efficiency: max_force / power;
				      power if max_force is 0;
			* Value: max_force / price;
				 power * 1000 / price if max_force is 0;
			* Rating: max_force / (power - 100) * Value / 1000 (where 100 is the standard
			  thrust recharge rate).
		*/

		power_usage_calc := thruster.PowerUsage
		if power_usage_calc == 0 {
			power_usage_calc = 1
		}
		if thruster.MaxForce > 0 {
			thruster.Efficiency = float64(thruster.MaxForce) / float64(power_usage_calc)
		} else {
			thruster.Efficiency = float64(thruster.PowerUsage)
		}

		price_calc := thruster.Price
		if price_calc == 0 {
			price_calc = 1
		}
		if thruster.MaxForce > 0 {
			thruster.Value = float64(thruster.MaxForce) / float64(price_calc)
		} else {
			thruster.Value = float64(thruster.Price) * 1000
		}

		e.exportInfocards(InfocardKey(thruster.Nickname), thruster.InfoID)
		thruster.DiscoveryTechCompat = CalculateTechCompat(e.Configs.Discovery, ids, thruster.Nickname)
		thrusters = append(thrusters, thruster)
	}
	return thrusters
}

func (e *Exporter) FilterToUsefulThrusters(thrusters []Thruster) []Thruster {
	var items []Thruster = make([]Thruster, 0, len(thrusters))
	for _, item := range thrusters {
		if !e.Buyable(item.Bases) {
			continue
		}
		items = append(items, item)
	}
	return items
}
