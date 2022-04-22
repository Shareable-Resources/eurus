package kyc_model

import (
	"eurus-backend/foundation/network"
	"reflect"
)

type KYCSubmitImageTask struct {
	network.MQEvent
	Data *KYCSubmitImage `json:"data"`
}

type KYCSubmitImage struct {
	KYCStatusId   uint64       `json:"userKYCStatusId"`
	UserId        uint64       `json:"userId"`
	ImageType     KYCImageType `json:"imageType"`
	FileExtension string       `json:"fileExtension"`
}

func NewKycSubmitImageTask() *KYCSubmitImageTask {
	task := new(KYCSubmitImageTask)
	task.EventType = reflect.TypeOf(task).Name()
	task.Data = new(KYCSubmitImage)
	return task
}
