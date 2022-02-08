package bot

import (
	"chainrunner/internal/graph"
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

// Creates edges given reserves and pairs
func CreateEdges(reserves map[common.Address]*PoolReserve, pairInfos util.UniswapPairs, tokenNameToId map[string]int) []*graph.Edge {
	var wg sync.WaitGroup
	defer util.Duration(util.Track("CreateEdges-1000"))

	var edges []*graph.Edge
	log.Info("Creating edges")

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
		}) {
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

			// create two quotes
			firstEdge := graph.NewEdge(tokenNameToId[pair.Token0.Symbol], tokenNameToId[pair.Token1.Symbol], p0_neg_log)
			secondEdge := graph.NewEdge(tokenNameToId[pair.Token1.Symbol], tokenNameToId[pair.Token0.Symbol], p1_neg_log)

			edges = append(edges, firstEdge, secondEdge)
			wg.Done()
		}(pair)
	}
	wg.Wait()
	return edges
}

// Run the bot
func (b *Bot) Run() (e error) {
	var g errgroup.Group

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

			reserves := make(map[common.Address]*PoolReserve)

			for poolAddress, rawReserve := range res {
				res0, res1 := util.DeriveReservesFromSlot(rawReserve.String())
				reserves[poolAddress] = &PoolReserve{
					reserve0: res0,
					reserve1: res1,
				}
			}

			edges := CreateEdges(reserves, pairInfos, tokenNameToId)

			log.Info("Got edges", "count", len(edges))
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
