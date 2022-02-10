package bot

import (
	"chainrunner/internal/global"
	"chainrunner/internal/graph"
	"chainrunner/internal/util"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/inconshreveable/log15"
	"github.com/syndtr/goleveldb/leveldb"
	"golang.org/x/sync/errgroup"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
)

const (
	ESTIMATE_GAS_FAIL    = "estimate_gas_fail.json"
	FOUND_ARBS           = "found_arbs.json"
	REASONS_NOT_WORTH_IT = "not_worth.json"
)

type clients struct {
	primary   *ethclient.Client
	uniswap   *ethclient.Client
	sushiswap *ethclient.Client
	otherwise *ethclient.Client
}

type ReportMessage struct {
	Error    error
	Role     string
	When     time.Time
	ArgsUsed interface{}
	Extra    interface{}
}

type report_logs struct {
	found_arb         chan *ReportMessage
	estimate_fail_log chan *ReportMessage
	not_worth_it_log  chan *ReportMessage
}

type Bot struct {
	db                  *leveldb.DB
	clients             clients
	shutdown            chan struct{}
	log_update_incoming report_logs
}

func NewBot(db_path, client_path string) (*Bot, error) {
	client, err := ethclient.Dial(client_path)

	if err != nil {
		return nil, err
	}

	primaryClient, err := ethclient.Dial(client_path)

	if err != nil {
		return nil, err
	}

	uniswapClient, err := ethclient.Dial(client_path)

	if err != nil {
		return nil, err
	}

	sushiSwapClient, err := ethclient.Dial(client_path)

	if err != nil {
		return nil, err
	}

	db, err := leveldb.RecoverFile(db_path, nil)

	if err != nil {
		return nil, err
	}

	return &Bot{
		clients: clients{
			primary:   primaryClient,
			uniswap:   uniswapClient,
			sushiswap: sushiSwapClient,
			otherwise: client,
		},
		log_update_incoming: report_logs{
			make(chan *ReportMessage),
			make(chan *ReportMessage),
			make(chan *ReportMessage),
		},
		shutdown: make(chan struct{}),
		db:       db,
	}, nil
}

func (b *Bot) CloseResources() error {
	b.clients.primary.Close()
	b.clients.otherwise.Close()
	b.clients.uniswap.Close()
	b.clients.sushiswap.Close()

	return b.db.Close()
}

func (b *Bot) KickoffFailureLogs(file_used string, failures chan *ReportMessage) {
	f, err := os.OpenFile(file_used, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		log.Debug("come back to this")
		return
	}

	for {
		select {
		case <-b.shutdown:
			f.Close()
			return
		case msg := <-failures:
			s, err := json.MarshalIndent(msg, "", "\t")
			if err != nil {
				log.Error("come back to this")
				return
			}

			if _, err := f.Write(s); err != nil {
				log.Error("some error on writing to "+file_used, err)
				f.Close()
				return
			}

			f.WriteString("\n")
		}
	}
}

// Adjacency list implementation

// TODO: Function that looks at reserves (in terms of tokenAddress) that are changing and if we want that to be starting path
// Maybe need to hardcode WETH and WMATIC trades

func FindArb(pairInfos util.UniswapPairs, tokenHelper *util.TokenHelper) {}

// helper to create an array with incremental range
func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

