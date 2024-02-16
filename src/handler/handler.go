package handler

import (
	"thor/src/server/container"
)

type Handler struct {
	UserHandler           *userHandler
	AuthHandler           *authHandler
	LovHandler            *lovHandler
	RolePermissionHandler *rolePermissionHandler
	ModuleMenuHandler     *moduleMenuHandler
}

func InitializeHandler(container *container.DefaultContainer) *Handler {
	return &Handler{
		UserHandler:           NewUserHandler(container.UserService),
		AuthHandler:           NewAuthHandler(container.AuthService),
		LovHandler:            NewLovHandler(container.LovService),
		RolePermissionHandler: NewRolePermissionHandler(container.RolePermissionService),
		ModuleMenuHandler:     NewModuleMenuHandler(container.ModuleMenuService),
	}
}
