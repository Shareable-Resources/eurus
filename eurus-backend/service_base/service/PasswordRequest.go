package service

type PasswordRequest struct {
	Password    string `json:"password"`
	IsEncrypted bool   `json:"isEncrypted"`
}

func (me *PasswordRequest) MethodName() string {
	return "passwordRequest"
}