// created directed graph (used to find simple cycles)
func BuildDirectedGraph(reserves map[common.Address]*global.PoolReserve, pairInfos util.UniswapPairs, tokenHelper *util.TokenHelper) *simple.DirectedGraph {
	defer util.Duration(util.Track("CREATE GONUM EDGES"))
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

// Run the bot
func (b *Bot) Run() (e error) {
	var g errgroup.Group
	// uncomment while inspecting simulation stuff
	// gwei, _, _ := big.ParseFloat("1e9", 10, 0, big.ToNearestEven)

	go b.KickoffFailureLogs(ESTIMATE_GAS_FAIL, b.log_update_incoming.estimate_fail_log)
	go b.KickoffFailureLogs(REASONS_NOT_WORTH_IT, b.log_update_incoming.not_worth_it_log)

	g.Go(func() error {
		interrupt := make(chan os.Signal)
		defer signal.Stop(interrupt)
		defer close(interrupt)
		signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-interrupt
		close(b.shutdown)
		time.Sleep(time.Second * 2)
		return nil
	})

	// start here
	g.Go(func() error {
		// Create token helper struct
		tokenHelper := util.NewTokenHelper()
		addresses, pairInfos := util.GetDemoPairs(b.clients.primary)
		index := 0

		// create unique indexes / id for tokens and populate tokenHelper struct
		// (populates tokenHelper)
		for _, pair := range pairInfos.Data.Pairs {
			// int -> symbol & symbol -> int
			// symbol -> id
			_, ok := tokenHelper.TokenNameToId[pair.Token0.Symbol]
			if !ok {
				tokenHelper.TokenIdToName[index] = pair.Token0.Symbol
				tokenHelper.TokenNameToId[pair.Token0.Symbol] = index
				index++
			}

			// symbol -> id
			_, notexis := tokenHelper.TokenNameToId[pair.Token1.Symbol]
			if !notexis {
				tokenHelper.TokenIdToName[index] = pair.Token1.Symbol
				tokenHelper.TokenNameToId[pair.Token1.Symbol] = index
				index++
			}

			// symbol1 -> addr
			_, exists := tokenHelper.TokenToAddr[pair.Token0.Symbol]
			if !exists {
				tokenHelper.TokenToAddr[pair.Token0.Symbol] = common.HexToAddress(pair.Token0.Address)
			}

			// symbol2 -> addr
			_, err := tokenHelper.TokenToAddr[pair.Token0.Symbol]
			if !err {
				tokenHelper.TokenToAddr[pair.Token1.Symbol] = common.HexToAddress(pair.Token1.Address)
			}
		}

		log.Info("Number of Tokens", "TokenNameToId", len(tokenHelper.TokenNameToId))

		// loop analytics
		i := float64(0)
		avg := float64(0)

		for {
			// LOOP START
			start := time.Now()
			// get reserves slots

			// 1 - Get reserves
			res, err := b.clients.primary.GetReservesSlots(context.Background(), addresses, nil)

			if err != nil {
				log.Error("ERROR Getting reserves", "error", err)
				panic("exiting")
			}

			log.Info("Reserves", "Got pairs", len(res))

			// 2 - Get simulation
			simulation, err := b.clients.otherwise.SimulateMempool(context.Background(), 5000)

			if err != nil {
				log.Error("error simulating mempool", "error", err)
				panic("exiting")
			}

			log.Info("Simulation", "simulation length", len(simulation))

			// create reserves mapping (address -> reserve)
			reserves := make(map[common.Address]*global.PoolReserve)

			// Below we 'process' the simulation
			// TODO: We need to store backrunnable shit here somewhere

			// for each  pool address in the simulated data
			for poolAddress, rawReserve := range res {
				// if pairs + simulation are mergable (i.e: simulation on one of the pairs the bot supports)
				if val, ok := simulation[poolAddress]; ok {
					// log.Info("Overriding pool", "address", poolAddress)
					fmt.Println(val)
					for _, simulatedReserves := range val {
						reserves[poolAddress] = &global.PoolReserve{
							Reserve0: simulatedReserves[0].Reserve0,
							Reserve1: simulatedReserves[0].Reserve1,
						}
					}
				}
				res0, res1 := util.DeriveReservesFromSlot(rawReserve.String())
				reserves[poolAddress] = &global.PoolReserve{
					Reserve0: res0,
					Reserve1: res1,
				}
			}

			// instead of calling CreateEdges() I'm trying the github analysis repo's method
			// create graph using adjacency list
			//  IterateAndGetPaths(reserves, pairInfos, tokenHelper, "WETH", "WETH", 5, []string, [][]string)

			// testing gonum stuff
			graph := graph.BuildDirectedGraph(reserves, pairInfos, tokenHelper)

			// manually getting token id's
			WETH, found := tokenHelper.TokenNameToId["WETH"]
			if !found {
				log.Info("unable to find token id", "token", "WETH")
			}
			MATIC, found := tokenHelper.TokenNameToId["MATIC"]
			if !found {
				log.Info("unable to find token id", "token", "MATIC")
			}

			var tokens []int64
			tokens = append(tokens, int64(WETH), int64(MATIC))
			// TODO: Fix me P0

			preCycles := time.Now()
			cycles := topo.DirectedCyclesOfMaxLenContainingAnyOf(graph, 5, tokens)

			postCycles := time.Since(preCycles)

			log.Info("Time to find cycles", "time", postCycles, "cycles", len(cycles))

			// visited := make(map[int]bool)
			// Perform DFS on graph starting with WETH. lets see
			// get weth token id
			// params: graph, initial token id, visited, tokenHelper

			// DFS(grap, WETH, visited, tokenHelper, true)

			// log.Info("Got edges", "count", len(edges))
			log.Info("Finished loop")

			// LOOP END
			elapsed := time.Since(start)
			log.Info("Time Elapsed", "iteration", elapsed)
			i++
			avg = avg*(i-1)/i + elapsed.Seconds()*1/i
			log.Info("Time Elapsed", "average", avg)
			log.Info("---------------------------------------")
		}

	})

	return g.Wait()
}
