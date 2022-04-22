package kyc_model

import (
	"eurus-backend/foundation/database"
)

type KYCImageStatusType int

const (
	KYCImageStatusReceived           KYCImageStatusType = iota // User Server Received the image
	KYCImageStatusUploaded                                     // the image is uploaded to S3
	KYCImageStatusWaitingForResubmit                           // the image is rejected and waiting for resubmit
	KYCImageStatusApproved                                     // the image is approved by CS admin
	KYCImageStatusVoided                                       // the image is no longer being used (replace by new image)
)

type UserKYCImage struct {
	database.DbModel
	UserKYCId    uint64              `gorm:"column:user_kyc_id" json:"userKYCId"`
	DocType      KYCImageType        `gorm:"column:doc_type" json:"docType"`
	ImageSeq     int                 `gorm:"column:image_seq" json:"imageSeq"`
	Status       KYCImageStatusType  `gorm:"column:status" json:"status"`
	ImagePath    string              `gorm:"column:image_path" json:"imagePath"`
	RejectReason string              `gorm:"column:reject_reason" json:"rejectReason"`
	OperatorId   database.NullString `gorm:"column:operator_id" json:"operatorId"`
}

func (t UserKYCImage) TableName() string {
	return "user_kyc_images"
}

func NewUserKYCImage(
	userKYCId uint64,
	imageType KYCImageType,
	imageSeq int,
	status KYCImageStatusType,
	imagePath string,
	rejectReason string,
	operatorId database.NullString,

) *UserKYCImage {
	obj := new(UserKYCImage)
	obj.UserKYCId = userKYCId
	obj.DocType = imageType
	obj.ImageSeq = imageSeq
	obj.Status = status
	obj.ImagePath = imagePath
	obj.RejectReason = rejectReason
	obj.OperatorId = operatorId
	obj.InitDate()
	return obj
}
