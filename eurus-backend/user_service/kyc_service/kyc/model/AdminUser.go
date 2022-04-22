package kyc_model

type AdminUser struct {
	Username string `gorm:"column:username" json:"username"`
	Password string `gorm:"column:password" json:"password"`
}

func (t AdminUser) TableName() string {
	return "admin_users"
}
