package configs_export

import (
	"github.com/darklab8/fl-configs/configs/cfgtype"
	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/equipment_mapped/equip_mapped"
	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/initialworld/flhash"
)

func (g Engine) GetNickname() string                 { return g.Nickname }
func (g Engine) GetTechCompat() *DiscoveryTechCompat { return g.DiscoveryTechCompat }

type Engine struct {
	Name  string
	Price int

	CruiseSpeed      int
	CruiseChargeTime int
	LinearDrag       int
	MaxForce         int
	ReverseFraction  float64
	ImpulseSpeed     float64

	HpType          string
	HpTypeHash      flhash.HashCode
	FlameEffect     string
	FlameEffectHash flhash.HashCode
	TrailEffect     string
	TrailEffectHash flhash.HashCode

	Nickname     string
	NicknameHash flhash.HashCode
	NameID       int
	InfoID       int

	Bases map[cfgtype.BaseUniNick]*GoodAtBase
	*DiscoveryTechCompat
	Mass float64
}

func (e *Exporter) GetEngineSpeed(engine_info *equip_mapped.Engine) int {
	if cruise_speed, ok := engine_info.CruiseSpeed.GetValue(); ok {
		return cruise_speed
	} else {
		if cruise_speed, ok := e.Configs.Consts.EngineEquipConsts.CRUISING_SPEED.GetValue(); ok {
			return cruise_speed
		}
	}
	return 350
}

func (e *Exporter) GetEngines(ids []Tractor) []Engine {
	var engines []Engine

	for _, engine_info := range e.Configs.Equip.Engines {
		engine := Engine{
			Bases: make(map[cfgtype.BaseUniNick]*GoodAtBase),
		}
		engine.Mass, _ = engine_info.Mass.GetValue()

		engine.Nickname = engine_info.Nickname.Get()
		engine.CruiseSpeed = e.GetEngineSpeed(engine_info)
		engine.CruiseChargeTime, _ = engine_info.CruiseChargeTime.GetValue()
		engine.LinearDrag = engine_info.LinearDrag.Get()
		engine.MaxForce = engine_info.MaxForce.Get()
		engine.ReverseFraction = engine_info.ReverseFraction.Get()
		engine.ImpulseSpeed = float64(engine.MaxForce) / float64(engine.LinearDrag)

		engine.HpType, _ = engine_info.HpType.GetValue()
		engine.FlameEffect, _ = engine_info.FlameEffect.GetValue()
		engine.TrailEffect, _ = engine_info.TrailEffect.GetValue()

		engine.NameID = engine_info.IdsName.Get()
		engine.InfoID = engine_info.IdsInfo.Get()

		if good_info, ok := e.Configs.Goods.GoodsMap[engine.Nickname]; ok {
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

		engine.DiscoveryTechCompat = CalculateTechCompat(e.Configs.Discovery, ids, engine.Nickname)
		engine.NicknameHash = flhash.HashNickname(engine.Nickname)
		engine.HpTypeHash = flhash.HashNickname(engine.HpType)
		engine.FlameEffectHash = flhash.HashNickname(engine.FlameEffect)
		engine.TrailEffectHash = flhash.HashNickname(engine.TrailEffect)

		e.Hashes[engine.Nickname] = engine.NicknameHash
		e.Hashes[engine.HpType] = engine.HpTypeHash
		e.Hashes[engine.FlameEffect] = engine.FlameEffectHash
		e.Hashes[engine.TrailEffect] = engine.TrailEffectHash

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
