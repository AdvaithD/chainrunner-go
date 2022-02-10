package util

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
)

type TokenHelper struct {
	// id -> 'token name' mapping
	TokenIdToName map[int]string
	// name -> id
	TokenNameToId map[string]int
	// pair name -> address
	TokenToAddr map[string]common.Address
}

func NewTokenHelper() *TokenHelper {
	tokenHelper := &TokenHelper{
		TokenIdToName: make(map[int]string),
		TokenNameToId: make(map[string]int),
		TokenToAddr:   make(map[string]common.Address),
	}
	return tokenHelper
}

// given an id, gets symbol
func (t *TokenHelper) GetTokenSymbol(source int) (string, error) {
	if val, ok := t.TokenIdToName[source]; ok {
		return val, nil
	}
	return "", errors.New("token not found")
}

func (t *TokenHelper) GetTokenId(source string) (int, error) {
	if val, ok := t.TokenNameToId[source]; ok {
		return val, nil
	}
	return 0, errors.New("token not found")
}
