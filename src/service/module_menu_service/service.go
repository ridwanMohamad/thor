package module_menu_service

import (
	"errors"
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
	"strings"
	"thor/src/constanta"
	"thor/src/constanta/enum"
	"thor/src/domain/tbl_menu_permission"
	"thor/src/domain/tbl_menus"
	"thor/src/domain/tbl_modules"
	"thor/src/domain/tbl_roles"
	"thor/src/domain/tbl_users"
	"thor/src/payload/module_menu_req"
	"thor/src/payload/module_menu_resp"
	"thor/src/payload/response"
	"thor/src/payload/role_permission_resp"
	"thor/src/repository/modules_menus"
	"thor/src/repository/roles_permissions"
	"thor/src/repository/users"
	"thor/src/util"
	"thor/src/util/date_util"
	"thor/src/util/string_util"
	"time"
)

type service struct {
	moduleMenuRepo modules_menus.IModulesMenusRepository
	roleRepo       roles_permissions.IRolesPermissionRepository
	userRepo       users.IUsersRepository
}

func NewMenuService(moduleMenuRepo modules_menus.IModulesMenusRepository,
	roleRepo roles_permissions.IRolesPermissionRepository,
	userRepo users.IUsersRepository) IModuleMenuService {

	return &service{moduleMenuRepo: moduleMenuRepo, roleRepo: roleRepo, userRepo: userRepo}
}

func (s service) CreateModule(data module_menu_req.CreateModuleReq) (resp response.GlobalResponse) {
	//TODO implement me
	var id uuid.UUID
	err := errors.New("")

	moduleCode := strings.ToUpper(string_util.TrimSpace(data.Name))

	if dt, err := s.moduleMenuRepo.FindModuleByCode(moduleCode); dt != nil && err == nil {
		return util.CreateGlobalResponse(constanta.ErrorFailedToModuleExists, nil)
	}

	if id, err = uuid.NewUUID(); err != nil {
		return util.CreateGlobalResponse(constanta.ErrorGeneral, nil)
	}

	var mtcStart null.Time
	var mtcEnd null.Time

	if data.Maintenance {
		mtcStart = date_util.StringToNilTime(data.MtcStart)
		mtcEnd = date_util.StringToNilTime(data.MtcEnd)
	}

	if dt, err := s.moduleMenuRepo.SaveModule(&tbl_modules.Modules{
		ModuleId:        id.String(),
		Name:            data.Name,
		Description:     data.Description,
		Status:          data.Status,
		MaintenanceMode: data.Maintenance,
		MtcStart:        mtcStart,
		MtcEnd:          mtcEnd,
		ModuleCode:      moduleCode,
	}); err == nil {
		return util.CreateGlobalResponse(constanta.Success, dt)
	}
	return util.CreateGlobalResponse(constanta.ErrorFailedToCreateModule, nil)
}

func (s service) UpdateModule(moduleId string, data module_menu_req.CreateModuleReq) (resp response.GlobalResponse) {
	//TODO implement me
	var dt *tbl_modules.Modules
	err := errors.New("")

	if dt, err = s.moduleMenuRepo.FindModuleById(moduleId); err != nil || dt == nil {
		return util.CreateGlobalResponse(constanta.ErrorModuleDataNotFound, nil)
	}

	moduleCode := strings.ToUpper(string_util.TrimSpace(data.Name))

	if dt.Name != data.Name {
		if cd, err := s.moduleMenuRepo.FindModuleByCode(moduleCode); cd != nil && err == nil {
			return util.CreateGlobalResponse(constanta.ErrorFailedToModuleExists, nil)
		}
	}

	var mtcStart null.Time
	var mtcEnd null.Time

	if data.Maintenance {
		mtcStart = date_util.StringToNilTime(data.MtcStart)
		mtcEnd = date_util.StringToNilTime(data.MtcEnd)
	}

	if err := s.moduleMenuRepo.UpdateModule(&tbl_modules.Modules{
		ModuleId:        dt.ModuleId,
		Name:            data.Name,
		Description:     data.Description,
		Status:          data.Status,
		MaintenanceMode: data.Maintenance,
		MtcStart:        mtcStart,
		MtcEnd:          mtcEnd,
		ModuleCode:      moduleCode,
	}); err == nil {
		return util.CreateGlobalResponse(constanta.Success, nil)
	}
	return util.CreateGlobalResponse(constanta.ErrorFailedToUpdateModule, nil)
}

