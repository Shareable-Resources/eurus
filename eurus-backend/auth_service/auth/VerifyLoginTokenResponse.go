package auth

type VerifyLoginTokenResponse struct { //implements network.ILoginToken
	ExpiredTime      int64  `json:"expiredTime"`
	CreatedDate      int64  `json:"createdDate"`
	LastModifiedDate int64  `json:"lastModifiedDate"`
	UserId           string `json:"userId"`
	token            string
	Type 		int16 `json:"type"`
}

func (me *VerifyLoginTokenResponse) GetToken() string {
	return me.token
}

func (me *VerifyLoginTokenResponse) GetCreatedDate() int64 {
	return me.CreatedDate
}

func (me *VerifyLoginTokenResponse) GetExpiredTime() int64 {
	return me.ExpiredTime
}

func (me *VerifyLoginTokenResponse) GetLastModifiedDate() int64 {
	return me.LastModifiedDate
}

func (me *VerifyLoginTokenResponse) GetUserId() string {
	return me.UserId
}

func (me *VerifyLoginTokenResponse) SetUserId(userId string) {
	me.UserId = userId
}

func (me *VerifyLoginTokenResponse) SetToken(token string) {
	me.token = token
}

func (me *VerifyLoginTokenResponse) GetTokenType() int16 {
	return me.Type
}