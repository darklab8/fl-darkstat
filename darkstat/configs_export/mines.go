package configs_export

import (
	"math"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/initialworld/flhash"
	"github.com/darklab8/go-utils/utils/ptr"
)

type Mine struct {
	Name                string          `json:"name"`
	Price               int             `json:"price"`
	AmmoPrice           int             `json:"ammo_price"`
	Nickname            string          `json:"nickname"`
	MineDropperHash     flhash.HashCode `json:"mine_dropper_hash" format:"int64"`
	ProjectileArchetype string          `json:"projectyle_archetype"`
	MineHash            flhash.HashCode `json:"mine_hash" format:"int64"`
	IdsName             int             `json:"ids_name"`
	IdsInfo             int             `json:"ids_info"`

	HullDamage    int     `json:"hull_damage"`
	EnergyDamange int     `json:"energy_damage"`
	ShieldDamage  int     `json:"shield_damage"`
	PowerUsage    float64 `json:"power_usage"`

	Value              float64 `json:"value"`
	Refire             float64 `json:"refire"`
	DetonationDistance float64 `json:"detonation_distance"`
	Radius             float64 `json:"radius"`
	SeekDistance       int     `json:"seek_distance"`
	TopSpeed           int     `json:"top_speed"`
	Acceleration       int     `json:"acceleration"`
	LinearDrag         float64 `json:"linear_drag"`
	LifeTime           float64 `json:"life_time"`
	OwnerSafe          int     `json:"owner_safe"`
	Toughness          float64 `json:"toughness"`

	HitPts   int  `json:"hit_pts"`
	Lootable bool `json:"lootable"`

	Bases map[cfg.BaseUniNick]*MarketGood `json:"-" swaggerignore:"true"`

	*DiscoveryTechCompat `json:"-" swaggerignore:"true"`

	AmmoLimit AmmoLimit `json:"ammo_limit"`
	Mass      float64   `json:"mass"`
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

	for _, mine_dropper := range e.Configs.Equip.MineDroppers {
		mine := Mine{
			Bases: make(map[cfg.BaseUniNick]*MarketGood),
		}
		mine.Mass, _ = mine_dropper.Mass.GetValue()

		mine.Nickname = mine_dropper.Nickname.Get()
		mine.MineDropperHash = flhash.HashNickname(mine.Nickname)
		e.Hashes[mine.Nickname] = mine.MineDropperHash

		mine.IdsInfo = mine_dropper.IdsInfo.Get()
		mine.IdsName = mine_dropper.IdsName.Get()
		mine.PowerUsage = mine_dropper.PowerUsage.Get()
		mine.Lootable = mine_dropper.Lootable.Get()
		mine.Toughness = mine_dropper.Toughness.Get()
		mine.HitPts = mine_dropper.HitPts.Get()

		if good_info, ok := e.Configs.Goods.GoodsMap[mine.Nickname]; ok {
			if price, ok := good_info.Price.GetValue(); ok {
				mine.Price = price
				mine.Bases = e.GetAtBasesSold(GetCommodityAtBasesInput{
					Nickname: good_info.Nickname.Get(),
					Price:    price,
				})
			}
		}

		mine.Name = e.GetInfocardName(mine.IdsName, mine.Nickname)

		mine_info := e.Configs.Equip.MinesMap[mine_dropper.ProjectileArchetype.Get()]
		mine.ProjectileArchetype = mine_info.Nickname.Get()
		mine.MineHash = flhash.HashNickname(mine.ProjectileArchetype)
		e.Hashes[mine.ProjectileArchetype] = mine.MineHash

		explosion := e.Configs.Equip.ExplosionMap[mine_info.ExplosionArch.Get()]

		mine.HullDamage = explosion.HullDamage.Get()
		mine.EnergyDamange = explosion.EnergyDamange.Get()
		mine.ShieldDamage = int(float64(mine.HullDamage)*float64(e.Configs.Consts.ShieldEquipConsts.HULL_DAMAGE_FACTOR.Get()) + float64(mine.EnergyDamange))

		mine.Radius = float64(explosion.Radius.Get())

		mine.Refire = float64(1 / mine_dropper.RefireDelay.Get())

		mine.DetonationDistance = float64(mine_info.DetonationDistance.Get())
		mine.OwnerSafe = mine_info.OwnerSafeTime.Get()
		mine.SeekDistance = mine_info.SeekDist.Get()
		mine.TopSpeed = mine_info.TopSpeed.Get()
		mine.Acceleration = mine_info.Acceleration.Get()
		mine.LifeTime = mine_info.Lifetime.Get()
		mine.LinearDrag = mine_info.LinearDrag.Get()

		if mine_good_info, ok := e.Configs.Goods.GoodsMap[mine_info.Nickname.Get()]; ok {
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

		e.exportInfocards(InfocardKey(mine.Nickname), mine.IdsInfo)
		mine.DiscoveryTechCompat = CalculateTechCompat(e.Configs.Discovery, ids, mine.Nickname)

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
