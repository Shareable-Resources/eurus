package admin

import "time"

type ElasticAdminLoginData struct {
	AdminId   uint64    `json:"adminId"`
	UserName  string    `json:"userName"`
	LoginTime time.Time `json:"loginTime"`
	LoginIp   string    `json:"loginIp"`
}
