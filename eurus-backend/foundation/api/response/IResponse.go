package response

type IResponse interface {
	GetReturnCode() int64
	GetMessage() string
	GetNonce() string
	GetData() interface{}
	IsInterfaceNil() bool
}
