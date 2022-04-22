//go:build prod
// +build prod

package env

import (
	"math/big"

	"github.com/ethereum/go-ethereum/params"
)

const Tag = "prod"

const IsConfigEncrypted bool = true
const IsDeleteConfigAfterUsed bool = false
const DefaultEurusChainId int64 = 1008
const DefaultMainnetChainId int64 = 1

func CreateEurusChainConfig() *params.ChainConfig {
	return &params.ChainConfig{
		ChainID:             big.NewInt(DefaultEurusChainId),
		HomesteadBlock:      nil,
		DAOForkBlock:        nil,
		DAOForkSupport:      true,
		EIP150Block:         nil,
		EIP155Block:         big.NewInt(0),
		EIP158Block:         nil,
		ByzantiumBlock:      nil,
		ConstantinopleBlock: nil,
		PetersburgBlock:     nil,
		IstanbulBlock:       nil,
		MuirGlacierBlock:    nil,
		BerlinBlock:         nil,
		LondonBlock:         nil,
		Clique: &params.CliqueConfig{
			Period: 30,
			Epoch:  30000,
		},
	}
}
