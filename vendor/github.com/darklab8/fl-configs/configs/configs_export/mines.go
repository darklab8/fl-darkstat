package configs_export

import "math"

type Mine struct {
	Name      string
	Price     int
	AmmoPrice int
	Nickname  string
	IdsName   int
	IdsInfo   int

	HullDamage    int
	EnergyDamange int
	ShieldDamage  int
	PowerUsage    float64

	Value              float64
	Refire             float64
	DetonationDistance float64
	Radius             float64
	SeekDistance       int
	TopSpeed           int
	Acceleration       int
	LinearDrag         float64
	LifeTime           float64
	OwnerSafe          int
	Toughness          float64

	HitPts   int
	Lootable bool

	Bases []*GoodAtBase

	*DiscoveryTechCompat
}

func (e *Exporter) GetMines(ids []Tractor) []Mine {
	var mines []Mine

	for _, mine_dropper := range e.configs.Equip.MineDroppers {
		mine := Mine{}

		mine.Nickname = mine_dropper.Nickname.Get()
		mine.IdsInfo = mine_dropper.IdsInfo.Get()
		mine.IdsName = mine_dropper.IdsName.Get()
		mine.PowerUsage = mine_dropper.PowerUsage.Get()
		mine.Lootable = mine_dropper.Lootable.Get()
		mine.Toughness = mine_dropper.Toughness.Get()
		mine.HitPts = mine_dropper.HitPts.Get()

		if good_info, ok := e.configs.Goods.GoodsMap[mine.Nickname]; ok {
			if price, ok := good_info.Price.GetValue(); ok {
				mine.Price = price
				mine.Bases = e.GetAtBasesSold(GetAtBasesInput{
					Nickname: good_info.Nickname.Get(),
					Price:    price,
				})
			}
		}

		mine.Name = e.GetInfocardName(mine.IdsName, mine.Nickname)

		mine_info := e.configs.Equip.MinesMap[mine_dropper.ProjectileArchetype.Get()]
		explosion := e.configs.Equip.ExplosionMap[mine_info.ExplosionArch.Get()]

		mine.HullDamage = explosion.HullDamage.Get()
		mine.EnergyDamange = explosion.EnergyDamange.Get()
		mine.ShieldDamage = int(float64(mine.HullDamage)*float64(e.configs.Consts.ShieldEquipConsts.HULL_DAMAGE_FACTOR.Get()) + float64(mine.EnergyDamange))

		mine.Radius = float64(explosion.Radius.Get())

		mine.Refire = float64(1 / mine_dropper.RefireDelay.Get())

		mine.DetonationDistance = float64(mine_info.DetonationDistance.Get())
		mine.OwnerSafe = mine_info.OwnerSafeTime.Get()
		mine.SeekDistance = mine_info.SeekDist.Get()
		mine.TopSpeed = mine_info.TopSpeed.Get()
		mine.Acceleration = mine_info.Acceleration.Get()
		mine.LifeTime = mine_info.Lifetime.Get()
		mine.LinearDrag = mine_info.LinearDrag.Get()

		if mine_good_info, ok := e.configs.Goods.GoodsMap[mine_info.Nickname.Get()]; ok {
			if price, ok := mine_good_info.Price.GetValue(); ok {
				mine.AmmoPrice = price
				mine.Value = math.Max(float64(mine.HullDamage), float64(mine.ShieldDamage)) / float64(mine.AmmoPrice)
			}
		}

		e.exportInfocards(InfocardKey(mine.Nickname), mine.IdsInfo)
		mine.DiscoveryTechCompat = CalculateTechCompat(e.configs.Discovery, ids, mine.Nickname)
		mines = append(mines, mine)
	}

	return mines
}

func (e *Exporter) FilterToUsefulMines(mines []Mine) []Mine {
	var items []Mine = make([]Mine, 0, len(mines))
	for _, item := range mines {
		if !e.Buyable(item.Bases) {
			continue
		}
		items = append(items, item)
	}
	return items
}
