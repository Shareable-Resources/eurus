package auth

type RefreshLoginTokenResponse struct { //implements network.ILoginToken
	ExpiredTime      int64  `json:"expiredTime"`
	CreatedDate      int64  `json:"createdDate"`
	LastModifiedDate int64  `json:"lastModifiedDate"`
	UserId           string `json:"userId"`
	token            string
	Type 		int16 `json:"type"`
}

func (me *RefreshLoginTokenResponse) GetToken() string {
	return me.token
}

func (me *RefreshLoginTokenResponse) GetExpiredTime() int64 {
	return me.ExpiredTime
}

func (me *RefreshLoginTokenResponse) GetCreatedDate() int64 {
	return me.CreatedDate
}

func (me *RefreshLoginTokenResponse) GetLastModifiedDate() int64 {
	return me.LastModifiedDate
}

func (me *RefreshLoginTokenResponse) GetUserId() string {
	return me.UserId
}

func (me *RefreshLoginTokenResponse) SetUserId(userId string) {
	me.UserId = userId
}

func (me *RefreshLoginTokenResponse) SetToken(token string) {
	me.token = token
}

func (me *RefreshLoginTokenResponse) GetTokenType() int16 {
	return me.Type
}