func (s service) DeleteModule(moduleId string) (resp response.GlobalResponse) {
	//TODO implement me
	panic("implement me")
}

func (s service) GetAllModule() (resp response.GlobalResponse) {
	//TODO implement me
	if dt, err := s.moduleMenuRepo.FindAllModules(); err == nil && dt != nil {
		return util.CreateGlobalResponse(constanta.Success, dt)
	}
	return util.CreateGlobalResponse(constanta.ErrorModuleDataNotFound, nil)
}

func (s service) GetModuleById(moduleId string) (resp response.GlobalResponse) {
	if module, err := s.moduleMenuRepo.FindModuleById(moduleId); module != nil && err == nil {
		return util.CreateGlobalResponse(constanta.Success, module)
	}
	return util.CreateGlobalResponse(constanta.ErrorModuleDataNotFound, nil)
}

func (s service) CreateMenu(data module_menu_req.CreateMenuReq) (resp response.GlobalResponse) {
	//TODO implement me
	var module *tbl_modules.Modules
	var parent *tbl_menus.Menus
	err := errors.New("")

	if module, err = s.moduleMenuRepo.FindModuleById(data.ModuleId); err != nil || module == nil {
		return util.CreateGlobalResponse(constanta.ErrorModuleDataNotFound, nil)
	}

	if data.ParentId != "" {
		if parent, err = s.moduleMenuRepo.FindMenuById(data.ParentId); err != nil || parent == nil {
			return util.CreateGlobalResponse(constanta.ErrorMenuDataNotFound, nil)
		}

		if parent.PkModuleId != module.Pk {
			return util.CreateGlobalResponse(constanta.ErrorMenuParentIdBelongTo, nil)
		}
	}

	var menuType = enum.Parent
	var parentId = "00000000-0000-0000-0000-000000000000"

	if data.ParentId != "" {
		menuType = enum.Child
		parentId = parent.MenuId.String

		if len(data.Permission) <= 0 {
			return util.CreateGlobalResponse(constanta.ErrorFailedToCreateMenu, nil)
		} else {
			for _, val := range data.Permission {
				if val.PermissionName == "" {
					return util.CreateGlobalResponse(constanta.ErrorFailedToCreateMenu, nil)
				}
			}
		}
	}

	menuCode := strings.ToUpper(string_util.TrimSpace(data.Name))

	if dt, err := s.moduleMenuRepo.FindMenuByCode(menuCode, parentId, module.Pk, menuType); dt != nil && err == nil {
		return util.CreateGlobalResponse(constanta.ErrorFailedToMenuExists, nil)
	}

	var id uuid.UUID

	if id, err = uuid.NewUUID(); err != nil {
		return util.CreateGlobalResponse(constanta.ErrorGeneral, nil)
	}

	if dt, err := s.moduleMenuRepo.SaveMenu(&tbl_menus.Menus{
		PkModuleId: module.Pk,
		MenuId:     null.StringFrom(id.String()),
		ParentId:   null.StringFrom(parentId),
		Name:       data.Name,
		Path:       data.Path,
		Type:       menuType,
		MenuCode:   menuCode,
		MenuIcon:   data.Icon,
	}); err == nil {

		var tempPerm []tbl_menu_permission.MenuPermission

		for _, val := range data.Permission {
			permCode := strings.ToUpper(string_util.TrimSpace(val.PermissionName))
			if permCode != "" {
				tempPerm = append(tempPerm, tbl_menu_permission.MenuPermission{PkMenu: dt.Pk, Name: val.PermissionName, PermCode: permCode})
			}
		}

		if sav, err := s.moduleMenuRepo.SaveMenuPermissionByMenuId(tempPerm); err == nil {
			result := module_menu_resp.MenuResponse{
				MenuId:     dt.MenuId,
				ParentId:   dt.ParentId,
				Name:       dt.Name,
				Path:       dt.Path,
				Type:       dt.Type,
				MenuCode:   dt.MenuCode,
				MenuIcon:   dt.MenuIcon,
				CreatedAt:  dt.CreatedAt,
				UpdatedAt:  dt.UpdatedAt,
				Permission: *sav,
			}

			return util.CreateGlobalResponse(constanta.Success, result)
		}
		return util.CreateGlobalResponse(constanta.Success, dt)
	}
	return util.CreateGlobalResponse(constanta.ErrorFailedToCreateMenu, nil)
}

