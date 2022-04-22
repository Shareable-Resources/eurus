package admin

import (
	"encoding/base64"
	"eurus-backend/admin_service/admin_common"
	"eurus-backend/foundation"
	"eurus-backend/foundation/crypto"
	"eurus-backend/foundation/database"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type AdminAccountDBProcessor struct {
	AdminDBProcessor
}

func NewAdminAccountDBProcessor(config *AdminServerConfig, db *database.Database) *AdminAccountDBProcessor {
	processor := new(AdminAccountDBProcessor)
	processor.config = config
	processor.db = db
	return processor
}

func (me *AdminAccountDBProcessor) DbQueryAccountList(userName string, roleName string, state admin_common.AccountState) ([]*AdminUserEx, error) {
	conn, err := me.db.GetConn()
	if err != nil {
		return nil, err
	}

	query := conn.Session(&gorm.Session{}).Model(AdminUser{}).
		Distinct("admin_users.*", "m.username as modified_by_user").
		Order(" admin_users.id").
		Joins(" LEFT JOIN admin_users as m ON admin_users.modified_by = m.id").
		Where("admin_users.status <> ?", AdminDeleted)

	if userName != "" {
		query = query.Where("admin_users.username = ?", userName)
	}
	if roleName != "" {
		query = query.Joins(" LEFT JOIN admin_user_role_relations as rel ON rel.admin_id = admin_users.id").
			Joins(" INNER JOIN admin_roles as r ON rel.role_id = r.id").
			Where("lower(r.role_name) = lower(?)", roleName)
	}
	if state != admin_common.AccountAll {
		query = query.Where("admin_users.status = ?", state)
	}
	var accountList []*AdminUserEx = make([]*AdminUserEx, 0)
	tx := query.Find(&accountList)
	if tx.Error != nil {
		return nil, tx.Error
	}

	var adminIdList []uint64 = make([]uint64, 0)
	var adminIdMap map[uint64]*AdminUserEx = make(map[uint64]*AdminUserEx) // Map for indexing by admin id
	for _, admin := range accountList {
		adminIdList = append(adminIdList, admin.Id)
		adminIdMap[admin.Id] = admin
	}

	//Query all admin related roles
	var simplifiedRoleList []*SimplifiedAdminRole = make([]*SimplifiedAdminRole, 0)
	tx = conn.Session(&gorm.Session{}).Model(AdminUser{}).
		Select("rel.admin_id", "rel.role_id", "r.role_name").
		Joins(" LEFT JOIN admin_user_role_relations as rel ON rel.admin_id = admin_users.id").
		Joins(" INNER JOIN admin_roles as r ON rel.role_id = r.id").
		Where("admin_users.id IN ? ", adminIdList).
		Order("rel.admin_id").
		Find(&simplifiedRoleList)
	if tx.Error != nil {
		return nil, tx.Error
	}

	var currentAdminId uint64 = 0
	var startIndex int = 0
	var endIndex int = 0
	size := len(simplifiedRoleList)
	if size > 0 {
		currentAdminId = simplifiedRoleList[0].AdminId
		//Assing roles to all admin respectively
		for index, role := range simplifiedRoleList {
			if currentAdminId != role.AdminId {
				endIndex = index
				adminIdMap[currentAdminId].RoleList = simplifiedRoleList[startIndex:endIndex]
				currentAdminId = role.AdminId
				startIndex = endIndex
			}
			if index == size-1 {
				adminIdMap[currentAdminId].RoleList = simplifiedRoleList[startIndex:]
			}
		}
	}
	return accountList, nil

}
func (me *AdminAccountDBProcessor) DbCreateAccount(req *admin_common.CreateAccountRequest, createdByAdminId uint64) (uint64, error) {
	conn, err := me.db.GetConn()
	if err != nil {
		return 0, err
	}

	account := new(AdminUser)
	account.Username = req.UserName
	account.ModifiedBy = createdByAdminId
	account.Status = AdminWaitForBindGA
	salt := reverse(account.Username)

	hash, err := crypto.GenerateBCryptHash([]byte(req.Password), []byte(salt))
	if err != nil {
		return 0, errors.Wrap(err, "Unable to encrypt password")
	}
	account.Password = base64.StdEncoding.EncodeToString(hash)
	tx := conn.Create(&account)
	if tx.Error != nil {
		return 0, errors.Wrap(tx.Error, "Insert admin user failed")
	}

	if account.Id == 0 {
		return 0, errors.New("Cannot generate admin id")
	}

	return account.Id, nil
}

func (me *AdminAccountDBProcessor) DbQueryAdminById(adminId uint64) (*AdminUser, *foundation.ServerError) {
	conn, err := me.db.GetConn()
	if err != nil {
		return nil, foundation.NewErrorWithMessage(foundation.DatabaseError, err.Error())
	}

	account := new(AdminUser)
	tx := conn.Session(&gorm.Session{}).Model(AdminUser{}).Where("id = ? and status <> ?", adminId, AdminDeleted).FirstOrInit(&account)
	if tx.Error != nil {
		return nil, foundation.NewErrorWithMessage(foundation.DatabaseError, err.Error())
	}
	if account.Id == 0 {
		return nil, foundation.NewErrorWithMessage(foundation.UserNotFound, "Admin not found")
	}

	return account, nil
}

func (me *AdminAccountDBProcessor) DbUpdateAccount(req *admin_common.UpdateAccountRequest, salt string, updatedByAdminId uint64) error {
	conn, err := me.db.GetConn()
	if err != nil {
		return err
	}
	err = conn.Transaction(func(db *gorm.DB) error {
		var updateMap map[string]interface{} = make(map[string]interface{})
		updateMap["last_modified_date"] = time.Now()
		updateMap["modified_by"] = updatedByAdminId

		if req.UpdateField&admin_common.AccountFieldPassword > 0 {

			hash, err := crypto.GenerateBCryptHash([]byte(req.Password), []byte(salt))
			if err != nil {
				return errors.Wrap(err, "Unable to generate password hash")
			}
			updateMap["password"] = base64.StdEncoding.EncodeToString(hash)
		}

		if req.UpdateField&admin_common.AccountFieldState > 0 {
			if req.State == int(AdminWaitForBindGA) {
				updateMap["secret"] = ""
			}
			updateMap["status"] = req.State
		}

		tx := db.Session(&gorm.Session{}).Model(AdminUser{}).Where("id = ?", req.AdminId).Updates(updateMap)
		if tx.Error != nil {
			return errors.Wrap(tx.Error, "Update field failed")
		}

		if req.UpdateField&admin_common.AccountFieldRole > 0 {

			tx = db.Session(&gorm.Session{}).Model(AdminUserRoleRelation{}).Delete(nil, "admin_id = ?", req.AdminId)
			if tx.Error != nil {
				return errors.Wrap(tx.Error, "Delete  admin role failed")
			}

			for _, roleId := range req.RoleIdList {
				var checkId uint64
				tx = db.Session(&gorm.Session{}).Model(AdminRole{}).Select("id").Where("id = ?", roleId).First(&checkId)
				if tx.Error != nil {
					return errors.Wrap(tx.Error, "Role Id not found: Role id: "+strconv.FormatUint(roleId, 10))
				}
				roleRelation := new(AdminUserRoleRelation)
				roleRelation.InitDate()
				roleRelation.AdminId = req.AdminId
				roleRelation.RoleId = roleId
				roleRelation.CreatedBy = updatedByAdminId
				tx = db.Session(&gorm.Session{}).Model(roleRelation).Create(&roleRelation)
				if tx.Error != nil {
					return errors.Wrap(tx.Error, "Insert admin role failed. role id: "+strconv.FormatUint(roleId, 10))
				}
			}
		}

		return nil
	})

	return err
}

func (me *AdminAccountDBProcessor) DbDeleteAccount(adminId uint64, updatedByAdminId uint64) error {
	conn, err := me.db.GetConn()
	if err != nil {
		return err
	}

	var updateMap map[string]interface{} = make(map[string]interface{})
	updateMap["username"] = gorm.Expr("'_' || id::text || '_' || username")
	updateMap["modified_by"] = updatedByAdminId
	updateMap["last_modified_date"] = time.Now()
	updateMap["status"] = AdminDeleted
	tx := conn.Session(&gorm.Session{}).Model(AdminUser{}).
		Where("id = ?", adminId).
		Updates(updateMap)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "Unable to update record")
	}
	return nil
}
