package conf_api

type AuthService struct {
	Id             int    `json:"id"`
	ServiceName    string `json:"serviceName"`
	PubKey         string `json:"pubKey"`
	WalletAddress  string `json:"walletAddress"`
	ServiceGroupId int    `json:"serviceGroupId"`
}
