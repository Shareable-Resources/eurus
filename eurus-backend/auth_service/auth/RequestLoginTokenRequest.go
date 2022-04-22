package auth

type RequestLoginTokenRequest struct {
	UserId          string `json:"userId"`
	ExpiredDuration int64  `json:"expiredDuration"`
}

func (me *RequestLoginTokenRequest) MethodName() string {
	return "requestLoginToken"
}
