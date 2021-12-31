package main

import (
	"chainrunner/bindings/uniquery"
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

// 1. Get GraphQL pair data
// 2. Create edges
// 3. Perform graph search algorithm
// 4. log it if possible (dry run)

// TODO: Finish stack code to trace a negative cycle
// func TraceNegativeCycle(pre map[string]string, string v) ([]string) {
//         for !Stack.contains(v) {
//                 Stack.push(v)
//                 v = pre[v]
//         }

//         cycle := make([]string)
//         cycle = append(cycle, v)

//         for Stack.top() != v {
//                 cycle = append(Stack.pop())
//         }
//         cycle = append(cycle, v)

//         return cycle
// }


// PSUEDOCODE TO DETECT CYCLE
// 
// bool dfs(int u) {
//         vis[u] = true;
//         on_stk[u] = true;
//         for (int v : adj[u]) {
//           if ((!vis[v] && dfs(v)) || on_stk[v])
//             return true;
//         on_stk[u] = false;
//         return false;
//       }
//
// PSUEDOCODE TO DETECT CYCLE

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

        type pairData struct {
                Address common.Address
                Token0  struct {
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

        reserves, pairs, pairInfos := getUniswapPairs(uniquery)

        // we measure timefrom here (post data collection)

        logger.Printf("reserves: %v  pairs: %v \n", len(reserves), len(pairs))

        // id -> 'token name' mapping
        tokenIdMap := make(map[int]string)

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

                quotes = append(quotes, firstQuote, secondQuote)
        }

        fmt.Printf("[Create Edges]: Took %v to create edges for %v pairs \n", time.Since(now), len(pairs))
        fmt.Printf("[EDGE] Edge Count: %v, nodes: %v tokenToName: %v\n", len(quotes), nodes, len(tokenToName))
        fmt.Printf("[EDGE] Quotescount: %v, edgesFromTo: %v \n", len(quotes), len(edgesFromTo))


        index := 0
        loop := time.Now()
        for key := range tokenToName {
                tokenIdMap[index] = key
                index++
        }


        fmt.Println("tokenIdMap", tokenIdMap)
	// data, _ := json.MarshalIndent(edgesFromTo["WETH"], "", " ")

	// fmt.Println("DATA", string(data))
       
}