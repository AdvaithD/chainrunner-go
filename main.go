package main

import (
	"chainrunner/testpackage"
	"fmt"
)



func main() {
	// godotenv.Load(".env")

	// string client = flag.String("client", "xxx", "Gateway to the bsc protocol. Available options:\n\t-bsc_testnet\n\t-bsc\n\t-geth_http\n\t-geth_ipc")
	fmt.Println(testpackage.NumberOfTransactions)

}