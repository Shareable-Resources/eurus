package banner

type MarketingBannerDetail struct {
	BannerId        uint64 `json:"bannerId"`
	LangCode        string `json:"langCode"`
	BannerUrlMobile string `json:"bannerUrlMobile"`
	Content         string `json:"content"`
}

type LangCode string

const (
	English            LangCode = "en"
	TraditionalChinese LangCode = "zh"
	SimplifiedChinese  LangCode = "cn"
)
