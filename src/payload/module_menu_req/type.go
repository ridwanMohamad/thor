package module_menu_req

import (
	"thor/src/constanta/enum"
)

type CreateModuleReq struct {
	Name        string          `json:"name" validate:"required"`
	Status      enum.StatusEnum `json:"status" validate:"required"`
	Description string          `json:"description"`
	Maintenance bool            `json:"maintenance"`
	MtcStart    string          `json:"mtcStart" validate:"ISO8601date"`
	MtcEnd      string          `json:"mtcEnd" validate:"ISO8601date"`
}

type CreateMenuReq struct {
	ModuleId   string            `json:"moduleId" validate:"required"`
	ParentId   string            `json:"parentId"`
	Name       string            `json:"name" validate:"required"`
	Path       string            `json:"path" validate:"required"`
	Icon       string            `json:"icon"`
	Permission []PermissionShard `json:"menuPermission"`
}

type PermissionShard struct {
	PermissionName string             `json:"name"`
	State          enum.AddRemoveEnum `json:"state"`
}
