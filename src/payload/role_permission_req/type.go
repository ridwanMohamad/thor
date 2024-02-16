package role_permission_req

import "thor/src/constanta/enum"

type RoleWithPermissionReq struct {
	Name       string            `json:"name" validate:"required"`
	Status     string            `json:"status" validate:"required"`
	Permission []ShardPermission `json:"permission" validate:"required"`
}

type AdditionalPrivilegeReq struct {
	AdditionalPrivilege ShardAdditionalPrivilege `json:"privilege" validate:"required"`
}

type ShardPermission struct {
	MenuId     string                `json:"menuId" validate:"required"`
	State      enum.AddRemoveEnum    `json:"state" validate:"required"`
	Permission []ShardMenuPermission `json:"permission" validate:"required"`
}

type ShardMenuPermission struct {
	PermId int64              `json:"permId"`
	State  enum.AddRemoveEnum `json:"state"`
}

type ShardAdditionalPrivilege struct {
	UserId     int64                 `json:"userId" validate:"required"`
	MenuId     int64                 `json:"menuId" validate:"required"`
	Permission []ShardMenuPermission `json:"permission" validate:"required"`
}
