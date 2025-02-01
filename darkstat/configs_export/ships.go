package configs_export

import (
	"math"
	"sort"
	"strings"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/equipment_mapped/equip_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/initialworld/flhash"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/ship_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
	"github.com/darklab8/go-typelog/typelog"
)

func (g Ship) GetTechCompat() *DiscoveryTechCompat { return g.DiscoveryTechCompat }

type ShipPackage struct {
	Nickname           string
	equipped_thrusters []*equip_mapped.Thruster
}

type Ship struct {
	Nickname     string          `json:"nickname"`
	NicknameHash flhash.HashCode `json:"nickname_hash"`

	Name      string  `json:"name"`
	Class     int     `json:"class"`
	Type      string  `json:"type"`
	Price     int     `json:"price"`
	Armor     int     `json:"armor"`
	HoldSize  int     `json:"hold_size"`
	Nanobots  int     `json:"nanobots"`
	Batteries int     `json:"batteries"`
	Mass      float64 `json:"mass"`

	PowerCapacity     int     `json:"power_capacity"`
	PowerRechargeRate int     `json:"power_recharge_rate"`
	CruiseSpeed       int     `json:"cruise_speed"`
	LinearDrag        float64 `json:"linear_drag"`
	EngineMaxForce    int     `json:"engine_max_force"`
	ImpulseSpeed      float64 `json:"impulse_speed"`
	ThrusterSpeed     []int   `json:"thruster_speed"`
	ReverseFraction   float64 `json:"reverse_fraction"`
	ThrustCapacity    int     `json:"thrust_capacity"`
	ThrustRecharge    int     `json:"thrust_recharge"`

	MaxAngularSpeedDegS           float64 `json:"max_ansgular_speed"`
	AngularDistanceFrom0ToHalfSec float64 `json:"angular_distance_from_0_to_halfsec"`
	TimeTo90MaxAngularSpeed       float64 `json:"time_to_90_max_angular_speed"`

	NudgeForce  float64 `json:"nudge_force"`
	StrafeForce float64 `json:"strafe_force"`
	NameID      int     `json:"name_id"`
	InfoID      int     `json:"info_id"`

	Bases            map[cfg.BaseUniNick]*MarketGood `json:"_" swaggerignore:"true"`
	Slots            []EquipmentSlot                 `json:"equipment_slots"`
	BiggestHardpoint []string                        `json:"biggest_hardpoint"`
	ShipPackages     []ShipPackage                   `json:"ship_packages"`

	*DiscoveryTechCompat `json:"_" swaggerignore:"true"`

	DiscoShip *DiscoShip `json:"discovery_ship"`
}

func (b Ship) GetNickname() string { return string(b.Nickname) }

func (b Ship) GetBases() map[cfg.BaseUniNick]*MarketGood { return b.Bases }

func (b Ship) GetDiscoveryTechCompat() *DiscoveryTechCompat { return b.DiscoveryTechCompat }

