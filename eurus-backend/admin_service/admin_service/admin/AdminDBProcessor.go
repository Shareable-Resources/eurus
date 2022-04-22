package admin

import (
	"encoding/base64"
	"errors"
	"eurus-backend/admin_service/admin_common"
	"eurus-backend/foundation/crypto"
	"eurus-backend/foundation/database"
	"time"

	"gorm.io/gorm"
)

type AdminDBProcessor struct {
	db     *database.Database
	config *AdminServerConfig
}

func NewAdminDBProcessor(config *AdminServerConfig, db *database.Database) *AdminDBProcessor {
	adminDBProcessor := new(AdminDBProcessor)
	adminDBProcessor.config = config
	adminDBProcessor.db = db
	return adminDBProcessor
}

func (me *AdminDBProcessor) DbVerifyAdminPassword(userName string, password string) (*AdminUser, error) {

	conn, err := me.db.GetConn()
	if err != nil {
		return nil, err
	}
	adminUser := new(AdminUser)
	tx := conn.Model(adminUser).Where("username = ?", userName).FirstOrInit(&adminUser)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, errors.New("User not found")
	}
	hash, err := base64.StdEncoding.DecodeString(adminUser.Password)
	if err != nil {
		return nil, err
	}
	err = crypto.VerifyBCryptHash([]byte(password), []byte(reverse(userName)), []byte(hash))
	if err != nil {
		return nil, errors.New("Invalid password")
	}

	return adminUser, nil
}

func (me *AdminDBProcessor) DbUpdateAdminSecret(userId uint64, secret string) error {
	conn, err := me.db.GetConn()
	if err != nil {
		return err
	}

	tx := conn.Model(AdminUser{}).Where("id = ?", userId).Updates(map[string]interface{}{
		"secret": secret,
		"status": AdminWaitForBindGA})
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("User not found")
	}

	return nil
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func (me *AdminDBProcessor) DbGetAdminUserById(userId uint64) (*AdminUser, error) {
	conn, err := me.db.GetConn()
	if err != nil {
		return nil, err
	}

	adminUser := new(AdminUser)

	tx := conn.Model(adminUser).Where("id = ?", userId).First(&adminUser)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return adminUser, nil
}

func (me *AdminDBProcessor) DbUpdateAdminUserStatus(userId uint64, status AdminUserStatus) error {
	conn, err := me.db.GetConn()
	if err != nil {
		return err
	}

	tx := conn.Model(AdminUser{}).Where("id = ?", userId).Updates(map[string]interface{}{
		"status": status,
	})
	return tx.Error
}

func (me *AdminDBProcessor) DbUpdateAdminUserLoginInfo(userId uint64, ip string, loginTime time.Time) error {
	conn, err := me.db.GetConn()
	if err != nil {
		return err
	}

	tx := conn.Model(AdminUser{}).Where("id = ?", userId).Updates(map[string]interface{}{
		"login_ip":        ip,
		"last_login_time": loginTime,
	})
	return tx.Error
}

