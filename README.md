## Logs

- Able to get 10k uniswap pairs over ipc in ~900ms
- Able to get 100k uniswap pairs over ipc in ~40s

## Research doc (11 Jan 2022)

- Connect to client, uniswapFlashbotsQuery, get reserves, pairs, and pairInfos
- Create util mappings TokenIdToName, TokenNameToId, TokenToAddr
- edges are graph.Edge, and vertices are int

- for each pair in pairinfos we apply them to TokenIdToName TokenNameToId and TokenToAddr mappings

## Snippets

- Code to inspect simulation data (reserves for now)

```
   // Code that inspects simulation data
   for address, gasReserveMap := range res {
    for gasPrice, reserves := range gasReserveMap {
     log.Info("Possible Backrun", "Address", address, "Gas Price", new(big.Float).Quo(new(big.Float).SetUint64(uint64(gasPrice)), gwei), "Reserve", reserves) // "Reserve0", reserve0Float, "Reserve1", reserve1Float)
     log.Info("Possible Backrun", "Gas Price", new(big.Float).Quo(new(big.Float).SetUint64(uint64(gasPrice)), gwei), "Reserve", reserves) // "Reserve0", reserve0Float, "Reserve1", reserve1Float)
    }
   }
```
