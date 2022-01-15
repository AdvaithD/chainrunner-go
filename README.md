## Logs

- Able to get 10k uniswap pairs over ipc in ~900ms
- Able to get 100k uniswap pairs over ipc in ~40s


## Research doc (11 Jan 2022)

- Connect to client, uniswapFlashbotsQuery, get reserves, pairs, and pairInfos
- Create util mappings tokenIdToName, tokenNameToId, tokenToAddr
- edges are graph.Edge, and vertices are int

- for each pair in pairinfos we apply them to tokenIdToName tokenNameToId and tokenToAddr mappings

