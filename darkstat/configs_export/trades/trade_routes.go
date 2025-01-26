package trades

import (
	"math"
	"strings"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/initialworld/flhash"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped/systems_mapped"
	"github.com/darklab8/fl-darkstat/darkstat/settings"
	"github.com/darklab8/go-utils/utils/ptr"
)

type SystemObject struct {
	nickname string
	pos      cfg.Vector
}

func DistanceForVecs(Pos1 cfg.Vector, Pos2 cfg.Vector) float64 {
	// if _, ok := Pos1.X.GetValue(); !ok {
	// 	return 0, errors.New("no x")
	// }
	// if _, ok := Pos2.X.GetValue(); !ok {
	// 	return 0, errors.New("no x")
	// }

	x_dist := math.Pow((Pos1.X - Pos2.X), 2)
	y_dist := math.Pow((Pos1.Y - Pos2.Y), 2)
	z_dist := math.Pow((Pos1.Z - Pos2.Z), 2)
	distance := math.Pow((x_dist + y_dist + z_dist), 0.5)
	return distance
}

type WithFreighterPaths bool

type RouteShipType int64

const (
	RouteTransport RouteShipType = iota
	RouteFrigate
	RouteFreighter
)

type ShipSpeeds struct {
	AvgTransportCruiseSpeed int
	AvgFrigateCruiseSpeed   int
	AvgFreighterCruiseSpeed int
}

var VanillaSpeeds ShipSpeeds = ShipSpeeds{
	AvgTransportCruiseSpeed: 350,
	AvgFrigateCruiseSpeed:   350,
	AvgFreighterCruiseSpeed: 350,
}

var DiscoverySpeeds ShipSpeeds = ShipSpeeds{
	AvgTransportCruiseSpeed: 350, // TODO You should grab those speeds from some ship example
	AvgFrigateCruiseSpeed:   500, // TODO You should grab those speeds from some ship example
	AvgFreighterCruiseSpeed: 500, // TODO You should grab those speeds from some ship example
}

var FLSRSpeeds ShipSpeeds = ShipSpeeds{
	AvgTransportCruiseSpeed: 500, // TODO You should grab those speeds from some ship example
	AvgFrigateCruiseSpeed:   500, // TODO You should grab those speeds from some ship example
	AvgFreighterCruiseSpeed: 500, // TODO You should grab those speeds from some ship example
}

const (
	// already accounted for

	// Add for every pair of jumphole in path
	JumpHoleDelaySec = 15 // and jump gate
	// add for every tradelane vertex pair in path
	TradeLaneDockingDelaySec = 10
	// add just once
	BaseDockingDelay = 10
)

type ExtraBase struct {
	Pos      cfg.Vector
	Nickname cfg.BaseUniNick
}

/*
Algorithm should be like this:
We iterate through list of Systems:
Adding all bases, jump gates, jump holes, tradelanes as Vertexes.
We scan in advance nicknames for object on another side of jump gate/hole and add it as vertix
We calculcate distances between them. Distance between jump connections is 0 (or time to wait measured in distance)
We calculate distances between trade lanes as shorter than real distance for obvious reasons.
The matrix built on a fight run will be having connections between vertixes as hashmaps of possible edges? For optimized memory consumption in a sparse matrix.

Then on second run, knowing amount of vertixes
We build Floyd matrix? With allocating memory in bulk it should be rather rapid may be.
And run Floud algorithm.
Thus we have stuff calculated for distances between all possible trading locations. (edited)
[6:02 PM]
====
Then we build table of Bases as starting points.
And on click we show proffits of delivery to some location. With time of delivery. And profit per time.
[6:02 PM]
====
Optionally print sum of two best routes that can be started within close range from each other.
*/
type MappingOptions struct {
	TradeRoutesDetailedTradeLane *bool
}

