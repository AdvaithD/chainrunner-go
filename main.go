package main

import (
	"chainrunner/uniquery"
	"chainrunner/util"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

// get uniswap pairs to bootstrap reserves data
func getUniswapPairs(query *uniquery.FlashBotsUniswapQuery) {
	defer util.Duration(util.Track("getUniswapPairs - 100k"))
	factory := common.HexToAddress("0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f")
	pairs, err := query.GetPairsByIndexRange(&bind.CallOpts{}, factory, big.NewInt(0), big.NewInt(100000))
	if err != nil {
		log.Fatalf("err getting data")
	}
	fmt.Println("Got pairs", len(pairs), pairs)
}


func main() {
	// init .env into program context
	godotenv.Load(".env")

	// string client = flag.String("client", "xxx", "Gateway to the bsc protocol. Available options:\n\t-bsc_testnet\n\t-bsc\n\t-geth_http\n\t-geth_ipc")

	// create client
	// rpcClient := services.InitRPCClient()
	conn, err := ethclient.Dial("ws://157.90.35.22:8545")

	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// address where multicall contract is located
	uniqueryAddr := common.HexToAddress("0x5EF1009b9FCD4fec3094a5564047e190D72Bd511")
	// uniswap factory address
	// uniFactoryAddr := "0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f"

	// query the contract
	uniquery, err := uniquery.NewFlashBotsUniswapQuery(uniqueryAddr, conn)

	if err != nil {
		fmt.Println("error initiating contract to query mass")
	}

	getUniswapPairs(uniquery)
}