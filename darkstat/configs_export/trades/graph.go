package trades

/*
Game graph simplifies for us conversion of data from Freelancer space simulator to different graph algorithms.
*/

import (
	"errors"
	"math"
	"reflect"

	"github.com/darklab8/fl-darkstat/configs/cfg"
)

type VertexName string

type Intg = int32

const intgmax = Intg(math.MaxInt32)

type GameGraph struct {
	matrix                    map[VertexName]map[VertexName]float64
	IndexByNick               map[VertexName]Intg `json:"index_by_nickname" validate:"required"`
	NicknameByIndex           map[Intg]VertexName `json:"nickname_by_index" validate:"required"`
	AllowedVertixesForCalcs   map[VertexName]bool // Consider deleting this
	AvgCruiseSpeed            Intg
	idsNamesByNick            map[VertexName]int
	IsTradelane               map[VertexName]bool
	CanVisitFreightersOnlyJHs WithFreighterPaths
}

func (g *GameGraph) WipeMatrix() {
	g.matrix = nil
}

func NewGameGraph(avgCruiseSpeed int, canVisitFreighterOnlyJHs WithFreighterPaths) *GameGraph {
	return &GameGraph{
		matrix:                    make(map[VertexName]map[VertexName]float64),
		IndexByNick:               map[VertexName]Intg{},
		NicknameByIndex:           make(map[Intg]VertexName),
		AllowedVertixesForCalcs:   make(map[VertexName]bool),
		AvgCruiseSpeed:            Intg(avgCruiseSpeed),
		idsNamesByNick:            make(map[VertexName]int),
		IsTradelane:               make(map[VertexName]bool),
		CanVisitFreightersOnlyJHs: canVisitFreighterOnlyJHs,
	}
}

func (graph *GameGraph) GetVertexIndexByName(key VertexName) Intg {

	value, ok := graph.IndexByNick[key]

	if !ok {
		value = Intg(len(graph.IndexByNick))
		graph.IndexByNick[key] = value
		graph.NicknameByIndex[value] = key
	}

	return value
}

func (f *GameGraph) SetEdge(keya string, keyb string, distance float64) {
	vertex, vertex_exists := f.matrix[VertexName(keya)]
	if !vertex_exists {
		vertex = make(map[VertexName]float64)
		f.matrix[VertexName(keya)] = vertex
	}

	if _, vert_target_exists := f.matrix[VertexName(keyb)]; !vert_target_exists {
		f.matrix[VertexName(keyb)] = make(map[VertexName]float64)
	}

	_, already_set := vertex[VertexName(keyb)]
	if already_set {
		return // otherwise u will overwrite tradelane distances.
	}

	vertex[VertexName(keyb)] = distance
}

func (f *GameGraph) SetIdsName(keya string, ids_name int) {
	f.idsNamesByNick[VertexName(keya)] = ids_name
}

func (f *GameGraph) SetIstRadelane(keya string) {
	f.IsTradelane[VertexName(keya)] = true
}

func GetTimeMs(f *GameGraph, dist [][]Intg, keya string, keyb string) (cfg.MillisecondsI, error) {
	sourse_index, source_found := f.IndexByNick[VertexName(keya)]
	target_index, target_found := f.IndexByNick[VertexName(keyb)]
	_ = source_found
	if !source_found && !target_found {
		return INF, errors.New("both source and destination are not found")
	}
	if !source_found {
		return INF, errors.New("source is not found")
	}
	if !target_found {
		return INF, errors.New("destination is not found")
	}
	return dist[sourse_index][target_index], nil
}
func GetTimeMs2(f *GameGraph, dist [][]Intg, keya string, keyb string) cfg.MillisecondsI {
	result, _ := GetTimeMs(f, dist, keya, keyb)
	return result
}

type Path struct {
	Node     Intg
	NextNode Intg
	Dist     Intg
}

func GetPath(graph *GameGraph, parents [][]Parent, dist [][]Intg, source_key string, target_key string) []Path {
	// fmt.Println("get_path", source_key, target_key)
	S := []Path{}
	u, found_u := graph.IndexByNick[VertexName(target_key)] // target
	if !found_u {
		return []Path{}
	}
	_ = found_u
	source := graph.IndexByNick[VertexName(source_key)]

	distance_skipped_buffer := Intg(0)

	add_node := func(parent Parent) {
		path_to_add := Path{
			Node: u,
		}
		if len(S) > 0 {
			path_to_add.NextNode = S[len(S)-1].Node
		} else {
			path_to_add.NextNode = NO_PARENT
		}
		if path_to_add.Node != NO_PARENT && path_to_add.NextNode != NO_PARENT {
			path_to_add.Dist = parent.weight + distance_skipped_buffer //  dist[path_to_add.Node][path_to_add.NextNode]
			distance_skipped_buffer = 0
		}

		S = append(S, path_to_add)
	}
	add_node(Parent{node: u})

	if parents[source][u].node != NO_PARENT || u == source {
		for {
			parent := parents[source][u]
			u = parent.node

			nickname := graph.NicknameByIndex[u]
			if _, ok := graph.IsTradelane[nickname]; ok {
				distance_skipped_buffer += parent.weight
				continue
			}

			add_node(parent)

			if u == NO_PARENT {
				break
			}
		}
	}
	ReverseSlice(S)
	return S
}

// panic if s is not a slice
func ReverseSlice(s interface{}) {
	size := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, size-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

type DetailedPath struct {
	PrevName    string
	NextName    string
	PrevIdsName int
	NextIdsName int
	PrevNode    Intg
	NextNode    Intg
	Dist        Intg
	TimeMinutes int
	TimeSeconds int
}

func (graph *GameGraph) GetPaths(parents [][]Parent, dist [][]Intg, source_key string, target_key string) []DetailedPath {
	var detailed_paths []DetailedPath

	paths := GetPath(graph, parents, dist, source_key, target_key)
	for _, path := range paths {
		minutes := int(math.Floor(graph.GetTimeForDist(float64(path.Dist)) / 60))
		detailed_path := DetailedPath{
			PrevName:    string(graph.NicknameByIndex[path.Node]),
			NextName:    string(graph.NicknameByIndex[path.NextNode]),
			PrevIdsName: graph.idsNamesByNick[graph.NicknameByIndex[path.Node]],
			NextIdsName: graph.idsNamesByNick[graph.NicknameByIndex[path.NextNode]],
			PrevNode:    path.Node,
			NextNode:    path.NextNode,
			Dist:        path.Dist,
			TimeMinutes: int(minutes),
			TimeSeconds: int(int(graph.GetTimeForDist(float64(path.Dist))) - minutes*60),
		}

		detailed_paths = append(detailed_paths, detailed_path)
	}

	return detailed_paths
}
