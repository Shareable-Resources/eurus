package test

import (
	"eurus-backend/user_service/user_service/user"
	"fmt"
	"testing"
	"time"
)

func TestCheckTimestamp(t *testing.T) {
	importWalletRequest := user.NewImportWalletRequest()
	nowTime := time.Now().UnixNano()
	fmt.Println("current time: ", nowTime)
	importWalletRequest.Timestamp = nowTime
	err := importWalletRequest.CheckTimestamp()
	if err != nil {
		t.Errorf("Invalid Checking Error: %v", err.Error())
	}
}