func (s service) UpdateMenu(menuId string, data module_menu_req.CreateMenuReq) (resp response.GlobalResponse) {
	//TODO implement me
	var menu *tbl_menus.Menus
	err := errors.New("")

	if menu, err = s.moduleMenuRepo.FindMenuById(menuId); err != nil || menu == nil {
		return util.CreateGlobalResponse(constanta.ErrorMenuDataNotFound, nil)
	}

	var permAdd []tbl_menu_permission.MenuPermission
	var permRemove []tbl_menu_permission.MenuPermission

	if menuPermission, err := s.moduleMenuRepo.FindAllMenuPermissionByMenuPk(menu.Pk); err == nil {
		size := len(*menuPermission)

		for _, val := range data.Permission {
			permCode := strings.ToUpper(string_util.TrimSpace(val.PermissionName))

			if size > 0 {
				// check permission exist or no
				check := linq.From(*menuPermission).Where(func(i interface{}) bool {
					t := i.(tbl_menu_permission.MenuPermission)

					return t.PermCode == permCode
				}).Count()

				//add permission to add
				if check == 0 && val.State == enum.Add {
					permAdd = append(permAdd, tbl_menu_permission.MenuPermission{PkMenu: menu.Pk, Name: val.PermissionName, PermCode: permCode})
				}

				//add permission to remove
				if check > 0 && val.State == enum.Removed {
					permRemove = append(permRemove, tbl_menu_permission.MenuPermission{PkMenu: menu.Pk, PermCode: permCode})
				}

			} else {
				if val.State == enum.Add {
					permAdd = append(permAdd, tbl_menu_permission.MenuPermission{PkMenu: menu.Pk, Name: val.PermissionName, PermCode: permCode})
				}
			}
		}
	}

	var defParentId = "00000000-0000-0000-0000-000000000000"

	var chName = data.Name
	var chParentId = data.ParentId
	var chMenuType = menu.Type
	var chMenuIcon = data.Icon
	var chMenuPath = data.Path
	var menuCode = menu.MenuCode

	var flagCoreChange = false

	//change name
	if strings.ToUpper(menu.Name) != strings.ToUpper(chName) {
		chName = data.Name
		flagCoreChange = true
	}
	//change parent
	if menu.ParentId.String != chParentId {
		flagCoreChange = true
		if chParentId == "" {
			chParentId = defParentId
			chMenuType = enum.Parent
			if chParentId == menu.ParentId.String && chMenuType == menu.Type {
				flagCoreChange = false
			}
		}
		if chParentId != "" && chParentId != defParentId {
			chParentId = data.ParentId
			chMenuType = enum.Child
		}
	}
	//change icon
	if menu.MenuIcon != chMenuIcon {
		chMenuIcon = data.Icon
	}
	//change path
	if menu.Path != chMenuPath {
		chMenuPath = data.Path
	}

	//if (chName != data.Name) || (chParentId != data.ParentId) {
	//	if chName != data.Name {
	//		chName = data.Name
	//	}
	//	if chParentId != data.ParentId {
	//		chParentId = data.ParentId
	//		if chParentId == "" {
	//			chParentId = defParentId
	//			chMenuType = enum.Parent
	//		}
	//	}
	//}
	if flagCoreChange {
		menuCode = strings.ToUpper(string_util.TrimSpace(chName))

		if dt, err := s.moduleMenuRepo.FindMenuByCode(menuCode, chParentId, menu.PkModuleId, chMenuType); dt != nil && err == nil {
			return util.CreateGlobalResponse(constanta.ErrorFailedToMenuExists, nil)
		}
	}

	if err := s.moduleMenuRepo.UpdateMenu(&tbl_menus.Menus{
		MenuId:    menu.MenuId,
		ParentId:  null.StringFrom(chParentId),
		Name:      chName,
		Path:      chMenuPath,
		Type:      chMenuType,
		MenuCode:  menuCode,
		MenuIcon:  chMenuIcon,
		UpdatedAt: null.TimeFrom(time.Now()),
	}); err == nil {
		if len(permRemove) > 0 {
			if err = s.moduleMenuRepo.RemoveMenuPermissionByMenuId(permRemove); err == nil {

			}
		}
		if len(permAdd) > 0 {
			if dta, err := s.moduleMenuRepo.SaveMenuPermissionByMenuId(permAdd); err == nil && dta != nil {

			}
		}
		return util.CreateGlobalResponse(constanta.Success, nil)
	}
	return util.CreateGlobalResponse(constanta.ErrorFailedToUpdateMenu, nil)
}

