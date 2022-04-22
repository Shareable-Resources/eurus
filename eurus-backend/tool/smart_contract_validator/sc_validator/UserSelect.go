package sc_validator

type UserSelect struct{
	Network 	NetworkSelect
}
type NetworkSelect int

const (
	StatusError           NetworkSelect = iota
	MainNet
	SideChain
	AllNetwork
	Exit
)
