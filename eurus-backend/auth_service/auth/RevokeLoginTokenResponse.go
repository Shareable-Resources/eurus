package auth

type RevokeLoginTokenResponse struct {
	ExpiredTime int64  `json:"expiredTime"`
	UserId      string `json:"userId"`
	token            string
	Type 		int16 `json:"type"`
	CreatedDate      int64  `json:"createdDate"`
	LastModifiedDate int64  `json:"lastModifiedDate"`
}

func (me *RevokeLoginTokenResponse) GetToken() string {
	return me.token
}

func (me *RevokeLoginTokenResponse) GetExpiredTime() int64 {
	return me.ExpiredTime
}

func (me *RevokeLoginTokenResponse) GetCreatedDate() int64 {
	return me.CreatedDate
}

func (me *RevokeLoginTokenResponse) GetLastModifiedDate() int64 {
	return me.LastModifiedDate
}

func (me *RevokeLoginTokenResponse) GetUserId() string {
	return me.UserId
}

func (me *RevokeLoginTokenResponse) SetUserId(userId string) {
	me.UserId = userId
}

func (me *RevokeLoginTokenResponse) SetToken(token string) {
	me.token = token
}

func (me *RevokeLoginTokenResponse) GetTokenType() int16 {
	return me.Type
}