func (s service) DeleteMenu(menuId string) (resp response.GlobalResponse) {
	//TODO implement me
	panic("implement me")
}

func (s service) GetAllMenu() (resp response.GlobalResponse) {
	//TODO implement me
	if dt, err := s.moduleMenuRepo.FindAllMenus(); err == nil {
		return util.CreateGlobalResponse(constanta.Success, dt)
	}
	return util.CreateGlobalResponse(constanta.ErrorMenuDataNotFound, nil)
}

func (s service) GetAllMenuByModuleId(moduleId string) (resp response.GlobalResponse) {
	//TODO implement me
	if modules, _ := s.moduleMenuRepo.FindModuleById(moduleId); modules != nil {
		var parentMenu []role_permission_resp.ShardMenu

		if menu, _ := s.moduleMenuRepo.FindParentMenuByModulePk(modules.Pk); menu != nil {
			for _, valMenu := range *menu {
				var shardMenu []role_permission_resp.ShardSubMenu

				//add parent permission
				var mastPerm []role_permission_resp.JoinedUserMenuPermissionMatrix
				if tempMasterPerm, _ := s.moduleMenuRepo.FindAllMenuPermissionByMenuPk(valMenu.Pk); tempMasterPerm != nil {
					for _, val := range *tempMasterPerm {
						mastPerm = append(mastPerm, role_permission_resp.JoinedUserMenuPermissionMatrix{
							Pk:       val.Pk,
							PkMenu:   val.PkMenu,
							Name:     val.Name,
							PermCode: val.PermCode,
						})
					}
				}

				//check if parent has any child
				if subMenu, _ := s.moduleMenuRepo.FindChildMenuByParentId(valMenu.MenuId.String); subMenu != nil {
					for _, val := range *subMenu {
						//add child permission
						var childPerm []role_permission_resp.JoinedUserMenuPermissionMatrix
						if tempChildPerm, _ := s.moduleMenuRepo.FindAllMenuPermissionByMenuPk(val.Pk); tempChildPerm != nil {
							for _, val3 := range *tempChildPerm {
								childPerm = append(childPerm, role_permission_resp.JoinedUserMenuPermissionMatrix{
									Pk:       val3.Pk,
									PkMenu:   val3.PkMenu,
									Name:     val3.Name,
									PermCode: val3.PermCode,
								})
							}
						}
						shardMenu = append(shardMenu, role_permission_resp.ShardSubMenu{
							ParentId:   val.ParentId.String,
							MenuId:     val.MenuId.String,
							Type:       string(val.Type),
							MenuName:   val.Name,
							Path:       val.Path,
							Icon:       val.MenuIcon,
							Permission: childPerm,
						})
					}
				}

				parentMenu = append(parentMenu, role_permission_resp.ShardMenu{
					ParentId:   valMenu.ParentId.String,
					MenuId:     valMenu.MenuId.String,
					Type:       string(valMenu.Type),
					Path:       valMenu.Path,
					MenuName:   valMenu.Name,
					Icon:       valMenu.MenuIcon,
					Permission: mastPerm,
					SubMenu:    shardMenu,
				})
			}
		}
		return util.CreateGlobalResponse(constanta.Success, parentMenu)
	}
	return util.CreateGlobalResponse(constanta.ErrorMenuDataNotFound, nil)
}

