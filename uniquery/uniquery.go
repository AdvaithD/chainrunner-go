// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package uniquery

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// FlashBotsUniswapQueryABI is the input ABI used to generate the binding from.
const FlashBotsUniswapQueryABI = "[{\"inputs\":[{\"internalType\":\"contractUniswapV2Factory\",\"name\":\"_uniswapFactory\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_stop\",\"type\":\"uint256\"}],\"name\":\"getPairsByIndexRange\",\"outputs\":[{\"internalType\":\"address[3][]\",\"name\":\"\",\"type\":\"address[3][]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIUniswapV2Pair[]\",\"name\":\"_pairs\",\"type\":\"address[]\"}],\"name\":\"getReservesByPairs\",\"outputs\":[{\"internalType\":\"uint256[3][]\",\"name\":\"\",\"type\":\"uint256[3][]\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// FlashBotsUniswapQueryFuncSigs maps the 4-byte function signature to its string representation.
var FlashBotsUniswapQueryFuncSigs = map[string]string{
	"ab2217e4": "getPairsByIndexRange(address,uint256,uint256)",
	"4dbf0f39": "getReservesByPairs(address[])",
}

// FlashBotsUniswapQueryBin is the compiled bytecode used for deploying new contracts.
var FlashBotsUniswapQueryBin = "0x608060405234801561001057600080fd5b5061092e806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80634dbf0f391461003b578063ab2217e414610064575b600080fd5b61004e610049366004610654565b610084565b60405161005b91906107f7565b60405180910390f35b6100776100723660046106e6565b610252565b60405161005b9190610784565b606060008267ffffffffffffffff8111156100a1576100a16108ca565b6040519080825280602002602001820160405280156100da57816020015b6100c76105f6565b8152602001906001900390816100bf5790505b50905060005b8381101561024a578484828181106100fa576100fa6108b4565b905060200201602081019061010f91906106c9565b6001600160a01b0316630902f1ac6040518163ffffffff1660e01b815260040160606040518083038186803b15801561014757600080fd5b505afa15801561015b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061017f919061071b565b826001600160701b03169250816001600160701b031691508063ffffffff1690508484815181106101b2576101b26108b4565b60200260200101516000600381106101cc576101cc6108b4565b602002018585815181106101e2576101e26108b4565b60200260200101516001600381106101fc576101fc6108b4565b60200201868681518110610212576102126108b4565b602002602001015160026003811061022c5761022c6108b4565b6020020192909252919052528061024281610883565b9150506100e0565b509392505050565b60606000846001600160a01b031663574f2ba36040518163ffffffff1660e01b815260040160206040518083038186803b15801561028f57600080fd5b505afa1580156102a3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102c7919061076b565b9050808311156102d5578092505b838310156103295760405162461bcd60e51b815260206004820181905260248201527f73746172742063616e6e6f7420626520686967686572207468616e2073746f70604482015260640160405180910390fd5b6000610335858561086c565b905060008167ffffffffffffffff811115610352576103526108ca565b60405190808252806020026020018201604052801561038b57816020015b6103786105f6565b8152602001906001900390816103705790505b50905060005b828110156105eb5760006001600160a01b038916631e3dd18b6103b4848b610854565b6040518263ffffffff1660e01b81526004016103d291815260200190565b60206040518083038186803b1580156103ea57600080fd5b505afa1580156103fe573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104229190610630565b9050806001600160a01b0316630dfe16816040518163ffffffff1660e01b815260040160206040518083038186803b15801561045d57600080fd5b505afa158015610471573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104959190610630565b8383815181106104a7576104a76108b4565b60200260200101516000600381106104c1576104c16108b4565b60200201906001600160a01b031690816001600160a01b031681525050806001600160a01b031663d21220a76040518163ffffffff1660e01b815260040160206040518083038186803b15801561051757600080fd5b505afa15801561052b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061054f9190610630565b838381518110610561576105616108b4565b602002602001015160016003811061057b5761057b6108b4565b60200201906001600160a01b031690816001600160a01b031681525050808383815181106105ab576105ab6108b4565b60200260200101516002600381106105c5576105c56108b4565b6001600160a01b03909216602092909202015250806105e381610883565b915050610391565b509695505050505050565b60405180606001604052806003906020820280368337509192915050565b80516001600160701b038116811461062b57600080fd5b919050565b60006020828403121561064257600080fd5b815161064d816108e0565b9392505050565b6000806020838503121561066757600080fd5b823567ffffffffffffffff8082111561067f57600080fd5b818501915085601f83011261069357600080fd5b8135818111156106a257600080fd5b8660208260051b85010111156106b757600080fd5b60209290920196919550909350505050565b6000602082840312156106db57600080fd5b813561064d816108e0565b6000806000606084860312156106fb57600080fd5b8335610706816108e0565b95602085013595506040909401359392505050565b60008060006060848603121561073057600080fd5b61073984610614565b925061074760208501610614565b9150604084015163ffffffff8116811461076057600080fd5b809150509250925092565b60006020828403121561077d57600080fd5b5051919050565b602080825282518282018190526000919084820190604085019084805b828110156107ea57845184835b60038110156107d45782516001600160a01b0316825291880191908801906001016107ae565b50505093850193606093909301926001016107a1565b5091979650505050505050565b602080825282518282018190526000919084820190604085019084805b828110156107ea57845184835b600381101561083e57825182529188019190880190600101610821565b5050509385019360609390930192600101610814565b600082198211156108675761086761089e565b500190565b60008282101561087e5761087e61089e565b500390565b60006000198214156108975761089761089e565b5060010190565b634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052604160045260246000fd5b6001600160a01b03811681146108f557600080fd5b5056fea26469706673582212205f75e3c05c4fc88ca24506b22a6628b12311bbebe4994beb183bc3aa6510638f64736f6c63430008060033"

