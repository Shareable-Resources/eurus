package kyc_model

import (
	"database/sql"
	"eurus-backend/foundation/database"
)

type KYCStatusType int

const (
	KYCStatusPending            KYCStatusType = iota //kyc status just created
	KYCStatusWaitingForApproval                      //kyc status is submitted for approval
	KYCStatusWaitingForResubmit                      //kyc status is rejected and waiting to resubmit again
	KYCStatusApproved                                //kyc status is approved
	KYCStatusRejected                                //kyc status is rejected and no longer can submit (need reset)
)

type UserKYCStatus struct {
	database.DbModel
	Id             uint64              `gorm:"column:id;autoIncrement" json:"id"`
	UserId         uint64              `gorm:"column:user_id" json:"userId"`
	KYCLevel       int                 `gorm:"column:kyc_level" json:"kycLevel"`
	ApprovalDate   database.NullTime   `gorm:"column:approval_date" json:"approvalDate"`
	OperatorId     database.NullString `gorm:"column:operator_id" json:"operatorId"`
	KYCRetryCount  int                 `gorm:"column:kyc_retry_count" json:"kycRetryCount"`
	KYCCountryCode string              `gorm:"column:kyc_country_code" json:"kycCountryCode"`
	KYCStatus      KYCStatusType       `gorm:"column:kyc_status" json:"kycStatus"`
	KYCDoc         KYCDocType          `gorm:"column:kyc_doc" json:"kycDoc"`
	Images         []*UserKYCImage     `gorm:"foreignKey:UserKYCId" json:"images"`
}

func (t UserKYCStatus) TableName() string {
	return "user_kyc_statuses"
}

func NewUserKYCStatus(id string,
	userId string,
	KYCLevel int,
	approvalDate database.NullTime,
	operatorId database.NullString,
	KYCRetryCount int,
	KYCCountryCode string,
	KYCStatus KYCStatusType,
	KYCDoc KYCDocType) *UserKYCStatus {
	obj := new(UserKYCStatus)
	obj.KYCLevel = KYCLevel
	obj.ApprovalDate = database.NullTime{sql.NullTime{Valid: false}}
	obj.OperatorId = database.NullString{sql.NullString{Valid: false}}
	obj.KYCRetryCount = KYCRetryCount
	obj.KYCCountryCode = KYCCountryCode
	obj.KYCStatus = KYCStatus
	obj.KYCDoc = KYCDoc
	obj.InitDate()
	return obj
}
