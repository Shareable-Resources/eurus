package foundation

import "time"

type IRetrySetting interface {
	GetRetryCount() int
	GetRetryInterval() time.Duration
}
