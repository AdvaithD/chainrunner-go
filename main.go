package main

import (
	"chainrunner/uniquery"
	"chainrunner/util"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

// get uniswap pairs to bootstrap reserves data
func getUniswapPairs(query *uniquery.FlashBotsUniswapQuery) {
	defer util.Duration(util.Track("getUniswapPairs - 00k"))
	factory := common.HexToAddress("0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f")
	pairs, err := query.GetPairsByIndexRange(&bind.CallOpts{}, factory, big.NewInt(0), big.NewInt(10000))
	if err != nil {
		log.Fatalf("err getting data", err)
	}
	fmt.Println("Got pairs", len(pairs))

	unipairs := make([]common.Address, len(pairs))
	for key := range pairs {
		unipairs = append(unipairs, pairs[key][2])
	}

	fmt.Println("unipairs", unipairs)
	res, err := query.GetReservesByPairs(&bind.CallOpts{}, unipairs)

	if err != nil {
		fmt.Println("err getting reserves", err)
	}

	fmt.Println(res)

}


func main() {
	// init .env into program context
	godotenv.Load(".env")

	// create client
	// rpcClient := services.InitRPCClient()
	conn, err := ethclient.Dial(os.Getenv("GETH_IPC_URL"))

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
