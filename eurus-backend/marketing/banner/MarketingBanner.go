package banner

import "time"

type MarketingBanner struct {
	Id            uint64                `json:"id" gorm:"primaryKey"`
	Position      uint64                `json:"position"`
	Seq           uint64                `json:"seq"`
	BannerType    uint64                `json:"bannerType"`
	IconUrlMobile string                `json:"iconUrlMobile"`
	LinkMobile    string                `json:"linkMobile"`
	Status        MarketingBannerStatus `json:"status"`
	StartDate     *time.Time            `json:"startDate"`
	EndDate       *time.Time            `json:"endDate"`
	IsDefault     bool                  `json:"isDefault"`
}

type MarketingBannerStatus int16

const (
	BannerDisabled MarketingBannerStatus = iota
	BannerEnabled
)
