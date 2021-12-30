package main

import (
	"chainrunner/bindings/uniquery"
	"chainrunner/internal/mainnet"
	"chainrunner/internal/memory"
	"chainrunner/internal/util"
	"container/list"
	"fmt"
	"log"
	"math/big"
	"os"

	"strconv"
	"time"

	"github.com/ALTree/bigfloat"
	logger "github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

var (
	ten       = new(big.Int).SetInt64(10)
	zero       = new(big.Int).SetInt64(0)
	neg_one   = new(big.Float).SetFloat64(-1)
)
// Struct for id -> token (or) id -> pair address
type Arber struct {
	tokens   map[uint]common.Address
	pairs    map[uint]common.Address
	pairInfo map[common.Address]*memory.Pair
}

// get uniswap pairs to bootstrap reserves data
func getUniswapPairs(query *uniquery.FlashBotsUniswapQuery) ([][3]*big.Int, []common.Address, util.UniswapPairs) {
	defer util.Duration(util.Track("getUniswapPairsAndReserves-5000"))

	pairInfos, err:= util.GetUniswapPairs()

	if err != nil {
		fmt.Println("err getting graphql pairdata")
	}

	// code to get n pairs from index 0 to n (not being used in favour of graphql rn)
	// factory := common.HexToAddress("0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f")
	// pairs, err := query.GetPairsByIndexRange(&bind.CallOpts{}, factory, big.NewInt(0), big.NewInt(5000))
	// if err != nil {
	// 	log.Fatalf("err getting data", err)
	// }
	// logger.Printf("Got %v pairs \n", len(pairs))

	pairAddresses := make([]common.Address, 0)

	for _, pair := range pairInfos.Data.Pairs {
		pairAddresses = append(pairAddresses, common.HexToAddress(pair.Address))
	}

	reserves, err := query.GetReservesByPairs(&bind.CallOpts{Context: nil}, pairAddresses)

	if err != nil {
		logger.Error("err getting reserves", err)
	}

	// fmt.Printf("%V \n", res)
	return reserves, pairAddresses, pairInfos
}

// TODO: calculate amount out
func GetAmountOut(amountIn *big.Int, reserve0 *big.Int, reserve1 *big.Int) (*big.Int, error) {
	amountInWithFee := amountIn.Mul(amountIn, new(big.Int).SetInt64(997))
	var numerator = new(big.Int)
	var denominator = new(big.Int)
	var amountOut = new(big.Int)

	numerator = numerator.Mul(amountIn, new(big.Int).SetInt64(997))
	denominator = denominator.Add(reserve0.Mul(reserve0, new(big.Int).SetInt64(1000)), amountInWithFee)

	amountOut = numerator.Div(numerator, denominator)

	return amountOut, nil
}

// is this an edge?
type price_quote struct {
	TokenIn       string
	TokenOut      string
	PriceInToOut  *big.Float
	PriceNegOfLog *big.Float
}

// 1. Get GraphQL pair data
// 2. Create edges
// 3. Perform graph search algorithm
// 4. log it if possible (dry run)

func main() {
	// init .env into program context
	godotenv.Load(".env")
	
	conn, err := ethclient.Dial(os.Getenv("INFURA_WS_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// query the contract
	uniquery, err := uniquery.NewFlashBotsUniswapQuery(mainnet.UNIQUERY_ADDR, conn)
	if err != nil {
		fmt.Println("error initiating contract to query mass")
	}

	// arbExplore(reserves, pairs, pairInfos)
	type pairData struct {
		Address common.Address 
		Token0 struct {
			Decimals uint8
			Address  *common.Address
			Symbol   *string
		}
		Token1 struct {
			Decimals *uint8
			Address  *common.Address 
			Symbol   *string
		}
	} 

	// pairAddress -> pairInfo
	// var pairInfoMappingmap map[common.Address]pairData
	// [reserv0, reserve1, blockTimestampLast]

	reserves, pairs, pairInfos := getUniswapPairs(uniquery)

	logger.Printf("reserves: %v  pairs: %v \n", len(reserves), len(pairs))

	// pair name -> address
	tokenToName := make(map[string]common.Address)

	// price quotes for each pair 0 -> 1 and 1 -> 0 included
	var quotes []price_quote
	// loop over pairs
	now := time.Now()

    logger.Printf("tokenToName has %v\n", len(tokenToName))
    logger.Printf("pairs user were %v\n", len(pairs))

    nodes := make([]string, len(tokenToName))    
    edgesFromTo := make(map[string][]price_quote)

	// for each pair, create edges for all the pairs that we have
	for key, pair := range pairInfos.Data.Pairs {
		tokenToName[pair.Token0.Symbol] = common.HexToAddress(pair.Token0.Address)
		tokenToName[pair.Token1.Symbol] = common.HexToAddress(pair.Token1.Address)

		reserve0 := reserves[key][0]
		reserve1 := reserves[key][1]

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

		price_0_to_1, err := GetAmountOut(one_token0, reserve0, reserve1)

		if err != nil {
			fmt.Println("comeback")		
		}

		price_1_to_0, err := GetAmountOut(one_token1, reserve1, reserve0)


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
		firstQuote := price_quote{
			TokenIn: pair.Token0.Symbol, TokenOut: pair.Token1.Symbol,
			PriceInToOut: p0, PriceNegOfLog: p0_neg_log,
		}

		secondQuote := price_quote{
			TokenIn: pair.Token1.Symbol, TokenOut: pair.Token0.Symbol,
			PriceInToOut: p1, PriceNegOfLog: p1_neg_log,
		}

		// edges from a node mapping store
		edgesFromTo[firstQuote.TokenIn] = append(edgesFromTo[firstQuote.TokenIn], firstQuote)
		edgesFromTo[secondQuote.TokenIn] = append(edgesFromTo[secondQuote.TokenIn], secondQuote)

		quotes = append(quotes, firstQuote)
		quotes = append(quotes, secondQuote)
	}
    fmt.Printf("[Create Edges]: Took %v to create edges for %v pairs \n", time.Since(now), len(pairs))
    fmt.Printf("[EDGE] Edge Count: %v, nodes: %v tokenToName: %v\n", len(quotes), nodes, len(tokenToName))
    fmt.Printf("[EDGE] Quotescount: %v, edgesFromTo: %v \n", len(quotes), len(edgesFromTo))



    // length (in amount of edges) of current shortest path from the source to u
    length := make(map[string]int64)

    // distance is the weight of the current shortest path from source to u 
    distances := make(map[string]*big.Float)

    // Notation:
    // weight is price, u and v are tokenin and tokenout

    // FIFO Queue
    queue := list.New()

    // queue is not empty condition


    // SFPA - START
    // for each vertex, set initial distances to 0
    for token := range tokenToName {
	    length[token] = 0
	    distances[token] = new(big.Float).SetInt(zero)
        queue.PushBack(token)
    }

    for queue.Len() > 0 {
        u := queue.Front()

        // fmt.Printf("u, %+v %T \n", u, u.Value)
        queue.Remove(u)
        // now, loop  over each edge (u,v) in Edges of the graph

	for _, v := range edgesFromTo[u.Value.(string)] {
          //fmt.Printf("key, %+v\n", v)
	//   fmt.Printf("u: %+v \n", u.Value)
	//   fmt.Printf("v: %+v \n", v)
	  // if sum of (distance of u, weight w(u, v)) is less than distance[v]

	  if (distances[u.Value.(string)].Add(distances[u.Value.(string)], v.PriceNegOfLog)).Cmp(distances[v.TokenOut]) < 0 {
		length[v.TokenOut] = length[u.Value.(string)] + 1
		if length[v.TokenOut] < 0 {
			logger.Warn("Negative cycle!")
		}

		distances[v.TokenOut] = distances[u.Value.(string)].Add(distances[u.Value.(string)], v.PriceNegOfLog)
		//TODO: if Queue not containts v push it to queue

	  }
	}
    }
}
