package kyc_model

type KYCCountryCode struct {
	Code          string `gorm:"column:code" json:"code"`
	Name          string `gorm:"column:name" json:"name"`
	FullName      string `gorm:"column:full_name" json:"fullName"`
	Iso3          string `gorm:"column:iso3" json:"iso3"`
	Number        string `gorm:"column:number" json:"number"`
	ContinentCode string `gorm:"column:continent_code" json:"continentCode"`
}

func (t KYCCountryCode) TableName() string {
	return "kyc_country_codes"
}

func NewKYCCountryCode(code string, name string, fullName string, iso3 string, number string, continentCode string) *KYCCountryCode {
	obj := new(KYCCountryCode)
	obj.Code = code
	obj.Name = name
	obj.FullName = fullName
	obj.Iso3 = iso3
	obj.Number = number
	obj.ContinentCode = continentCode
	return obj
}
