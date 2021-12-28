package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type UniswapPairs struct {
	Data struct {
		Pairs []struct {
			Address     string `json:"id"`
			Token0 struct {
				Decimals string `json:"decimals"`
				Address       string `json:"id"`
				Symbol   string `json:"symbol"`
			} `json:"token0"`
			Token1 struct {
				Decimals string `json:"decimals"`
				Address       string `json:"id"`
				Symbol   string `json:"symbol"`
			} `json:"token1"`
		} `json:"pairs"`
	} `json:"data"`
}

func GetUniswapPairs() (UniswapPairs, error) {
    jsonData := map[string]string{
        "query": `
        {
          pairs(first: 500, skip: 0, where: {volumeUSD_gt: "10000000"}, orderBy: reserveUSD, orderDirection: desc) {
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