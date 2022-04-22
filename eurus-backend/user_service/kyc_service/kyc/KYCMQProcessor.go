package kyc

import "eurus-backend/foundation/log"

type KYCMQProcessor struct {
}

func failOnError(err error, msg string) {
	if err != nil {
		log.GetLogger(log.Name.Root).Errorf("%s: %s\r\n", msg, err)
	}
}
