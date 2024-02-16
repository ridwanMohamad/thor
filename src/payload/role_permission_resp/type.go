package role_permission_resp

import (
	"gopkg.in/guregu/null.v4"
	"thor/src/constanta/enum"
)

type AfterCreateRole struct {
	RoleId string `json:"roleId"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type RolePermissionResp struct {
	RoleId              string      `json:"roleId"`
	Name                string      `json:"name"`
	Status              string      `json:"status,omitempty"`
	PermissionMenu      interface{} `json:"permissionMenu,omitempty"`
	AdditionalPrivilege interface{} `json:"additionalPrivilege,omitempty"`
}

type ShardPermission struct {
	Pk              int64       `json:"-"`
	ModuleId        string      `json:"moduleId"`
	Name            string      `json:"moduleName"`
	MaintenanceMode bool        `json:"maintenanceMode"`
	MtcStart        null.Time   `json:"mtcStart"`
	MtcEnd          null.Time   `json:"mtcEnd"`
	Menu            interface{} `json:"menu,omitempty"`
}

type ShardMenu struct {
	Pk         int64                            `json:"-"`
	ParentId   string                           `json:"parentId,omitempty"`
	MenuId     string                           `json:"menuId"`
	Type       string                           `json:"menuType"`
	MenuName   string                           `json:"name"`
	Path       string                           `json:"path"`
	Icon       string                           `json:"icon"`
	Registered enum.RegisteredEnum              `json:"registeredStatus,omitempty"`
	Permission []JoinedUserMenuPermissionMatrix `json:"permission"`
	SubMenu    interface{}                      `json:"subMenu,omitempty"`
}

type ShardSubMenu struct {
	Pk         int64                            `json:"-"`
	ParentId   string                           `json:"parentId,omitempty"`
	MenuId     string                           `json:"menuId"`
	MenuName   string                           `json:"name"`
	Type       string                           `json:"type"`
	Path       string                           `json:"path"`
	Icon       string                           `json:"icon"`
	Registered enum.RegisteredEnum              `json:"registeredStatus,omitempty"`
	Permission []JoinedUserMenuPermissionMatrix `json:"permission"`
}

type ShardAdditionalPrivilege struct {
	MenuId   string `json:"menuId"`
	MenuName string `json:"name"`
}

type JoinedModules struct {
	Pk              int64
	ModuleId        string
	Name            string
	MaintenanceMode bool
	MtcStart        null.Time
	MtcEnd          null.Time
}

type JoinedMenu struct {
	Pk       int64
	ParentId string
	MenuId   string
	Type     string
	Name     string
	Path     string
	MenuIcon string
}

type NewJoinedMenu struct {
	Pk               int64               `json:"pk"`
	ParentId         string              `json:"parentId"`
	MenuId           string              `json:"menuId"`
	Type             string              `json:"type"`
	Name             string              `json:"name"`
	Path             string              `json:"path"`
	MenuIcon         string              `json:"menuIcon"`
	RegisteredStatus enum.RegisteredEnum `json:"registeredStatus"`
}

type JoinedPrivilegeMenu struct {
	Pk       int64
	ParentId string
	MenuId   string
	Type     string
	Name     string
	Path     string
	MenuIcon string
}

type JoinedRoleMenu struct {
	PkRole int64
	PkMenu int64
	MenuId string
	Name   string
}

type FinalJoinedRoleMenu struct {
	PkRole     int64
	PkMenu     int64
	MenuId     string
	Name       string
	Permission []JoinedRoleMenuPermissionMatrix `json:"MenuPermission,omitempty"`
}

type JoinedRoleMenuPermissionMatrix struct {
	Pk       int64
	PkMenu   int64
	Name     string
	PermCode string
}

type FinalJoinedUserMenu struct {
	PkUser     int64
	PkMenu     int64
	MenuId     string
	Name       string
	Permission []JoinedUserMenuPermissionMatrix `json:"UserPermission,omitempty"`
}

type JoinedUserMenuPermissionMatrix struct {
	Pk               int64               `json:"pk"`
	PkMenu           int64               `json:"pkMenu"`
	Name             string              `json:"name"`
	PermCode         string              `json:"permCode"`
	RegisteredStatus enum.RegisteredEnum `json:"registeredStatus,omitempty"`
}

type JoinedPermittedMenu struct {
	PkModuleId int64  `json:"pkModuleId"`
	ModuleId   string `json:"moduleId"`
	ModuleName string `json:"moduleName"`
	ParentId   string `json:"parentId"`
	MenuId     string `json:"menuId"`
	PkMenu     int64  `json:"pkMenu"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Path       string `json:"path"`
	MenuIcon   string `json:"menuIcon"`
}

//modules
type PermittedPermissionRes struct {
	PkModuleId int64           `json:"pkModuleId"`
	ModuleId   string          `json:"moduleId"`
	ModuleName string          `json:"moduleName"`
	Menu       []PermittedMenu `json:"listMenu"`
}

//menu
type PermittedMenu struct {
	ParentId   string                           `json:"parentId"`
	MenuId     string                           `json:"menuId"`
	PkMenu     int64                            `json:"pkMenu"`
	Name       string                           `json:"name"`
	Type       string                           `json:"type"`
	Path       string                           `json:"path"`
	MenuIcon   string                           `json:"menuIcon"`
	Permission []JoinedRoleMenuPermissionMatrix `json:"permission"`
	SubMenu    []PermittedSubMenu               `json:"listSubMenu"`
}

//sub menu
type PermittedSubMenu struct {
	ParentId   string                           `json:"parentId"`
	MenuId     string                           `json:"menuId"`
	PkMenu     int64                            `json:"pkMenu"`
	Name       string                           `json:"name"`
	Type       string                           `json:"type"`
	Path       string                           `json:"path"`
	MenuIcon   string                           `json:"menuIcon"`
	Permission []JoinedRoleMenuPermissionMatrix `json:"permission"`
}
