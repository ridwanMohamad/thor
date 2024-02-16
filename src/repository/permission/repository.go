package permission

import (
	"errors"
	"gorm.io/gorm"
	"thor/src/constanta"
	"thor/src/domain/tbl_additional_privileges"
	"thor/src/server/database"
)

type permission struct {
	db *database.Database
}

func NewPermissionRepository(DB *database.Database) IPermissionRepository {

	return &permission{db: DB}
}

func (p permission) SaveAddPriv(data tbl_additional_privileges.AdditionalPrivileges) (resp *tbl_additional_privileges.AdditionalPrivileges, err error) {
	//TODO implement me
	tx := p.db

	if err = tx.Create(data).Error; err != nil {
		return nil, constanta.DbFailedToInsertData
	}

	if err = tx.First(&resp).Where("pk_user_id = ? and pk_menu_id", data.PkUserId, data.PkMenuId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (p permission) FindAddPrivByUserId(userId int64) (resp *tbl_additional_privileges.AdditionalPrivileges, err error) {
	//TODO implement me
	tx := p.db.Begin()

	if err = tx.Model(&resp).Where("pk_user_id = ?", userId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return nil, nil
		}
		tx.Rollback()
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (p permission) DeleteAddPriv(data tbl_additional_privileges.AdditionalPrivileges) (err error) {
	//TODO implement me
	panic("implement me")
}
