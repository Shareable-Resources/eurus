package sc_validator

import (
	"encoding/json"
)

type SmartContractDictionary struct {
	SmartContractMap map[string]*SmartContractInfo `json:"smartContract"`
}

func (me *SmartContractDictionary) UnmarshalJSON(data []byte) error {
	if me.SmartContractMap == nil {
		me.SmartContractMap = make(map[string]*SmartContractInfo)
	}

	type cloneType SmartContractDictionary
	err := json.Unmarshal(data, (*cloneType)(me))
	if err != nil {
		var strVar string
		return json.Unmarshal(data, &strVar)
	}

	return nil
}
