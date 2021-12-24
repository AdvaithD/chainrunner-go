package main

import (
	"log"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)





func main() {
	// init .env into program context
	godotenv.Load(".env")

	// create client
	// rpcClient := services.InitRPCClient()
	_, err := ethclient.Dial(os.Getenv("INFURA_WS_URL"))

	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// address where multicall contract is located
	// uniqueryAddr := common.HexToAddress("0x5EF1009b9FCD4fec3094a5564047e190D72Bd511")
	// uniswap factory address
	// uniFactoryAddr := "0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f"

	// query the contract
	// uniquery, err := uniquery.NewFlashBotsUniswapQuery(uniqueryAddr, conn)

	// if err != nil {
	// 	fmt.Println("error initiating contract to query mass")
	// }

	// getUniswapPairs(uniquery)
}
