package graph

import (
	"chainrunner/internal/global"
	"chainrunner/internal/util"
	"fmt"
	"math/big"
	"strconv"
	"sync"

	"github.com/ALTree/bigfloat"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

type AdjGraph struct {
	vertices int
	adgeList [][]int
	weights  map[int]map[int]*big.Float
}

func GetAdjGraph(vertices int) *AdjGraph {
	var me *AdjGraph = &AdjGraph{}
	me.vertices = vertices
	me.adgeList = make([][]int, vertices)
	me.weights = make(map[int]map[int]*big.Float)
	for i := 0; i < me.vertices; i++ {
		me.adgeList = append(me.adgeList)
	}
	return me
}

func (this *AdjGraph) addEdge(u, v int, w *big.Float) {
	if u < 0 || u >= this.vertices || v < 0 || v >= this.vertices {
		return
	}
	// add node edge
	this.adgeList[u] = append(this.adgeList[u], v)
	// add node weight
	if this.weights[u][v] == nil {
		this.weights[u] = make(map[int]*big.Float)
		this.weights[u][v] = w
	}
}
func (this *AdjGraph) printGraph() {
	fmt.Print("\n Graph Adjacency List ")
	for i := 0; i < this.vertices; i++ {
		fmt.Print(" \n [", i, "] :")
		// iterate edges of i node
		for j := 0; j < len(this.adgeList[i]); j++ {
			fmt.Print("  ", this.adgeList[i][j])
		}
	}
}

// Creates edges given reserves, pairInfo and token helper
// @returns Graph adjacency list
func CreateEdges(reserves map[common.Address]*global.PoolReserve, pairInfos util.UniswapPairs, tokenHelper *util.TokenHelper) *AdjGraph {
	var wg sync.WaitGroup
	var mu = &sync.Mutex{}
	defer util.Duration(util.Track("CreateEdges-300"))
	log.Info("Creating edges")

	graph := GetAdjGraph(len(tokenHelper.TokenNameToId))

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
		}, graphWrapper *AdjGraph) {
			defer wg.Done()
			reserve0 := reserves[common.HexToAddress(pair.Address)].Reserve0
			reserve1 := reserves[common.HexToAddress(pair.Address)].Reserve1

			token0Decimals, err := strconv.ParseInt(pair.Token0.Decimals, 10, 64)
			if err != nil {
				fmt.Println("comeback")
			}

			token1Decimals, err := strconv.ParseInt(pair.Token1.Decimals, 10, 64)
			if err != nil {
				fmt.Println("comeback")
			}

			one_token0 := new(big.Int).Exp(global.Ten, big.NewInt(token0Decimals), nil)
			one_token1 := new(big.Int).Exp(global.Ten, big.NewInt(token1Decimals), nil)

			price_0_to_1, err := util.GetAmountOut(one_token0, reserve0, reserve1)
			if err != nil {
				fmt.Println("comeback")
			}

			price_1_to_0, err := util.GetAmountOut(one_token1, reserve1, reserve0)
			if err != nil {
				fmt.Println("comeback")
			}

			// applying negative log
			p0 := new(big.Float).SetInt(price_0_to_1)
			p0.Quo(p0, new(big.Float).SetInt(one_token1))

			p1 := new(big.Float).SetInt(price_1_to_0)
			p1.Quo(p1, new(big.Float).SetInt(one_token0))

			p0_neg_log := bigfloat.Log(p0)
			p0_neg_log.Mul(p0_neg_log, global.Neg_one)

			p1_neg_log := bigfloat.Log(p1)
			p1_neg_log.Mul(p1_neg_log, global.Neg_one)

			// create two quotes (u, v, w) two vertices by names and w is weigth
			token0Id := tokenHelper.TokenNameToId[pair.Token0.Symbol]
			token1Id := tokenHelper.TokenNameToId[pair.Token1.Symbol]

			mu.Lock()

			graph.addEdge(token0Id, token1Id, p0_neg_log)
			graph.addEdge(token1Id, token0Id, p1_neg_log)

			mu.Unlock()
		}(pair, graph)
	}
	wg.Wait()
	return graph
}
