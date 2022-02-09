package bot

import (
	"chainrunner/internal/util"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/ALTree/bigfloat"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/inconshreveable/log15"
	"github.com/syndtr/goleveldb/leveldb"
	"golang.org/x/sync/errgroup"
)

var (
	ten       = new(big.Int).SetInt64(10)
	zero      = new(big.Int).SetInt64(0)
	neg_one   = new(big.Float).SetFloat64(-1)
	inf       = new(big.Float).SetInf(true)
	repl_mode = flag.Bool(
		"repl_mode", false, "repl mode to inspect DB",
	)
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

type PoolReserve struct {
	reserve0 *big.Int
	reserve1 *big.Int
}

func (p *PoolReserve) IncreaseR0(amount *big.Int) {
	newReserve := new(big.Int)
	newReserve.Add(p.reserve0, amount)
	p.reserve0 = newReserve
}

func (p *PoolReserve) DecreaseR0(amount *big.Int) {
	newReserve := new(big.Int)
	newReserve.Sub(p.reserve0, amount)
	p.reserve0 = newReserve
}

func (p *PoolReserve) IncreaseR1(amount *big.Int) {
	newReserve := new(big.Int)
	newReserve.Add(p.reserve0, amount)
	p.reserve0 = newReserve
}

func (p *PoolReserve) DecreaseR1(amount *big.Int) {
	newReserve := new(big.Int)
	newReserve.Sub(p.reserve0, amount)
	p.reserve0 = newReserve
}

// Adjacency list implementation
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

// TODO: Function that looks at reserves (in terms of tokenAddress) that are changing and if we want that to be starting path
// Maybe need to hardcode WETH and WMATIC trades

// Creates edges given reserves and pairs
func CreateEdges(reserves map[common.Address]*PoolReserve, pairInfos util.UniswapPairs, tokenNameToId map[string]int) *AdjGraph {
	var wg sync.WaitGroup
	var mu = &sync.Mutex{}
	defer util.Duration(util.Track("CreateEdges-300"))
	log.Info("Creating edges")

	graph := GetAdjGraph(len(tokenNameToId))

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
			reserve0 := reserves[common.HexToAddress(pair.Address)].reserve0
			reserve1 := reserves[common.HexToAddress(pair.Address)].reserve1

			token0Decimals, err := strconv.ParseInt(pair.Token0.Decimals, 10, 64)
			if err != nil {
				fmt.Println("comeback")
			}

			token1Decimals, err := strconv.ParseInt(pair.Token1.Decimals, 10, 64)
			if err != nil {
				fmt.Println("comeback")
			}

			one_token0 := new(big.Int).Exp(ten, big.NewInt(token0Decimals), nil)
			one_token1 := new(big.Int).Exp(ten, big.NewInt(token1Decimals), nil)

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
			p0_neg_log.Mul(p0_neg_log, neg_one)

			p1_neg_log := bigfloat.Log(p1)
			p1_neg_log.Mul(p1_neg_log, neg_one)

			// create two quotes (u, v, w) two vertices by names and w is weigth
			token0Id := tokenNameToId[pair.Token0.Symbol]
			token1Id := tokenNameToId[pair.Token1.Symbol]

			mu.Lock()

			graph.addEdge(token0Id, token1Id, p0_neg_log)
			graph.addEdge(token1Id, token0Id, p1_neg_log)

			mu.Unlock()
		}(pair, graph)
	}
	wg.Wait()
	return graph
}

func DFS(g *AdjGraph, source int) {
	defer util.Duration(util.Track("DFS"))
}

type TokenHelper struct {
		// id -> 'token name' mapping
		tokenIdToName map[int]string
		// name -> id
		tokenNameToId map[string]int
		// pair name -> address
		tokenToAddr map[string]common.Address
}

