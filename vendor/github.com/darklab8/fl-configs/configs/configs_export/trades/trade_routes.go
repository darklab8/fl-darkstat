package trades

import (
	"math"
	"strings"

	"github.com/darklab8/fl-configs/configs/cfgtype"
	"github.com/darklab8/fl-configs/configs/configs_mapped"
	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/initialworld/flhash"
	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped/systems_mapped"
	"github.com/darklab8/fl-configs/configs/configs_settings"
)

type SystemObject struct {
	nickname string
	pos      cfgtype.Vector
}

func DistanceForVecs(Pos1 cfgtype.Vector, Pos2 cfgtype.Vector) float64 {
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

const (
	// already accounted for
	AvgTradeLaneSpeed = 2250

	// Add for every pair of jumphole in path
	JumpHoleDelaySec = 15 // and jump gate
	// add for every tradelane vertex pair in path
	TradeLaneDockingDelaySec = 10
	// add just once
	BaseDockingDelay = 10
)

type ExtraBase struct {
	Pos      cfgtype.Vector
	Nickname string
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
func MapConfigsToFGraph(
	configs *configs_mapped.MappedConfigs,
	avgCruiseSpeed int,
	with_freighter_paths WithFreighterPaths,
	extra_bases_by_system map[string][]ExtraBase,
) *GameGraph {
	graph := NewGameGraph(avgCruiseSpeed, with_freighter_paths)
	for _, system := range configs.Systems.Systems {

		var system_objects []SystemObject = make([]SystemObject, 0, 50)

		if bases, ok := extra_bases_by_system[system.Nickname]; ok {
			for _, base := range bases {
				object := SystemObject{
					nickname: base.Nickname,
					pos:      base.Pos,
				}
				graph.SetIdsName(object.nickname, 0)

				for _, existing_object := range system_objects {
					distance := DistanceForVecs(object.pos, existing_object.pos) + graph.GetDistForTime(BaseDockingDelay)
					graph.SetEdge(object.nickname, existing_object.nickname, distance)
					graph.SetEdge(existing_object.nickname, object.nickname, distance)
				}

				graph.AllowedVertixesForCalcs[VertexName(object.nickname)] = true

				system_objects = append(system_objects, object)
			}
		}

		for _, system_obj := range system.Bases {
			object := SystemObject{
				nickname: system_obj.Base.Get(),
				pos:      system_obj.Pos.Get(),
			}
			graph.SetIdsName(object.nickname, system_obj.IdsName.Get())

			if system_obj.Archetype.Get() == systems_mapped.BaseArchetypeInvisible {
				continue
			}

			goods, goods_defined := configs.Market.GoodsPerBase[object.nickname]
			if !goods_defined {
				continue
			}

			if len(goods.MarketGoods) == 0 {
				continue
			}

			for _, existing_object := range system_objects {
				distance := DistanceForVecs(object.pos, existing_object.pos) + graph.GetDistForTime(BaseDockingDelay)
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
				_, dockable := solar.DockingSphere.GetValue()

				if !dockable {
					continue
				}
			}

			// Check locked_gate if it is enterable.
			hash_id := flhash.HashNickname(object.nickname)
			if _, ok := configs.InitialWorld.LockedGates[hash_id]; ok {
				continue
			}

			// Condition is taken from FLCompanion
			// https://github.com/Corran-Raisu/FLCompanion/blob/021159e3b3a1b40188c93064f1db136780424ea9/Datas.cpp#L585
			// Check Aingar Fork for Disco version if necessary.
			if strings.Contains(jh_archetype, "_fighter") ||
				strings.Contains(jh_archetype, "_notransport") ||
				jh_archetype == "dsy_comsat_planetdock" ||
				jh_archetype == "dsy_hypergate_all" {
				if !with_freighter_paths {
					continue
				}
			}

			for _, existing_object := range system_objects {
				distance := DistanceForVecs(object.pos, existing_object.pos) + graph.GetDistForTime(JumpHoleDelaySec)
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

			if configs_settings.Env.IsDevEnv {
				// for dev env purposes to speed up test execution, we treat tradelanes as single entity
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
				distance_inside_tradelane := distance * float64(graph.AvgCruiseSpeed) / float64(AvgTradeLaneSpeed)
				graph.SetEdge(object.nickname, last_tradelane.Nickname.Get(), distance_inside_tradelane)
			} else {
				// in production every trade lane ring will work as separate entity
				if next_exists {
					if last_tradelane, ok := system.TradelaneByNick[next_tradelane]; ok {
						distance := DistanceForVecs(object.pos, last_tradelane.Pos.Get())
						distance_inside_tradelane := distance * float64(graph.AvgCruiseSpeed) / float64(AvgTradeLaneSpeed)
						graph.SetEdge(object.nickname, last_tradelane.Nickname.Get(), distance_inside_tradelane)
					}
				}

				if prev_exists {
					if last_tradelane, ok := system.TradelaneByNick[prev_tradelane]; ok {
						distance := DistanceForVecs(object.pos, last_tradelane.Pos.Get())
						distance_inside_tradelane := distance * float64(graph.AvgCruiseSpeed) / float64(AvgTradeLaneSpeed)
						graph.SetEdge(object.nickname, last_tradelane.Nickname.Get(), distance_inside_tradelane)
					}
				}
			}

			for _, existing_object := range system_objects {
				distance := DistanceForVecs(object.pos, existing_object.pos) + graph.GetDistForTime(TradeLaneDockingDelaySec)
				graph.SetEdge(object.nickname, existing_object.nickname, distance)
				graph.SetEdge(existing_object.nickname, object.nickname, distance)
			}

			system_objects = append(system_objects, object)
		}
	}
	return graph
}

func (graph *GameGraph) GetDistForTime(time int) float64 {
	return float64(time * graph.AvgCruiseSpeed)
}

func (graph *GameGraph) GetTimeForDist(dist float64) int {
	return int(dist / float64(graph.AvgCruiseSpeed))
}
