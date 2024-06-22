package trades

/*
Game graph simplifies for us conversion of data from Freelancer space simulator to different graph algorithms.
*/

import "math"

const INF = math.MaxFloat32

type VertexName string

type GameGraph struct {
	matrix                        map[VertexName]map[VertexName]float64
	index_by_nickname             map[VertexName]int
	vertex_to_calculate_paths_for map[VertexName]bool
}

func NewGameGraph() *GameGraph {
	return &GameGraph{
		matrix:                        make(map[VertexName]map[VertexName]float64),
		index_by_nickname:             map[VertexName]int{},
		vertex_to_calculate_paths_for: make(map[VertexName]bool),
	}
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
	vertex[VertexName(keyb)] = distance
}

func GetDist[T any](f *GameGraph, dist [][]T, keya string, keyb string) T {
	return dist[f.index_by_nickname[VertexName(keya)]][f.index_by_nickname[VertexName(keyb)]]
}
