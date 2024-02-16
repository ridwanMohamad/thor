package router

import (
	"net/http"
	"thor/src/handler"

	"github.com/labstack/echo/v4"
)

func InitializeRouter(server *echo.Echo, handler *handler.Handler) {
	server.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "service up")
	})

	g := server.Group("/api")

	//g = server.Group("/v1/users")
	//user routing
	g.GET("/v1/users", handler.UserHandler.GetAllUsers)
	g.GET("/v1/users/username/:username", handler.UserHandler.GetAllUsers)
	g.GET("/v1/users/:userId", handler.UserHandler.GetUserByUserId)
	g.GET("/v1/users/:username/mapping-location", handler.UserHandler.GetMappingLocationUser)
	g.PUT("/v1/users/:username/mapping-location/update", handler.UserHandler.UpdateMappingLocationUser)
	g.GET("/v1/users/:username/mapping-roles", handler.UserHandler.GetMappingRoleUser)
	g.PUT("/v1/users/:username/mapping-roles/update", handler.UserHandler.UpdateMappingRoleUser)
	g.POST("/v1/users", handler.UserHandler.CreateNewUser)
	g.PUT("/v1/users", handler.UserHandler.UpdateUser)
	g.PUT("/v1/users/:userId/update-profile", handler.UserHandler.UpdateUserProfile)
	g.PUT("/v1/users/change-password", handler.UserHandler.ChangeUserPassword)
	g.DELETE("/v1/users/:userId/delete", handler.UserHandler.InActiveUser)
	g.DELETE("/v1/users/:id/mapping-location/delete", handler.UserHandler.RemoveMappingLocationUser)
	g.DELETE("/v1/users/:id/mapping-roles/delete", handler.UserHandler.RemoveMappingRoleUser)

	//auth routing
	g.POST("/v1/login", handler.AuthHandler.Login)
	g.POST("/v2/login", handler.AuthHandler.LoginV2)
	g.POST("/v1/logout", handler.AuthHandler.Logout)
	g.POST("/v1/check-token", handler.AuthHandler.CheckToken)
	g.POST("/v1/forgot-password", handler.AuthHandler.ForgotPassword)

	//Module routing
	g.POST("/v1/modules", handler.ModuleMenuHandler.CreateNewModule)
	g.PUT("/v1/modules/:moduleId/update", handler.ModuleMenuHandler.UpdateModuleData)
	g.GET("/v1/modules", handler.ModuleMenuHandler.GetAllModule)
	g.GET("/v1/modules/:moduleId", handler.ModuleMenuHandler.GetModuleByIdModule)

	//Menu routing
	g.POST("/v1/menus", handler.ModuleMenuHandler.CreateNewMenu)
	g.GET("/v1/test/menus", handler.ModuleMenuHandler.GetTestAllMenu)
	g.PUT("/v1/menus/:menuId/update", handler.ModuleMenuHandler.UpdateMenuData)
	g.GET("/v1/menus", handler.ModuleMenuHandler.GetAllMenu)
	g.GET("/v1/menus/by-module/:moduleId", handler.ModuleMenuHandler.GetAllMenuByModuleId)
	g.GET("/v1/menus/:menuId", handler.ModuleMenuHandler.GetDetailMenuByMenuId)
	g.GET("/v1/menus/unregistered/:roleId", handler.ModuleMenuHandler.GetMenuNotExistsInRoleMenu)
	g.GET("/v1/menus/unregistered/by-user/:userPk", handler.ModuleMenuHandler.GetAllMenuAndPermissionByUserId)

	//Role routing
	g.POST("/v1/role", handler.RolePermissionHandler.CreateNewRoleWithPermission)
	g.PUT("/v1/role/:roleId/update", handler.RolePermissionHandler.UpdateRoleData)
	g.DELETE("/v1/role/:roleId/delete", handler.RolePermissionHandler.DeleteRole)
	g.GET("/v1/role", handler.RolePermissionHandler.GetAllRole)
	g.GET("/v1/permission/:roleId", nil)

	//LOV routing
	g.POST("/v1/lov", handler.LovHandler.CreateNewLov)
	g.POST("/v1/lov/:headerId/detail", handler.LovHandler.CreateNewLovDetail)
	g.GET("/v1/lov", handler.LovHandler.GetAllLovHeader)
	g.PUT("/v1/lov/:headerId", handler.LovHandler.UpdateLovHeader)
	g.DELETE("/v1/lov/:headerId", handler.LovHandler.RemoveLovHeader)
	g.GET("/v1/lov/:headerId", handler.LovHandler.GetAllLovByHeaderIdOrName)
	g.GET("/v2/lov/", handler.LovHandler.GetAllLovByHeaderIdOrNameTest)
	g.PUT("/v1/lov/update/detail", handler.LovHandler.UpdateLovDetail)
	g.DELETE("/v1/lov/delete/detail/:detailId", handler.LovHandler.RemoveLovDetail)

	//LOV Permission
	g.POST("/v1/lov/permission", handler.LovHandler.CreateNewListOfPermission)
	g.PUT("/v1/lov/permission", handler.LovHandler.UpdateListOfPermission)
	g.DELETE("/v1/lov/permission/:code", handler.LovHandler.RemoveListOfPermission)
	g.GET("/v1/lov/permission", handler.LovHandler.GetAllListOfPermission)

	//permission role routing
	g.GET("/v1/permission/role/:roleId", handler.RolePermissionHandler.GetPermissionByRoleId)
	g.GET("/v1/permission/user/:userId", handler.RolePermissionHandler.GetPermissionByUserId)
	g.POST("/v1/additional/privilege", handler.RolePermissionHandler.CreateNewAdditionalPrivilege)
	g.DELETE("/v1/additional/privilege", handler.RolePermissionHandler.RemoveAdditionalPrivilege)

	g.GET("/test", func(c echo.Context) error {
		name := c.QueryParam("name")
		return c.JSON(http.StatusOK, name)
	})
	//list of enum
	g.GET("/v1/enum", func(c echo.Context) error {
		//var listEnum [3]map[string]string
		//var list []map[string]string
		//list = append(list, map[string]string{"enum_add_or_remove": string(enum.Add)})
		//list = append(list,map[string]string{"enum_add_or_remove": string(enum.Exist)}
		//list = append(list,map[string]string{"enum_add_or_remove": string(enum.Modified)}
		//list = append(list,map[string]string{"enum_add_or_remove": string(enum.Removed)}
		//listEnum[0] = list
		//listEnum[1] = map[string]string{"enum_status": string(enum.Active)}
		//listEnum[1] = map[string]string{"enum_status": string(enum.InActive)}
		//
		//listEnum[2] = map[string]string{"enum_menu_type": string(enum.Child)}
		//listEnum[2] = map[string]string{"enum_menu_type": string(enum.Parent)}

		return c.JSON(http.StatusOK, "List of enum")
	})
}
