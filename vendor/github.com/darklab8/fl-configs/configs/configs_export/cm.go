package configs_export

type CounterMeasure struct {
	Name  string
	Price int

	HitPts        int
	AIRange       int
	AmmoLimit     int
	Lifetime      int
	Range         int
	DiversionPctg int

	Lootable bool
	Nickname string
	NameID   int
	InfoID   int

	Bases []GoodAtBase
}

func (e *Exporter) GetCounterMeasures() []CounterMeasure {
	var tractors []CounterMeasure

	for _, cm_info := range e.configs.Equip.CounterMeasureDroppers {
		cm := CounterMeasure{}
		cm.Nickname = cm_info.Nickname.Get()
		cm.HitPts = cm_info.HitPts.Get()
		cm.AIRange = cm_info.AIRange.Get()
		cm.Lootable = cm_info.Lootable.Get()
		cm.NameID = cm_info.IdsName.Get()
		cm.InfoID = cm_info.IdsInfo.Get()

		if good_info, ok := e.configs.Goods.GoodsMap[cm.Nickname]; ok {
			if price, ok := good_info.Price.GetValue(); ok {
				cm.Price = price
				cm.Bases = e.GetAtBasesSold(GetAtBasesInput{
					Nickname:       good_info.Nickname.Get(),
					Price:          price,
					PricePerVolume: -1,
				})
			}
		}

		if name, ok := e.configs.Infocards.Infonames[cm.NameID]; ok {
			cm.Name = string(name)
		}

		infocards := []int{cm.InfoID}
		if ammo_info, ok := e.configs.Equip.CounterMeasureMap[cm_info.ProjectileArchetype.Get()]; ok {
			cm.AmmoLimit = ammo_info.AmmoLimit.Get()
			cm.Lifetime = ammo_info.Lifetime.Get()
			cm.Range = ammo_info.Range.Get()
			cm.DiversionPctg = ammo_info.DiversionPctg.Get()

			if id, ok := ammo_info.IdsInfo.GetValue(); ok {
				infocards = append(infocards, id)
			}
		}

		e.infocards_parser.Set(InfocardKey(cm.Nickname), infocards...)
		tractors = append(tractors, cm)
	}
	return tractors
}

func FilterToUsefulCounterMeasures(cms []CounterMeasure) []CounterMeasure {
	var useful_items []CounterMeasure = make([]CounterMeasure, 0, len(cms))
	for _, item := range cms {
		if !Buyable(item.Bases) {
			continue
		}
		useful_items = append(useful_items, item)
	}
	return useful_items
}
