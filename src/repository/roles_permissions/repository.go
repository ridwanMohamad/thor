package roles_permissions

import (
	"errors"
	"thor/src/constanta"
	"thor/src/domain/tbl_additional_privileges"
	"thor/src/domain/tbl_menu_permission"
	"thor/src/domain/tbl_menus"
	"thor/src/domain/tbl_role_menu"
	"thor/src/domain/tbl_role_menu_matrix"
	"thor/src/domain/tbl_roles"
	"thor/src/domain/tbl_user_menu_matrix"
	"thor/src/domain/tbl_users"
	"thor/src/payload/role_permission_resp"
	"thor/src/server/database"
	"time"

	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type rolesPermissions struct {
	db *database.Database
}

func NewRolesRepository(DB *database.Database) IRolesPermissionRepository {

	return &rolesPermissions{db: DB}
}

func (r rolesPermissions) SaveRole(data *tbl_roles.Roles) (resp *tbl_roles.Roles, err error) {
	//TODO implement me
	tx := r.db

	if err = tx.Create(data).Error; err != nil {
		return nil, constanta.DbFailedToInsertData
	}

	if err = tx.Where("role_id = ?", data.RoleId).First(&resp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (r rolesPermissions) FindAllRole() (resp *[]tbl_roles.Roles, err error) {
	//TODO implement me
	tx := r.db

	if err = tx.Order("name asc").Find(&resp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return nil, nil
		}
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (r rolesPermissions) FindRoleById(roleId string) (resp *tbl_roles.Roles, err error) {
	//TODO implement me
	tx := r.db

	if err = tx.Where("role_id = ? ", roleId).First(&resp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (r rolesPermissions) FindRoleByPkId(roleId int64) (resp *tbl_roles.Roles, err error) {
	//TODO implement me
	tx := r.db

	if err = tx.Where("pk = ? ", roleId).First(&resp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (r rolesPermissions) UpdateRole(data tbl_roles.Roles) (err error) {
	//TODO implement me
	tx := r.db.Begin()

	if err = tx.Where("role_id = ? ", data.RoleId).
		//Select("name, status, updated_at").
		Updates(tbl_roles.Roles{
			Name:      data.Name,
			Status:    data.Status,
			UpdatedAt: null.TimeFrom(time.Now())}).
		Error; err != nil {

		tx.Rollback()
		return constanta.DbFailedToUpdateData
	}
	tx.Commit()
	return
}

func (r rolesPermissions) SaveRolesMenus(data *[]tbl_role_menu.RoleMenu) (resp *[]role_permission_resp.JoinedRoleMenu, err error) {
	//TODO implement me
	tx := r.db

	if err = tx.Create(data).Error; err != nil {
		return nil, constanta.DbFailedToInsertData
	}

	var pkRole int64

	for _, val := range *data {
		pkRole = val.PkRole
		break
	}

	if err = tx.Model(&tbl_menus.Menus{}).
		Select("role_menu.pk_role, menus.pk as pk_menu, menus.menu_id, menus.name, NULL as permission").
		Joins("inner join role_menu on role_menu.pk_menu = menus.pk").
		Where("role_menu.pk_role = ?", pkRole).
		Scan(&resp).Error; err != nil {
		return nil, constanta.DbFailedToExecuteQuery
	}
	//
	//if err = tx.Where("pk_role = ?", pkRole).Find(&resp).Error; err != nil {
	//	if errors.Is(err, gorm.ErrRecordNotFound) {
	//		return nil, nil
	//	}
	//	return nil, constanta.DbFailedToExecuteQuery
	//}
	return
}

func (r rolesPermissions) FindRolesMenusByRoleId(roleId string) (resp *[]tbl_role_menu.RoleMenu, err error) {
	//TODO implement me
	tx := r.db

	if err = tx.Where("pk = ?", roleId).Find(&resp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, constanta.DbFailedToExecuteQuery
	}

	return
}

func (r rolesPermissions) FindRolesMenusByRoleIdAndMenuId(roleId string, menuId string) (resp *tbl_role_menu.RoleMenu, err error) {
	tx := r.db

	if err = tx.Model(&tbl_role_menu.RoleMenu{}).
		Select("role_menu.role_menu_id, role_menu.pk_role, role_menu.pk_menu  ").
		Joins("inner join roles on roles.pk  = role_menu.pk_role").
		Joins("inner join menus on menus.pk = role_menu.pk_menu").
		Where("roles.role_id = ? and menus.menu_id = ?", roleId, menuId).
		Scan(&resp).Error; err != nil {
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (r rolesPermissions) FindRolesMenusByPkRole(pkRole int64) (resp *[]tbl_role_menu.RoleMenu, err error) {
	//TODO implement me
	tx := r.db

	if err = tx.Where("pk_role = ?", pkRole).Find(&resp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (r rolesPermissions) DeleteRolesMenus(data *[]tbl_role_menu.RoleMenu) (err error) {
	//TODO implement me
	tx := r.db.Begin()

	for _, val := range *data {
		if err = tx.Where("pk_role = ? and pk_menu = ?", val.PkRole, val.PkMenu).Delete(&tbl_role_menu.RoleMenu{}).Error; err != nil {
			tx.Rollback()
			return constanta.DbFailedToDeleteData
		}
	}
	tx.Commit()
	return
}

func (r rolesPermissions) SaveRoleMatrix(data *[]tbl_role_menu_matrix.RoleMenuMatrix) (resp *[]role_permission_resp.JoinedRoleMenuPermissionMatrix, err error) {
	tx := r.db

	if err = tx.Create(data).Error; err != nil {
		return nil, constanta.DbFailedToInsertData
	}

	var pkRoleId int64

	for _, val := range *data {
		pkRoleId = val.PkRole
		break
	}

	if err = tx.Model(&tbl_menu_permission.MenuPermission{}).
		Select("menu_permission.pk, menu_permission.pk_menu, menu_permission.name, menu_permission.perm_code").
		Joins("inner join role_menu_matrix on role_menu_matrix.pk_menu = menu_permission.pk_menu and role_menu_matrix.pk_menu_perm = menu_permission.pk").
		Where("role_menu_matrix.pk_role = ?", pkRoleId).
		Scan(&resp).Error; err != nil {
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (r rolesPermissions) FindRoleMatrixByPkRoleAndPkMenu(pkRole int64, pkMenu int64) (resp *[]role_permission_resp.JoinedRoleMenuPermissionMatrix, err error) {
	tx := r.db

	if err = tx.Model(&tbl_menu_permission.MenuPermission{}).
		Select("menu_permission.pk, menu_permission.name, menu_permission.perm_code").
		Joins("inner join role_menu_matrix on role_menu_matrix.pk_menu = menu_permission.pk_menu and role_menu_matrix.pk_menu_perm = menu_permission.pk").
		//Joins("inner join menus on menus.pk = role_menu.pk_menu").
		Where("role_menu_matrix.pk_role = ? and role_menu_matrix.pk_menu = ?", pkRole, pkMenu).
		Scan(&resp).Error; err != nil {
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}
func (r rolesPermissions) RemoveRoleMatrix(data *[]tbl_role_menu_matrix.RoleMenuMatrix) (err error) {
	tx := r.db.Begin()

	for _, val := range *data {
		if err = tx.Where("pk_menu = ? and pk_menu_perm = ? and pk_role = ?", val.PkMenu, val.PkMenuPerm, val.PkRole).
			Delete(&tbl_role_menu_matrix.RoleMenuMatrix{}).
			Error; err != nil {

			tx.Rollback()
			return constanta.DbFailedToDeleteData
		}
	}
	tx.Commit()
	return
}

func (r rolesPermissions) SaveAddPriv(data *[]tbl_additional_privileges.AdditionalPrivileges) (resp *[]tbl_additional_privileges.AdditionalPrivileges, err error) {
	//TODO implement me
	tx := r.db

	if err = tx.Create(data).Error; err != nil {
		return nil, constanta.DbFailedToInsertData
	}

	var pkUserId int64

	for _, val := range *data {
		pkUserId = val.PkUserId
		return
	}

	if err = tx.Where("pk_user_id = ?", pkUserId).Find(&resp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (r rolesPermissions) FindAddPrivByUserId(userId int64) (resp *[]tbl_additional_privileges.AdditionalPrivileges, err error) {
	//TODO implement me
	tx := r.db

	if err = tx.Model(&resp).Where("pk_user_id = ?", userId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (r rolesPermissions) DeleteAddPriv(data *[]tbl_additional_privileges.AdditionalPrivileges) (err error) {
	tx := r.db.Begin()

	for _, val := range *data {
		if err = tx.Where("pk_user_id = ? and pk_menu_id = ?", val.PkUserId, val.PkMenuId).
			Delete(&tbl_additional_privileges.AdditionalPrivileges{}).
			Error; err != nil {

			tx.Rollback()
			return constanta.DbFailedToDeleteData
		}
	}
	tx.Commit()
	return
}

func (r rolesPermissions) DeleteAddPrivByUserPk(userPk int64) (err error) {
	tx := r.db

	if err = tx.Where("pk_user_id = ?", userPk).Delete(&tbl_additional_privileges.AdditionalPrivileges{}).Error; err != nil {
		return constanta.DbFailedToDeleteData
	}
	return
}

func (r rolesPermissions) SaveUserMatrix(data *[]tbl_user_menu_matrix.UserMenuMatrix) (resp *[]role_permission_resp.JoinedUserMenuPermissionMatrix, err error) {
	tx := r.db

	if err = tx.Create(data).Error; err != nil {
		return nil, constanta.DbFailedToInsertData
	}

	var pkUserId int64

	for _, val := range *data {
		pkUserId = val.PkUser
		break
	}

	if err = tx.Model(&tbl_menu_permission.MenuPermission{}).
		Select("menu_permission.pk, menu_permission.pk_menu, menu_permission.name, menu_permission.perm_code").
		Joins("inner join role_menu_matrix on role_menu_matrix.pk_menu = menu_permission.pk_menu and role_menu_matrix.pk_menu_perm = menu_permission.pk").
		Where("role_menu_matrix.pk_role = ?", pkUserId).
		Scan(&resp).Error; err != nil {
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (r rolesPermissions) FindUserMatrixByPkUserAndPkMenu(pkUser int64, pkMenu int64) (resp *[]role_permission_resp.JoinedUserMenuPermissionMatrix, err error) {
	tx := r.db

	if err = tx.Model(&tbl_menu_permission.MenuPermission{}).
		Select("menu_permission.pk, additional_privileges.pk_menu_id as pk_menu, menu_permission.name, menu_permission.perm_code").
		Joins("inner join additional_privileges on additional_privileges.pk_menu_id = menu_permission.pk_menu and additional_privileges.pk_menu_perm_id = menu_permission.pk").
		//Joins("inner join menus on menus.pk = role_menu.pk_menu").
		Where("additional_privileges.pk_user_id = ? and additional_privileges.pk_menu_id = ?", pkUser, pkMenu).
		Scan(&resp).Error; err != nil {
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (r rolesPermissions) RemoveUserMatrix(data *[]tbl_user_menu_matrix.UserMenuMatrix) (err error) {
	tx := r.db.Begin()

	for _, val := range *data {
		if err = tx.Where("pk_menu = ? and pk_menu_perm = ? and pk_role = ?", val.PkMenu, val.PkMenuPerm, val.PkUser).
			Delete(&tbl_role_menu_matrix.RoleMenuMatrix{}).
			Error; err != nil {

			tx.Rollback()
			return constanta.DbFailedToDeleteData
		}
	}
	tx.Commit()
	return
}

func (r rolesPermissions) FindModuleByPkRole(id int64) (resp *[]role_permission_resp.JoinedModules, err error) {
	tx := r.db

	if err = tx.Model(&tbl_role_menu.RoleMenu{}).
		Select("distinct modules.pk, modules.module_id, modules.name, modules.maintenance_mode, modules.mtc_start, modules.mtc_end").
		Joins("inner join menus on menus.pk = role_menu.pk_menu and menus.type = 'parent'").
		Joins("inner join modules on modules.pk = menus.pk_module_id").
		Where("role_menu.pk_role = ? and modules.status = 'active'", id).
		Scan(&resp).Error; err != nil {
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (r rolesPermissions) FindMenuParentByPkRole(pkRole int64, pkModules int64) (resp *[]role_permission_resp.JoinedMenu, err error) {
	tx := r.db

	if err = tx.Model(&tbl_role_menu.RoleMenu{}).
		Select("menus.pk, menus.parent_id, menus.menu_id, menus.type, menus.name, menus.path, menus.menu_icon").
		Joins("inner join menus on menus.pk = role_menu.pk_menu and menus.type = 'parent'").
		Joins("inner join modules on modules.pk = menus.pk_module_id").
		Where("role_menu.pk_role = ? and modules.pk = ?", pkRole, pkModules).
		Scan(&resp).Error; err != nil {
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (r rolesPermissions) FindMenuChildByParentId(pkRole int64, parentId string) (resp *[]role_permission_resp.JoinedMenu, err error) {
	tx := r.db

	if err = tx.Model(&tbl_role_menu.RoleMenu{}).
		Select("menus.pk, menus.parent_id, menus.menu_id, menus.type, menus.name, menus.path, menus.menu_icon").
		Joins("inner join menus on menus.pk = role_menu.pk_menu and menus.type = 'child'").
		Joins("inner join modules on modules.pk = menus.pk_module_id").
		Where("role_menu.pk_role = ? and menus.parent_id = ?", pkRole, parentId).
		Scan(&resp).Error; err != nil {
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (r rolesPermissions) FindMenuPrivByUserId(pkUser int64) (resp *[]role_permission_resp.JoinedPrivilegeMenu, err error) {
	//TODO implement me
	tx := r.db

	if err = tx.Model(&tbl_additional_privileges.AdditionalPrivileges{}).
		Select("distinct menus.pk, menus.parent_id, menus.menu_id, menus.type, menus.name, menus.path, menus.menu_icon").
		Joins("inner join menus on menus.pk = additional_privileges.pk_menu_id").
		Joins("inner join users on users.pk = additional_privileges.pk_user_id").
		Where("additional_privileges.pk_user_id = ?", pkUser).
		Scan(&resp).Error; err != nil {
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (r rolesPermissions) FindPermissionWithRegisteredStatus(pkRole int64, pkMenu int64) (resp *[]role_permission_resp.JoinedUserMenuPermissionMatrix, err error) {
	tx := r.db

	var sel = "menu_permission.pk, menu_permission.pk_menu, menu_permission.name, menu_permission.perm_code, (case when coalesce(role_menu_matrix.pk_menu_perm,0) <> 0 then 'registered' else 'unregistered' end) as registered_status"
	var join = "left join role_menu_matrix on role_menu_matrix.pk_menu = menu_permission.pk_menu and role_menu_matrix.pk_menu_perm = menu_permission.pk and role_menu_matrix.pk_role = ?"

	if err = tx.Model(&tbl_menu_permission.MenuPermission{}).
		Select(sel).
		Joins("inner join menus on menus.pk = menu_permission.pk_menu").
		Joins(join, pkRole).
		Where("menus.pk = ?", pkMenu).Scan(&resp).Error; err != nil {
	}

	return
}

func (r rolesPermissions) FindAllPermissionWithRegisteredStatus(pkUser int64, pkMenu int64) (resp *[]role_permission_resp.JoinedUserMenuPermissionMatrix, err error) {
	tx := r.db

	if err = tx.Model(&tbl_menu_permission.MenuPermission{}).
		Select("menu_permission.pk, menu_permission.pk_menu, menu_permission.name, menu_permission.perm_code,"+
			"(case when (coalesce(rm.pk_menu, 0) <> 0 or coalesce(adp.pk_menu_id, 0) <> 0) then 'registered' else 'unregistered' end) as registered_status").
		Joins("inner join menus on menus.pk = menu_permission.pk_menu").
		Joins("inner join users on users.pk = ?", pkUser).
		Joins("left join thor.role_menu_matrix rm on rm.pk_menu = menu_permission.pk_menu and rm.pk_menu_perm = menu_permission.pk and rm.pk_role = users.fk_role").
		Joins("left join thor.additional_privileges adp on adp.pk_menu_id = menu_permission.pk_menu and adp.pk_menu_perm_id = menu_permission.pk and adp.pk_user_id = users.pk").
		Where("menus.pk = ?", pkMenu).Scan(&resp).Error; err != nil {
	}

	return
}

func (r rolesPermissions) FindPermittedPermission(pkUser int64) (resp *[]role_permission_resp.JoinedRoleMenuPermissionMatrix, err error) {
	tx := r.db

	if err = tx.Raw("? UNION ?",
		tx.Model(&tbl_users.Users{}).
			Select("mp.pk, mp.pk_menu, mp.name, mp.perm_code").
			Joins("inner join roles on roles.pk = users.fk_role").
			Joins("inner join role_menu_matrix rmm on rmm.pk_role = roles.pk").
			Joins("inner join menus on menus.pk = rmm.pk_menu").
			Joins("inner join menu_permission as mp on mp.pk = rmm.pk_menu_perm").
			Where("users.pk = ?", pkUser),
		tx.Model(&tbl_users.Users{}).
			Select("mp.pk, mp.pk_menu, mp.name, mp.perm_code").
			Joins("inner join additional_privileges ap on ap.pk_user_id = users.pk").
			Joins("inner join menus on menus.pk = ap.pk_menu_id").
			Joins("inner join menu_permission as mp on mp.pk = ap.pk_menu_perm_id").
			Where(" users.pk = ?", pkUser)).Scan(&resp).Error; err != nil {

	}
	return
}
