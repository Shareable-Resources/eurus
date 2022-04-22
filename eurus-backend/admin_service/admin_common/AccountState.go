package admin_common

type AccountState uint16

const (
	AccountDisabled AccountState = 0
	AccountEnabled  AccountState = 1
	AccountAll      AccountState = 255
)