func (me *AdminDBProcessor) DbQueryAllFeaturePermission() ([]*AdminFeature, error) {
	conn, err := me.db.GetConn()
	if err != nil {
		return nil, err
	}
	var rootFeatureList []*AdminFeature = make([]*AdminFeature, 0)
	var featureMap map[uint64]*AdminFeature = make(map[uint64]*AdminFeature)

	tx := conn.Session(&gorm.Session{}).Model(AdminFeature{}).Where("is_enabled = ? AND parent_feature_id IS NULL", true).Order("parent_feature_id").Find(&rootFeatureList)
	if tx.Error != nil {
		return nil, tx.Error
	}

	for _, feature := range rootFeatureList {
		featureMap[feature.Id] = feature
	}

	var subFeatureList []*AdminFeature = make([]*AdminFeature, 0)
	tx = conn.Session(&gorm.Session{}).Model(AdminFeature{}).Where("is_enabled = ? AND parent_feature_id IS NOT NULL", true).Order("parent_feature_id").Find(&subFeatureList)
	if tx.Error != nil {
		return nil, tx.Error
	}

	var featureIdList []uint64 = make([]uint64, 0)

	for _, feature := range subFeatureList {
		featureMap[feature.Id] = feature
		featureIdList = append(featureIdList, feature.Id)
	}

	for _, feature := range subFeatureList {
		parent, ok := featureMap[feature.ParentFeatureId]
		if ok {
			if parent.SubFeatureList == nil {
				parent.SubFeatureList = make([]*admin_common.Feature, 0)
			}
			parent.SubFeatureList = append(parent.SubFeatureList, (*admin_common.Feature)(feature))
		}
	}
	var permissionList []*AdminFeaturePermission = make([]*AdminFeaturePermission, 0)
	tx = conn.Session(&gorm.Session{}).Raw("SELECT r.feature_id, p.id, p.name FROM admin_feature_permission_relations as r INNER JOIN admin_feature_permissions as p ON r.permission_id = p.id WHERE r.feature_id IN (?)", featureIdList).Order("r.feature_id, p.id").Find(&permissionList)
	if tx.Error != nil {
		return nil, tx.Error
	}

	for _, permission := range permissionList {
		feature := featureMap[permission.FeatureId]
		if feature.AvailablePermission == nil {
			feature.AvailablePermission = make([]*admin_common.FeaturePermission, 0)
		}
		feature.AvailablePermission = append(feature.AvailablePermission, (*admin_common.FeaturePermission)(permission))
	}

	return rootFeatureList, nil

}

func (me *AdminDBProcessor) DbQueryAdminEffectivePermission(adminUserId uint64) ([]*admin_common.FeaturePermissionPair, error) {
	conn, err := me.db.GetConn()
	if err != nil {
		return nil, err
	}

	var permissionList []*admin_common.FeaturePermissionPair = make([]*admin_common.FeaturePermissionPair, 0)
	tx := conn.Session(&gorm.Session{}).Raw(
		`SELECT DISTINCT p.feature_id, p.permission_id 
			FROM admin_role_permissions as p 
			INNER JOIN admin_user_role_relations as r on r.role_id = p.role_id 
			INNER JOIN admin_roles as k ON r.role_id = k.id 
			WHERE r.admin_id = ? AND k.state = ? ORDER BY p.feature_id, p.permission_id`, adminUserId, admin_common.RoleEnabled).Find(&permissionList)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return permissionList, nil
}

func (me *AdminDBProcessor) DbQueryAdminEffectiveRole(adminUserId uint64) ([]*AdminRole, error) {
	conn, err := me.db.GetConn()
	if err != nil {
		return nil, err
	}

	var roleList []*AdminRole = make([]*AdminRole, 0)

	tx := conn.Session(&gorm.Session{}).Raw(`SELECT r.* FROM admin_roles as r 
	INNER JOIN admin_user_role_relations as a ON a.role_id = r.id 
	WHERE a.admin_id = ? AND r.state = ?`, adminUserId, admin_common.RoleEnabled).Find(&roleList)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return roleList, nil
}

func (me *AdminDBProcessor) DbAdminHasPermission(adminId uint64, featureId FeatureId, permissionId PermissionId) (bool, error) {
	conn, err := me.db.GetConn()
	if err != nil {
		return false, err
	}
	var found int
	tx := conn.Raw(`SELECT 1 FROM admin_user_role_relations as a 
	INNER JOIN admin_users as u ON a.admin_id = u.id
	INNER JOIN admin_roles as r ON a.role_id = r.id 
	INNER JOIN admin_role_permissions as p ON r.id = p.role_id
	INNER JOIN admin_feature_permission_relations AS fp ON fp.feature_id = p.feature_id AND fp.permission_id = p.permission_id
	INNER JOIN admin_features AS f ON f.id = fp.feature_id
	WHERE r.state = ?
	AND f.is_enabled = ?
	AND a.admin_id = ?
	AND u.status = ?
	AND fp.feature_id = ?
	AND fp.permission_id = ?`, admin_common.RoleEnabled, true, adminId, AdminNormal, featureId, permissionId).Find(&found)

	if tx.Error != nil {
		return false, tx.Error
	}

	if found == 1 {
		return true, nil
	}
	return false, nil

}
