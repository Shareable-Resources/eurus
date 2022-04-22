package patch

import "eurus-backend/foundation/server"

type PatchConfigBase struct {
	server.ServerConfigBase
	GasLimit uint64 `json:"gasLimit"`
}

func (me *PatchConfigBase) GetServerConfigBase() *server.ServerConfigBase {
	return &me.ServerConfigBase
}

func (me *PatchConfigBase) GetParent() interface{} {
	return &me.ServerConfigBase
}
