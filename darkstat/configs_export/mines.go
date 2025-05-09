package configs_export

import (
	"math"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export/infocarder"
	"github.com/darklab8/go-utils/utils/ptr"
)

type Mine struct {
	Name                string `json:"name" validate:"required"`
	Price               int    `json:"price" validate:"required"`
	AmmoPrice           int    `json:"ammo_price" validate:"required"`
	Nickname            string `json:"nickname" validate:"required"`
	ProjectileArchetype string `json:"projectyle_archetype" validate:"required"`
	IdsName             int    `json:"ids_name" validate:"required"`
	IdsInfo             int    `json:"ids_info" validate:"required"`

	HullDamage    int     `json:"hull_damage" validate:"required"`
	EnergyDamange int     `json:"energy_damage" validate:"required"`
	ShieldDamage  int     `json:"shield_damage" validate:"required"`
	PowerUsage    float64 `json:"power_usage" validate:"required"`

	Value              float64 `json:"value" validate:"required"`
	Refire             float64 `json:"refire" validate:"required"`
	DetonationDistance float64 `json:"detonation_distance" validate:"required"`
	Radius             float64 `json:"radius" validate:"required"`
	SeekDistance       int     `json:"seek_distance" validate:"required"`
	TopSpeed           int     `json:"top_speed" validate:"required"`
	Acceleration       int     `json:"acceleration" validate:"required"`
	LinearDrag         float64 `json:"linear_drag" validate:"required"`
	LifeTime           float64 `json:"life_time" validate:"required"`
	OwnerSafe          int     `json:"owner_safe" validate:"required"`
	Toughness          float64 `json:"toughness" validate:"required"`

	HitPts   int  `json:"hit_pts" validate:"required"`
	Lootable bool `json:"lootable" validate:"required"`

	Bases map[cfg.BaseUniNick]*MarketGood `json:"-" swaggerignore:"true"`

	*DiscoveryTechCompat `json:"-" swaggerignore:"true"`

	AmmoLimit AmmoLimit `json:"ammo_limit" validate:"required"`
	Mass      float64   `json:"mass" validate:"required"`
}

func (b Mine) GetNickname() string { return string(b.Nickname) }

func (b Mine) GetBases() map[cfg.BaseUniNick]*MarketGood { return b.Bases }

func (b Mine) GetDiscoveryTechCompat() *DiscoveryTechCompat { return b.DiscoveryTechCompat }

type AmmoLimit struct {
	// Disco stuff
	AmountInCatridge *int
	MaxCatridges     *int
}

func (e *Exporter) GetMines(ids []*Tractor) []Mine {
	var mines []Mine

	for _, mine_dropper := range e.Mapped.Equip().MineDroppers {
		mine := Mine{
			Bases: make(map[cfg.BaseUniNick]*MarketGood),
		}
		mine.Mass, _ = mine_dropper.Mass.GetValue()

		mine.Nickname = mine_dropper.Nickname.Get()

		mine.IdsInfo = mine_dropper.IdsInfo.Get()
		mine.IdsName = mine_dropper.IdsName.Get()
		mine.PowerUsage = mine_dropper.PowerUsage.Get()
		mine.Lootable = mine_dropper.Lootable.Get()
		mine.Toughness = mine_dropper.Toughness.Get()
		mine.HitPts = mine_dropper.HitPts.Get()

		if good_info, ok := e.Mapped.Goods.GoodsMap[mine.Nickname]; ok {
			if price, ok := good_info.Price.GetValue(); ok {
				mine.Price = price
				mine.Bases = e.GetAtBasesSold(GetCommodityAtBasesInput{
					Nickname: good_info.Nickname.Get(),
					Price:    price,
				})
			}
		}

		mine.Name = e.GetInfocardName(mine.IdsName, mine.Nickname)

		mine_info := e.Mapped.Equip().MinesMap[mine_dropper.ProjectileArchetype.Get()]
		mine.ProjectileArchetype = mine_info.Nickname.Get()

		explosion := e.Mapped.Equip().ExplosionMap[mine_info.ExplosionArch.Get()]

		mine.HullDamage = explosion.HullDamage.Get()
		mine.EnergyDamange = explosion.EnergyDamange.Get()
		mine.ShieldDamage = int(float64(mine.HullDamage)*float64(e.Mapped.Consts.ShieldEquipConsts.HULL_DAMAGE_FACTOR.Get()) + float64(mine.EnergyDamange))

		mine.Radius = float64(explosion.Radius.Get())

		mine.Refire = float64(1 / mine_dropper.RefireDelay.Get())

		mine.DetonationDistance = float64(mine_info.DetonationDistance.Get())
		mine.OwnerSafe = mine_info.OwnerSafeTime.Get()
		mine.SeekDistance = mine_info.SeekDist.Get()
		mine.TopSpeed = mine_info.TopSpeed.Get()
		mine.Acceleration = mine_info.Acceleration.Get()
		mine.LifeTime = mine_info.Lifetime.Get()
		mine.LinearDrag = mine_info.LinearDrag.Get()

		if mine_good_info, ok := e.Mapped.Goods.GoodsMap[mine_info.Nickname.Get()]; ok {
			if price, ok := mine_good_info.Price.GetValue(); ok {
				mine.AmmoPrice = price
				mine.Value = math.Max(float64(mine.HullDamage), float64(mine.ShieldDamage)) / float64(mine.AmmoPrice)
			}
		}

		if value, ok := mine_info.AmmoLimitAmountInCatridge.GetValue(); ok {
			mine.AmmoLimit.AmountInCatridge = ptr.Ptr(value)
		}
		if value, ok := mine_info.AmmoLimitMaxCatridges.GetValue(); ok {
			mine.AmmoLimit.MaxCatridges = ptr.Ptr(value)
		}

		e.exportInfocards(infocarder.InfocardKey(mine.Nickname), mine.IdsInfo)
		mine.DiscoveryTechCompat = CalculateTechCompat(e.Mapped.Discovery, ids, mine.Nickname)

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
