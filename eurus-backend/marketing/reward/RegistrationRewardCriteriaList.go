package reward

import (
	"encoding/json"
	"math/big"
	"strings"
)

type RegistrationRewardCriteriaList struct {
	CenteralizeUserSetting   *RegistrationUserRewardSetting `json:"centralizedUser"`
	DecentralizedUserSetting *RegistrationUserRewardSetting `json:"decentralizedUser"`
}

type LowerCaseStringList []string

func (me *LowerCaseStringList) UnmarshalJSON(data []byte) error {
	var addrList []string = make([]string, 0)
	err := json.Unmarshal(data, &addrList)
	if err != nil {
		return err
	}
	for _, addr := range addrList {
		*me = append(*me, strings.ToLower(addr))
	}
	return nil
}

type RegistrationUserRewardSetting struct {
	ExcludeSenderList LowerCaseStringList                    `json:"excludeSenderList"`
	ExcludeSenderMap  map[string]string                      `json:"-"`
	Reward            *RegistrationReward                    `json:"reward"`
	Criteria          map[string]*RegistrationRewardCriteria `json:"criteria"` //Asset name to criteria mapping
}

func (me *RegistrationUserRewardSetting) UnmarshalJSON(data []byte) error {
	type cloneType RegistrationUserRewardSetting
	err := json.Unmarshal(data, (*cloneType)(me))
	if err != nil {
		return err
	}
	if me.ExcludeSenderList != nil {
		me.ExcludeSenderMap = make(map[string]string)
		for _, addr := range me.ExcludeSenderList {
			me.ExcludeSenderMap[addr] = addr
		}
	}
	return nil
}

type RegistrationReward struct {
	AssetName string   `json:"assetName"`
	Amount    *big.Int `json:"amount"`
}

type RegistrationRewardCriteria struct {
	AssetName       string          `json:"assetName"`
	CompareCriteria CompareOperator `json:"compareCriteria"` // supported operator >, >= , < , <= , ==
	CompareAmount   *big.Int        `json:"compareAmount"`
}

type CompareOperator string

func (me *CompareOperator) Compare(left *big.Int, right *big.Int) bool {
	switch strings.Trim(string(*me), " ") {
	case "==":
		return left.Cmp(right) == 0
	case ">":
		return left.Cmp(right) > 0
	case "<":
		return left.Cmp(right) < 0
	case ">=":
		return left.Cmp(right) >= 0
	case "<=":
		return left.Cmp(right) <= 0
	default:
		return false

	}
}
