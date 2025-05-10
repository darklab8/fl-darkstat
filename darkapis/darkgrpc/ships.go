package darkgrpc

import (
	"context"

	pb "github.com/darklab8/fl-darkstat/darkapis/darkgrpc/statproto"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
)

func (s *Server) GetShips(_ context.Context, in *pb.GetEquipmentInput) (*pb.GetShipsReply, error) {
	if s.app_data != nil {
		s.app_data.RLock()
		defer s.app_data.RUnlock()
	}

	var input []configs_export.Ship
	if in.FilterToUseful {
		input = s.app_data.Configs.FilterToUsefulShips(s.app_data.Configs.Ships)
	} else {
		input = s.app_data.Configs.Ships
	}
	input = FilterNicknames(in.FilterNicknames, input)

	var items []*pb.Ship
	for _, item := range input {
		result := &pb.Ship{
			Nickname:                      item.Nickname,
			Name:                          item.Name,
			Class:                         int64(item.Class),
			Type:                          item.Type,
			Price:                         int64(item.Price),
			Armor:                         int64(item.Armor),
			HoldSize:                      int64(item.HoldSize),
			Nanobots:                      int64(item.Nanobots),
			Batteries:                     int64(item.Batteries),
			Mass:                          item.Mass,
			PowerCapacity:                 int64(item.PowerCapacity),
			PowerRechargeRate:             int64(item.PowerRechargeRate),
			CruiseSpeed:                   int64(item.CruiseSpeed),
			LinearDrag:                    item.LinearDrag,
			EngineMaxForce:                int64(item.EngineMaxForce),
			ImpulseSpeed:                  item.ImpulseSpeed,
			ReverseFraction:               item.ReverseFraction,
			ThrustCapacity:                int64(item.ThrustCapacity),
			ThrustRecharge:                int64(item.ThrustRecharge),
			MaxAngularSpeedDegS:           item.MaxAngularSpeedDegS,
			AngularDistanceFrom0ToHalfSec: item.AngularDistanceFrom0ToHalfSec,
			TimeTo90MaxAngularSpeed:       item.TimeTo90MaxAngularSpeed,
			NudgeForce:                    item.NudgeForce,
			StrafeForce:                   item.StrafeForce,
			NameId:                        int64(item.NameID),
			InfoId:                        int64(item.InfoID),
		}
		for _, thruster_speed := range item.ThrusterSpeed {
			result.ThrusterSpeed = append(result.ThrusterSpeed, int64(thruster_speed))
		}
		for _, slot := range item.Slots {
			result.Slots = append(result.Slots, &pb.EquipmentSlot{
				SlotName:     slot.SlotName,
				AllowedEquip: slot.AllowedEquip,
			})
		}
		result.BiggestHardpoint = item.BiggestHardpoint
		for _, ship_package := range item.ShipPackages {
			result.ShipPackages = append(result.ShipPackages, &pb.ShipPackage{
				Nickname: ship_package.Nickname,
			})
		}
		if item.DiscoShip != nil {
			result.DiscoShip = &pb.DiscoShip{
				ArmorMult: item.DiscoShip.ArmorMult,
			}
		}

		if in.IncludeMarketGoods {
			result.Bases = NewBases(item.Bases)
		}
		if in.IncludeTechCompat {
			result.DiscoveryTechCompat = NewTechCompat(item.DiscoveryTechCompat)
		}
		items = append(items, result)
	}
	return &pb.GetShipsReply{Items: items}, nil
}
