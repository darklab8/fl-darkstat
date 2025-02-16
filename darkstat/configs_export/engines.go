package configs_export

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/equipment_mapped/equip_mapped"
)

func (g Engine) GetTechCompat() *DiscoveryTechCompat { return g.DiscoveryTechCompat }

type Engine struct {
	Name  string `json:"name"  validate:"required"`
	Price int    `json:"price"  validate:"required"`

	CruiseSpeed      int     `json:"cruise_speed"  validate:"required"`
	CruiseChargeTime int     `json:"cruise_charge_time"  validate:"required"`
	LinearDrag       int     `json:"linear_drag"  validate:"required"`
	MaxForce         int     `json:"max_force"  validate:"required"`
	ReverseFraction  float64 `json:"reverse_fraction"  validate:"required"`
	ImpulseSpeed     float64 `json:"impulse_speed"  validate:"required"`

	HpType      string `json:"hp_type"  validate:"required"`
	FlameEffect string `json:"flame_effect"  validate:"required"`
	TrailEffect string `json:"trail_effect"  validate:"required"`

	Nickname string `json:"nickname"  validate:"required"`
	NameID   int    `json:"name_id"  validate:"required"`
	InfoID   int    `json:"info_id"  validate:"required"`

	Bases                map[cfg.BaseUniNick]*MarketGood `json:"-" swaggerignore:"true"`
	*DiscoveryTechCompat `json:"-" swaggerignore:"true"`
	Mass                 float64 `json:"mass"  validate:"required"`
}

func (b Engine) GetNickname() string { return string(b.Nickname) }

func (b Engine) GetBases() map[cfg.BaseUniNick]*MarketGood { return b.Bases }

func (b Engine) GetDiscoveryTechCompat() *DiscoveryTechCompat { return b.DiscoveryTechCompat }

func (e *Exporter) GetEngineSpeed(engine_info *equip_mapped.Engine) int {
	if cruise_speed, ok := engine_info.CruiseSpeed.GetValue(); ok {
		return cruise_speed
	} else {
		if cruise_speed, ok := e.mapped.Consts.EngineEquipConsts.CRUISING_SPEED.GetValue(); ok {
			return cruise_speed
		}
	}
	return 350
}

func (e *Exporter) GetEngines(ids []*Tractor) []Engine {
	var engines []Engine

	for _, engine_info := range e.mapped.Equip.Engines {
		engine := Engine{
			Bases: make(map[cfg.BaseUniNick]*MarketGood),
		}
		engine.Mass, _ = engine_info.Mass.GetValue()

		engine.Nickname = engine_info.Nickname.Get()
		engine.CruiseSpeed = e.GetEngineSpeed(engine_info)
		engine.CruiseChargeTime, _ = engine_info.CruiseChargeTime.GetValue()
		engine.LinearDrag = engine_info.LinearDrag.Get()
		engine.MaxForce = engine_info.MaxForce.Get()
		engine.ReverseFraction = engine_info.ReverseFraction.Get()
		linear_drag_for_calc := engine.LinearDrag
		if linear_drag_for_calc == 0 {
			linear_drag_for_calc = 1
		}
		engine.ImpulseSpeed = float64(engine.MaxForce) / float64(linear_drag_for_calc)

		engine.HpType, _ = engine_info.HpType.GetValue()
		engine.FlameEffect, _ = engine_info.FlameEffect.GetValue()
		engine.TrailEffect, _ = engine_info.TrailEffect.GetValue()

		engine.NameID = engine_info.IdsName.Get()
		engine.InfoID = engine_info.IdsInfo.Get()

		if good_info, ok := e.mapped.Goods.GoodsMap[engine.Nickname]; ok {
			if price, ok := good_info.Price.GetValue(); ok {
				engine.Price = price
				engine.Bases = e.GetAtBasesSold(GetCommodityAtBasesInput{
					Nickname: good_info.Nickname.Get(),
					Price:    price,
				})
			}
		}

		engine.Name = e.GetInfocardName(engine.NameID, engine.Nickname)

		e.exportInfocards(InfocardKey(engine.Nickname), engine.InfoID)

		engine.DiscoveryTechCompat = CalculateTechCompat(e.mapped.Discovery, ids, engine.Nickname)

		engines = append(engines, engine)
	}
	return engines
}

func (e *Exporter) FilterToUsefulEngines(engines []Engine) []Engine {
	var buyable_engines []Engine = make([]Engine, 0, len(engines))
	for _, engine := range engines {
		if !e.Buyable(engine.Bases) {
			continue
		}
		buyable_engines = append(buyable_engines, engine)
	}

	return buyable_engines
}
