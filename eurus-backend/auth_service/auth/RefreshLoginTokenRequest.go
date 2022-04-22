package auth

type RefreshLoginTokenRequest struct {
	Token string `json:"token"`
}

func (me *RefreshLoginTokenRequest) MethodName() string {
	return "refreshLoginToken"
}
