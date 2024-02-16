package permission_req

type AdditionalPrivilegeReq struct {
	UserId int64 `json:"userId"`
	MenuId int64 `json:"menuId"`
}

type MenuPermissionReq struct {
	PermissionName string `json:"permissionName"`
}
