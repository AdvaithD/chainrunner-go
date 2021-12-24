package main

import (
	"chainrunner/services"

	"github.com/joho/godotenv"
)

// get uniswap pairs to bootstrap reserves data
func getUniswapPairs() {

}


func main() {
	// init .env into program context
	godotenv.Load(".env")

	// string client = flag.String("client", "xxx", "Gateway to the bsc protocol. Available options:\n\t-bsc_testnet\n\t-bsc\n\t-geth_http\n\t-geth_ipc")

	// create client
	rpcClient := services.InitRPCClient()
}