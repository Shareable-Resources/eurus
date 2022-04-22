package admin_common

import "strings"

type RoleState int16

const (
	RoleDisabled RoleState = iota
	RoleEnabled
	RoleDeleted
	RoleAll RoleState = 255
)

func (me *RoleState) UnmarshalJSON(b []byte) error {
	str := string(b)
	if str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}
	str = strings.Trim(str, " ")
	switch str {
	case "0":
		*me = RoleDisabled
	case "1":
		*me = RoleEnabled
	default:
		*me = RoleAll
	}
	return nil
}
