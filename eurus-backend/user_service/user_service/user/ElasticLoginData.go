package user

import (
	"eurus-backend/foundation/elastic"
)

type elasticLoginData struct {
	elastic.ElasticSearchDataBase
	UserId             uint64 `json:"userId"`
	Ip                 string `json:"ip"`
	AppVersion         string `json:"appVersion"`
	Os                 string `json:"os"`
	RegistrationSource string `json:"registrationSource"`
}

func newElasticLoginData() *elasticLoginData {
	data := new(elasticLoginData)
	data.Path = "login"
	data.Index = "loginData"

	return data
}
