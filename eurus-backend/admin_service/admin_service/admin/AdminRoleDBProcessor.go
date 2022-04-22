package admin

import (
	"eurus-backend/admin_service/admin_common"
	"eurus-backend/foundation/database"
	"eurus-backend/foundation/log"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type AdminRoleDBProcessor struct {
	AdminDBProcessor
	config *AdminServerConfig
}

func NewAdminRoleDBProcessor(config *AdminServerConfig, db *database.Database) *AdminRoleDBProcessor {
	dbProcessor := new(AdminRoleDBProcessor)
	dbProcessor.config = config
	dbProcessor.db = db
	return dbProcessor
}

func (me *AdminRoleDBProcessor) DbQueryAdminRoleList(roleName string, roleState admin_common.RoleState) ([]*AdminRoleEx, error) {
	conn, err := me.db.GetConn()
	if err != nil {
		return nil, err
	}
	var roleList []*AdminRoleEx = make([]*AdminRoleEx, 0)
	sess := conn.Session(&gorm.Session{}).Model(AdminRole{}).Select("admin_roles.*, u.username")
	sess = sess.Joins("LEFT JOIN admin_users as u ON admin_roles.modified_by = u.id")
	if roleName != "" {
		sess = sess.Where(" admin_roles.role_name = ?", roleName)
	}

	if roleState != admin_common.RoleAll {
		sess = sess.Where(" admin_roles.state = ?", roleState)
	}

	tx := sess.Find(&roleList)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return roleList, nil
}

func (me *AdminRoleDBProcessor) DbQueryAdminRolePermission(roleId uint64) ([]*AdminRolePermission, error) {
	conn, err := me.db.GetConn()
	if err != nil {
		return nil, err
	}

	var rolePermissionList []*AdminRolePermission = make([]*AdminRolePermission, 0)
	tx := conn.Session(&gorm.Session{}).Where("role_id = ?", roleId).Find(&rolePermissionList)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return rolePermissionList, nil
}

func (me *AdminRoleDBProcessor) DbQueryRoleDetail(roleId uint64) (*AdminRole, []*admin_common.FeaturePermissionPair, error) {

	conn, err := me.db.GetConn()
	if err != nil {
		return nil, nil, err
	}
	var featurePermissionList []*admin_common.FeaturePermissionPair = make([]*admin_common.FeaturePermissionPair, 0)
	tx := conn.Session(&gorm.Session{}).Model(AdminRole{}).
		Select("p.feature_id", "p.permission_id").
		Joins("INNER JOIN admin_role_permissions as p ON p.role_id = admin_roles.id").
		Joins("INNER JOIN admin_features as f ON f.id = p.feature_id").
		Joins("INNER JOIN admin_feature_permission_relations as fp ON f.id = fp.feature_id AND p.permission_id = fp.permission_id").
		Where("admin_roles.id = ?", roleId).
		Find(&featurePermissionList)
	if tx.Error != nil {
		return nil, nil, tx.Error
	}

	role := new(AdminRole)
	role.Id = roleId
	tx = conn.Session(&gorm.Session{}).Find(&role)
	if tx.Error != nil {
		return nil, nil, tx.Error
	}

	return role, featurePermissionList, nil
}

func (me *AdminRoleDBProcessor) DbCreateRole(roleName string, desc string, adminId uint64, permission []*admin_common.FeaturePermissionPair) (uint64, error) {
	conn, err := me.db.GetConn()
	if err != nil {
		return 0, err
	}
	role := new(AdminRole)

	err = conn.Transaction(func(db *gorm.DB) error {
		role.RoleName = roleName
		role.Description = desc
		role.ModifiedBy = adminId
		role.State = admin_common.RoleEnabled
		role.InitDate()

		tx := db.Session(&gorm.Session{}).Create(&role)
		if tx.Error != nil {
			log.GetLogger(log.Name.Root).Errorln("Unable to create role: ", err, " role name: ", roleName)
			return errors.Wrap(tx.Error, "Unable to create role")
		}

		if role.Id == 0 {
			log.GetLogger(log.Name.Root).Errorln("Unable to generate role Id. Role name: ", roleName)
			return errors.New("Unable to generate role Id. Role name: " + roleName)
		}

		for _, item := range permission {
			rolePermission := new(AdminRolePermission)
			rolePermission.RoleId = role.Id
			rolePermission.PermissionId = item.PermissionId
			rolePermission.FeatureId = item.FeatureId

			var permissionList []*admin_common.FeaturePermissionPair = make([]*admin_common.FeaturePermissionPair, 0)
			tx2 := db.Session(&gorm.Session{}).Table("admin_feature_permission_relations").
				Where("feature_id = ? AND permission_id = ?", rolePermission.FeatureId, rolePermission.PermissionId).
				Find(&permissionList)
			if tx2.Error != nil {
				log.GetLogger(log.Name.Root).Errorln("Unable to query feature permission: ", tx2.Error, " permission id: ", item.PermissionId, " feature id: ", item.FeatureId, " role name: ", roleName)
				return errors.Wrap(tx2.Error, "Unable to query feature permission. Permission id: "+strconv.FormatUint(item.PermissionId, 10)+" feature id: "+strconv.FormatUint(item.FeatureId, 10))
			}

			if len(permissionList) == 0 {
				log.GetLogger(log.Name.Root).Errorln("Feature permission not found: ", tx2.Error, " permission id: ", item.PermissionId, " feature id: ", item.FeatureId, " role name: ", roleName)
				return errors.New("Feature permission not found: " + strconv.FormatUint(item.PermissionId, 10) + " feature id: " + strconv.FormatUint(item.FeatureId, 10))
			}

			rolePermission.InitDate()
			tx1 := db.Session(&gorm.Session{}).Create(&rolePermission)
			if tx1.Error != nil {
				log.GetLogger(log.Name.Root).Errorln("Unable to insert role permission: ", tx1.Error, " permission id: ", item.PermissionId, " feature id: ", item.FeatureId, " role name: ", roleName)
				return errors.Wrap(tx1.Error, "Unable to insert role permission. Permission id: "+strconv.FormatUint(item.PermissionId, 10)+" feature id: "+strconv.FormatUint(item.FeatureId, 10))
			}
		}
		return nil
	})

	return role.Id, err
}

func (me *AdminRoleDBProcessor) DbUpdateRole(roleId uint64, adminId uint64, updateField admin_common.UpdateRoleField, roleName string, description string, state admin_common.RoleState, permissionList []*admin_common.FeaturePermissionPair) error {
	conn, err := me.db.GetConn()
	if err != nil {
		return err
	}

	err = conn.Transaction(func(db *gorm.DB) error {
		role := new(AdminRole)
		role.Id = roleId
		tx := db.Session(&gorm.Session{}).Where("id = ?", role.Id).First(&role)
		if tx.Error != nil {
			log.GetLogger(log.Name.Root).Errorln("Admin role not found. role id: ", roleId)
			return errors.Wrap(tx.Error, "Query admin role failed")
		}
		var updateMap map[string]interface{} = make(map[string]interface{})

		if updateField&admin_common.RoleFieldName > 0 {
			updateMap["role_name"] = roleName
		}
		if updateField&admin_common.RoleFieldDescription > 0 {
			updateMap["description"] = description
		}
		if updateField&admin_common.RoleFieldState > 0 {
			updateMap["state"] = state
		}

		updateMap["last_modified_date"] = time.Now()
		updateMap["modified_by"] = adminId

		tx = db.Session(&gorm.Session{}).Model(AdminRole{}).
			Where("id = ?", role.Id).
			Updates(updateMap)

		if tx.Error != nil {
			log.GetLogger(log.Name.Root).Errorln("Admin role update failed. role id: ", roleId)
			return errors.Wrap(tx.Error, "Admin role update failed.")
		}

		if updateField&admin_common.RoleFieldPermission > 0 {
			tx1 := db.Session(&gorm.Session{}).Model(AdminRolePermission{}).Delete(AdminRolePermission{RoleId: roleId}, "role_id = ?", roleId)
			if tx1.Error != nil {
				log.GetLogger(log.Name.Root).Errorln("Delete exists role permission failed: ", tx1.Error)
				return errors.Wrap(tx1.Error, "Delete existing role permission failed")
			}

			for _, item := range permissionList {
				rolePermission := new(AdminRolePermission)
				rolePermission.RoleId = roleId
				rolePermission.FeatureId = item.FeatureId
				rolePermission.PermissionId = item.PermissionId

				tx2 := db.Session(&gorm.Session{}).Create(&rolePermission)
				if tx2.Error != nil {
					log.GetLogger(log.Name.Root).Errorln("Insert role permission failed: ", tx2.Error, " feature id: ", item.FeatureId, " permission id: ", item.PermissionId)
					msg := fmt.Sprintf("Insert role permission failed. Feature id: %d, permission id: %d", item.FeatureId, item.PermissionId)
					return errors.Wrap(tx2.Error, msg)
				}
			}
		}
		return nil
	})

	return err
}

func (me *AdminRoleDBProcessor) DbDeleteRole(roleId uint64, adminId uint64) error {
	conn, err := me.db.GetConn()
	if err != nil {
		return err
	}

	tx := conn.Model(AdminRole{}).Where("id = ? AND state <> ?", roleId, admin_common.RoleDeleted).Updates(map[string]interface{}{
		"state": admin_common.RoleDeleted,
	})

	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("Role Id not found")
	}
	return nil
}
