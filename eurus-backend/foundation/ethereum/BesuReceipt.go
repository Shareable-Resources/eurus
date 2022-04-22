package ethereum

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
)

type BesuReceipt struct {
	types.Receipt
	RevertReason      string
	EffectiveGasPrice *big.Int
}

func (me *BesuReceipt) UnmarshalJSON(b []byte) error {
	type cloneType BesuReceipt

	type unmarshalType struct {
		EffectiveGasPrice string `json:"effectiveGasPrice"`
		RevertReason      string `json:"revertReason"`
	}
	var obj unmarshalType
	err := json.Unmarshal(b, &obj)
	if err != nil {
		return err
	}
	me.RevertReason = obj.RevertReason
	if obj.EffectiveGasPrice != "" && strings.HasPrefix(obj.EffectiveGasPrice, "0x") {
		hexStr := obj.EffectiveGasPrice[2:]
		if len(hexStr)%2 == 1 {
			hexStr = "0" + hexStr
		}
		data, err := hex.DecodeString(hexStr)
		if err != nil {
			return err
		}
		var paddingBytes []byte = make([]byte, 8)
		valBytes := append(paddingBytes[:], data[:]...)
		length := len(valBytes)
		me.EffectiveGasPrice = big.NewInt(0).SetUint64(binary.BigEndian.Uint64(valBytes[length-8:]))
	}

	return json.Unmarshal(b, (*cloneType)(me))
}

func (me *BesuReceipt) MarshalJSON() ([]byte, error) {
	if me.RevertReason != "" {
		var marshalMap map[string]interface{}

		parentByte, err := json.Marshal(me.Receipt)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(parentByte, &marshalMap)
		if err != nil {
			return nil, err
		}
		marshalMap["revertReason"] = me.RevertReason
		return json.Marshal(marshalMap)
	}

	type cloneType BesuReceipt
	return json.Marshal((*cloneType)(me))
}
