package password_api

const PasswordRequestMethodName = "passwordRequest"

type PasswordRequest struct {
}

func (me *PasswordRequest) MethodName() string {
	return PasswordRequestMethodName
}

type PasswordResponse struct {
	EncryptedPassword string `json:"encryptedPassword"`
}