// DeployFlashBotsUniswapQuery deploys a new Ethereum contract, binding an instance of FlashBotsUniswapQuery to it.
func DeployFlashBotsUniswapQuery(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *FlashBotsUniswapQuery, error) {
	parsed, err := abi.JSON(strings.NewReader(FlashBotsUniswapQueryABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(FlashBotsUniswapQueryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &FlashBotsUniswapQuery{FlashBotsUniswapQueryCaller: FlashBotsUniswapQueryCaller{contract: contract}, FlashBotsUniswapQueryTransactor: FlashBotsUniswapQueryTransactor{contract: contract}, FlashBotsUniswapQueryFilterer: FlashBotsUniswapQueryFilterer{contract: contract}}, nil
}

// FlashBotsUniswapQuery is an auto generated Go binding around an Ethereum contract.
type FlashBotsUniswapQuery struct {
	FlashBotsUniswapQueryCaller     // Read-only binding to the contract
	FlashBotsUniswapQueryTransactor // Write-only binding to the contract
	FlashBotsUniswapQueryFilterer   // Log filterer for contract events
}

// FlashBotsUniswapQueryCaller is an auto generated read-only Go binding around an Ethereum contract.
type FlashBotsUniswapQueryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FlashBotsUniswapQueryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FlashBotsUniswapQueryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FlashBotsUniswapQueryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FlashBotsUniswapQueryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FlashBotsUniswapQuerySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FlashBotsUniswapQuerySession struct {
	Contract     *FlashBotsUniswapQuery // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// FlashBotsUniswapQueryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FlashBotsUniswapQueryCallerSession struct {
	Contract *FlashBotsUniswapQueryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// FlashBotsUniswapQueryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FlashBotsUniswapQueryTransactorSession struct {
	Contract     *FlashBotsUniswapQueryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// FlashBotsUniswapQueryRaw is an auto generated low-level Go binding around an Ethereum contract.
type FlashBotsUniswapQueryRaw struct {
	Contract *FlashBotsUniswapQuery // Generic contract binding to access the raw methods on
}

// FlashBotsUniswapQueryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FlashBotsUniswapQueryCallerRaw struct {
	Contract *FlashBotsUniswapQueryCaller // Generic read-only contract binding to access the raw methods on
}

// FlashBotsUniswapQueryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FlashBotsUniswapQueryTransactorRaw struct {
	Contract *FlashBotsUniswapQueryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFlashBotsUniswapQuery creates a new instance of FlashBotsUniswapQuery, bound to a specific deployed contract.
func NewFlashBotsUniswapQuery(address common.Address, backend bind.ContractBackend) (*FlashBotsUniswapQuery, error) {
	contract, err := bindFlashBotsUniswapQuery(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FlashBotsUniswapQuery{FlashBotsUniswapQueryCaller: FlashBotsUniswapQueryCaller{contract: contract}, FlashBotsUniswapQueryTransactor: FlashBotsUniswapQueryTransactor{contract: contract}, FlashBotsUniswapQueryFilterer: FlashBotsUniswapQueryFilterer{contract: contract}}, nil
}

// NewFlashBotsUniswapQueryCaller creates a new read-only instance of FlashBotsUniswapQuery, bound to a specific deployed contract.
func NewFlashBotsUniswapQueryCaller(address common.Address, caller bind.ContractCaller) (*FlashBotsUniswapQueryCaller, error) {
	contract, err := bindFlashBotsUniswapQuery(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FlashBotsUniswapQueryCaller{contract: contract}, nil
}

// NewFlashBotsUniswapQueryTransactor creates a new write-only instance of FlashBotsUniswapQuery, bound to a specific deployed contract.
func NewFlashBotsUniswapQueryTransactor(address common.Address, transactor bind.ContractTransactor) (*FlashBotsUniswapQueryTransactor, error) {
	contract, err := bindFlashBotsUniswapQuery(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FlashBotsUniswapQueryTransactor{contract: contract}, nil
}

// NewFlashBotsUniswapQueryFilterer creates a new log filterer instance of FlashBotsUniswapQuery, bound to a specific deployed contract.
func NewFlashBotsUniswapQueryFilterer(address common.Address, filterer bind.ContractFilterer) (*FlashBotsUniswapQueryFilterer, error) {
	contract, err := bindFlashBotsUniswapQuery(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FlashBotsUniswapQueryFilterer{contract: contract}, nil
}

// bindFlashBotsUniswapQuery binds a generic wrapper to an already deployed contract.
func bindFlashBotsUniswapQuery(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(FlashBotsUniswapQueryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FlashBotsUniswapQuery *FlashBotsUniswapQueryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FlashBotsUniswapQuery.Contract.FlashBotsUniswapQueryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FlashBotsUniswapQuery *FlashBotsUniswapQueryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FlashBotsUniswapQuery.Contract.FlashBotsUniswapQueryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FlashBotsUniswapQuery *FlashBotsUniswapQueryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FlashBotsUniswapQuery.Contract.FlashBotsUniswapQueryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FlashBotsUniswapQuery *FlashBotsUniswapQueryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FlashBotsUniswapQuery.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FlashBotsUniswapQuery *FlashBotsUniswapQueryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FlashBotsUniswapQuery.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FlashBotsUniswapQuery *FlashBotsUniswapQueryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FlashBotsUniswapQuery.Contract.contract.Transact(opts, method, params...)
}

// GetPairsByIndexRange is a free data retrieval call binding the contract method 0xab2217e4.
//
// Solidity: function getPairsByIndexRange(address _uniswapFactory, uint256 _start, uint256 _stop) view returns(address[3][])
func (_FlashBotsUniswapQuery *FlashBotsUniswapQueryCaller) GetPairsByIndexRange(opts *bind.CallOpts, _uniswapFactory common.Address, _start *big.Int, _stop *big.Int) ([][3]common.Address, error) {
	var out []interface{}
	err := _FlashBotsUniswapQuery.contract.Call(opts, &out, "getPairsByIndexRange", _uniswapFactory, _start, _stop)

	if err != nil {
		return *new([][3]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([][3]common.Address)).(*[][3]common.Address)

	return out0, err

}

// GetPairsByIndexRange is a free data retrieval call binding the contract method 0xab2217e4.
//
// Solidity: function getPairsByIndexRange(address _uniswapFactory, uint256 _start, uint256 _stop) view returns(address[3][])
func (_FlashBotsUniswapQuery *FlashBotsUniswapQuerySession) GetPairsByIndexRange(_uniswapFactory common.Address, _start *big.Int, _stop *big.Int) ([][3]common.Address, error) {
	return _FlashBotsUniswapQuery.Contract.GetPairsByIndexRange(&_FlashBotsUniswapQuery.CallOpts, _uniswapFactory, _start, _stop)
}

// GetPairsByIndexRange is a free data retrieval call binding the contract method 0xab2217e4.
//
// Solidity: function getPairsByIndexRange(address _uniswapFactory, uint256 _start, uint256 _stop) view returns(address[3][])
func (_FlashBotsUniswapQuery *FlashBotsUniswapQueryCallerSession) GetPairsByIndexRange(_uniswapFactory common.Address, _start *big.Int, _stop *big.Int) ([][3]common.Address, error) {
	return _FlashBotsUniswapQuery.Contract.GetPairsByIndexRange(&_FlashBotsUniswapQuery.CallOpts, _uniswapFactory, _start, _stop)
}

// GetReservesByPairs is a free data retrieval call binding the contract method 0x4dbf0f39.
//
// Solidity: function getReservesByPairs(address[] _pairs) view returns(uint256[3][])
func (_FlashBotsUniswapQuery *FlashBotsUniswapQueryCaller) GetReservesByPairs(opts *bind.CallOpts, _pairs []common.Address) ([][3]*big.Int, error) {
	var out []interface{}
	err := _FlashBotsUniswapQuery.contract.Call(opts, &out, "getReservesByPairs", _pairs)

	if err != nil {
		return *new([][3]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([][3]*big.Int)).(*[][3]*big.Int)

	return out0, err

}

// GetReservesByPairs is a free data retrieval call binding the contract method 0x4dbf0f39.
//
// Solidity: function getReservesByPairs(address[] _pairs) view returns(uint256[3][])
func (_FlashBotsUniswapQuery *FlashBotsUniswapQuerySession) GetReservesByPairs(_pairs []common.Address) ([][3]*big.Int, error) {
	return _FlashBotsUniswapQuery.Contract.GetReservesByPairs(&_FlashBotsUniswapQuery.CallOpts, _pairs)
}

// GetReservesByPairs is a free data retrieval call binding the contract method 0x4dbf0f39.
//
// Solidity: function getReservesByPairs(address[] _pairs) view returns(uint256[3][])
func (_FlashBotsUniswapQuery *FlashBotsUniswapQueryCallerSession) GetReservesByPairs(_pairs []common.Address) ([][3]*big.Int, error) {
	return _FlashBotsUniswapQuery.Contract.GetReservesByPairs(&_FlashBotsUniswapQuery.CallOpts, _pairs)
}

// IUniswapV2PairABI is the input ABI used to generate the binding from.
const IUniswapV2PairABI = "[{\"inputs\":[],\"name\":\"getReserves\",\"outputs\":[{\"internalType\":\"uint112\",\"name\":\"reserve0\",\"type\":\"uint112\"},{\"internalType\":\"uint112\",\"name\":\"reserve1\",\"type\":\"uint112\"},{\"internalType\":\"uint32\",\"name\":\"blockTimestampLast\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"token0\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"token1\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// IUniswapV2PairFuncSigs maps the 4-byte function signature to its string representation.
var IUniswapV2PairFuncSigs = map[string]string{
	"0902f1ac": "getReserves()",
	"0dfe1681": "token0()",
	"d21220a7": "token1()",
}

// IUniswapV2Pair is an auto generated Go binding around an Ethereum contract.
type IUniswapV2Pair struct {
	IUniswapV2PairCaller     // Read-only binding to the contract
	IUniswapV2PairTransactor // Write-only binding to the contract
	IUniswapV2PairFilterer   // Log filterer for contract events
}

// IUniswapV2PairCaller is an auto generated read-only Go binding around an Ethereum contract.
type IUniswapV2PairCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUniswapV2PairTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IUniswapV2PairTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUniswapV2PairFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IUniswapV2PairFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUniswapV2PairSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IUniswapV2PairSession struct {
	Contract     *IUniswapV2Pair   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IUniswapV2PairCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IUniswapV2PairCallerSession struct {
	Contract *IUniswapV2PairCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// IUniswapV2PairTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IUniswapV2PairTransactorSession struct {
	Contract     *IUniswapV2PairTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// IUniswapV2PairRaw is an auto generated low-level Go binding around an Ethereum contract.
type IUniswapV2PairRaw struct {
	Contract *IUniswapV2Pair // Generic contract binding to access the raw methods on
}

// IUniswapV2PairCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IUniswapV2PairCallerRaw struct {
	Contract *IUniswapV2PairCaller // Generic read-only contract binding to access the raw methods on
}

// IUniswapV2PairTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IUniswapV2PairTransactorRaw struct {
	Contract *IUniswapV2PairTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIUniswapV2Pair creates a new instance of IUniswapV2Pair, bound to a specific deployed contract.
func NewIUniswapV2Pair(address common.Address, backend bind.ContractBackend) (*IUniswapV2Pair, error) {
	contract, err := bindIUniswapV2Pair(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IUniswapV2Pair{IUniswapV2PairCaller: IUniswapV2PairCaller{contract: contract}, IUniswapV2PairTransactor: IUniswapV2PairTransactor{contract: contract}, IUniswapV2PairFilterer: IUniswapV2PairFilterer{contract: contract}}, nil
}

// NewIUniswapV2PairCaller creates a new read-only instance of IUniswapV2Pair, bound to a specific deployed contract.
func NewIUniswapV2PairCaller(address common.Address, caller bind.ContractCaller) (*IUniswapV2PairCaller, error) {
	contract, err := bindIUniswapV2Pair(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IUniswapV2PairCaller{contract: contract}, nil
}

// NewIUniswapV2PairTransactor creates a new write-only instance of IUniswapV2Pair, bound to a specific deployed contract.
func NewIUniswapV2PairTransactor(address common.Address, transactor bind.ContractTransactor) (*IUniswapV2PairTransactor, error) {
	contract, err := bindIUniswapV2Pair(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IUniswapV2PairTransactor{contract: contract}, nil
}

// NewIUniswapV2PairFilterer creates a new log filterer instance of IUniswapV2Pair, bound to a specific deployed contract.
func NewIUniswapV2PairFilterer(address common.Address, filterer bind.ContractFilterer) (*IUniswapV2PairFilterer, error) {
	contract, err := bindIUniswapV2Pair(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IUniswapV2PairFilterer{contract: contract}, nil
}

// bindIUniswapV2Pair binds a generic wrapper to an already deployed contract.
func bindIUniswapV2Pair(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IUniswapV2PairABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IUniswapV2Pair *IUniswapV2PairRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IUniswapV2Pair.Contract.IUniswapV2PairCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IUniswapV2Pair *IUniswapV2PairRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IUniswapV2Pair.Contract.IUniswapV2PairTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IUniswapV2Pair *IUniswapV2PairRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IUniswapV2Pair.Contract.IUniswapV2PairTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IUniswapV2Pair *IUniswapV2PairCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IUniswapV2Pair.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IUniswapV2Pair *IUniswapV2PairTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IUniswapV2Pair.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IUniswapV2Pair *IUniswapV2PairTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IUniswapV2Pair.Contract.contract.Transact(opts, method, params...)
}

// GetReserves is a free data retrieval call binding the contract method 0x0902f1ac.
//
// Solidity: function getReserves() view returns(uint112 reserve0, uint112 reserve1, uint32 blockTimestampLast)
func (_IUniswapV2Pair *IUniswapV2PairCaller) GetReserves(opts *bind.CallOpts) (struct {
	Reserve0           *big.Int
	Reserve1           *big.Int
	BlockTimestampLast uint32
}, error) {
	var out []interface{}
	err := _IUniswapV2Pair.contract.Call(opts, &out, "getReserves")

	outstruct := new(struct {
		Reserve0           *big.Int
		Reserve1           *big.Int
		BlockTimestampLast uint32
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Reserve0 = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Reserve1 = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.BlockTimestampLast = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

// GetReserves is a free data retrieval call binding the contract method 0x0902f1ac.
//
// Solidity: function getReserves() view returns(uint112 reserve0, uint112 reserve1, uint32 blockTimestampLast)
func (_IUniswapV2Pair *IUniswapV2PairSession) GetReserves() (struct {
	Reserve0           *big.Int
	Reserve1           *big.Int
	BlockTimestampLast uint32
}, error) {
	return _IUniswapV2Pair.Contract.GetReserves(&_IUniswapV2Pair.CallOpts)
}

// GetReserves is a free data retrieval call binding the contract method 0x0902f1ac.
//
// Solidity: function getReserves() view returns(uint112 reserve0, uint112 reserve1, uint32 blockTimestampLast)
func (_IUniswapV2Pair *IUniswapV2PairCallerSession) GetReserves() (struct {
	Reserve0           *big.Int
	Reserve1           *big.Int
	BlockTimestampLast uint32
}, error) {
	return _IUniswapV2Pair.Contract.GetReserves(&_IUniswapV2Pair.CallOpts)
}

// Token0 is a free data retrieval call binding the contract method 0x0dfe1681.
//
// Solidity: function token0() view returns(address)
func (_IUniswapV2Pair *IUniswapV2PairCaller) Token0(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IUniswapV2Pair.contract.Call(opts, &out, "token0")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token0 is a free data retrieval call binding the contract method 0x0dfe1681.
//
// Solidity: function token0() view returns(address)
func (_IUniswapV2Pair *IUniswapV2PairSession) Token0() (common.Address, error) {
	return _IUniswapV2Pair.Contract.Token0(&_IUniswapV2Pair.CallOpts)
}

// Token0 is a free data retrieval call binding the contract method 0x0dfe1681.
//
// Solidity: function token0() view returns(address)
func (_IUniswapV2Pair *IUniswapV2PairCallerSession) Token0() (common.Address, error) {
	return _IUniswapV2Pair.Contract.Token0(&_IUniswapV2Pair.CallOpts)
}

// Token1 is a free data retrieval call binding the contract method 0xd21220a7.
//
// Solidity: function token1() view returns(address)
func (_IUniswapV2Pair *IUniswapV2PairCaller) Token1(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IUniswapV2Pair.contract.Call(opts, &out, "token1")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token1 is a free data retrieval call binding the contract method 0xd21220a7.
//
// Solidity: function token1() view returns(address)
func (_IUniswapV2Pair *IUniswapV2PairSession) Token1() (common.Address, error) {
	return _IUniswapV2Pair.Contract.Token1(&_IUniswapV2Pair.CallOpts)
}

// Token1 is a free data retrieval call binding the contract method 0xd21220a7.
//
// Solidity: function token1() view returns(address)
func (_IUniswapV2Pair *IUniswapV2PairCallerSession) Token1() (common.Address, error) {
	return _IUniswapV2Pair.Contract.Token1(&_IUniswapV2Pair.CallOpts)
}

// UniswapV2FactoryABI is the input ABI used to generate the binding from.
const UniswapV2FactoryABI = "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"allPairs\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"allPairsLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"getPair\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// UniswapV2FactoryFuncSigs maps the 4-byte function signature to its string representation.
var UniswapV2FactoryFuncSigs = map[string]string{
	"1e3dd18b": "allPairs(uint256)",
	"574f2ba3": "allPairsLength()",
	"e6a43905": "getPair(address,address)",
}

// UniswapV2Factory is an auto generated Go binding around an Ethereum contract.
type UniswapV2Factory struct {
	UniswapV2FactoryCaller     // Read-only binding to the contract
	UniswapV2FactoryTransactor // Write-only binding to the contract
	UniswapV2FactoryFilterer   // Log filterer for contract events
}

// UniswapV2FactoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type UniswapV2FactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapV2FactoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UniswapV2FactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapV2FactoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UniswapV2FactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapV2FactorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UniswapV2FactorySession struct {
	Contract     *UniswapV2Factory // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UniswapV2FactoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UniswapV2FactoryCallerSession struct {
	Contract *UniswapV2FactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// UniswapV2FactoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UniswapV2FactoryTransactorSession struct {
	Contract     *UniswapV2FactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// UniswapV2FactoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type UniswapV2FactoryRaw struct {
	Contract *UniswapV2Factory // Generic contract binding to access the raw methods on
}

// UniswapV2FactoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UniswapV2FactoryCallerRaw struct {
	Contract *UniswapV2FactoryCaller // Generic read-only contract binding to access the raw methods on
}

// UniswapV2FactoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UniswapV2FactoryTransactorRaw struct {
	Contract *UniswapV2FactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUniswapV2Factory creates a new instance of UniswapV2Factory, bound to a specific deployed contract.
func NewUniswapV2Factory(address common.Address, backend bind.ContractBackend) (*UniswapV2Factory, error) {
	contract, err := bindUniswapV2Factory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UniswapV2Factory{UniswapV2FactoryCaller: UniswapV2FactoryCaller{contract: contract}, UniswapV2FactoryTransactor: UniswapV2FactoryTransactor{contract: contract}, UniswapV2FactoryFilterer: UniswapV2FactoryFilterer{contract: contract}}, nil
}

// NewUniswapV2FactoryCaller creates a new read-only instance of UniswapV2Factory, bound to a specific deployed contract.
func NewUniswapV2FactoryCaller(address common.Address, caller bind.ContractCaller) (*UniswapV2FactoryCaller, error) {
	contract, err := bindUniswapV2Factory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UniswapV2FactoryCaller{contract: contract}, nil
}

// NewUniswapV2FactoryTransactor creates a new write-only instance of UniswapV2Factory, bound to a specific deployed contract.
func NewUniswapV2FactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*UniswapV2FactoryTransactor, error) {
	contract, err := bindUniswapV2Factory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UniswapV2FactoryTransactor{contract: contract}, nil
}

// NewUniswapV2FactoryFilterer creates a new log filterer instance of UniswapV2Factory, bound to a specific deployed contract.
func NewUniswapV2FactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*UniswapV2FactoryFilterer, error) {
	contract, err := bindUniswapV2Factory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UniswapV2FactoryFilterer{contract: contract}, nil
}

// bindUniswapV2Factory binds a generic wrapper to an already deployed contract.
func bindUniswapV2Factory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(UniswapV2FactoryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UniswapV2Factory *UniswapV2FactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UniswapV2Factory.Contract.UniswapV2FactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UniswapV2Factory *UniswapV2FactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapV2Factory.Contract.UniswapV2FactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UniswapV2Factory *UniswapV2FactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UniswapV2Factory.Contract.UniswapV2FactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UniswapV2Factory *UniswapV2FactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UniswapV2Factory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UniswapV2Factory *UniswapV2FactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapV2Factory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UniswapV2Factory *UniswapV2FactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UniswapV2Factory.Contract.contract.Transact(opts, method, params...)
}

// AllPairs is a free data retrieval call binding the contract method 0x1e3dd18b.
//
// Solidity: function allPairs(uint256 ) view returns(address)
func (_UniswapV2Factory *UniswapV2FactoryCaller) AllPairs(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _UniswapV2Factory.contract.Call(opts, &out, "allPairs", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AllPairs is a free data retrieval call binding the contract method 0x1e3dd18b.
//
// Solidity: function allPairs(uint256 ) view returns(address)
func (_UniswapV2Factory *UniswapV2FactorySession) AllPairs(arg0 *big.Int) (common.Address, error) {
	return _UniswapV2Factory.Contract.AllPairs(&_UniswapV2Factory.CallOpts, arg0)
}

// AllPairs is a free data retrieval call binding the contract method 0x1e3dd18b.
//
// Solidity: function allPairs(uint256 ) view returns(address)
func (_UniswapV2Factory *UniswapV2FactoryCallerSession) AllPairs(arg0 *big.Int) (common.Address, error) {
	return _UniswapV2Factory.Contract.AllPairs(&_UniswapV2Factory.CallOpts, arg0)
}

// AllPairsLength is a free data retrieval call binding the contract method 0x574f2ba3.
//
// Solidity: function allPairsLength() view returns(uint256)
func (_UniswapV2Factory *UniswapV2FactoryCaller) AllPairsLength(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UniswapV2Factory.contract.Call(opts, &out, "allPairsLength")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AllPairsLength is a free data retrieval call binding the contract method 0x574f2ba3.
//
// Solidity: function allPairsLength() view returns(uint256)
func (_UniswapV2Factory *UniswapV2FactorySession) AllPairsLength() (*big.Int, error) {
	return _UniswapV2Factory.Contract.AllPairsLength(&_UniswapV2Factory.CallOpts)
}

// AllPairsLength is a free data retrieval call binding the contract method 0x574f2ba3.
//
// Solidity: function allPairsLength() view returns(uint256)
func (_UniswapV2Factory *UniswapV2FactoryCallerSession) AllPairsLength() (*big.Int, error) {
	return _UniswapV2Factory.Contract.AllPairsLength(&_UniswapV2Factory.CallOpts)
}

// GetPair is a free data retrieval call binding the contract method 0xe6a43905.
//
// Solidity: function getPair(address , address ) view returns(address)
func (_UniswapV2Factory *UniswapV2FactoryCaller) GetPair(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (common.Address, error) {
	var out []interface{}
	err := _UniswapV2Factory.contract.Call(opts, &out, "getPair", arg0, arg1)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetPair is a free data retrieval call binding the contract method 0xe6a43905.
//
// Solidity: function getPair(address , address ) view returns(address)
func (_UniswapV2Factory *UniswapV2FactorySession) GetPair(arg0 common.Address, arg1 common.Address) (common.Address, error) {
	return _UniswapV2Factory.Contract.GetPair(&_UniswapV2Factory.CallOpts, arg0, arg1)
}

// GetPair is a free data retrieval call binding the contract method 0xe6a43905.
//
// Solidity: function getPair(address , address ) view returns(address)
func (_UniswapV2Factory *UniswapV2FactoryCallerSession) GetPair(arg0 common.Address, arg1 common.Address) (common.Address, error) {
	return _UniswapV2Factory.Contract.GetPair(&_UniswapV2Factory.CallOpts, arg0, arg1)
}
