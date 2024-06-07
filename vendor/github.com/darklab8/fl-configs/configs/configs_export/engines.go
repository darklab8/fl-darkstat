package configs_export

import "github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/equipment_mapped/equip_mapped"

type Engine struct {
	Name  string
	Price int

	CruiseSpeed      int
	CruiseChargeTime int
	LinearDrag       int
	MaxForce         int
	ReverseFraction  float64
	ImpulseSpeed     float64

	HpType      string
	FlameEffect string
	TrailEffect string

	Nickname string
	NameID   int
	InfoID   int

	Bases []GoodAtBase
	*DiscoveryTechCompat
}

func (e *Exporter) GetEngineSpeed(engine_info *equip_mapped.Engine) int {
	if cruise_speed, ok := engine_info.CruiseSpeed.GetValue(); ok {
		return cruise_speed
	} else {
		if cruise_speed, ok := e.configs.Consts.EngineEquipConsts.CRUISING_SPEED.GetValue(); ok {
			return cruise_speed
		}
	}
	return 350
}

func (e *Exporter) GetEngines(ids []Tractor) []Engine {
	var engines []Engine

	for _, engine_info := range e.configs.Equip.Engines {
		engine := Engine{}
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

		if good_info, ok := e.configs.Goods.GoodsMap[engine.Nickname]; ok {
			if price, ok := good_info.Price.GetValue(); ok {
				engine.Price = price
				engine.Bases = e.GetAtBasesSold(GetAtBasesInput{
					Nickname: good_info.Nickname.Get(),
					Price:    price,
				})
			}
		}

		engine.Name = e.GetInfocardName(engine.NameID, engine.Nickname)

		e.exportInfocards(InfocardKey(engine.Nickname), engine.InfoID)

		if engine.HpType == "" {
			continue
		}

		engine.DiscoveryTechCompat = CalculateTechCompat(e.configs.Discovery, ids, engine.Nickname)
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
