package test

import (
	"eurus-backend/foundation/elastic"
	"fmt"
	"testing"
	"time"
)

func TestInsertDataToElastic (t *testing.T) {
	path := "http://18.166.159.174:9200"
	data := &elastic.LoginData{UserId: 12456, Ip : "127.0.0.1", AppVersion : "V1.1111", Os : "osos",
		Timestamp : time.Now(), RegistrationSource : "testing by ryan"}
	err := elastic.InsertLoginLog(path,data)
	if err!= nil {
		fmt.Println(err)
	}
}