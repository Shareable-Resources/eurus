package admin_common

//DB model base class
type FeaturePermission struct {
	Id        uint64 `json:"id" gorm:"primaryKey"`
	Name      string `json:"name"`
	FeatureId uint64 `json:"-"`
}
