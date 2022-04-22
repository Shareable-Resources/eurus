package conf_api

type FaucetConfig struct {
	Key          string `json:"key"`
	Amount       uint   `json:"amount"`
	Decimal      uint   `json:"decimal"`
	DateInterval uint   `json:"dateInterval"`
}
