package roles_permissions

import (
	"thor/src/domain/tbl_additional_privileges"
	"thor/src/domain/tbl_role_menu"
	"thor/src/domain/tbl_role_menu_matrix"
	"thor/src/domain/tbl_roles"
	"thor/src/domain/tbl_user_menu_matrix"
	"thor/src/payload/role_permission_resp"
)

type IRolesPermissionRepository interface {
	// SaveRole role
	SaveRole(data *tbl_roles.Roles) (resp *tbl_roles.Roles, err error)
	FindAllRole() (resp *[]tbl_roles.Roles, err error)
	FindRoleById(roleId string) (resp *tbl_roles.Roles, err error)
	FindRoleByPkId(id int64) (resp *tbl_roles.Roles, err error)
	UpdateRole(data tbl_roles.Roles) (err error)

	// SaveRolesMenus role menu
	SaveRolesMenus(data *[]tbl_role_menu.RoleMenu) (resp *[]role_permission_resp.JoinedRoleMenu, err error)
	FindRolesMenusByPkRole(pkRole int64) (resp *[]tbl_role_menu.RoleMenu, err error)
	FindRolesMenusByRoleId(roleId string) (resp *[]tbl_role_menu.RoleMenu, err error)
	FindRolesMenusByRoleIdAndMenuId(roleId string, menuId string) (resp *tbl_role_menu.RoleMenu, err error)
	DeleteRolesMenus(data *[]tbl_role_menu.RoleMenu) (err error)

	// SaveRoleMatrix saving permission base on menu and role
	SaveRoleMatrix(data *[]tbl_role_menu_matrix.RoleMenuMatrix) (resp *[]role_permission_resp.JoinedRoleMenuPermissionMatrix, err error)
	FindRoleMatrixByPkRoleAndPkMenu(pkRole int64, pkMenu int64) (resp *[]role_permission_resp.JoinedRoleMenuPermissionMatrix, err error)
	RemoveRoleMatrix(data *[]tbl_role_menu_matrix.RoleMenuMatrix) (err error)

	// SaveAddPriv save additional privilege
	SaveAddPriv(data *[]tbl_additional_privileges.AdditionalPrivileges) (resp *[]tbl_additional_privileges.AdditionalPrivileges, err error)
	FindAddPrivByUserId(userId int64) (resp *[]tbl_additional_privileges.AdditionalPrivileges, err error)
	DeleteAddPriv(data *[]tbl_additional_privileges.AdditionalPrivileges) (err error)
	DeleteAddPrivByUserPk(userPk int64) (err error)

	// SaveUserMatrix saving permission base on menu and user for additional privilege
	SaveUserMatrix(data *[]tbl_user_menu_matrix.UserMenuMatrix) (resp *[]role_permission_resp.JoinedUserMenuPermissionMatrix, err error)
	FindUserMatrixByPkUserAndPkMenu(pkUser int64, pkMenu int64) (resp *[]role_permission_resp.JoinedUserMenuPermissionMatrix, err error)
	RemoveUserMatrix(data *[]tbl_user_menu_matrix.UserMenuMatrix) (err error)

	// FindModuleByPkRole select menu and module
	FindModuleByPkRole(id int64) (resp *[]role_permission_resp.JoinedModules, err error)
	FindMenuParentByPkRole(pkRole int64, pkModules int64) (resp *[]role_permission_resp.JoinedMenu, err error)
	FindMenuChildByParentId(pkRole int64, parentId string) (resp *[]role_permission_resp.JoinedMenu, err error)
	FindMenuPrivByUserId(pkUser int64) (resp *[]role_permission_resp.JoinedPrivilegeMenu, err error)

	// FindPermissionWithRegisteredStatus select permission with registered status by pk role and pk menu
	FindPermissionWithRegisteredStatus(pkRole int64, pkMenu int64) (resp *[]role_permission_resp.JoinedUserMenuPermissionMatrix, err error)
	FindAllPermissionWithRegisteredStatus(pkUser int64, pkMenu int64) (resp *[]role_permission_resp.JoinedUserMenuPermissionMatrix, err error)

	FindPermittedPermission(pkUser int64) (resp *[]role_permission_resp.JoinedRoleMenuPermissionMatrix, err error)
}
