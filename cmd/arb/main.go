package main

import (
	"chainrunner/bindings/uniquery"
	"chainrunner/memory"
	"chainrunner/util"
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

// get uniswap pairs to bootstrap reserves data
func getUniswapPairs(query *uniquery.FlashBotsUniswapQuery) ([][3]*big.Int, [][3]common.Address) {
	// re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	// numberOfPairs := 5000
	defer util.Duration(util.Track("getUniswapPairsAndReserves-5000"))
	factory := common.HexToAddress("0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f")


	// get pairs
	pairs, err := query.GetPairsByIndexRange(&bind.CallOpts{}, factory, big.NewInt(0), big.NewInt(5000))
	if err != nil {
		log.Fatalf("err getting data", err)
	}
	logger.Printf("Got %v pairs \n", len(pairs))

	unipairs := make([]common.Address, 0)

	for _, pair := range pairs {
		unipairs = append(unipairs, pair[2])
	}

	// fmt.Println("unipairs", "%v", unipairs)

	// var finalPairs []common.Address
	res, err := query.GetReservesByPairs(&bind.CallOpts{Context: nil}, unipairs)

	if err != nil {
		logger.Error("err getting reserves", err)
	}

	// fmt.Printf("%V \n", res)
	return res, pairs
}


func main() {
	// init .env into program context
	godotenv.Load(".env")

	database := memory.NewUniswapV2()
	// create client
	// rpcClient := services.InitRPCClient()
	
	conn, err := ethclient.Dial(os.Getenv("INFURA_WS_URL"))

	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// address where multicall contract is located
	uniqueryAddr := common.HexToAddress("0x5EF1009b9FCD4fec3094a5564047e190D72Bd511")

	// query the contract
	uniquery, err := uniquery.NewFlashBotsUniswapQuery(uniqueryAddr, conn)
	if err != nil {
		fmt.Println("error initiating contract to query mass")
	}

	// [reserv0, reserve1, blockTimestampLast]
	reserves, pairs := getUniswapPairs(uniquery)

	logger.Printf("reserves: %v  pairs: %v \n", len(reserves), len(pairs))
	logger.Printf("reserves: %T  pairs: %T \n", reserves, pairs)

	// loop over pairs
	for index, pair := range pairs {
		// logger.Printf("%v | %v | %T", index, , val)

		database.CreatePair(pair[0], pair[1], pair[2], reserves[index][0], reserves[index][1])
	}
	logger.Info("Finished writing to db")
}
