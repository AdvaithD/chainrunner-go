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

	logger "github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

var (
	ten       = new(big.Int).SetInt64(10)
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
	logger.Printf("reserves: %T  pairs: %T \n", reserves, pairs)
	logger.Printf("pairinfos: %+v \n", pairInfos)
	// pair name -> address
	tokenToName := make(map[string]common.Address)

	for _, pair := range pairInfos.Data.Pairs {

		fmt.Printf("%+v \n", pair)

		tokenToName[pair.Token0.Symbol] = common.HexToAddress(pair.Token0.Address)
		tokenToName[pair.Token1.Symbol] = common.HexToAddress(pair.Token1.Address)


		one_token0 = new(big.Int).Exp(ten)

	}
}

// type price_quote struct {
// 	TokenIn       string
// 	TokenOut      string
// 	PriceInToOut  *big.Float
// 	PriceNegOfLog *big.Float
// }

// func CreateEdges(reserves [][3]*big.Int,pairs []common.Address) error {
// 	var quotes []price_quote
// }

// func arbExplore(reserves [][3]*big.Int,pairs []common.Address,uniswapInfos *util.UniswapPairs) {
// 	logger.Info("Starting Arb Explore")

// 	// var quotes []price_quote
// 	// loop over each pair

// 	for key, pair := range pairs {
// 		// create edge
// 		token0 := uniswapInfos.Data.Pairs[key].Token0.Address
// 		token1 := uniswapInfos.Data.Pairs[key].Token0.Address

// 		// one_toke0 := new(big.Int).Exp(10, big.NewInt(int64(token0.Decimals)), nil)
// 		// one_toke0 := new(big.Int).Exp(10, big.NewInt(int64(token0.Decimals)), nil)
// 	}

// }
