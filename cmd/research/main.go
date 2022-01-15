package main

import (
	"chainrunner/bindings/uniquery"
	"chainrunner/internal/graph"
	"chainrunner/internal/mainnet"
	"chainrunner/internal/memory"
	"chainrunner/internal/util"
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
        ten     = new(big.Int).SetInt64(10)
        zero    = new(big.Int).SetInt64(0)
        neg_one = new(big.Float).SetFloat64(-1)
        inf       = new(big.Float).SetInf(true)
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

        pairInfos, err := util.GetUniswapPairs()

        if err != nil {
                fmt.Println("err getting graphql pairdata")
        }
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

// Definition of an edge in uniswapv2 terms
type price_quote struct {
        TokenIn       string
        TokenOut      string
        PriceInToOut  *big.Float
        PriceNegOfLog *big.Float
}

// helper to create an array with incremental range
func makeRange(min, max int) []int {
        a := make([]int, max-min+1)
        for i := range a {
            a[i] = min + i
        }
        return a
    }

// Get reserves, pairs, pairInfos
// create tokenIdToName mapping
// create tokenNameToId mapping
// tokenToAddr mapping
// edges variable to store all edges
// verticed to store all vertices
// for each pair 
// check if toeknNamToId exists

func main() {
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

        reserves, pairs, pairInfos := getUniswapPairs(uniquery)

        // we measure timefrom here (post data collection)

        logger.Printf("reserves: %v  pairs: %v \n", len(reserves), len(pairs))

        // id -> 'token name' mapping
        tokenIdToName := make(map[int]string)
        // name -> id
        tokenNameToId := make(map[string]int)
        // pair name -> address
        tokenToAddr := make(map[string]common.Address)

        // bellman edges
        var edges []*graph.Edge

        // bellman vertices
        var vertices []int
        // loop over pairs
        now := time.Now()

        logger.Printf("tokenToName has %v\n", len(tokenToAddr))
        logger.Printf("pairs user were %v\n", len(pairs))

        // create necessary token mappings (id to symbol, symbol to addr)
        utiltime := time.Now()
        // id counter
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

        // vertices start from 0, 1,2, 3,....
        vertices = makeRange(0, len(tokenIdToName)-1)
        fmt.Println("util mapping creation time: ", time.Since(utiltime))
        // for each pair, create edges for all the pairs that we have
        for key, pair := range pairInfos.Data.Pairs {
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
        }

        fmt.Printf("[Create Edges]: Took %v to create edges for %v pairs \n", time.Since(now), len(pairs))
        fmt.Printf("[EDGE] Edge Count: %v, vertices: %v tokenToName: %v\n", len(edges), len(vertices), len(tokenToAddr))
        fmt.Printf("[EDGE] tokenIdToName: %v, tokenNameToId: %v, tokenToName: %v\n", len(tokenIdToName), len(tokenNameToId), len(tokenToAddr))
        fmt.Printf("[EDGE] TokenNameToId: %+v \n", tokenNameToId)

        // PRINT VALUES FOR MAPPINGS
        // fmt.Println("tokenIdToName: ", tjokenIdToName)
        // fmt.Println("tokenNameToId: ", tokenNameToId)
        // fmt.Println("tokenToName: ", tokenToName)

        inputTokens := []string{"WETH", "USDT", "WBTC", "USDC"}

        for _, token := range inputTokens {
                arber := graph.NewGraph(edges, vertices, tokenIdToName, tokenToAddr, tokenNameToId)
                fmt.Println("token routes for: ", token)
                tokenId := arber.GetTokenId(token)
                fmt.Println("id: ", tokenId)

                loop := arber.FindArbitrageLoop(tokenId)

                for _, key := range loop {
                        fmt.Printf("%v -> ", arber.GetTokenName(key))
                }

                fmt.Printf("\n\n %v loop: %v \n\n", token, loop)
        }
}
