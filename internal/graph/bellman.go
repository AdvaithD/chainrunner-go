package graph

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Graph represents a graph consisting of edges and vertices
type Graph struct {
	edges    []*Edge
	vertices []int
	tokenIdToName map[int]string //  0 -> eth, 1 -> wbtc
	tokenNameToAddress map[string]common.Address // eth -> 0xabc, wbtc -> 0xbtc
}

// get token name given an id
func (g *Graph) GetTokenName(id int) string{
	return g.tokenIdToName[id]
}

// get token address given symbol
func (g *Graph) GetTokenAddr(name string) common.Address {
	return g.tokenNameToAddress[name]
}

var infinity = new(big.Float).SetInf(true)

// Edge represents a weighted line between two nodes
type Edge struct {
	From, To int
	Weight   *big.Float
}

// NewEdge returns a pointer to a new Edge
func NewEdge(from, to int, weight *big.Float) *Edge {
	return &Edge{From: from, To: to, Weight: weight}
}

// NewGraph returns a graph consisting of given edges and vertices (vertices must count from 0 upwards)
func NewGraph(edges []*Edge, vertices []int) *Graph {
	return &Graph{edges: edges, vertices: vertices}
}

// FindArbitrageLoop returns either an arbitrage loop or a nil map
func (g *Graph) FindArbitrageLoop(source int) []int {
	predecessors, distances := g.BellmanFord(source)
	return g.FindNegativeWeightCycle(predecessors, distances, source)
}

// BellmanFord determines the shortest path and returns the predecessors and distances
func (g *Graph) BellmanFord(source int) ([]int, []*big.Float) {
	size := len(g.vertices)
	distances := make([]*big.Float, size)
	predecessors := make([]int, size)
	for _, v := range g.vertices {
		distances[v] = infinity
	}
	distances[source] = new(big.Float).SetInt64(0)

	for i, changes := 0, 0; i < size-1; i, changes = i+1, 0 {
		for _, edge := range g.edges {
			var tempDist *big.Float
			if tempDist := tempDist.Add(distances[edge.From], edge.Weight); tempDist.Cmp(distances[edge.To]) == -1 {
				distances[edge.To] = tempDist
				predecessors[edge.To] = edge.From
				changes++
			}
		}
		if changes == 0 {
			break
		}
	}
	return predecessors, distances
}

// FindNegativeWeightCycle finds a negative weight cycle from predecessors and a source
func (g *Graph) FindNegativeWeightCycle(predecessors []int, distances []*big.Float, source int) []int {
	for _, edge := range g.edges {
		var tempBigFloat *big.Float
		if tempBigFloat.Add(distances[edge.From], edge.Weight).Cmp(distances[edge.To]) == -1 {
			return arbitrageLoop(predecessors, source)
		}
	}
	return nil
}

func arbitrageLoop(predecessors []int, source int) []int {
	size := len(predecessors)
	loop := make([]int, size)
	loop[0] = source

	exists := make([]bool, size)
	exists[source] = true

	indices := make([]int, size)

	var index, next int
	for index, next = 1, source; ; index++ {
		next = predecessors[next]
		loop[index] = next
		if exists[next] {
			return loop[indices[next] : index+1]
		}
		indices[next] = index
		exists[next] = true
	}
}