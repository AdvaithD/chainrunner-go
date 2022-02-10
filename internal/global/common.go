package global

import "math/big"

var (
	Ten     = new(big.Int).SetInt64(10)
	Zero    = new(big.Int).SetInt64(0)
	Neg_one = new(big.Float).SetFloat64(-1)
	Inf     = new(big.Float).SetInf(true)
)
