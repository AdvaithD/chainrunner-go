package global

import "github.com/ethereum/go-ethereum/common"

var ROUTERS = map[string]common.Address{
	"quickswap": common.HexToAddress("0xa5E0829CaCEd8fFDD4De3c43696c57F7D7A678ff"),
	"sushiswap": common.HexToAddress("0x1b02dA8Cb0d097eB8D57A175b88c7D8b47997506"),
}
