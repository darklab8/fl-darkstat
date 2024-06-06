package configs_export

type Thruster struct {
	Name       string
	Price      int
	MaxForce   int
	PowerUsage int
	Efficiency float64
	Value      float64
	Rating     float64
	HitPts     int
	Lootable   bool
	Nickname   string
	NameID     int
	InfoID     int

	Bases []GoodAtBase

	*DiscoveryTechCompat
}

func (e *Exporter) GetThrusters(ids []Tractor) []Thruster {
	var thrusters []Thruster

	for _, thruster_info := range e.configs.Equip.Thrusters {
		thruster := Thruster{}
		thruster.Nickname = thruster_info.Nickname.Get()
		thruster.MaxForce = thruster_info.MaxForce.Get()
		thruster.PowerUsage = thruster_info.PowerUsage.Get()
		thruster.HitPts = thruster_info.HitPts.Get()
		thruster.Lootable = thruster_info.Lootable.Get()
		thruster.NameID = thruster_info.IdsName.Get()
		thruster.InfoID = thruster_info.IdsInfo.Get()

		if good_info, ok := e.configs.Goods.GoodsMap[thruster.Nickname]; ok {
			if price, ok := good_info.Price.GetValue(); ok {
				thruster.Price = price
				thruster.Bases = e.GetAtBasesSold(GetAtBasesInput{
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
		thruster.DiscoveryTechCompat = CalculateTechCompat(e.configs.Discovery, ids, thruster.Nickname)
		thrusters = append(thrusters, thruster)
	}
	return thrusters
}

func FilterToUsefulThrusters(thrusters []Thruster) []Thruster {
	var items []Thruster = make([]Thruster, 0, len(thrusters))
	for _, item := range thrusters {
		if !Buyable(item.Bases) {
			continue
		}
		items = append(items, item)
	}
	return items
}
