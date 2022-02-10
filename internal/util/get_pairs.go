package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type UniswapPairs struct {
	Data struct {
		Pairs []struct {
			Address string `json:"id"`
			Token0  struct {
				Decimals string `json:"decimals"`
				Address  string `json:"id"`
				Symbol   string `json:"symbol"`
			} `json:"token0"`
			Token1 struct {
				Decimals string `json:"decimals"`
				Address  string `json:"id"`
				Symbol   string `json:"symbol"`
			} `json:"token1"`
		} `json:"pairs"`
	} `json:"data"`
}

// get uniswap pairs to bootstrap reserves data
func GetDemoPairs(client *ethclient.Client) ([]common.Address, UniswapPairs) {
	pairInfos, err := GetUniswapPairs()

	if err != nil {
		fmt.Println("err getting graphql pairdata")
	}

	pairAddresses := make([]common.Address, 0)

	for _, pair := range pairInfos.Data.Pairs {
		pairAddresses = append(pairAddresses, common.HexToAddress(pair.Address))
	}

	// fmt.Printf("%V \n", res)
	return pairAddresses, pairInfos
}

func GetUniswapPairs() (UniswapPairs, error) {
	jsonData := map[string]string{
		"query": `
        {
          pairs(first: 1000, skip: 0, where: {volumeUSD_gt: "10000000"}, orderBy: reserveUSD, orderDirection: desc) {
            id
            token0 {
              id
              decimals
              symbol
            }
            token1 {
              id
              decimals
              symbol
            }
        }
      }`,
	}
	jsonValue, _ := json.Marshal(jsonData)
	request, err := http.NewRequest("POST", "https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v2", bytes.NewBuffer(jsonValue))

	if err != nil {
		log.Fatalf("error requesting graphql data", err)
	}
	client := &http.Client{Timeout: time.Second * 10}
	response, err := client.Do(request)
	defer response.Body.Close()
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	data, _ := ioutil.ReadAll(response.Body)
	// fmt.Println(string(data))

	var pairs UniswapPairs

	json.Unmarshal(data, &pairs)

	fmt.Printf("Got %v pairs\n", len(pairs.Data.Pairs))

	return pairs, nil
}

// returns 1000 uniswap pair addresses
func Get1000PairAddresses() []common.Address {
	pairInfos, err := GetUniswapPairs()

	if err != nil {
		fmt.Println("err getting graphql pairdata")
	}

	pairAddresses := make([]common.Address, 0)

	for _, pair := range pairInfos.Data.Pairs {
		pairAddresses = append(pairAddresses, common.HexToAddress(pair.Address))
	}

	return pairAddresses
}

// calculate amount out given amount in, reserve0 and reserve1
func GetAmountOut(amountIn *big.Int, reserve0 *big.Int, reserve1 *big.Int) (*big.Int, error) {
	amountInWithFee := amountIn.Mul(amountIn, new(big.Int).SetInt64(997))

	var numerator = new(big.Int)
	var denominator = new(big.Int)
	var amountOut = new(big.Int)

	numerator = numerator.Mul(amountInWithFee, reserve1)
	denominator = denominator.Add(reserve0.Mul(reserve0, new(big.Int).SetInt64(1000)), amountInWithFee)

	amountOut = numerator.Div(numerator, denominator)

	return amountOut, nil
}

func GetAmountsOut(amountIn *big.Int) {}
