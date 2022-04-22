package auth

type RequestLoginTokenResponse struct { //implements network.ILoginToken
	Token            string `json:"token"`
	ExpiredTime      int64  `json:"expiredTime"`
	LastModifiedDate int64  `json:"lastModifiedDate"`
	CreatedDate int64 `json:"createdDate"`
	userId           string
	Type 		int16 `json:"type"`
}

func (me *RequestLoginTokenResponse) GetToken() string {
	return me.Token
}

func (me *RequestLoginTokenResponse) GetExpiredTime() int64 {
	return me.ExpiredTime
}

func (me *RequestLoginTokenResponse) GetCreatedDate() int64 {
	return me.CreatedDate
}

func (me *RequestLoginTokenResponse) GetLastModifiedDate() int64 {
	return me.LastModifiedDate
}

func (me *RequestLoginTokenResponse) GetUserId() string {
	return me.userId
}

func (me *RequestLoginTokenResponse) SetUserId(userId string) {
	me.userId = userId
}

func (me *RequestLoginTokenResponse) SetToken(token string) {
	me.Token = token
}

func (me *RequestLoginTokenResponse) GetTokenType() int16 {
	return me.Type
}
