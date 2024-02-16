package modules_menus

import (
	"thor/src/constanta/enum"
	"thor/src/domain/tbl_menu_permission"
	"thor/src/domain/tbl_menus"
	"thor/src/domain/tbl_modules"
	"thor/src/payload/role_permission_resp"
)

type IModulesMenusRepository interface {
	SaveModule(data *tbl_modules.Modules) (resp *tbl_modules.Modules, err error)
	FindAllModules() (resp *[]tbl_modules.Modules, err error)
	FindModuleById(moduleId string) (resp *tbl_modules.Modules, err error)
	FindModuleByCode(moduleCode string) (resp *tbl_modules.Modules, err error)
	UpdateModule(data *tbl_modules.Modules) (err error)
	DeleteModule(moduleId string) (err error)

	SaveMenu(data *tbl_menus.Menus) (resp *tbl_menus.Menus, err error)
	FindAllParentMenus() (resp *[]tbl_menus.Menus, err error)
	FindAllMenus() (resp *[]tbl_menus.Menus, err error)
	FindMenuById(menuId string) (resp *tbl_menus.Menus, err error)
	FindMenuByCode(menuCode string, parentId string, pkModule int64, menuType enum.MenuTypeEnum) (resp *tbl_menus.Menus, err error)
	FindParentMenuByModulePk(modulePk int64) (resp *[]tbl_menus.Menus, err error)
	FindMenuByIdAndModulePk(menuId string, modulePk int64) (resp *tbl_menus.Menus, err error)
	FindMenuByPk(id int64) (resp *tbl_menus.Menus, err error)
	FindParentMenuByPk(id int64) (resp *[]tbl_menus.Menus, err error)
	FindChildMenuByParentId(parentId string) (resp *[]tbl_menus.Menus, err error)
	UpdateMenu(data *tbl_menus.Menus) (err error)
	DeleteMenu(menuId string) (err error)
	FindMenuAndRoleMenuByPkRole(pkRole int64, parentId string) (resp *[]role_permission_resp.JoinedMenu, err error)

	SaveMenuPermissionByMenuId(data []tbl_menu_permission.MenuPermission) (resp *[]tbl_menu_permission.MenuPermission, err error)
	UpdateMenuPermissionByMenuId(data []tbl_menu_permission.MenuPermission) (resp *[]tbl_menu_permission.MenuPermission, err error)
	RemoveMenuPermissionByMenuId(data []tbl_menu_permission.MenuPermission) (err error)
	FindAllMenuPermissionByMenuPk(menuId int64) (resp *[]tbl_menu_permission.MenuPermission, err error)
	FindAllMenuPermByMenuPkAndPermPk(menuId int64, permPk int64) (resp *tbl_menu_permission.MenuPermission, err error)

	FindParentMenuWithRegisteredStatus(pkRole int64) (resp *[]role_permission_resp.NewJoinedMenu, err error)
	FindChildMenuWithRegisteredStatus(pkRole int64, menuId string) (resp *[]role_permission_resp.NewJoinedMenu, err error)

	FindAllParentMenuWithRegisteredStatus(pkUser int64) (resp *[]role_permission_resp.NewJoinedMenu, err error)
	FindAllChildMenuWithRegisteredStatus(pkUser int64, menuId string) (resp *[]role_permission_resp.NewJoinedMenu, err error)

	FindALlPermittedMenu(pkUser int64) (resp *[]role_permission_resp.JoinedPermittedMenu, err error)
}
