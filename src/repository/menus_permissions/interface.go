package menus_permissions

import (
	"thor/src/domain/tbl_menu_permission"
)

type IMenusPermissions interface {
	SaveMenuPermission(data tbl_menu_permission.MenuPermission) (resp *tbl_menu_permission.MenuPermission, err error)
	FindAllMenuPermissionByMenuId(menuId int64) (resp *[]tbl_menu_permission.MenuPermission, err error)
	UpdatePermissionByPermissionId(data tbl_menu_permission.MenuPermission) (err error)
	RemovePermissionByPermissionId(permId int64) (err error)
}
