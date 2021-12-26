package util

import (
	"chainrunner/bindings/erc20"
	"chainrunner/util/decimal"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var ten = new(big.Int).SetInt64(10)


func PrettyToken(
	tk common.Address,
	value *big.Int,
	cl *ethclient.Client,
) (string, decimal.Dec) {

	token, err := erc20.NewErc20(tk, cl)
	if err != nil {
		return "", decimal.NewDec(int64(0))
	}

	dec, err := token.Decimals(nil)
	if err != nil {
		return "", decimal.NewDec(int64(0))
	}

	symbol, err := token.Symbol(nil)
	if err != nil {
		return "", decimal.NewDec(int64(0))
	}

	one_token := decimal.NewDecFromBigInt(new(big.Int).Exp(
		ten,
		big.NewInt(int64(dec)),
		nil,
	))

	return symbol, decimal.NewDecFromBigInt(
		value,
	).Quo(one_token)
}