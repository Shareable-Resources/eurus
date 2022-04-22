package foundation

type IServerError interface {
	error
	GetReturnCode() ServerReturnCode
}
type ServerError struct {
	ReturnCode ServerReturnCode
	Message    string
}

func NewError(code ServerReturnCode) *ServerError {
	return NewErrorWithMessage(code, code.String())
}

func NewErrorWithMessage(code ServerReturnCode, message string) *ServerError {
	err := new(ServerError)
	err.ReturnCode = code
	err.Message = message

	return err
}

func (me *ServerError) Error() string {
	return me.Message
}

func (me *ServerError) GetReturnCode() ServerReturnCode {
	return me.ReturnCode
}
