package main

import (
	log "github.com/sirupsen/logrus"

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
		log.Info("Failed to connect to the Ethereum client: %v", err)
	}

	log.Info("Working..")
}
