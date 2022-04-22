package auth

type NonRefreshableLoginTokenRequest struct {
	UserId    string    `json:"userId"`
	Token     string    `json:"token"`
	Duration  int64     `json:"duration"`
	TokenType TokenType `json:"tokenType"`
}

func (me *NonRefreshableLoginTokenRequest) MethodName() string {
	return "requestPaymentLoginToken"
}
