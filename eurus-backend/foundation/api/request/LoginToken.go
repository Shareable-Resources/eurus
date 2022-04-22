package request

import "time"

type LoginToken struct {
	UserInfo   string
	Token      string
	ExpiryTime time.Time
}
