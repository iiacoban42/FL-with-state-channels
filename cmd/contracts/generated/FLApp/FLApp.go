// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package FLApp

import (
	"errors"
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
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// ChannelAllocation is an auto generated low-level Go binding around an user-defined struct.
type ChannelAllocation struct {
	Assets   []common.Address
	Balances [][]*big.Int
	Locked   []ChannelSubAlloc
}

// ChannelParams is an auto generated low-level Go binding around an user-defined struct.
type ChannelParams struct {
	ChallengeDuration *big.Int
	Nonce             *big.Int
	Participants      []common.Address
	App               common.Address
	LedgerChannel     bool
	VirtualChannel    bool
}

// ChannelState is an auto generated low-level Go binding around an user-defined struct.
type ChannelState struct {
	ChannelID [32]byte
	Version   uint64
	Outcome   ChannelAllocation
	AppData   []byte
	IsFinal   bool
}

// ChannelSubAlloc is an auto generated low-level Go binding around an user-defined struct.
type ChannelSubAlloc struct {
	ID       [32]byte
	Balances []*big.Int
	IndexMap []uint16
}

// FLAppMetaData contains all meta data concerning the FLApp contract.
var FLAppMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"challengeDuration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"app\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"ledgerChannel\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"virtualChannel\",\"type\":\"bool\"}],\"internalType\":\"structChannel.Params\",\"name\":\"params\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"channelID\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"address[]\",\"name\":\"assets\",\"type\":\"address[]\"},{\"internalType\":\"uint256[][]\",\"name\":\"balances\",\"type\":\"uint256[][]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"ID\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"balances\",\"type\":\"uint256[]\"},{\"internalType\":\"uint16[]\",\"name\":\"indexMap\",\"type\":\"uint16[]\"}],\"internalType\":\"structChannel.SubAlloc[]\",\"name\":\"locked\",\"type\":\"tuple[]\"}],\"internalType\":\"structChannel.Allocation\",\"name\":\"outcome\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structChannel.State\",\"name\":\"from\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"channelID\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"address[]\",\"name\":\"assets\",\"type\":\"address[]\"},{\"internalType\":\"uint256[][]\",\"name\":\"balances\",\"type\":\"uint256[][]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"ID\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"balances\",\"type\":\"uint256[]\"},{\"internalType\":\"uint16[]\",\"name\":\"indexMap\",\"type\":\"uint16[]\"}],\"internalType\":\"structChannel.SubAlloc[]\",\"name\":\"locked\",\"type\":\"tuple[]\"}],\"internalType\":\"structChannel.Allocation\",\"name\":\"outcome\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structChannel.State\",\"name\":\"to\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"signerIdx\",\"type\":\"uint256\"}],\"name\":\"validTransition\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506119b1806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c80630d1feb4f14610030575b600080fd5b61004361003e366004611294565b610045565b005b60026100546040860186611759565b90501461007c5760405162461bcd60e51b8152600401610073906115ed565b60405180910390fd5b600061008b60608501856117a7565b60008161009457fe5b919091013560f81c9150508181146100be5760405162461bcd60e51b81526004016100739061151d565b6100cb60608401846117a7565b6000816100d457fe5b9091013560f81c6001838101161490506101005760405162461bcd60e51b81526004016100739061156e565b61010d60608501856117a7565b600481811061011857fe5b919091013560f81c1590506101985761013460608401846117a7565b600281811061013f57fe5b909101356001600160f81b031916905061015c60608601866117a7565b600281811061016757fe5b9050013560f81c60f81b6001600160f81b031916146101985760405162461bcd60e51b815260040161007390611547565b60006101a760608501856117a7565b60028181106101b257fe5b919091013560f81c9150606190508082016101d060608701876117a7565b60028181106101db57fe5b919091013560f81c90506101f260608801886117a7565b60038181106101fd57fe5b9050013560f81c60f81b60f81c60ff16111561022b5760405162461bcd60e51b8152600401610073906114c2565b60ff841661065c5761024060608801886117a7565b600481811061024b57fe5b919091013560f81c15905061065c5761030461026a60608901896117a7565b6102799160339160059161184d565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506102bb9250505060608901896117a7565b6102ca9160339160059161184d565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250610bdc92505050565b156103215760405162461bcd60e51b815260040161007390611395565b61039161033160608901896117a7565b6103409160619160339161184d565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506103829250505060608901896117a7565b6102ca9160619160339161184d565b6103ad5760405162461bcd60e51b8152600401610073906113cc565b6104856103bd60608901896117a7565b6103d29160ff8688018116929087169161184d565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506104149250505060608901896117a7565b6104299160ff8789018116929088169161184d565b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061046b9250505060608b018b6117a7565b600381811061047657fe5b919091013560f81c9050610c5a565b6104a15760405162461bcd60e51b815260040161007390611654565b61051d6104b160608901896117a7565b6104c69160ff8588018116929086169161184d565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506105089250505060608901896117a7565b6104299160ff8689018116929087169161184d565b6105395760405162461bcd60e51b815260040161007390611471565b600161054860608901896117a7565b600381811061055357fe5b60ff92013560f81c9290920116905061056f60608801886117a7565b600381811061057a57fe5b9050013560f81c60f81b60f81c60ff16146105a75760405162461bcd60e51b81526004016100739061135e565b6105b460608701876117a7565b6105c160608a018a6117a7565b60038181106105cc57fe5b919091013560f81c850160ff1690508181106105e457fe5b919091013560f81c1515905080610640575061060360608701876117a7565b61061060608a018a6117a7565b600381811061061b57fe5b919091013560f81c840160ff16905081811061063357fe5b919091013560f81c151590505b61065c5760405162461bcd60e51b815260040161007390611716565b8360ff16600114156108a95761067861026a60608901896117a7565b6106945760405162461bcd60e51b815260040161007390611327565b6106a160608701876117a7565b60038181106106ac57fe5b919091013560f81c90506106c360608901896117a7565b60038181106106ce57fe5b9050013560f81c60f81b60f81c60ff16146106fb5760405162461bcd60e51b815260040161007390611403565b61077761070b60608901896117a7565b6107209160ff8688018116929087169161184d565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506107629250505060608901896117a7565b6102ca9160ff8789018116929088169161184d565b6107935760405162461bcd60e51b8152600401610073906116df565b61080f6107a360608901896117a7565b6107b89160ff8588018116929086169161184d565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506107fa9250505060608901896117a7565b6102ca9160ff8689018116929087169161184d565b61082b5760405162461bcd60e51b81526004016100739061143a565b61088c61083b60608801886117a7565b61084a9160619160339161184d565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506103829250505060608a018a6117a7565b156108a95760405162461bcd60e51b81526004016100739061161d565b60006108f56108bb60608901896117a7565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250610ce992505050565b905080151561090a60a0890160808a01611274565b1515146109295760405162461bcd60e51b815260040161007390611592565b6109c361093960408901896117ec565b6109439080611759565b808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152506109829250505060408b018b6117ec565b61098c9080611759565b80806020026020016040519081016040528093929190818152602001838360200280828437600092019190915250610d6592505050565b610a136109d360408901896117ec565b6109e1906040810190611759565b6109ea916118c3565b6109f760408b018b6117ec565b610a05906040810190611759565b610a0e916118c3565b610e60565b6000610a2260408a018a6117ec565b610a30906020810190611759565b610a3991611875565b9050610a85610a4b60608a018a6117a7565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250610ec592505050565b15610ba357670de0b6b3a764000060016000805b8451811015610b9e5783610ab060408f018f6117ec565b610abe906020810190611759565b83818110610ac857fe5b9050602002810190610ada9190611759565b8560ff16818110610ae757fe5b9050602002013501858281518110610afb57fe5b60200260200101518460ff1681518110610b1157fe5b602090810291909101015283610b2a60408f018f6117ec565b610b38906020810190611759565b83818110610b4257fe5b9050602002810190610b549190611759565b8460ff16818110610b6157fe5b9050602002013503858281518110610b7557fe5b60200260200101518360ff1681518110610b8b57fe5b6020908102919091010152600101610a99565b505050505b610bd0610bb360408a018a6117ec565b610bc1906020810190611759565b610bca91611875565b82610f11565b50505050505050505050565b60008151835114610bef57506000610c54565b60005b8351811015610c4e57828181518110610c0757fe5b602001015160f81c60f81b6001600160f81b031916848281518110610c2857fe5b01602001516001600160f81b03191614610c46576000915050610c54565b600101610bf2565b50600190505b92915050565b60008251845114610c6d57506000610ce2565b60005b8451811015610cdc578260ff16811415610c8957610cd4565b838181518110610c9557fe5b602001015160f81c60f81b6001600160f81b031916858281518110610cb657fe5b01602001516001600160f81b03191614610cd4576000915050610ce2565b600101610c70565b50600190505b9392505050565b805160009082906003908110610cfb57fe5b016020015182516001600160f81b03199091169083906002908110610d1c57fe5b01602001516001600160f81b031916148015610d4f5750815182906004908110610d4257fe5b60209101015160f81c6003145b15610d5c57506001610d60565b5060005b919050565b8051825114610dbb576040805162461bcd60e51b815260206004820152601960248201527f616464726573735b5d3a20756e657175616c206c656e67746800000000000000604482015290519081900360640190fd5b60005b8251811015610e5b57818181518110610dd357fe5b60200260200101516001600160a01b0316838281518110610df057fe5b60200260200101516001600160a01b031614610e53576040805162461bcd60e51b815260206004820152601760248201527f616464726573735b5d3a20756e657175616c206974656d000000000000000000604482015290519081900360640190fd5b600101610dbe565b505050565b8051825114610e815760405162461bcd60e51b8152600401610073906116a8565b60005b8251811015610e5b57610ebd838281518110610e9c57fe5b6020026020010151838381518110610eb057fe5b6020026020010151610f76565b600101610e84565b805160009060029083906004908110610eda57fe5b016020015160f81c148015610d4f5750815160009083906033908110610efc57fe5b016020015160f81c14610d5c57506001610d60565b8051825114610f325760405162461bcd60e51b8152600401610073906115b6565b60005b8251811015610e5b57610f6e838281518110610f4d57fe5b6020026020010151838381518110610f6157fe5b6020026020010151610fbf565b600101610f35565b8051825114610f975760405162461bcd60e51b8152600401610073906114ef565b610fa982602001518260200151610fbf565b610fbb826040015182604001516110a3565b5050565b8051825114611015576040805162461bcd60e51b815260206004820152601960248201527f75696e743235365b5d3a20756e657175616c206c656e67746800000000000000604482015290519081900360640190fd5b60005b8251811015610e5b5781818151811061102d57fe5b602002602001015183828151811061104157fe5b60200260200101511461109b576040805162461bcd60e51b815260206004820152601760248201527f75696e743235365b5d3a20756e657175616c206974656d000000000000000000604482015290519081900360640190fd5b600101611018565b80518251146110f9576040805162461bcd60e51b815260206004820152601860248201527f75696e7431365b5d3a20756e657175616c206c656e6774680000000000000000604482015290519081900360640190fd5b60005b8251811015610e5b5781818151811061111157fe5b602002602001015161ffff1683828151811061112957fe5b602002602001015161ffff1614611180576040805162461bcd60e51b815260206004820152601660248201527575696e7431365b5d3a20756e657175616c206974656d60501b604482015290519081900360640190fd5b6001016110fc565b600082601f830112611198578081fd5b813560206111ad6111a88361182f565b61180b565b82815281810190858301838502870184018810156111c9578586fd5b855b858110156111f657813561ffff811681146111e4578788fd5b845292840192908401906001016111cb565b5090979650505050505050565b600082601f830112611213578081fd5b813560206112236111a88361182f565b828152818101908583018385028701840188101561123f578586fd5b855b858110156111f657813584529284019290840190600101611241565b600060a0828403121561126e578081fd5b50919050565b600060208284031215611285578081fd5b81358015158114610ce2578182fd5b600080600080608085870312156112a9578283fd5b843567ffffffffffffffff808211156112c0578485fd5b9086019060c082890312156112d3578485fd5b909450602086013590808211156112e8578485fd5b6112f48883890161125d565b94506040870135915080821115611309578384fd5b506113168782880161125d565b949793965093946060013593505050565b60208082526019908201527f6163746f722063616e6e6f74206368616e6765206d6f64656c00000000000000604082015260600190565b6020808252601a908201527f6163746f72206d75737420696e6372656d656e7420726f756e64000000000000604082015260600190565b60208082526017908201527f6163746f722063616e6e6f7420736b6970206d6f64656c000000000000000000604082015260600190565b6020808252601d908201527f6163746f722063616e6e6f74206f766572726964652077656967687473000000604082015260600190565b6020808252601c908201527f6163746f722063616e6e6f7420696e6372656d656e7420726f756e6400000000604082015260600190565b6020808252601a908201527f6163746f722063616e6e6f74206f76657272696465206c6f7373000000000000604082015260600190565b60208082526031908201527f6163746f722063616e6e6f74206f76657272696465206c6f7373206f75747369604082015270032329031bab93932b73a103937bab7321607d1b606082015260800190565b602080825260139082015272726f756e64206f7574206f6620626f756e647360681b604082015260600190565b60208082526014908201527314dd58905b1b1bd8ce881d5b995c5d585b08125160621b604082015260600190565b60208082526010908201526f30b1ba37b9103737ba1039b4b3b732b960811b604082015260600190565b6020808252600d908201526c1c9bdd5b990818da185b99d959609a1b604082015260600190565b6020808252600a90820152693732bc3a1030b1ba37b960b11b604082015260600190565b6020808252600a908201526966696e616c20666c616760b01b604082015260600190565b6020808252601b908201527f75696e743235365b5d5b5d3a20756e657175616c206c656e6774680000000000604082015260600190565b6020808252601690820152756e756d626572206f66207061727469636970616e747360501b604082015260600190565b60208082526018908201527f6163746f722063616e6e6f7420736b6970207765696768740000000000000000604082015260600190565b60208082526034908201527f6163746f722063616e6e6f74206f76657272696465206163637572616379206f6040820152731d5d1cda59194818dd5c9c995b9d081c9bdd5b9960621b606082015260800190565b6020808252601a908201527f537562416c6c6f635b5d3a20756e657175616c206c656e677468000000000000604082015260600190565b6020808252601e908201527f6163746f722063616e6e6f74206f766572726964652061636375726163790000604082015260600190565b60208082526023908201527f6163746f722063616e6e6f7420736b697020616363757261637920616e64206c6040820152626f737360e81b606082015260800190565b6000808335601e1984360301811261176f578283fd5b83018035915067ffffffffffffffff821115611789578283fd5b60209081019250810236038213156117a057600080fd5b9250929050565b6000808335601e198436030181126117bd578283fd5b83018035915067ffffffffffffffff8211156117d7578283fd5b6020019150368190038213156117a057600080fd5b60008235605e19833603018112611801578182fd5b9190910192915050565b60405181810167ffffffffffffffff8111828210171561182757fe5b604052919050565b600067ffffffffffffffff82111561184357fe5b5060209081020190565b6000808585111561185c578182fd5b83861115611868578182fd5b5050820193919092039150565b60006118836111a88461182f565b8381526020808201919084845b878110156118b7576118a53683358901611203565b85529382019390820190600101611890565b50919695505050505050565b60006118d16111a88461182f565b8381526020808201919084845b878110156118b7578135870160608082360312156118fa578788fd5b604080519182019167ffffffffffffffff808411828510171561191957fe5b92825283358152868401359280841115611931578a8bfd5b61193d36858701611203565b8883015282850135935080841115611953578a8bfd5b5061196036848601611188565b918101919091528752505093820193908201906001016118de56fea26469706673582212206bef8fb4de6471d199ff8541cabe181bf92e8f12eb4ff1ebdbd779db9597f5c264736f6c63430007060033",
}

