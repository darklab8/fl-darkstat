package configs_export

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
)

type Thruster struct {
	Name       string  `json:"name" validate:"required"`
	Price      int     `json:"price" validate:"required"`
	MaxForce   int     `json:"max_force" validate:"required"`
	PowerUsage int     `json:"power_usage" validate:"required"`
	Efficiency float64 `json:"efficiency" validate:"required"`
	Value      float64 `json:"value" validate:"required"`
	HitPts     int     `json:"hit_pts" validate:"required"`
	Lootable   bool    `json:"lootable" validate:"required"`
	Nickname   string  `json:"nickname" validate:"required"`
	NameID     int     `json:"name_id" validate:"required"`
	InfoID     int     `json:"info_id" validate:"required"`

	Bases map[cfg.BaseUniNick]*MarketGood `json:"-" swaggerignore:"true"`

	*DiscoveryTechCompat `json:"-" swaggerignore:"true"`
	Mass                 float64 `json:"mass" validate:"required"`
}

func (b Thruster) GetNickname() string { return string(b.Nickname) }

func (b Thruster) GetBases() map[cfg.BaseUniNick]*MarketGood { return b.Bases }

func (b Thruster) GetDiscoveryTechCompat() *DiscoveryTechCompat { return b.DiscoveryTechCompat }

func (e *Exporter) GetThrusters(ids []*Tractor) []Thruster {
	var thrusters []Thruster

	for _, thruster_info := range e.Mapped.Equip().Thrusters {
		thruster := Thruster{
			Bases: make(map[cfg.BaseUniNick]*MarketGood),
		}
		thruster.Mass, _ = thruster_info.Mass.GetValue()

		thruster.Nickname = thruster_info.Nickname.Get()
		thruster.MaxForce = thruster_info.MaxForce.Get()
		thruster.PowerUsage = thruster_info.PowerUsage.Get()
		thruster.HitPts = thruster_info.HitPts.Get()
		thruster.Lootable = thruster_info.Lootable.Get()
		thruster.NameID = thruster_info.IdsName.Get()
		thruster.InfoID = thruster_info.IdsInfo.Get()

		if good_info, ok := e.Mapped.Goods.GoodsMap[thruster.Nickname]; ok {
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
		thruster.DiscoveryTechCompat = CalculateTechCompat(e.Mapped.Discovery, ids, thruster.Nickname)
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
