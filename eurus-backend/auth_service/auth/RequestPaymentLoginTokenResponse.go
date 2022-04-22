package auth

type NonRefreshableLoginTokenResponse struct { //implements network.ILoginToken
	ExpiredTime      int64     `json:"expiredTime"`
	CreatedDate      int64     `json:"createdDate"`
	LastModifiedDate int64     `json:"lastModifiedDate"`
	UserId           string    `json:"userId"`
	Type             TokenType `json:"type"`
	Token            string    `json:"token"`
}

func (me *NonRefreshableLoginTokenResponse) GetToken() string {
	return me.Token
}

func (me *NonRefreshableLoginTokenResponse) GetExpiredTime() int64 {
	return me.ExpiredTime
}

func (me *NonRefreshableLoginTokenResponse) GetCreatedDate() int64 {
	return me.CreatedDate
}

func (me *NonRefreshableLoginTokenResponse) GetLastModifiedDate() int64 {
	return me.LastModifiedDate
}

func (me *NonRefreshableLoginTokenResponse) GetUserId() string {
	return me.UserId
}

func (me *NonRefreshableLoginTokenResponse) SetUserId(userId string) {
	me.UserId = userId
}

func (me *NonRefreshableLoginTokenResponse) SetToken(token string) {
	me.Token = token
}

func (me *NonRefreshableLoginTokenResponse) GetTokenType() int16 {
	return int16(me.Type)
}
