package sc_validator

import (
	"encoding/json"
	"fmt"
)

type SmartContractInfo struct {
	Address string `json:"address"`
}

func (me *SmartContractInfo) UnmarshalJSON(data []byte) error {
	type cloneInfo SmartContractInfo
	err := json.Unmarshal(data, (*cloneInfo)(me))
	if err != nil {
		fmt.Println("cloneInfo error: ", err)
	}
	return nil
}