func (s service) GetMenuById(menuId string) (resp response.GlobalResponse) {
	//TODO implement me
	if dt, err := s.moduleMenuRepo.FindMenuById(menuId); err == nil {
		var perm *[]tbl_menu_permission.MenuPermission

		perm, err = s.moduleMenuRepo.FindAllMenuPermissionByMenuPk(dt.Pk)

		return util.CreateGlobalResponse(constanta.Success, module_menu_resp.MenuResponse{
			MenuId:     dt.MenuId,
			ParentId:   dt.ParentId,
			Name:       dt.Name,
			Path:       dt.Path,
			Type:       dt.Type,
			MenuCode:   dt.MenuCode,
			MenuIcon:   dt.MenuIcon,
			CreatedAt:  dt.CreatedAt,
			UpdatedAt:  dt.UpdatedAt,
			Permission: *perm,
		})
	}
	return util.CreateGlobalResponse(constanta.ErrorMenuDataNotFound, nil)

}

func (s service) GetMenuByPermission(roleId string) (resp response.GlobalResponse) {
	var role *tbl_roles.Roles

	err := errors.New("")

	if role, err = s.roleRepo.FindRoleById(roleId); err == nil && role == nil {
		return util.CreateGlobalResponse(constanta.ErrorRoleDataNotFound, nil)
	}

	var parentMenu []role_permission_resp.ShardMenu
	if tst, _ := s.moduleMenuRepo.FindParentMenuWithRegisteredStatus(role.Pk); tst != nil {
		for _, val := range *tst {
			var shardMenu []role_permission_resp.ShardSubMenu
			//add parent permission
			var mastPerm []role_permission_resp.JoinedUserMenuPermissionMatrix
			if tempMasterPerm, _ := s.roleRepo.FindPermissionWithRegisteredStatus(role.Pk, val.Pk); tempMasterPerm != nil {
				for _, valT := range *tempMasterPerm {
					mastPerm = append(mastPerm, role_permission_resp.JoinedUserMenuPermissionMatrix{
						Pk:               valT.Pk,
						PkMenu:           valT.PkMenu,
						Name:             valT.Name,
						PermCode:         valT.PermCode,
						RegisteredStatus: valT.RegisteredStatus,
					})
				}
			}
			if temp, _ := s.moduleMenuRepo.FindChildMenuWithRegisteredStatus(role.Pk, val.MenuId); temp != nil {
				for _, valT := range *temp {
					//add child permission
					var childPerm []role_permission_resp.JoinedUserMenuPermissionMatrix

					//if tempChildPerm, _ := s.moduleMenuRepo.FindAllMenuPermissionByMenuPk(val.Pk); tempChildPerm != nil {
					if tempChildPerm, _ := s.roleRepo.FindPermissionWithRegisteredStatus(role.Pk, valT.Pk); tempChildPerm != nil {
						for _, val3 := range *tempChildPerm {
							childPerm = append(childPerm, role_permission_resp.JoinedUserMenuPermissionMatrix{
								Pk:               val3.Pk,
								PkMenu:           val3.PkMenu,
								Name:             val3.Name,
								PermCode:         val3.PermCode,
								RegisteredStatus: val3.RegisteredStatus,
							})
						}
					}
					shardMenu = append(shardMenu, role_permission_resp.ShardSubMenu{
						ParentId:   valT.ParentId,
						MenuId:     valT.MenuId,
						Type:       valT.Type,
						MenuName:   valT.Name,
						Path:       valT.Path,
						Icon:       valT.MenuIcon,
						Registered: valT.RegisteredStatus,
						Permission: childPerm,
					})
				}
			}

			parentMenu = append(parentMenu, role_permission_resp.ShardMenu{
				ParentId:   val.ParentId,
				MenuId:     val.MenuId,
				Type:       val.Type,
				Path:       val.Path,
				MenuName:   val.Name,
				Icon:       val.MenuIcon,
				Registered: val.RegisteredStatus,
				Permission: mastPerm,
				SubMenu:    shardMenu,
			})
		}
		return util.CreateGlobalResponse(constanta.Success, parentMenu)
	}

	return util.CreateGlobalResponse(constanta.ErrorMenuDataNotFound, nil)
}

