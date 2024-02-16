package role_permission_service

import (
	"thor/src/payload/response"
	"thor/src/payload/role_permission_req"
)

type IRoleService interface {
	CreateNewRole(req role_permission_req.RoleWithPermissionReq) (resp response.GlobalResponse)
	UpdateRoleAndPermission(roleId string, req role_permission_req.RoleWithPermissionReq) (resp response.GlobalResponse)
	DeleteRole(roleId string) (resp response.GlobalResponse)
	GetAllRole() (resp response.GlobalResponse)
	GetAllRoleById(roleId string) (resp response.GlobalResponse)

	GetPermissionByRoleId(roleId string) (resp response.GlobalResponse)
	GetPermissionByUserId(userId int64, roleId string) (resp response.GlobalResponse)

	AddAdditionalMenuToUser(req []role_permission_req.ShardAdditionalPrivilege) (resp response.GlobalResponse)
	RemoveAdditionalMenuFromUser(req []role_permission_req.ShardAdditionalPrivilege) (resp response.GlobalResponse)
}