// FLAppABI is the input ABI used to generate the binding from.
// Deprecated: Use FLAppMetaData.ABI instead.
var FLAppABI = FLAppMetaData.ABI

// FLAppBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use FLAppMetaData.Bin instead.
var FLAppBin = FLAppMetaData.Bin

// DeployFLApp deploys a new Ethereum contract, binding an instance of FLApp to it.
func DeployFLApp(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *FLApp, error) {
	parsed, err := FLAppMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FLAppBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &FLApp{FLAppCaller: FLAppCaller{contract: contract}, FLAppTransactor: FLAppTransactor{contract: contract}, FLAppFilterer: FLAppFilterer{contract: contract}}, nil
}

// FLApp is an auto generated Go binding around an Ethereum contract.
type FLApp struct {
	FLAppCaller     // Read-only binding to the contract
	FLAppTransactor // Write-only binding to the contract
	FLAppFilterer   // Log filterer for contract events
}

// FLAppCaller is an auto generated read-only Go binding around an Ethereum contract.
type FLAppCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FLAppTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FLAppTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FLAppFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FLAppFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FLAppSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FLAppSession struct {
	Contract     *FLApp            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FLAppCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FLAppCallerSession struct {
	Contract *FLAppCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// FLAppTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FLAppTransactorSession struct {
	Contract     *FLAppTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FLAppRaw is an auto generated low-level Go binding around an Ethereum contract.
type FLAppRaw struct {
	Contract *FLApp // Generic contract binding to access the raw methods on
}

// FLAppCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FLAppCallerRaw struct {
	Contract *FLAppCaller // Generic read-only contract binding to access the raw methods on
}

// FLAppTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FLAppTransactorRaw struct {
	Contract *FLAppTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFLApp creates a new instance of FLApp, bound to a specific deployed contract.
func NewFLApp(address common.Address, backend bind.ContractBackend) (*FLApp, error) {
	contract, err := bindFLApp(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FLApp{FLAppCaller: FLAppCaller{contract: contract}, FLAppTransactor: FLAppTransactor{contract: contract}, FLAppFilterer: FLAppFilterer{contract: contract}}, nil
}

// NewFLAppCaller creates a new read-only instance of FLApp, bound to a specific deployed contract.
func NewFLAppCaller(address common.Address, caller bind.ContractCaller) (*FLAppCaller, error) {
	contract, err := bindFLApp(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FLAppCaller{contract: contract}, nil
}

// NewFLAppTransactor creates a new write-only instance of FLApp, bound to a specific deployed contract.
func NewFLAppTransactor(address common.Address, transactor bind.ContractTransactor) (*FLAppTransactor, error) {
	contract, err := bindFLApp(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FLAppTransactor{contract: contract}, nil
}

// NewFLAppFilterer creates a new log filterer instance of FLApp, bound to a specific deployed contract.
func NewFLAppFilterer(address common.Address, filterer bind.ContractFilterer) (*FLAppFilterer, error) {
	contract, err := bindFLApp(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FLAppFilterer{contract: contract}, nil
}

// bindFLApp binds a generic wrapper to an already deployed contract.
func bindFLApp(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FLAppMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FLApp *FLAppRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FLApp.Contract.FLAppCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FLApp *FLAppRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FLApp.Contract.FLAppTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FLApp *FLAppRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FLApp.Contract.FLAppTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FLApp *FLAppCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FLApp.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FLApp *FLAppTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FLApp.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FLApp *FLAppTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FLApp.Contract.contract.Transact(opts, method, params...)
}

// ValidTransition is a free data retrieval call binding the contract method 0x0d1feb4f.
//
// Solidity: function validTransition((uint256,uint256,address[],address,bool,bool) params, (bytes32,uint64,(address[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool) from, (bytes32,uint64,(address[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool) to, uint256 signerIdx) pure returns()
func (_FLApp *FLAppCaller) ValidTransition(opts *bind.CallOpts, params ChannelParams, from ChannelState, to ChannelState, signerIdx *big.Int) error {
	var out []interface{}
	err := _FLApp.contract.Call(opts, &out, "validTransition", params, from, to, signerIdx)

	if err != nil {
		return err
	}

	return err

}

// ValidTransition is a free data retrieval call binding the contract method 0x0d1feb4f.
//
// Solidity: function validTransition((uint256,uint256,address[],address,bool,bool) params, (bytes32,uint64,(address[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool) from, (bytes32,uint64,(address[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool) to, uint256 signerIdx) pure returns()
func (_FLApp *FLAppSession) ValidTransition(params ChannelParams, from ChannelState, to ChannelState, signerIdx *big.Int) error {
	return _FLApp.Contract.ValidTransition(&_FLApp.CallOpts, params, from, to, signerIdx)
}

// ValidTransition is a free data retrieval call binding the contract method 0x0d1feb4f.
//
// Solidity: function validTransition((uint256,uint256,address[],address,bool,bool) params, (bytes32,uint64,(address[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool) from, (bytes32,uint64,(address[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool) to, uint256 signerIdx) pure returns()
func (_FLApp *FLAppCallerSession) ValidTransition(params ChannelParams, from ChannelState, to ChannelState, signerIdx *big.Int) error {
	return _FLApp.Contract.ValidTransition(&_FLApp.CallOpts, params, from, to, signerIdx)
}
