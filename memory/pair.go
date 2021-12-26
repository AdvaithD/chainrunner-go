package memory

import (
	"errors"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

type Token common.Address
type Address common.Address

// const addressZero Address = ""

type UniswapV2 struct {
	muPairs         sync.RWMutex
	pairs           map[common.Address]*Pair
	// keyPairs        []pairKey
	pairAddresses   []common.Address
	isDirtyKeyPairs bool
}

func NewUniswapV2() *UniswapV2 {
	return &UniswapV2{pairs: map[common.Address]*Pair{}}
}

// type pairKey struct {
// 	TokenA, TokenB Token
// }

type pairData struct {
	*sync.RWMutex
	reserve0    *big.Int
	reserve1    *big.Int
}

// no clue
type dirty struct {
	isDirty         bool
	isDirtyBalances bool
}

// individual pair
type Pair struct {
	*sync.Mutex
	// tokenAddresses
	token0 common.Address
	token1 common.Address
	// invalid??
	pairData
	muBalance *sync.RWMutex
	*dirty
}

func (pd *pairData) Reserves() (reserve0 *big.Int, reserve1 *big.Int) {
	pd.RLock()
	defer pd.RUnlock()
	return new(big.Int).Set(pd.reserve0), new(big.Int).Set(pd.reserve1)
}

// return a pairs reserves
func (pd *pairData) Revert() pairData {
	return pairData{
		RWMutex:     pd.RWMutex,
		reserve0:    pd.reserve1,
		reserve1:    pd.reserve0,
	}
}

// return each unique pair address as a slice
func (s *UniswapV2) Pairs() ([]common.Address, error) {
	s.muPairs.Lock()
	defer s.muPairs.Unlock()

	return s.pairAddresses, nil
}

func (s *UniswapV2) pair(address common.Address) (*Pair, bool) {
	pair, ok := s.pairs[address]
	if !ok {
		return nil, false
	}
	return &Pair{
		muBalance: pair.muBalance,
		pairData:  pair.pairData.Revert(),
		dirty:     pair.dirty,
	}, true
}

func (s *UniswapV2) Pair(pairAddress common.Address) *Pair {
	s.muPairs.Lock()
	defer s.muPairs.Unlock()
	// key := pairKey{TokenA: coinA, TokenB: coinB}
	pair, _ := s.pair(pairAddress)
	return pair
}

var (
	ErrorIdenticalAddresses = errors.New("IDENTICAL_ADDRESSES")
	ErrorPairExists         = errors.New("PAIR_EXISTS")
)

func (s *UniswapV2) CreatePair(coinA, coinB, pairAddress common.Address, reserveA *big.Int, reserveB *big.Int) (*Pair, error) {
	if coinA == coinB {
		return nil, ErrorIdenticalAddresses
	}

	// pair := s.Pair(coinA, coinB)
	pair := s.Pair(pairAddress)
	if pair != nil {
		return nil, ErrorPairExists
	}

	s.muPairs.Lock()
	defer s.muPairs.Unlock()

	// key := pairKey{coinA, coinB}
	key := pairAddress
	pair = s.addPair(key, pairData{reserve0: reserveA, reserve1: reserveB})
	s.addKeyPair(key)

	return pair, nil
}

func (s *UniswapV2) addPair(pairKey common.Address, data pairData) *Pair {
	data.RWMutex = &sync.RWMutex{}
	pair := &Pair{
		muBalance: &sync.RWMutex{},
		pairData:  data,
		dirty: &dirty{
			isDirty:         false,
			isDirtyBalances: false,
		},
	}
	s.pairs[pairKey] = pair
	return pair
}

func (s *UniswapV2) addKeyPair(key common.Address) {
	s.pairAddresses = append(s.pairAddresses, key)
	s.isDirtyKeyPairs = true
}

// use to update reserves on each block
func (p *Pair) update(newAmount0, newAmount1 *big.Int) {
	p.pairData.Lock()
	defer p.pairData.Unlock()

	p.isDirty = true
	p.pairData.reserve0 = newAmount0
	p.pairData.reserve1 = newAmount1
}

var (
	ErrorInsufficientLiquidityMinted = errors.New("INSUFFICIENT_LIQUIDITY_MINTED")
)
var (
	ErrorInsufficientLiquidityBurned = errors.New("INSUFFICIENT_LIQUIDITY_BURNED")
)

var (
	ErrorK                        = errors.New("K")
	ErrorInsufficientInputAmount  = errors.New("INSUFFICIENT_INPUT_AMOUNT")
	ErrorInsufficientOutputAmount = errors.New("INSUFFICIENT_OUTPUT_AMOUNT")
	ErrorInsufficientLiquidity    = errors.New("INSUFFICIENT_LIQUIDITY")
)

// func (p *Pair) Swap(amount0In, amount1In, amount0Out, amount1Out *big.Int) (amount0, amount1 *big.Int, err error) {
// 	if amount0Out.Sign() != 1 && amount1Out.Sign() != 1 {
// 		return nil, nil, ErrorInsufficientOutputAmount
// 	}

// 	reserve0, reserve1 := p.Reserves()

// 	if amount0Out.Cmp(reserve0) == 1 || amount1Out.Cmp(reserve1) == 1 {
// 		return nil, nil, ErrorInsufficientLiquidity
// 	}

// 	amount0 = new(big.Int).Sub(amount0In, amount0Out)
// 	amount1 = new(big.Int).Sub(amount1In, amount1Out)

// 	if amount0.Sign() != 1 && amount1.Sign() != 1 {
// 		return nil, nil, ErrorInsufficientInputAmount
// 	}

// 	balance0Adjusted := new(big.Int).Sub(new(big.Int).Mul(new(big.Int).Add(amount0, reserve0), big.NewInt(1000)), new(big.Int).Mul(amount0In, big.NewInt(3)))
// 	balance1Adjusted := new(big.Int).Sub(new(big.Int).Mul(new(big.Int).Add(amount1, reserve1), big.NewInt(1000)), new(big.Int).Mul(amount1In, big.NewInt(3)))

// 	if new(big.Int).Mul(balance0Adjusted, balance1Adjusted).Cmp(new(big.Int).Mul(new(big.Int).Mul(reserve0, reserve1), big.NewInt(1000000))) == -1 {
// 		return nil, nil, ErrorK
// 	}

// 	p.update(amount0, amount1)

// 	return amount0, amount1, nil
// }

// func (p *Pair) mint(address Address, value *big.Int) {
// 	p.pairData.Lock()
// 	defer p.pairData.Unlock()

// 	p.muBalance.Lock()
// 	defer p.muBalance.Unlock()

// 	p.isDirtyBalances = true
// 	p.isDirty = true
// 	p.totalSupply.Add(p.totalSupply, value)
// 	balance := p.balances[address]
// 	if balance == nil {
// 		p.balances[address] = big.NewInt(0)
// 	}
// 	p.balances[address].Add(p.balances[address], value)
// }

// func (p *Pair) burn(address Address, value *big.Int) {
// 	p.pairData.Lock()
// 	defer p.pairData.Unlock()
// 	p.muBalance.Lock()
// 	defer p.muBalance.Unlock()

// 	p.isDirtyBalances = true
// 	p.isDirty = true
// 	p.balances[address].Sub(p.balances[address], value)
// 	p.totalSupply.Sub(p.totalSupply, value)
// }



// func (p *Pair) Amounts(liquidity *big.Int) (amount0 *big.Int, amount1 *big.Int) {
// 	p.pairData.RLock()
// 	defer p.pairData.RUnlock()
// 	amount0 = new(big.Int).Div(new(big.Int).Mul(liquidity, p.reserve0), p.totalSupply)
// 	amount1 = new(big.Int).Div(new(big.Int).Mul(liquidity, p.reserve1), p.totalSupply)
// 	return amount0, amount1
// }

// func startingSupply(amount0 *big.Int, amount1 *big.Int) *big.Int {
// 	mul := new(big.Int).Mul(amount0, amount1)
// 	sqrt := new(big.Int).Sqrt(mul)
// 	return new(big.Int).Sub(sqrt, big.NewInt(minimumLiquidity))
// }