func NewTokenHelper() *TokenHelper {
	tokenHelper := &TokenHelper{
		tokenIdToName: make(map[int]string),
		tokenNameToId: make(map[string]int),
		tokenToAddr: make(map[string]common.Address),
	}
	return tokenHelper
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
		// id -> 'token name' mapping
		tokenIdToName := make(map[int]string)
		// name -> id
		tokenNameToId := make(map[string]int)
		// pair name -> address
		tokenToAddr := make(map[string]common.Address)
		// get top 1000 pairs on uniswapv2

		// Create token helper struct

		tokenHelper := NewTokenHelper()
		
		fmt.Println("tokenheper", tokenHelper)
		addresses, pairInfos := util.GetDemoPairs(b.clients.primary)
		index := 0
		// create unique indexes / id for tokens and populate mappings
		for _, pair := range pairInfos.Data.Pairs {
			// int -> symbol & symbol -> int
			// symbol -> id
			_, ok := tokenNameToId[pair.Token0.Symbol]
			if !ok {
				tokenIdToName[index] = pair.Token0.Symbol
				tokenNameToId[pair.Token0.Symbol] = index
				index++
			}

			// symbol -> id
			_, notexis := tokenNameToId[pair.Token1.Symbol]
			if !notexis {
				tokenIdToName[index] = pair.Token1.Symbol
				tokenNameToId[pair.Token1.Symbol] = index
				index++
			}

			// symbol1 -> addr
			_, exists := tokenToAddr[pair.Token0.Symbol]
			if !exists {
				tokenToAddr[pair.Token0.Symbol] = common.HexToAddress(pair.Token0.Address)
			}

			// symbol2 -> addr
			_, err := tokenToAddr[pair.Token0.Symbol]
			if !err {
				tokenToAddr[pair.Token1.Symbol] = common.HexToAddress(pair.Token1.Address)
			}
		}
		log.Info("Number of Tokens", "tokenNameToId", len(tokenNameToId))

		// loop analytics
		i := float64(0)
		avg := float64(0)

		for {
			// LOOP START
			start := time.Now()

			// get reserves slots
			res, err := b.clients.primary.GetReservesSlots(context.Background(), addresses, nil)

			if err != nil {
				log.Error("ERROR Getting reserves", "error", err)
				panic("exiting")
			}

			log.Info("Reserves", "Got pairs", len(res))

			simulation, err := b.clients.otherwise.SimulateMempool(context.Background(), 5000)

			if err != nil {
				log.Error("error simulating mempool", "error", err)
				panic("exiting")
			}

			log.Info("Simulation", "simulation length", len(simulation))

			reserves := make(map[common.Address]*PoolReserve)

			// TODO: We need to store backrunnable shit here somewhere
			for poolAddress, rawReserve := range res {
				if val, ok := simulation[poolAddress]; ok {
					log.Info("Overriding pool", "address", poolAddress)
					fmt.Println(val)
					for _, simulatedReserves := range val {
						reserves[poolAddress] = &PoolReserve{
							reserve0: simulatedReserves[0].Reserve0,
							reserve1: simulatedReserves[0].Reserve1,
						}
					}
				}
				res0, res1 := util.DeriveReservesFromSlot(rawReserve.String())
				reserves[poolAddress] = &PoolReserve{
					reserve0: res0,
					reserve1: res1,
				}
			}

			// create graph using adjacency list
			grap := CreateEdges(reserves, pairInfos, tokenNameToId)

			// Get node id's for two tokens (WETH and MATIC)
			// Note: These serve as initial starting points for arbs on polygon
			// TODO: Expand initial tokens for arb to others later
			WETH, found := tokenNameToId["WETH"]
			if !found {
				log.Info("unable to find token id", "token", "WETH")
			}
			WMATIC, found := tokenNameToId["MATIC"]
			if !found {
				log.Info("unable to find token id", "token", "MATIC")
			}

			fmt.Println(WETH, WMATIC)
			// Perform DFS on graph starting with WETH. lets see
			// get weth token id
			// params: graph, initial token id, tokenNameToId and tokenIdToName
			DFS(grap, WETH)
			// graph.printGraph()

			// Code that inspects simulation data
			// for address, gasReserveMap := range res {
			// 	for gasPrice, reserves := range gasReserveMap {
			// 		// log.Info("Possible Backrun", "Address", address, "Gas Price", new(big.Float).Quo(new(big.Float).SetUint64(uint64(gasPrice)), gwei), "Reserve", reserves) // "Reserve0", reserve0Float, "Reserve1", reserve1Float)
			// 		log.Info("Possible Backrun", "Gas Price", new(big.Float).Quo(new(big.Float).SetUint64(uint64(gasPrice)), gwei), "Reserve", reserves) // "Reserve0", reserve0Float, "Reserve1", reserve1Float)
			// 	}
			// }

			// log.Info("Got edges", "count", len(edges))
			log.Info("Finished creating latest reserves mapping")

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