func (s service) GetUnregisteredMenuByUserId(userId int64) (resp response.GlobalResponse) {
	var user *tbl_users.Users

	err := errors.New("")

	if user, err = s.userRepo.FindById(userId); err == nil && user == nil {
		return util.CreateGlobalResponse(constanta.ErrorUserNotFound, nil)
	}

	var parentMenu []role_permission_resp.ShardMenu
	if tst, _ := s.moduleMenuRepo.FindAllParentMenuWithRegisteredStatus(user.Pk); tst != nil {
		for _, val := range *tst {
			var shardMenu []role_permission_resp.ShardSubMenu
			//add parent permission
			var mastPerm []role_permission_resp.JoinedUserMenuPermissionMatrix
			if tempMasterPerm, _ := s.roleRepo.FindAllPermissionWithRegisteredStatus(user.Pk, val.Pk); tempMasterPerm != nil {
				for _, valT := range *tempMasterPerm {
					mastPerm = append(mastPerm, role_permission_resp.JoinedUserMenuPermissionMatrix{
						Pk:               valT.Pk,
						PkMenu:           valT.PkMenu,
						Name:             valT.Name,
						PermCode:         valT.PermCode,
						RegisteredStatus: valT.RegisteredStatus,
					})
				}
			}
			if temp, _ := s.moduleMenuRepo.FindAllChildMenuWithRegisteredStatus(user.Pk, val.MenuId); temp != nil {
				for _, valT := range *temp {
					//add child permission
					var childPerm []role_permission_resp.JoinedUserMenuPermissionMatrix

					//if tempChildPerm, _ := s.moduleMenuRepo.FindAllMenuPermissionByMenuPk(val.Pk); tempChildPerm != nil {
					if tempChildPerm, _ := s.roleRepo.FindAllPermissionWithRegisteredStatus(user.Pk, valT.Pk); tempChildPerm != nil {
						for _, val3 := range *tempChildPerm {
							childPerm = append(childPerm, role_permission_resp.JoinedUserMenuPermissionMatrix{
								Pk:               val3.Pk,
								PkMenu:           val3.PkMenu,
								Name:             val3.Name,
								PermCode:         val3.PermCode,
								RegisteredStatus: val3.RegisteredStatus,
							})
						}
					}
					shardMenu = append(shardMenu, role_permission_resp.ShardSubMenu{
						ParentId:   valT.ParentId,
						MenuId:     valT.MenuId,
						Type:       valT.Type,
						MenuName:   valT.Name,
						Path:       valT.Path,
						Icon:       valT.MenuIcon,
						Registered: valT.RegisteredStatus,
						Permission: childPerm,
					})
				}
			}

			parentMenu = append(parentMenu, role_permission_resp.ShardMenu{
				ParentId:   val.ParentId,
				MenuId:     val.MenuId,
				Type:       val.Type,
				Path:       val.Path,
				MenuName:   val.Name,
				Icon:       val.MenuIcon,
				Registered: val.RegisteredStatus,
				Permission: mastPerm,
				SubMenu:    shardMenu,
			})
		}
		return util.CreateGlobalResponse(constanta.Success, parentMenu)
	}

	return util.CreateGlobalResponse(constanta.ErrorMenuDataNotFound, nil)
}
func (s service) GetAllPermittedMenu(userId int64) (resp response.GlobalResponse) {

	return
}

