package trades

import (
	"math"
	"strings"

	"github.com/darklab8/fl-configs/configs/configs_mapped"
	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped/systems_mapped"
	"github.com/darklab8/fl-configs/configs/conftypes"
)

type SystemObject struct {
	nickname string
	pos      conftypes.Vector
}

func DistanceForVecs(Pos1 conftypes.Vector, Pos2 conftypes.Vector) float64 {
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

const (
	// already accounted for
	AvgTransportCruiseSpeed = 350
	AvgFreighterCruiseSpeed = 500
	// already accounted for
	AvgTradeLaneSpeed = 1900

	// Add for every pair of jumphole in path
	JumpHoleDelaySec = 15 // and jump gate
	// add for every tradelane vertex pair in path
	TradeLaneDockingDelaySec = 10
	// add just once
	BaseDockingDelay = 20
)

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
func MapConfigsToFGraph(configs *configs_mapped.MappedConfigs, avgCruiseSpeed int, with_freighter_paths WithFreighterPaths) *GameGraph {
	graph := NewGameGraph(avgCruiseSpeed)
	for _, system := range configs.Systems.Systems {

		var system_objects []SystemObject = make([]SystemObject, 0, 50)

		for _, system_obj := range system.Bases {
			object := SystemObject{
				nickname: system_obj.Base.Get(),
				pos:      system_obj.Pos.Get(),
			}
			graph.SetIdsName(object.nickname, system_obj.IdsName.Get())

			if strings.Contains(object.nickname, "proxy_") {
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

			// if strings.Contains(object.nickname, strings.ToLower("Rh02_to_Iw02_hole")) {
			// 	fmt.Println()
			// }

			jh_archetype := jumphole.Archetype.Get()

			// TODO Check Solar if this is Dockable
			if jh_archetype == "jumphole_noentry" { // hardcoded for now
				continue
			}
			// TODO Check locked_gate if it is enterable.

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

			next_tradelane, next_exists := tradelane.NextRing.GetValue()
			prev_tradelane, prev_exists := tradelane.PrevRing.GetValue()
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
