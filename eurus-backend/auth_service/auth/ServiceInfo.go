package auth

type ServiceInfo struct {
	ServiceId   int64  `json:"serviceId"`
	ServiceName string `json:"serviceName"`
	PublicKey   string `json:"pubKey"`
}
