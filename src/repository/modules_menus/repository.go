package modules_menus

import (
	"errors"
	"thor/src/constanta"
	"thor/src/constanta/enum"
	"thor/src/domain/tbl_menu_permission"
	"thor/src/domain/tbl_menus"
	"thor/src/domain/tbl_modules"
	"thor/src/domain/tbl_role_menu"
	"thor/src/domain/tbl_users"
	"thor/src/payload/role_permission_resp"
	"thor/src/server/database"
	"time"

	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type modulesMenus struct {
	db *database.Database
}

func NewModulesMenusRepository(DB *database.Database) IModulesMenusRepository {

	return &modulesMenus{db: DB}
}

func (m modulesMenus) SaveModule(data *tbl_modules.Modules) (resp *tbl_modules.Modules, err error) {
	//TODO implement me
	tx := m.db

	if err = tx.Create(data).Error; err != nil {
		return nil, constanta.DbFailedToInsertData
	}

	if err = tx.Where("module_id = ?", data.ModuleId).Find(&resp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return nil, nil
		}
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (m modulesMenus) FindAllModules() (resp *[]tbl_modules.Modules, err error) {
	//TODO implement me
	tx := m.db

	if err = tx.Order("name asc").Find(&resp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (m modulesMenus) FindModuleById(moduleId string) (resp *tbl_modules.Modules, err error) {
	//TODO implement me
	tx := m.db

	if err = tx.First(&resp, "module_id = ? ", moduleId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (m modulesMenus) FindModuleByCode(moduleCode string) (resp *tbl_modules.Modules, err error) {
	//TODO implement me
	tx := m.db

	if err = tx.First(&resp, "module_code = ?", moduleCode).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (m modulesMenus) UpdateModule(data *tbl_modules.Modules) (err error) {
	//TODO implement me
	tx := m.db.Begin()

	if err = tx.Model(&tbl_modules.Modules{}).
		Where("module_id = ? ", data.ModuleId).
		Updates(map[string]interface{}{
			"name":             data.Name,
			"status":           data.Status,
			"description":      data.Description,
			"maintenance_mode": data.MaintenanceMode,
			"mtc_start":        data.MtcStart,
			"mtc_end":          data.MtcEnd,
			"module_code":      data.ModuleCode,
			"updated_at":       null.TimeFrom(time.Now()),
		}).
		Error; err != nil {

		tx.Rollback()
		return constanta.DbFailedToUpdateData
	}
	tx.Commit()
	return
}

func (m modulesMenus) DeleteModule(moduleId string) (err error) {
	//TODO implement me
	panic("implement me")
}

func (m modulesMenus) SaveMenu(data *tbl_menus.Menus) (resp *tbl_menus.Menus, err error) {
	//TODO implement me
	tx := m.db

	if err = tx.Create(data).Error; err != nil {
		return nil, constanta.DbFailedToInsertData
	}

	if err = tx.Where("menu_id = ?", data.MenuId).First(&resp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (m modulesMenus) FindAllParentMenus() (resp *[]tbl_menus.Menus, err error) {
	//TODO implement me
	tx := m.db

	if err = tx.Find(&resp, "type = 'parent'").Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (m modulesMenus) FindAllMenus() (resp *[]tbl_menus.Menus, err error) {
	//TODO implement me
	tx := m.db

	if err = tx.Find(&resp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (m modulesMenus) FindMenuById(menuId string) (resp *tbl_menus.Menus, err error) {
	//TODO implement me
	tx := m.db

	if err = tx.Where("menu_id = ? ", menuId).First(&resp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (m modulesMenus) FindMenuByCode(menuCode string, parentId string, pkModule int64, menuType enum.MenuTypeEnum) (resp *tbl_menus.Menus, err error) {
	//TODO implement me
	tx := m.db

	if err = tx.First(&resp, "menu_code = ? and parent_id = ? and pk_module_id = ? and type = ?", menuCode, parentId, pkModule, menuType).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (m modulesMenus) FindParentMenuByModulePk(modulePk int64) (resp *[]tbl_menus.Menus, err error) {
	tx := m.db

	if err = tx.Find(&resp, "pk_module_id = ? and type = 'parent'", modulePk).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constanta.DbFailedToExecuteQuery
		}
	}
	return
}
func (m modulesMenus) FindMenuByIdAndModulePk(menuId string, modulePk int64) (resp *tbl_menus.Menus, err error) {

	return
}

func (m modulesMenus) FindMenuByPk(id int64) (resp *tbl_menus.Menus, err error) {
	//TODO implement me
	tx := m.db

	if err = tx.Where("pk = ? ", id).Find(&resp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (m modulesMenus) FindParentMenuByPk(id int64) (resp *[]tbl_menus.Menus, err error) {
	return
}

func (m modulesMenus) FindChildMenuByParentId(parentId string) (resp *[]tbl_menus.Menus, err error) {
	tx := m.db

	if err = tx.Find(&resp, "parent_id = ? and type = 'child'", parentId).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constanta.DbFailedToExecuteQuery
		}
	}
	return
}

func (m modulesMenus) UpdateMenu(data *tbl_menus.Menus) (err error) {
	//TODO implement me
	tx := m.db

	if err = tx.Model(&tbl_menus.Menus{}).
		Where("menu_id = ? ", data.MenuId).
		Updates(tbl_menus.Menus{
			ParentId:  data.ParentId,
			Name:      data.Name,
			Path:      data.Path,
			MenuCode:  data.MenuCode,
			Type:      data.Type,
			MenuIcon:  data.MenuIcon,
			UpdatedAt: null.TimeFrom(time.Now())}).
		Error; err != nil {
		return constanta.DbFailedToUpdateData
	}
	return
}

func (m modulesMenus) DeleteMenu(menuId string) (err error) {
	//TODO implement me
	panic("implement me")
}

func (m modulesMenus) FindMenuAndRoleMenuByPkRole(pkRole int64, parentId string) (resp *[]role_permission_resp.JoinedMenu, err error) {
	tx := m.db

	subQuery := tx.Distinct("pk_menu").Where("pk_role = ?", pkRole).Find(&tbl_role_menu.RoleMenu{})

	if err = tx.Model(&tbl_menus.Menus{}).
		Select("menus.pk, menus.parent_id, menus.menu_id, menus.type, menus.name, menus.path, menus.menu_icon").
		Where("menus.pk NOT IN (?) and type = 'child' and parent_id = ?", subQuery, parentId).Scan(&resp).Error; err != nil {
	}
	return
}

func (m modulesMenus) SaveMenuPermissionByMenuId(data []tbl_menu_permission.MenuPermission) (resp *[]tbl_menu_permission.MenuPermission, err error) {
	tx := m.db

	if err = tx.Create(data).Error; err != nil {
		return nil, constanta.DbFailedToInsertData
	}

	if err = tx.Find(&resp, "pk_menu = ?", data[0].PkMenu).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (m modulesMenus) UpdateMenuPermissionByMenuId(data []tbl_menu_permission.MenuPermission) (resp *[]tbl_menu_permission.MenuPermission, err error) {
	return
}

func (m modulesMenus) RemoveMenuPermissionByMenuId(data []tbl_menu_permission.MenuPermission) (err error) {
	tx := m.db.Begin()

	for _, val := range data {
		if err = tx.Where("pk_menu = ? and perm_code = ?", val.PkMenu, val.PermCode).
			Delete(&tbl_menu_permission.MenuPermission{}).
			Error; err != nil {

			tx.Rollback()
			return constanta.DbFailedToDeleteData
		}
	}
	tx.Commit()
	return
}

func (m modulesMenus) FindAllMenuPermissionByMenuPk(menuId int64) (resp *[]tbl_menu_permission.MenuPermission, err error) {
	tx := m.db

	if err = tx.Find(&resp, "pk_menu = ?", menuId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (m modulesMenus) FindAllMenuPermByMenuPkAndPermPk(menuId int64, permPk int64) (resp *tbl_menu_permission.MenuPermission, err error) {
	tx := m.db

	if err = tx.Find(&resp, "pk_menu = ? AND pk = ?", menuId, permPk).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, constanta.DbFailedToExecuteQuery
	}
	if resp.Pk == 0 {
		return nil, nil
	}
	return
}

func (m modulesMenus) FindParentMenuWithRegisteredStatus(pkRole int64) (resp *[]role_permission_resp.NewJoinedMenu, err error) {
	tx := m.db

	if err = tx.Model(&tbl_menus.Menus{}).
		Select("menus.pk, menus.parent_id, menus.menu_id, menus.type, menus.name, menus.path, menus.menu_icon,"+
			"(case when (SELECT count(1) FROM role_menu WHERE pk_role = ? and pk_menu = menus.pk) > 0 then 'registered' else 'unregistered' end) as registered_status", pkRole).
		Where("menus.type = 'parent'").Scan(&resp).Error; err != nil {
	}

	return
}

func (m modulesMenus) FindChildMenuWithRegisteredStatus(pkRole int64, menuId string) (resp *[]role_permission_resp.NewJoinedMenu, err error) {
	tx := m.db

	if err = tx.Model(&tbl_menus.Menus{}).
		Select("menus.pk, menus.parent_id, menus.menu_id, menus.type, menus.name, menus.path, menus.menu_icon,"+
			"(case when (SELECT count(1) FROM role_menu WHERE pk_role = ? and pk_menu = menus.pk) > 0 then 'registered' else 'unregistered' end) as registered_status", pkRole).
		Where("menus.type = 'child' and menus.parent_id = ?", menuId).Scan(&resp).Error; err != nil {
	}

	return
}

func (m modulesMenus) FindAllParentMenuWithRegisteredStatus(pkUser int64) (resp *[]role_permission_resp.NewJoinedMenu, err error) {
	tx := m.db

	if err = tx.Model(&tbl_menus.Menus{}).
		Select("menus.pk, menus.parent_id, menus.menu_id, menus.type, menus.name, menus.path, menus.menu_icon,"+
			"(case when (coalesce(role_menu.pk_menu, 0) <> 0 or coalesce(adds.pk_menu_id, 0) <> 0) then 'registered' else 'unregistered' end) as registered_status").
		Joins("inner join users on users.pk = ?", pkUser).
		Joins("left join role_menu on role_menu.pk_role = users.fk_role and role_menu.pk_menu = menus.pk").
		Joins("left join (select distinct pk_user_id, pk_menu_id from additional_privileges adp where adp.pk_user_id = ?) adds on adds.pk_user_id = users.pk and adds.pk_menu_id = menus.pk", pkUser).
		Where("menus.type = 'parent'").Scan(&resp).Error; err != nil {
	}

	return
}
func (m modulesMenus) FindAllChildMenuWithRegisteredStatus(pkUser int64, menuId string) (resp *[]role_permission_resp.NewJoinedMenu, err error) {
	tx := m.db

	if err = tx.Model(&tbl_menus.Menus{}).
		Select("menus.pk, menus.parent_id, menus.menu_id, menus.type, menus.name, menus.path, menus.menu_icon,"+
			"(case when (coalesce(role_menu.pk_menu, 0) <> 0 or coalesce(adds.pk_menu_id, 0) <> 0) then 'registered' else 'unregistered' end) as registered_status").
		Joins("inner join users on users.pk = ?", pkUser).
		Joins("left join role_menu on role_menu.pk_role = users.fk_role and role_menu.pk_menu = menus.pk").
		Joins("left join (select distinct pk_user_id, pk_menu_id from additional_privileges adp where adp.pk_user_id = ?) adds on adds.pk_user_id = users.pk and adds.pk_menu_id = menus.pk", pkUser).
		Where("menus.type = 'child' and menus.parent_id = ?", menuId).Scan(&resp).Error; err != nil {
	}
	return
}

func (m modulesMenus) FindALlPermittedMenu(pkUser int64) (resp *[]role_permission_resp.JoinedPermittedMenu, err error) {
	tx := m.db

	if err = tx.Raw("SELECT XX.* FROM (? UNION ?) AS XX ORDER BY 6, 3 ASC",
		tx.Model(&tbl_users.Users{}).Select("menus.pk_module_id, modules.module_id, modules.name as module_name, menus.parent_id, menus.menu_id, role_menu.pk_menu as pk_menu, menus.name, menus.type, menus.path, menus.menu_icon").
			Joins("inner join role_menu on role_menu.pk_role = users.fk_role").
			Joins("inner join menus on menus.pk = role_menu.pk_menu").
			Joins("inner join modules on modules.pk = menus.pk_module_id").
			Where("users.pk = ?", pkUser),
		tx.Model(&tbl_users.Users{}).Select("distinct menus.pk_module_id, modules.module_id, modules.name as module_name, menus.parent_id, menus.menu_id, adp.pk_menu_id as pk_menu, menus.name, menus.type, menus.path, menus.menu_icon").
			Joins("inner join additional_privileges adp on adp.pk_user_id = users.pk").
			Joins("inner join menus on menus.pk = adp.pk_menu_id").
			Joins("inner join modules on modules.pk = menus.pk_module_id").
			Where("users.pk = ?", pkUser)).Scan(&resp).Error; err != nil {
		return nil, err
	}
	return
}
