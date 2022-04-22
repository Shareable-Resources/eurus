package ethereum

import (
	"eurus-backend/env"
	"math/big"

	"github.com/ethereum/go-ethereum/params"
)

var EurusChainConfig *params.ChainConfig = env.CreateEurusChainConfig()

func GetChainConfigFromChainId(chainId *big.Int) *params.ChainConfig {
	switch chainId.Uint64() {
	case 1:
		return params.MainnetChainConfig
	case 3:
		return params.RopstenChainConfig
	case 4:
		return params.RinkebyChainConfig
	case 5:
		return params.GoerliChainConfig
	default:
		return EurusChainConfig
	}
}
