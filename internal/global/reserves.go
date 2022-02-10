package global

import "math/big"

type PoolReserve struct {
	Reserve0 *big.Int
	Reserve1 *big.Int
}

func (p *PoolReserve) IncreaseR0(amount *big.Int) {
	newReserve := new(big.Int)
	newReserve.Add(p.Reserve0, amount)
	p.Reserve0 = newReserve
}

func (p *PoolReserve) DecreaseR0(amount *big.Int) {
	newReserve := new(big.Int)
	newReserve.Sub(p.Reserve0, amount)
	p.Reserve0 = newReserve
}

func (p *PoolReserve) IncreaseR1(amount *big.Int) {
	newReserve := new(big.Int)
	newReserve.Add(p.Reserve1, amount)
	p.Reserve1 = newReserve
}

func (p *PoolReserve) DecreaseR1(amount *big.Int) {
	newReserve := new(big.Int)
	newReserve.Sub(p.Reserve1, amount)
	p.Reserve1 = newReserve
}
