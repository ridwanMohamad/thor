package container

import (
	"fmt"
	"os"
	"thor/src/repository/lov"
	"thor/src/repository/modules_menus"
	"thor/src/repository/roles_permissions"
	"thor/src/repository/sessions"
	"thor/src/repository/users"
	"thor/src/server/config"
	"thor/src/server/database"
	mail2 "thor/src/server/mail"
	"thor/src/server/service"
	"thor/src/service/auth_service"
	"thor/src/service/lov_service"
	"thor/src/service/module_menu_service"
	"thor/src/service/role_permission_service"
	"thor/src/service/user_service"
)

type DefaultContainer struct {
	//#register config
	Config                *config.DefaultConfig
	UserService           user_service.IUserService
	AuthService           auth_service.IAuthService
	LovService            lov_service.ILovService
	ModuleMenuService     module_menu_service.IModuleMenuService
	RolePermissionService role_permission_service.IRoleService
}

func IntializeContainer() *DefaultContainer {

	config := config.ConfigApps("./resources/")

	//initialize database
	db := database.InitializeDatabase(config.Database)

	//initialize mailer
	//mail := mail2.MailInitialize(config.Mail)
	goMail := mail2.GoMailInitialize(config.Mail)

	//initialize repository
	userRepo := users.NewUsersRepository(db)
	sessionRepo := sessions.NewSessionsRepository(db)
	lovRepo := lov.NewLovRepository(db)
	moduleMenuRepo := modules_menus.NewModulesMenusRepository(db)
	rolePermissionRepo := roles_permissions.NewRolesRepository(db)

	//attach repo to service
	userService := user_service.NewUserService(userRepo, moduleMenuRepo, rolePermissionRepo, lovRepo, sessionRepo, goMail)
	authService := auth_service.NewAuthService(sessionRepo, userRepo, goMail, moduleMenuRepo, rolePermissionRepo, lovRepo, config.AuthConfig)
	lovService := lov_service.NewLovService(lovRepo)
	moduleMenuService := module_menu_service.NewMenuService(moduleMenuRepo, rolePermissionRepo, userRepo)
	rolePermissionService := role_permission_service.NewRoleService(rolePermissionRepo, moduleMenuRepo, userRepo)

	row := db.Raw("SELECT 1").Row()
	service.InitializeGcp(config.Service)
	fmt.Println(row)

	os.Setenv("ENV", config.Apps.Env)

	return &DefaultContainer{
		Config:                config,
		UserService:           userService,
		AuthService:           authService,
		LovService:            lovService,
		ModuleMenuService:     moduleMenuService,
		RolePermissionService: rolePermissionService,
	}
}
