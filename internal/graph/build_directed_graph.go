package graph

import (
	"chainrunner/internal/util"
	"sync"

	"gonum.org/v1/gonum/graph/simple"
)

// created directed graph (used to find simple cycles)
func BuildDirectedGraph(pairInfos util.UniswapPairs, tokenHelper *util.TokenHelper) *simple.DirectedGraph {
	defer util.Duration(util.Track("BuildDirectedGraph"))
	var wg sync.WaitGroup
	var mu = &sync.Mutex{}
	graph := simple.NewDirectedGraph()

	// create the edges first
	for key := range tokenHelper.TokenIdToName {
		// log.Info("createfonum", "key", key, "value", value)
		if graph.Node(int64(key)) == nil {
			graph.AddNode(simple.Node(key))
		}
	}

	for _, pair := range pairInfos.Data.Pairs {
		wg.Add(1)
		go func(pair struct {
			Address string `json:"id"`
			Token0  struct {
				Decimals string `json:"decimals"`
				Address  string `json:"id"`
				Symbol   string `json:"symbol"`
			} `json:"token0"`
			Token1 struct {
				Decimals string `json:"decimals"`
				Address  string `json:"id"`
				Symbol   string `json:"symbol"`
			} `json:"token1"`
		}, graph *simple.DirectedGraph) {
			defer wg.Done()

			token0Id := tokenHelper.TokenNameToId[pair.Token0.Symbol]
			token1Id := tokenHelper.TokenNameToId[pair.Token1.Symbol]

			mu.Lock()
			defer mu.Unlock()
			graph.SetEdge(simple.Edge{F: simple.Node(int64(token0Id)), T: simple.Node(int64(token1Id))})
			graph.SetEdge(simple.Edge{F: simple.Node(int64(token1Id)), T: simple.Node(int64(token0Id))})

		}(pair, graph)
	}
	wg.Wait()

	return graph
}
