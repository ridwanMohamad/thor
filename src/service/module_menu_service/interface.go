package module_menu_service

import (
	"thor/src/payload/module_menu_req"
	"thor/src/payload/response"
)

type IModuleMenuService interface {
	CreateModule(data module_menu_req.CreateModuleReq) (resp response.GlobalResponse)
	UpdateModule(moduleId string, data module_menu_req.CreateModuleReq) (resp response.GlobalResponse)
	DeleteModule(moduleId string) (resp response.GlobalResponse)
	GetAllModule() (resp response.GlobalResponse)
	GetModuleById(moduleId string) (resp response.GlobalResponse)

	CreateMenu(data module_menu_req.CreateMenuReq) (resp response.GlobalResponse)
	UpdateMenu(menuId string, data module_menu_req.CreateMenuReq) (resp response.GlobalResponse)
	DeleteMenu(menuId string) (resp response.GlobalResponse)
	GetAllMenu() (resp response.GlobalResponse)
	GetAllMenuByModuleId(moduleId string) (resp response.GlobalResponse)
	GetMenuById(menuId string) (resp response.GlobalResponse)
	GetMenuByPermission(roleId string) (resp response.GlobalResponse)
	GetUnregisteredMenuByUserId(userId int64) (resp response.GlobalResponse)
	GetAllPermittedMenu(userId int64) (resp response.GlobalResponse)
	TestGetAllMenu() (resp response.GlobalResponse)
}