func (s service) TestGetAllMenu() (resp response.GlobalResponse) {
	//TODO implement me
	dtMenu, _ := s.moduleMenuRepo.FindALlPermittedMenu(52)
	dtPermission, _ := s.roleRepo.FindPermittedPermission(52)

	if dtMenu != nil {
		//distinct module
		var result []role_permission_resp.PermittedPermissionRes
		var tmpMod []role_permission_resp.JoinedPermittedMenu

		//distinct module
		linq.From(*dtMenu).DistinctBy(func(i interface{}) interface{} {
			return i.(role_permission_resp.JoinedPermittedMenu).PkModuleId
		}).ToSlice(&tmpMod)

		//loop the module
		for _, val := range tmpMod {
			var parentMenu []role_permission_resp.PermittedMenu
			fmt.Println(val.PkMenu)
			var tmpMenu []role_permission_resp.JoinedPermittedMenu

			//get menu by module id && parent type
			linq.From(*dtMenu).Where(func(i interface{}) bool {
				t := i.(role_permission_resp.JoinedPermittedMenu)
				return t.PkModuleId == val.PkModuleId && t.ModuleId == val.ModuleId && t.Type == string(enum.Parent)
			}).ToSlice(&tmpMenu)

			//loop the parent menu
			for _, val2 := range tmpMenu {
				var parentPerm []role_permission_resp.JoinedRoleMenuPermissionMatrix
				var subMenu []role_permission_resp.PermittedSubMenu

				//get permission for parent menu by pk menu
				linq.From(*dtPermission).Where(func(y interface{}) bool {
					t2 := y.(role_permission_resp.JoinedRoleMenuPermissionMatrix)
					return t2.PkMenu == val2.PkMenu
				}).ForEach(func(z interface{}) {
					t4 := z.(role_permission_resp.JoinedRoleMenuPermissionMatrix)
					parentPerm = append(parentPerm, role_permission_resp.JoinedRoleMenuPermissionMatrix{
						Pk:       t4.Pk,
						PkMenu:   t4.PkMenu,
						Name:     t4.Name,
						PermCode: t4.PermCode,
					})
				})
				//get menu by parent id && child type
				linq.From(*dtMenu).Where(func(i interface{}) bool {
					t := i.(role_permission_resp.JoinedPermittedMenu)

					return t.ParentId == val2.MenuId && t.Type == string(enum.Child)
				}).ForEach(func(x interface{}) {
					t1 := x.(role_permission_resp.JoinedPermittedMenu)
					var childPerm []role_permission_resp.JoinedRoleMenuPermissionMatrix

					//get permission for child menu by pk menu
					linq.From(*dtPermission).Where(func(y interface{}) bool {
						t2 := y.(role_permission_resp.JoinedRoleMenuPermissionMatrix)
						return t2.PkMenu == t1.PkMenu
					}).ForEach(func(z interface{}) {
						t4 := z.(role_permission_resp.JoinedRoleMenuPermissionMatrix)
						childPerm = append(childPerm, role_permission_resp.JoinedRoleMenuPermissionMatrix{
							Pk:       t4.Pk,
							PkMenu:   t4.PkMenu,
							Name:     t4.Name,
							PermCode: t4.PermCode,
						})
					})

					subMenu = append(subMenu, role_permission_resp.PermittedSubMenu{
						ParentId:   t1.ParentId,
						MenuId:     t1.MenuId,
						PkMenu:     t1.PkMenu,
						Name:       t1.Name,
						Type:       t1.Type,
						Path:       t1.Path,
						MenuIcon:   t1.MenuIcon,
						Permission: childPerm,
					})
				})
				//result menu
				parentMenu = append(parentMenu, role_permission_resp.PermittedMenu{
					ParentId:   val2.ParentId,
					MenuId:     val2.MenuId,
					PkMenu:     val2.PkMenu,
					Name:       val2.Name,
					Type:       val2.Type,
					Path:       val2.Path,
					MenuIcon:   val2.MenuIcon,
					Permission: parentPerm,
					SubMenu:    subMenu,
				})

			}
			//result permission
			result = append(result, role_permission_resp.PermittedPermissionRes{
				PkModuleId: val.PkModuleId,
				ModuleId:   val.ModuleId,
				ModuleName: val.ModuleName,
				Menu:       parentMenu,
			})
		}

		return util.CreateGlobalResponse(constanta.Success, result)
	}

	return util.CreateGlobalResponse(constanta.ErrorMenuDataNotFound, nil)
}
