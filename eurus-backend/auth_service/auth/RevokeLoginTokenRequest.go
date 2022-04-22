package auth

type RevokeLoginTokenRequest struct {
	Token string `json:"token"`
}

func (me *RevokeLoginTokenRequest) MethodName() string {
	return "revokeLoginToken"
}