package auth

import (
	"eurus-backend/foundation/database"
	"time"
)

type UserSession struct {
	database.DbModel
	Token       string
	ServiceId   int64
	ExpiredTime *time.Time
	UserId      string
	Disabled    bool
	Type        int16
}

type LoginRequestTokenMap struct {
	database.DbModel
	LoginRequestToken string
	Token             *string
	ExpiredTime       *time.Time
}
