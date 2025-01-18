package configs_export

import (
	"github.com/darklab8/fl-configs/configs/cfgtype"
	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/initialworld/flhash"
)

type Thruster struct {
	Name         string
	Price        int
	MaxForce     int
	PowerUsage   int
	Efficiency   float64
	Value        float64
	Rating       float64
	HitPts       int
	Lootable     bool
	Nickname     string
	NicknameHash flhash.HashCode
	NameID       int
	InfoID       int

	Bases map[cfgtype.BaseUniNick]*GoodAtBase

	*DiscoveryTechCompat
	Mass float64
}

func (e *Exporter) GetThrusters(ids []Tractor) []Thruster {
	var thrusters []Thruster

	for _, thruster_info := range e.Configs.Equip.Thrusters {
		thruster := Thruster{
			Bases: make(map[cfgtype.BaseUniNick]*GoodAtBase),
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

		if thruster.MaxForce > 0 {
			thruster.Efficiency = float64(thruster.MaxForce) / float64(thruster.PowerUsage)
		} else {
			thruster.Efficiency = float64(thruster.PowerUsage)
		}

		if thruster.MaxForce > 0 {
			thruster.Value = float64(thruster.MaxForce) / float64(thruster.Price)
		} else {
			thruster.Value = float64(thruster.Price) * 1000
		}

		thruster.Rating = float64(thruster.MaxForce) / float64(thruster.PowerUsage-100) * thruster.Value / 1000
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
