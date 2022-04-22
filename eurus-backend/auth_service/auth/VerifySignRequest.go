package auth

type VerifySignRequest struct {
	ServiceId int64  `json:"serviceId"`
	Data      []byte `json:"data"`
	Sign      string `json:"sign"`
}

func (me *VerifySignRequest) MethodName() string {
	return "verifySign"
}
