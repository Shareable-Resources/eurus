package ethereum

import (
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

var extractStateNameMap map[ExtractState]string = map[ExtractState]string{
	ExtractFailed:   "Failed",
	ExtractABIError: "ABIError",
	ExtractSuccess:  "Success",
}

var DefaultABIDecoder *ABIDecoder = NewABIDecoder()

var ABIMethodNotFoundError error = errors.New("Method not found")

type ABIDecoder struct {
	abiCache map[string]*abi.ABI
}

func NewABIDecoder() *ABIDecoder {
	decoder := new(ABIDecoder)
	decoder.abiCache = make(map[string]*abi.ABI)
	return decoder
}

func (me *ABIDecoder) ImportABIJson(contractName string, abiJsonStr string) error {
	loadedAbi, err := abi.JSON(strings.NewReader(abiJsonStr))
	if err != nil {
		return err
	}
	me.abiCache[contractName] = &loadedAbi
	return nil
}

func (me *ABIDecoder) GetABI(contractName string) *abi.ABI {
	abiObj := me.abiCache[contractName]
	return abiObj
}

func (me *ABIDecoder) DecodeABIInputArgument(abiInput []byte, contractName string, expectedMethodName string) (map[string]interface{}, error, ExtractState) {
	if len(abiInput) < 4 {
		return nil, ABIMethodNotFoundError, ExtractFailed
	}
	var contractAbi *abi.ABI
	var ok bool
	if contractAbi, ok = me.abiCache[contractName]; !ok {
		return nil, errors.New("Contract not found"), ExtractFailed
	}
	// decode txInput method signature
	method, err := contractAbi.MethodById(abiInput[:4])
	if err != nil {
		return nil, ABIMethodNotFoundError, ExtractFailed
	}
	if method.RawName != expectedMethodName {
		return nil, ABIMethodNotFoundError, ExtractFailed
	}
	meth, ok := contractAbi.Methods[expectedMethodName]
	if !ok {
		return nil, ABIMethodNotFoundError, ExtractFailed
	}
	dataMap := make(map[string]interface{})
	err = meth.Inputs.UnpackIntoMap(dataMap, abiInput[4:])
	if err != nil {
		return nil, err, ExtractABIError
	}
	return dataMap, nil, ExtractSuccess

}

func (me *ABIDecoder) DecodeABIEventData(abiEventData []byte, contractName string, eventName string) ([]interface{}, error) {
	var contractAbi *abi.ABI
	var ok bool
	if contractAbi, ok = me.abiCache[contractName]; !ok {
		return nil, errors.New("Contract not found")
	}
	var event abi.Event
	event, ok = contractAbi.Events[eventName]
	if !ok {
		return nil, errors.New("Event not found")
	}

	return event.Inputs.Unpack(abiEventData)
}

type ExtractState int16

const (
	ExtractABIError ExtractState = -2
	ExtractFailed   ExtractState = -1
	ExtractSuccess  ExtractState = 0
)

func (me ExtractState) String() string {
	val, ok := extractStateNameMap[me]
	if !ok {
		return ""
	}
	return val
}
