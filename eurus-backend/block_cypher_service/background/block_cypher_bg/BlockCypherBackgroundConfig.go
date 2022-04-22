package block_cypher_bg

import (
	"eurus-backend/foundation/server"
)

type BlockCypherBackgroundConfig struct {
	server.ServerConfigBase
	CoinListJson string   `json:"coinList"`
	CoinList     []string `json:"-"`
	Chain        string   `json:"chain"`
}

func (me *BlockCypherBackgroundConfig) GetServerConfigBase() *server.ServerConfigBase {
	return &me.ServerConfigBase
}

func (me *BlockCypherBackgroundConfig) GetParent() interface{} {
	return &me.ServerConfigBase
}
