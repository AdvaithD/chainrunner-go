package util

type JsonPair []struct {
	Index   int    `json:"index"`
	Address string `json:"address"`
	Token0  struct {
		Address string `json:"address"`
		Symbol  string `json:"symbol"`
		Decimal int    `json:"decimal"`
	} `json:"token0"`
	Token1 struct {
		Address string `json:"address"`
		Symbol  string `json:"symbol"`
		Decimal int    `json:"decimal"`
	} `json:"token1"`
	Reserve0 int64 `json:"reserve0"`
	Reserve1 int64 `json:"reserve1"`
}