/*
For each ship

	ship_packages = find buyable/craftable ship packages
*/
func remove(s []int, i int) []int {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func is_thruster_slot(slot EquipmentSlot) bool {
	for _, smth := range slot.AllowedEquip {
		if smth == "hp_thruster" {
			return true
		}
	}
	return false
}

func (s *Ship) getThrusterSpeed(
	e *Exporter,
	equipped_thrusters []*equip_mapped.Thruster,
	linear_drag float64,
	ship_info *ship_mapped.Ship,
	ThrustCapacity float64,
	ThrustRecharge float64,
	Slots []EquipmentSlot,
	ThrusterMap map[string]*Thruster,
) float64 {
	// find amount of thrusters
	thruster_amount := 0
	for _, slot := range Slots {
		if is_thruster_slot(slot) {
			thruster_amount++
		}
	}

	total_thruster_force := 0

	// find_max_forced_compatible_thruster
	max_thruster_force := 0
	// var found_thruster1 *equip_mapped.Thruster # debug data
	// var found_thruster2 *Thruster # debug data

	for _, thruster := range e.Configs.Equip.Thrusters {
		thrust_usage := thruster.PowerUsage.Get()

		seconds_thrust_usage := int(ThrustCapacity / (float64(thrust_usage*thruster_amount) - ThrustRecharge))
		// 2000 / (2*120000 - 200)
		if seconds_thrust_usage < 0 {
			seconds_thrust_usage = 9999
		}

		// exclude not usable. if they are usable less than 3 seconds
		if seconds_thrust_usage >= 0 && seconds_thrust_usage < 3 {
			continue
		}

		// add check if item is buyable or craftable
		thruster_info, found_thruster := ThrusterMap[thruster.Nickname.Get()]
		if !found_thruster {
			continue
		}
		if !e.Buyable(thruster_info.Bases) {
			continue
		}

		thruster_force := thruster.MaxForce.Get()

		// no point to select weak
		if thruster_force < max_thruster_force {
			continue
		}

		max_thruster_force = thruster_force
		// found_thruster1 = thruster
		// found_thruster2 = thruster_info
	}

	// _ = found_thruster2
	// _ = found_thruster1

	// for each thruster at a ship
	for i := 0; i < thruster_amount; i++ {

		//   if already installed zero price thrustre is installed and has zero price and it is disco
		// 	  add its force
		if i < len(equipped_thrusters) {
			thruster := equipped_thrusters[i]
			thruster_price := 0
			if good_info, ok := e.Configs.Goods.GoodsMap[thruster.Nickname.Get()]; ok {
				if price, ok := good_info.Price.GetValue(); ok {
					thruster_price = price
				}
			}
			if e.Configs.Discovery != nil && thruster_price == 0 {
				total_thruster_force += thruster.MaxForce.Get()
				continue
			}
		}

		//   else:
		// 	  add max forced compatible thruster
		total_thruster_force += max_thruster_force
	}

	return s.ImpulseSpeed + float64(total_thruster_force)/linear_drag
}

type DiscoShip struct {
	ArmorMult float64 `json:"armor_mult"`
}

func (e *Exporter) GetShips(ids []*Tractor, TractorsByID map[cfg.TractorID]*Tractor, Thrusters []Thruster) []Ship {
	var ships []Ship

	var ThrusterMap map[string]*Thruster = make(map[string]*Thruster)
	for _, thruster := range Thrusters {
		ThrusterMap[thruster.Nickname] = &thruster
	}

	for _, ship_info := range e.Configs.Shiparch.Ships {
		ship := Ship{
			Nickname: ship_info.Nickname.Get(),
			Bases:    make(map[cfg.BaseUniNick]*MarketGood),
		}
		ship.NicknameHash = flhash.HashNickname(ship.Nickname)
		e.Hashes[ship.Nickname] = ship.NicknameHash

		// defer func() {
		// 	if r := recover(); r != nil {
		// 		fmt.Println("Recovered in f", r)
		// 		fmt.Println("ship.Nickname", ship.Nickname)
		// 		panic(r)
		// 	}
		// }()

		ship.Class, _ = ship_info.ShipClass.GetValue()
		if _, ok := ship_info.Type.GetValue(); !ok {
			logus.Log.Warn("ship problem with type", typelog.Any("nickname", ship.Nickname))
		}
		ship.Type, _ = ship_info.Type.GetValue()
		ship.Type = strings.ToLower(ship.Type)

		if ship_name_id, ship_has_name := ship_info.IdsName.GetValue(); ship_has_name {
			ship.NameID = ship_name_id
		} else {
			logus.Log.Warn("WARNING, ship has no ItdsName", typelog.String("ship.Nickname", ship.Nickname))
		}

		ship.InfoID, _ = ship_info.IdsInfo.GetValue()

		if bots, ok := ship_info.Nanobots.GetValue(); ok {
			ship.Nanobots = bots
		} else {
			continue
		}
		ship.Batteries = ship_info.Batteries.Get()
		ship.Mass = ship_info.Mass.Get()
		ship.NudgeForce = ship_info.NudgeForce.Get()
		ship.StrafeForce, _ = ship_info.StrafeForce.GetValue()

		ship.Name = e.GetInfocardName(ship.NameID, ship.Nickname)

		if ship_hull_good, ok := e.Configs.Goods.ShipHullsMapByShip[ship.Nickname]; ok {
			ship.Price = ship_hull_good.Price.Get()

			ship_hull_nickname := ship_hull_good.Nickname.Get()
			if ship_package_goods, ok := e.Configs.Goods.ShipsMapByHull[ship_hull_nickname]; ok {

				for _, ship_package_good := range ship_package_goods {
					var equipped_thrusters []*equip_mapped.Thruster
					for _, addon := range ship_package_good.Addons {

						// can be Power or Engine or Smth else
						// addon = dsy_hessian_engine, HpEngine01, 1
						// addon = dsy_loki_core, internal, 1
						// addon = ge_s_scanner_01, internal, 1
						addon_nickname := addon.ItemNickname.Get()

						if good_info, ok := e.Configs.Goods.GoodsMap[addon_nickname]; ok {
							if addon_price, ok := good_info.Price.GetValue(); ok {
								ship.Price += addon_price
							}
						}
						if thruster, ok := e.Configs.Equip.ThrusterMap[addon_nickname]; ok {
							equipped_thrusters = append(equipped_thrusters, thruster)
						}
						if power, ok := e.Configs.Equip.PowersMap[addon_nickname]; ok {
							ship.PowerCapacity = power.Capacity.Get()
							ship.PowerRechargeRate = power.ChargeRate.Get()

							ship.ThrustCapacity = power.ThrustCapacity.Get()
							ship.ThrustRecharge = power.ThrustRecharge.Get()
						}
						if engine, ok := e.Configs.Equip.EnginesMap[addon_nickname]; ok {
							ship.CruiseSpeed = e.GetEngineSpeed(engine)
							engine_linear_drag, _ := engine.LinearDrag.GetValue()
							ship_linear_drag, _ := ship_info.LinearDrag.GetValue()
							ship.EngineMaxForce, _ = engine.MaxForce.GetValue()
							ship.LinearDrag = (float64(engine_linear_drag) + float64(ship_linear_drag))
							ship.ImpulseSpeed = float64(ship.EngineMaxForce) / ship.LinearDrag

							ship.ReverseFraction = engine.ReverseFraction.Get()

							ship.MaxAngularSpeedDegS = ship_info.SteeringTorque.X.Get() / ship_info.AngularDrag.X.Get()
							ship.TimeTo90MaxAngularSpeed = ship_info.RotationIntertia.X.Get() / (ship_info.AngularDrag.X.Get() * LogOgE)

							ship.MaxAngularSpeedDegS *= Pi180

							if ship.TimeTo90MaxAngularSpeed > 0.5 {
								ship.AngularDistanceFrom0ToHalfSec = ship.MaxAngularSpeedDegS * (0.5 / ship.TimeTo90MaxAngularSpeed) / 2
							} else {
								ship.AngularDistanceFrom0ToHalfSec = ship.MaxAngularSpeedDegS*(0.5-ship.TimeTo90MaxAngularSpeed) + ship.MaxAngularSpeedDegS*ship.TimeTo90MaxAngularSpeed/2
							}
						}
					}

					ships_at_bases := e.GetAtBasesSold(GetCommodityAtBasesInput{
						Nickname: ship_package_good.Nickname.Get(),
						Price:    ship.Price,
					})
					for key, value := range ships_at_bases {
						ship.Bases[key] = value
					}

					ship.ShipPackages = append(ship.ShipPackages,
						ShipPackage{
							Nickname:           ship_package_good.Nickname.Get(),
							equipped_thrusters: equipped_thrusters,
						},
					)

				}

			}

		}

		ship.HoldSize = ship_info.HoldSize.Get()
		ship.Armor = ship_info.HitPts.Get()

		var hardpoints map[string][]string = make(map[string][]string)
		for _, hp_type := range ship_info.HpTypes {
			for _, equipment := range hp_type.AllowedEquipments {
				equipment_slot := equipment.Get()
				hardpoints[equipment_slot] = append(hardpoints[equipment_slot], hp_type.Nickname.Get())
			}
		}

		for slot_name, allowed_equip := range hardpoints {
			ship.Slots = append(ship.Slots, EquipmentSlot{
				SlotName:     slot_name,
				AllowedEquip: allowed_equip,
			})
		}

		sort.Slice(ship.Slots, func(i, j int) bool {
			return ship.Slots[i].SlotName < ship.Slots[j].SlotName
		})
		for _, slot := range ship.Slots {
			sort.Slice(slot.AllowedEquip, func(i, j int) bool {
				return slot.AllowedEquip[i] < slot.AllowedEquip[j]
			})
		}

		for _, slot := range ship.Slots {
			if len(slot.AllowedEquip) > len(ship.BiggestHardpoint) {
				ship.BiggestHardpoint = slot.AllowedEquip
			}
		}

		var infocards []int
		if id, ok := ship_info.IdsInfo1.GetValue(); ok {
			infocards = append(infocards, id)
		}
		// if id, ok := ship_info.IdsInfo2.GetValue(); ok {
		// 	infocards = append(infocards, id)
		// }
		// Nobody uses it?
		// if id, ok := ship_info.IdsInfo3.GetValue(); ok {
		// 	infocards = append(infocards, id)
		// }
		if id, ok := ship_info.IdsInfo.GetValue(); ok {
			infocards = append(infocards, id)
		}
		e.exportInfocards(InfocardKey(ship.Nickname), infocards...)
		ship.DiscoveryTechCompat = CalculateTechCompat(e.Configs.Discovery, ids, ship.Nickname)

		if e.Configs.Discovery != nil {
			armor_mult, _ := ship_info.ArmorMult.GetValue()
			ship.DiscoShip = &DiscoShip{ArmorMult: armor_mult}
		}

		var thruster_speeds map[int]bool = make(map[int]bool)
		for _, ship_package := range ship.ShipPackages {

			thrust_speed := ship.getThrusterSpeed(e,
				ship_package.equipped_thrusters,
				ship.LinearDrag,
				ship_info,
				float64(ship.ThrustCapacity),
				float64(ship.ThrustRecharge),
				ship.Slots,
				ThrusterMap,
			)
			thruster_speeds[int(thrust_speed)] = true
		}
		for thrust_speed, _ := range thruster_speeds {
			ship.ThrusterSpeed = append(ship.ThrusterSpeed, thrust_speed)
		}

		ships = append(ships, ship)
	}

	return ships
}

type EquipmentSlot struct {
	SlotName     string
	AllowedEquip []string
}

var Pi180 = 180 / math.Pi // number turning radians to degrees
var LogOgE = math.Log10(math.E)

func (e *Exporter) FilterToUsefulShips(ships []Ship) []Ship {
	var items []Ship = make([]Ship, 0, len(ships))
	for _, item := range ships {
		if !e.Buyable(item.Bases) {
			continue
		}
		items = append(items, item)
	}
	return items
}

type CompatibleIDsForTractor struct {
	TechCompat float64
	Tractor    *Tractor
}