func MapConfigsToFGraph(
	configs *configs_mapped.MappedConfigs,
	avgCruiseSpeed int,
	with_freighter_paths WithFreighterPaths,
	extra_bases_by_system map[string][]ExtraBase,
	opts MappingOptions,
) *GameGraph {
	if opts.TradeRoutesDetailedTradeLane == nil {
		opts.TradeRoutesDetailedTradeLane = ptr.Ptr(settings.Env.TradeRoutesDetailedTradeLane)
	}
	average_trade_lane_speed := configs.GetAvgTradeLaneSpeed()

	graph := NewGameGraph(avgCruiseSpeed, with_freighter_paths)
	for _, system := range configs.Systems.Systems {
		system_speed_multiplier := configs.Overrides.GetSystemSpeedMultiplier(system.Nickname)

		var system_objects []SystemObject = make([]SystemObject, 0, 50)

		if bases, ok := extra_bases_by_system[system.Nickname]; ok {
			for _, base := range bases {
				object := SystemObject{
					nickname: base.Nickname.ToStr(),
					pos:      base.Pos,
				}
				graph.SetIdsName(object.nickname, int(flhash.HashNickname(object.nickname)))

				for _, existing_object := range system_objects {
					distance := graph.DistanceToTime(
						DistanceForVecs(object.pos, existing_object.pos),
						system_speed_multiplier,
					) + BaseDockingDelay*PrecisionMultipiler
					graph.SetEdge(object.nickname, existing_object.nickname, distance)
					graph.SetEdge(existing_object.nickname, object.nickname, distance)
				}

				graph.AllowedVertixesForCalcs[VertexName(object.nickname)] = true

				system_objects = append(system_objects, object)
			}
		}

		for _, system_obj := range system.Bases {
			// system_base_base := system_obj.Base.Get()
			system_base_base, dockable := system_obj.DockWith.GetValue()

			if !dockable {
				continue
			}
			object := SystemObject{
				nickname: system_base_base,
				pos:      system_obj.Pos.Get(),
			}
			graph.SetIdsName(object.nickname, system_obj.IdsName.Get())

			if system_obj.Archetype.Get() == systems_mapped.BaseArchetypeInvisible {
				continue
			}

			object_nickname := system_obj.Nickname.Get()
			if _, ok := configs.InitialWorld.LockedGates[flhash.HashNickname(object_nickname)]; ok {
				continue
			}

			// get all objects with same Base?
			// Check if any of them has docking sphere medium

			if configs.Discovery != nil {
				is_dockable_by_transports := false
				if bases, ok := system.AllBasesByDockWith[system_base_base]; ok {
					for _, base_obj := range bases {
						base_archetype := base_obj.Archetype.Get()
						if solar, ok := configs.Solararch.SolarsByNick[base_archetype]; ok {
							if solar.IsDockableByCaps() {
								is_dockable_by_transports = true
							}
						}
					}
				}
				if !is_dockable_by_transports && bool(!with_freighter_paths) {
					continue
				}
			}

			// Lets allow flying between all bases
			// goods, goods_defined := configs.Market.GoodsPerBase[object.nickname]
			// if !goods_defined {
			// 	continue
			// }

			// if len(goods.MarketGoods) == 0 {
			// 	continue
			// }

			for _, existing_object := range system_objects {
				distance := graph.DistanceToTime(
					DistanceForVecs(object.pos, existing_object.pos),
					system_speed_multiplier,
				) + BaseDockingDelay*PrecisionMultipiler
				graph.SetEdge(object.nickname, existing_object.nickname, distance)
				graph.SetEdge(existing_object.nickname, object.nickname, distance)
			}

			graph.AllowedVertixesForCalcs[VertexName(object.nickname)] = true

			system_objects = append(system_objects, object)
		}

		for _, jumphole := range system.Jumpholes {
			object := SystemObject{
				nickname: jumphole.Nickname.Get(),
				pos:      jumphole.Pos.Get(),
			}
			graph.SetIdsName(object.nickname, jumphole.IdsName.Get())

			jh_archetype := jumphole.Archetype.Get()

			// Check Solar if this is Dockable
			if solar, ok := configs.Solararch.SolarsByNick[jh_archetype]; ok {
				if len(solar.DockingSpheres) == 0 {
					continue
				}
			}

			// Check locked_gate if it is enterable.
			hash_id := flhash.HashNickname(object.nickname)
			if _, ok := configs.InitialWorld.LockedGates[hash_id]; ok {
				continue
			}

			if strings.Contains(jh_archetype, "invisible") {
				continue
			}

			if configs.Discovery != nil {
				is_dockable_by_transports := false
				if solar, ok := configs.Solararch.SolarsByNick[jh_archetype]; ok {
					// strings.Contains(jh_archetype, "_fighter") || // Atmospheric entry points. Dockable only by fighters/freighters
					// included into `IsDockableByCaps` as they don't have capital docking_sphere dockings
					if solar.IsDockableByCaps() {
						is_dockable_by_transports = true
					}
				}

				// Condition is initiallly taken from FLCompanion
				// https://github.com/Corran-Raisu/FLCompanion/blob/021159e3b3a1b40188c93064f1db136780424ea9/Datas.cpp#L585
				// but then rewritted to docking_sphere checks.
				// only with docking_sphere =jump, moor_large we can dock in disco by transports
				if strings.Contains(jh_archetype, "_notransport") { // jumphole_notransport Dockable only by ships with below 650 cargo on board
					// "dsy_hypergate_all" is one directional hypergate dockable by everything, no need to exclude for freighter only paths
					is_dockable_by_transports = false
				}
				if !is_dockable_by_transports && bool(!with_freighter_paths) {
					continue
				}
			}

			for _, existing_object := range system_objects {
				distance := graph.DistanceToTime(DistanceForVecs(object.pos, existing_object.pos),
					system_speed_multiplier,
				) + JumpHoleDelaySec*PrecisionMultipiler
				graph.SetEdge(object.nickname, existing_object.nickname, distance)
				graph.SetEdge(existing_object.nickname, object.nickname, distance)
			}

			jumphole_target_hole := jumphole.GotoHole.Get()
			graph.SetEdge(object.nickname, jumphole_target_hole, 0)
			system_objects = append(system_objects, object)
		}

		for _, tradelane := range system.Tradelanes {
			object := SystemObject{
				nickname: tradelane.Nickname.Get(),
				pos:      tradelane.Pos.Get(),
			}
			graph.SetIstRadelane(object.nickname)

			next_tradelane, next_exists := tradelane.NextRing.GetValue()
			prev_tradelane, prev_exists := tradelane.PrevRing.GetValue()

			if *opts.TradeRoutesDetailedTradeLane {
				// in production every trade lane ring will work as separate entity
				// CONSUMES A LOT OF RAM MEMORY.
				if next_exists {
					if last_tradelane, ok := system.TradelaneByNick[next_tradelane]; ok {
						distance := DistanceForVecs(object.pos, last_tradelane.Pos.Get())
						distance_inside_tradelane := distance * PrecisionMultipiler / float64(average_trade_lane_speed)
						graph.SetEdge(object.nickname, last_tradelane.Nickname.Get(), distance_inside_tradelane)
					}
				}

				if prev_exists {
					if last_tradelane, ok := system.TradelaneByNick[prev_tradelane]; ok {
						distance := DistanceForVecs(object.pos, last_tradelane.Pos.Get())
						distance_inside_tradelane := distance * PrecisionMultipiler / float64(average_trade_lane_speed)
						graph.SetEdge(object.nickname, last_tradelane.Nickname.Get(), distance_inside_tradelane)
					}
				}
			} else {
				// for dev env purposes to speed up test execution, we treat tradelanes as single entity
				// THIS CONSUMES FAR LESS RAM MEMORY. For this reason making it default.
				if next_exists && prev_exists {
					continue
				}

				// next or previous tradelane
				chained_tradelane := ""
				if next_exists {
					chained_tradelane = next_tradelane
				} else {
					chained_tradelane = prev_tradelane
				}
				var last_tradelane *systems_mapped.TradeLaneRing
				// iterate to last in a chain
				for {
					another_tradelane, ok := system.TradelaneByNick[chained_tradelane]
					if !ok {
						break
					}
					last_tradelane = another_tradelane

					if next_exists {
						chained_tradelane, _ = another_tradelane.NextRing.GetValue()
					} else {
						chained_tradelane, _ = another_tradelane.PrevRing.GetValue()
					}
					if chained_tradelane == "" {
						break
					}
				}

				if last_tradelane == nil {
					continue
				}
				distance := DistanceForVecs(object.pos, last_tradelane.Pos.Get())
				distance_inside_tradelane := distance * PrecisionMultipiler / float64(average_trade_lane_speed)
				graph.SetEdge(object.nickname, last_tradelane.Nickname.Get(), distance_inside_tradelane)
			}

			for _, existing_object := range system_objects {
				distance := graph.DistanceToTime(
					DistanceForVecs(object.pos, existing_object.pos),
					system_speed_multiplier,
				) + TradeLaneDockingDelaySec*PrecisionMultipiler
				graph.SetEdge(object.nickname, existing_object.nickname, distance)
				graph.SetEdge(existing_object.nickname, object.nickname, distance)
			}

			system_objects = append(system_objects, object)
		}
	}
	return graph
}

// func (graph *GameGraph) GetDistForTime(time int) float64 {
// 	return float64(time * graph.AvgCruiseSpeed)
// }

func (graph *GameGraph) DistanceToTime(distance float64, system_speed_multiplier float64) cfg.Milliseconds {
	// we assume graph.AvgCruiseSpeed is above zero smth. Not going to check correctness
	// lets try in milliseconds
	return distance * float64(PrecisionMultipiler) / (float64(graph.AvgCruiseSpeed) * system_speed_multiplier)
}

func (graph *GameGraph) GetTimeForDist(dist cfg.Milliseconds) cfg.Seconds {
	// Surprise ;) Distance is time now.
	return dist / PrecisionMultipiler
}

// makes time in ms. Higher int value help having better calcs.
const PrecisionMultipiler = cfg.Milliseconds(1000)
