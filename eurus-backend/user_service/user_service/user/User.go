package user

import (
	"eurus-backend/foundation/database"
	"time"
)

type UserId struct {
	NextVal uint64
}
type UserList struct {
	User []*User
}

type User struct {
	Id                         uint64
	LoginAddress               string
	WalletAddress              string
	OwnerWalletAddress         string
	MainnetWalletAddress       string
	Email                      string
	KycLevel                   int16
	CreatedDate                time.Time
	LastModifiedDate           time.Time
	IsMetamaskAddr             bool
	LastLoginTime              time.Time
	Status                     UserStatus
	Mnemonic                   string
	ChangeLoginPasswordCount   int
	ChangePaymentPasswordCount int
}

type Verification struct {
	UserId           uint64
	Code             string
	Type             int
	CreatedDate      time.Time
	LastModifiedDate time.Time
	ExpiredTime      time.Time
	Count            int
}

type LoginRequestTokenMap struct {
	database.DbModel
	Id                *uint64
	LoginRequestToken string
	Token             *string
	ExpiredTime       *time.Time
}

type UserPreferenceStorage struct {
	UserId   uint64
	Platform uint8
	Sequence uint64
	Storage  string
}

type UserDevice struct {
	database.DbModel
	CustomerId   uint64
	CustomerType int16
	DeviceId     string
	PubKey       string
}

type UserStatus int16
type VerificationType int

const (
	UserStatusUnknown UserStatus = iota
	UserStatusNormal
	UserStatusSuspended
	UserStatusVerifiedNotSetPaymentAddress
	UserStatusNotVerify
	UserStatusDeleted
)

const (
	VerificationRegistration VerificationType = iota
	VerificationForgetLoginPassword
	VerificationForgetPaymentPassword
	VerificationRegisterDevice
)

type UserType int16

const (
	DecentralizedUser UserType = iota
	CentralizedUser
)

type CustomerType int16

const (
	CustomerUser CustomerType = iota
	CustomerMerchant
)
