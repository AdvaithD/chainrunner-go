package graph

import (
	"math/big"

	"chainrunner/internal/global"

	"github.com/ethereum/go-ethereum/common"
)

// Naive calculation over a route given reserves
func EvaluatePath(route []common.Address, reserves map[common.Address]*global.PoolReserve, initialAmount *big.Int) *big.Int {
	amount := initialAmount
	for _, poolAddr := range route {
		res := reserves[poolAddr]
		if res == nil {
			return nil
		}
		// Naive swap: token0 -> token1
		numerator := new(big.Int).Mul(amount, res.Reserve1)
		denominator := new(big.Int).Add(res.Reserve0, amount)
		amount = new(big.Int).Quo(numerator, denominator)
	}
	return amount
}
