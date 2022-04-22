package auth

type AuthenticateResponse struct {
	SessionId      int64  `json:"sessionId"`
	Token          string `json:"token"`
	AdditionalInfo string `json:"additionalInfo"`
}
