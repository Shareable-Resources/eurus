package auth

type VerifyLoginTokenRequest struct {
	Token      string `json:"token"`
	VerifyMode int    `json:"verifyMode"`
}

func (me *VerifyLoginTokenRequest) MethodName() string {
	return "verifyLoginToken"
}
