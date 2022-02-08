package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/machinebox/graphql"
)

func main() {
	fmt.Println("Graphql reserves test")

	var pairs []common.Address
	pairs = []common.Address{
		common.HexToAddress("0x21b8065d10f73ee2e260e5b47d3344d3ced7596e"), common.HexToAddress("0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc"), common.HexToAddress("0x9928e4046d7c6513326ccea028cd3e7a91c7590a"),
	}
	// create a graphql client
	client := graphql.NewClient("http://localhost:8545/graphql")

	// make request to query reserves via storage slots
	req := graphql.NewRequest(`
                query getPairs($pairs: [Address!]) {
                        logs(filter: {addresses: $pairs}) {
                        account {
                                storage(slot: "0x0000000000000000000000000000000000000000000000000000000000000008")
                                }
                        }
                }
        `)

	req.Var("pairs", pairs)

	// set header fields
	req.Header.Set("Cache-Control", "no-cache")

	// define a Context for the request
	ctx := context.Background()

	type AutoGenerated struct {
		Data struct {
			Logs []struct {
				Account struct {
					Storage string `json:"storage"`
					Address string `json:"address"`
				} `json:"account"`
			} `json:"logs"`
		} `json:"data"`
	}

	// raw reserve
	type Reserve2 map[string]map[interface{}]interface{}

	// manual interface
	type Responser map[string]map[string]map[string][]Reserve2

	// run it and capture the response
	var respData map[string]interface{}
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatal(err)
	}

	fmt.Println(respData)

}