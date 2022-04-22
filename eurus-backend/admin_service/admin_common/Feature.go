package admin_common

//DB model base class
type Feature struct {
	Id                  uint64               `json:"id" gorm:"primaryKey"`
	Name                string               `json:"name"`
	SubFeatureList      []*Feature           `json:"subFeature" gorm:"-"`
	AvailablePermission []*FeaturePermission `json:"availablePermission" gorm:"many2many:admin_feature_permission_relations;foreignKey:Id;joinForeignKey:featureId;References:Id;joinReferences:PermissionId"`
	ParentFeatureId     uint64               `json:"-"`
	IsEnabled           bool                 `json:"isEnabled"`
}
