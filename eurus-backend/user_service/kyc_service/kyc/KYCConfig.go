package kyc

import (
	"eurus-backend/foundation/server"
)

// Extends ServerConfigBase struct, should add new attributes in here
type KYCConfig struct {
	server.ServerConfigBase
	S3BucketZone             string `json:"s3BucketZone"`
	AwsAccessKeyId           string `json:"awsAccessKeyId"`
	AwsAccessSecretAccessKey string `json:"awsAccessSecretAccessKey" eurus_conf:"noPrint"`
	RowLimit                 int    `json:"rowLimit"`
	AdminAESKey              string `json:"adminAESKey" eurus_conf:"noPrint"`
}

func NewKYCConfig() *KYCConfig {
	config := new(KYCConfig)
	config.ActualConfig = config
	return config
}

func (me *KYCConfig) GetServerConfigBase() *server.ServerConfigBase {
	return &me.ServerConfigBase
}

func (me *KYCConfig) GetParent() interface{} {
	return &me.ServerConfigBase
}

func (me *KYCConfig) GetConfigFileOnlyFieldList() []string {
	return []string{"adminAESKey"}
}
