package signature

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	signer "github.com/ethereum/go-ethereum/signer/core"
)

const EIP712DomainName string = "DApp Payment API"
const EIP712DomainVersion string = "1"

func ReadSignature(s string) ([]byte, error) {
	var ret []byte
	var err error

	// Normally we represent signature in string starts with 0x, but hex.DecodeString() only need the hex part
	// So write a simple function to handle both cases
	if strings.HasPrefix(s, "0x") {
		ret, err = hex.DecodeString(s[2:])
	} else {
		ret, err = hex.DecodeString(s)
	}

	if err != nil {
		return nil, err
	}

	return ret, nil
}

func GetSigner(chainID int, txHash common.Hash, coin string, amount *big.Int, merchant string, tag string, signature []byte) (common.Address, error) {
	// Louis TODO: handle v including chain id
	if len(signature) != 65 {
		return zeroAddr(), fmt.Errorf("invalid signature length: %d", len(signature))
	}

	if signature[64] != 27 && signature[64] != 28 {
		return zeroAddr(), fmt.Errorf("invalid recovery id: %d", signature[64])
	}
	signature[64] -= 27

	_chainID := math.HexOrDecimal256(*big.NewInt(int64(chainID)))
	_amount := math.HexOrDecimal256(*amount)

	data := &signer.TypedData{
		Types: signer.Types{
			"EIP712Domain": []signer.Type{
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
			},
			"Message": []signer.Type{
				{Name: "txHash", Type: "bytes32"},
				{Name: "coin", Type: "string"},
				{Name: "amount", Type: "uint256"},
				{Name: "merchant", Type: "string"},
				{Name: "tag", Type: "string"},
			},
		},
		PrimaryType: "Message",
		Domain: signer.TypedDataDomain{
			Name:    EIP712DomainName,
			Version: EIP712DomainVersion,
			ChainId: &_chainID,
		},
		Message: signer.TypedDataMessage{
			"txHash":   txHash.Hex(),
			"coin":     coin,
			"amount":   &_amount,
			"merchant": merchant,
			"tag":      tag,
		},
	}

	domainSeparator, err := data.HashStruct("EIP712Domain", data.Domain.Map())
	if err != nil {
		return zeroAddr(), err
	}

	typedDataHash, err := data.HashStruct(data.PrimaryType, data.Message)
	if err != nil {
		return zeroAddr(), err
	}

	rawData := []byte{0x19, 0x01}
	rawData = append(rawData, domainSeparator...)
	rawData = append(rawData, typedDataHash...)
	pubKeyRaw, err := crypto.Ecrecover(crypto.Keccak256Hash(rawData).Bytes(), signature)
	if err != nil {
		return zeroAddr(), err
	}

	pubKey, err := crypto.UnmarshalPubkey(pubKeyRaw)
	if err != nil {
		return zeroAddr(), err
	}

	return crypto.PubkeyToAddress(*pubKey), nil
}

func zeroAddr() common.Address {
	return *new(common.Address)
}
