package role_permission_service

import (
	"errors"
	"thor/src/constanta"
	"thor/src/constanta/enum"
	"thor/src/domain/tbl_additional_privileges"
	"thor/src/domain/tbl_menus"
	"thor/src/domain/tbl_role_menu"
	"thor/src/domain/tbl_role_menu_matrix"
	"thor/src/domain/tbl_roles"
	"thor/src/domain/tbl_users"
	"thor/src/payload/response"
	"thor/src/payload/role_permission_req"
	"thor/src/payload/role_permission_resp"
	"thor/src/repository/modules_menus"
	"thor/src/repository/roles_permissions"
	"thor/src/repository/users"
	"thor/src/util"

	"github.com/google/uuid"
)

type roleService struct {
	rolePermissionRepo roles_permissions.IRolesPermissionRepository
	menuRepo           modules_menus.IModulesMenusRepository
	userRepo           users.IUsersRepository
}

func NewRoleService(roleRepo roles_permissions.IRolesPermissionRepository,
	menuRepo modules_menus.IModulesMenusRepository,
	userRepo users.IUsersRepository) IRoleService {
	return &roleService{rolePermissionRepo: roleRepo, menuRepo: menuRepo, userRepo: userRepo}
}

func (r roleService) CreateNewRole(req role_permission_req.RoleWithPermissionReq) (resp response.GlobalResponse) {
	//TODO implement me
	var role *tbl_roles.Roles
	var menu []tbl_menus.Menus
	var permMenu []tbl_role_menu_matrix.RoleMenuMatrix

	var id uuid.UUID
	err := errors.New("")
	var menuId = ""
	for _, val := range req.Permission {
		if val.MenuId == menuId {
			return util.CreateGlobalResponse(constanta.ErrorMenuIdDuplicate, nil)
		}
		if val.MenuId == "" {
			return util.CreateGlobalResponse(constanta.ErrorMenuPartialDataNotFound, nil)
		}
		if (val.State == enum.Removed || val.State == enum.Add) == false {
			return util.CreateGlobalResponse(constanta.ErrorRoleMenuStateNotValid, nil)
		}

		if dt, err := r.menuRepo.FindMenuById(val.MenuId); err == nil && dt != nil {

			for _, val2 := range val.Permission {
				if tmp, err := r.menuRepo.FindAllMenuPermByMenuPkAndPermPk(dt.Pk, val2.PermId); err == nil && tmp != nil {
					permMenu = append(permMenu, tbl_role_menu_matrix.RoleMenuMatrix{
						PkMenu: tmp.PkMenu, PkMenuPerm: tmp.Pk,
					})
				} else {
					return util.CreateGlobalResponse(constanta.ErrorMenuPermissionNotExists, nil)
				}
			}

			menu = append(menu, tbl_menus.Menus{
				Pk: dt.Pk,
			})
		} else {
			return util.CreateGlobalResponse(constanta.ErrorMenuPartialDataNotFound, nil)
		}
		menuId = val.MenuId
	}

	if id, err = uuid.NewRandom(); err != nil {
		return util.CreateGlobalResponse(constanta.ErrorGeneral, nil)
	}

	var status enum.StatusEnum

	if req.Status == "inactive" {
		status = enum.InActive
	} else {
		status = enum.Active
	}

	if role, err = r.rolePermissionRepo.SaveRole(&tbl_roles.Roles{
		RoleId: id.String(),
		Name:   req.Name,
		Status: status,
	}); err == nil {

		var permission []tbl_role_menu.RoleMenu
		var permissionMenu []tbl_role_menu_matrix.RoleMenuMatrix

		for _, val := range menu {
			permission = append(permission, tbl_role_menu.RoleMenu{
				PkRole: role.Pk,
				PkMenu: val.Pk,
			})
		}

		for _, val := range permMenu {
			permissionMenu = append(permissionMenu, tbl_role_menu_matrix.RoleMenuMatrix{
				PkRole:     role.Pk,
				PkMenu:     val.PkMenu,
				PkMenuPerm: val.PkMenuPerm,
			})
		}

		if result, err := r.rolePermissionRepo.SaveRolesMenus(&permission); err == nil {
			var finalRoleMenu []role_permission_resp.FinalJoinedRoleMenu
			if resultMatrix, err := r.rolePermissionRepo.SaveRoleMatrix(&permissionMenu); err == nil {
				for i := 0; i < len(*result); i++ {
					var matrix []role_permission_resp.JoinedRoleMenuPermissionMatrix
					roleMenu := (*result)[i]
					for _, val := range *resultMatrix {
						if roleMenu.PkMenu == val.PkMenu {
							matrix = append(matrix, role_permission_resp.JoinedRoleMenuPermissionMatrix{
								Pk:       val.Pk,
								PkMenu:   val.PkMenu,
								Name:     val.Name,
								PermCode: val.PermCode,
							})
						}
					}
					finalRoleMenu = append(finalRoleMenu, role_permission_resp.FinalJoinedRoleMenu{
						PkRole:     roleMenu.PkRole,
						PkMenu:     roleMenu.PkMenu,
						MenuId:     roleMenu.MenuId,
						Name:       roleMenu.Name,
						Permission: matrix,
					})
				}
			}
			res := role_permission_resp.RolePermissionResp{
				RoleId:         role.RoleId,
				Name:           role.Name,
				Status:         string(role.Status),
				PermissionMenu: finalRoleMenu,
			}
			return util.CreateGlobalResponse(constanta.Success, res)
		}
	}
	return util.CreateGlobalResponse(constanta.ErrorFailedToCreateRole, nil)
}

func (r roleService) UpdateRoleAndPermission(roleId string, req role_permission_req.RoleWithPermissionReq) (resp response.GlobalResponse) {
	//TODO implement me
	var data *tbl_roles.Roles
	err := errors.New("")

	if data, err = r.rolePermissionRepo.FindRoleById(roleId); data == nil && (err != nil || err == nil) {
		return util.CreateGlobalResponse(constanta.ErrorRoleDataNotFound, nil)
	}

	var status enum.StatusEnum

	if req.Status == "inactive" {
		status = enum.InActive
	} else {
		status = enum.Active
	}

	var flagUpdate = false

	if data.Status != status || req.Name != data.Name {
		flagUpdate = true
	}
	if flagUpdate {
		if err = r.rolePermissionRepo.UpdateRole(tbl_roles.Roles{
			RoleId: data.RoleId,
			Name:   req.Name,
			Status: status,
		}); err != nil {
			return util.CreateGlobalResponse(constanta.ErrorFailedToUpdateRole, nil)
		}
	}

	var stateRemove = false
	var stateInsert = false
	var stateModifiedPermission = false
	var removedRoleMenu []tbl_role_menu.RoleMenu
	var insertRoleMenu []tbl_role_menu.RoleMenu
	var addPermission []tbl_role_menu_matrix.RoleMenuMatrix
	var removedPermission []tbl_role_menu_matrix.RoleMenuMatrix

	for _, val := range req.Permission {
		if val.State == "removed" {
			var perm *tbl_role_menu.RoleMenu

			if perm, err = r.rolePermissionRepo.FindRolesMenusByRoleIdAndMenuId(data.RoleId, val.MenuId); perm != nil && err == nil {
				removedRoleMenu = append(removedRoleMenu, tbl_role_menu.RoleMenu{
					PkRole: perm.PkRole,
					PkMenu: perm.PkMenu,
				})

				var menuPermission *[]role_permission_resp.JoinedRoleMenuPermissionMatrix

				if menuPermission, err = r.rolePermissionRepo.FindRoleMatrixByPkRoleAndPkMenu(perm.PkRole, perm.PkMenu); menuPermission != nil && err == nil {
					for _, valA := range *menuPermission {
						removedPermission = append(removedPermission, tbl_role_menu_matrix.RoleMenuMatrix{
							PkRole:     perm.PkRole,
							PkMenu:     perm.PkMenu,
							PkMenuPerm: valA.Pk,
						})
					}
				}
				stateRemove = true
			} else {
				return util.CreateGlobalResponse(constanta.ErrorPermissionNotFound, nil)
			}
		}
		if val.State == "add" {
			var menu *tbl_menus.Menus
			var perm *tbl_role_menu.RoleMenu

			if perm, err = r.rolePermissionRepo.FindRolesMenusByRoleIdAndMenuId(data.RoleId, val.MenuId); perm != nil && err == nil {
				continue
			}

			if menu, err = r.menuRepo.FindMenuById(val.MenuId); menu != nil || err == nil {
				insertRoleMenu = append(insertRoleMenu, tbl_role_menu.RoleMenu{
					PkRole: data.Pk,
					PkMenu: menu.Pk,
				})
				for _, valA := range val.Permission {
					if tmp, err := r.menuRepo.FindAllMenuPermByMenuPkAndPermPk(menu.Pk, valA.PermId); err == nil && tmp != nil {
						addPermission = append(addPermission, tbl_role_menu_matrix.RoleMenuMatrix{
							PkMenu: tmp.PkMenu, PkMenuPerm: tmp.Pk, PkRole: data.Pk,
						})
					} else {
						return util.CreateGlobalResponse(constanta.ErrorMenuPermissionNotExists, nil)
					}
				}
				stateInsert = true
			} else {
				return util.CreateGlobalResponse(constanta.ErrorMenuDataNotFound, nil)
			}
		}
		if val.State == "modified" {
			var menuPermission *[]role_permission_resp.JoinedRoleMenuPermissionMatrix
			var perm *tbl_role_menu.RoleMenu

			if perm, err = r.rolePermissionRepo.FindRolesMenusByRoleIdAndMenuId(data.RoleId, val.MenuId); perm != nil && err == nil {
				if menuPermission, err = r.rolePermissionRepo.FindRoleMatrixByPkRoleAndPkMenu(perm.PkRole, perm.PkMenu); menuPermission != nil && err == nil {
					for _, valA := range val.Permission {
						count := 0
						for _, valB := range *menuPermission {
							if valA.PermId == valB.Pk {
								count += 1
								break
							}
						}
						if valA.State == enum.Add && count == 0 {
							//add permission
							if tmp, err := r.menuRepo.FindAllMenuPermByMenuPkAndPermPk(perm.PkMenu, valA.PermId); err == nil && tmp != nil {
								addPermission = append(addPermission, tbl_role_menu_matrix.RoleMenuMatrix{
									PkRole:     perm.PkRole,
									PkMenu:     perm.PkMenu,
									PkMenuPerm: valA.PermId,
								})
							} else {
								return util.CreateGlobalResponse(constanta.ErrorMenuPermissionNotExists, nil)
							}
						}
						if valA.State == enum.Removed && count == 1 {
							//remove permission
							removedPermission = append(removedPermission, tbl_role_menu_matrix.RoleMenuMatrix{
								PkRole:     perm.PkRole,
								PkMenu:     perm.PkMenu,
								PkMenuPerm: valA.PermId,
							})
						}
					}

					if len(addPermission) > 0 || len(removedPermission) > 0 {
						stateModifiedPermission = true
					}
				} else {
					for _, valA := range val.Permission {
						if tmp, err := r.menuRepo.FindAllMenuPermByMenuPkAndPermPk(perm.PkMenu, valA.PermId); err == nil && tmp != nil {
							addPermission = append(addPermission, tbl_role_menu_matrix.RoleMenuMatrix{
								PkMenu: tmp.PkMenu, PkMenuPerm: tmp.Pk, PkRole: perm.PkRole,
							})
						} else {
							return util.CreateGlobalResponse(constanta.ErrorMenuPermissionNotExists, nil)
						}
					}
					if len(addPermission) > 0 {
						stateModifiedPermission = true
					}
				}
			} else {
				return util.CreateGlobalResponse(constanta.ErrorPermissionNotFound, nil)
			}
		}
	}

	if stateRemove {
		if err = r.rolePermissionRepo.DeleteRolesMenus(&removedRoleMenu); err != nil {
			return util.CreateGlobalResponse(constanta.ErrorFailedToRemovePermission, nil)
		}
		if len(removedPermission) > 0 && err == nil {
			if err = r.rolePermissionRepo.RemoveRoleMatrix(&removedPermission); err != nil {
				return util.CreateGlobalResponse(constanta.ErrorRemovedMenuPermission, nil)
			}
		}
	}
	var skip = false

	if stateInsert && stateModifiedPermission {
		skip = true
	}

	if stateInsert {
		if _, err = r.rolePermissionRepo.SaveRolesMenus(&insertRoleMenu); err != nil {
			return util.CreateGlobalResponse(constanta.ErrorFailedToAddPermission, nil)
		}
		if len(addPermission) > 0 && err == nil {
			if _, err = r.rolePermissionRepo.SaveRoleMatrix(&addPermission); err != nil {
				return util.CreateGlobalResponse(constanta.ErrorAddMenuPermission, nil)
			}
		}
	}

	if stateModifiedPermission {
		if len(removedPermission) > 0 {
			if err = r.rolePermissionRepo.RemoveRoleMatrix(&removedPermission); err != nil {
				return util.CreateGlobalResponse(constanta.ErrorRemovedMenuPermission, nil)
			}
		}
		if len(addPermission) > 0 && skip == false {
			if _, err = r.rolePermissionRepo.SaveRoleMatrix(&addPermission); err != nil {
				return util.CreateGlobalResponse(constanta.ErrorAddMenuPermission, nil)
			}
		}
	}

	return util.CreateGlobalResponse(constanta.Success, nil)
}

func (r roleService) DeleteRole(roleId string) (resp response.GlobalResponse) {
	//TODO implement me
	panic("implement me")
}

func (r roleService) GetAllRole() (resp response.GlobalResponse) {
	//TODO implement me
	var data *[]tbl_roles.Roles
	err := errors.New("")

	if data, err = r.rolePermissionRepo.FindAllRole(); err != nil {
		return util.CreateGlobalResponse(constanta.ErrorRoleDataNotFound, nil)
	}
	return util.CreateGlobalResponse(constanta.Success, data)
}

func (r roleService) GetAllRoleById(roleId string) (resp response.GlobalResponse) {
	//TODO implement me
	var data *tbl_roles.Roles
	err := errors.New("")

	if data, err = r.rolePermissionRepo.FindRoleById(roleId); err != nil {
		return util.CreateGlobalResponse(constanta.ErrorRoleDataNotFound, nil)
	}
	return util.CreateGlobalResponse(constanta.Success, data)
}

func (r roleService) GetPermissionByRoleId(roleId string) (resp response.GlobalResponse) {
	var role *tbl_roles.Roles
	err := errors.New("")

	if role, err = r.rolePermissionRepo.FindRoleById(roleId); role == nil && err == nil {
		return util.CreateGlobalResponse(constanta.ErrorRoleDataNotFound, nil)
	}

	if dt, err := r.rolePermissionRepo.FindRolesMenusByPkRole(role.Pk); err == nil && dt != nil {
		var res role_permission_resp.RolePermissionResp

		if role, err = r.rolePermissionRepo.FindRoleByPkId(role.Pk); err == nil && role != nil {
			//var menu []tbl_menus.Menus

			var perm []role_permission_resp.ShardPermission
			if modules, _ := r.rolePermissionRepo.FindModuleByPkRole(role.Pk); modules != nil {
				for _, valMod := range *modules {
					var parentMenu []role_permission_resp.ShardMenu
					if menu, _ := r.rolePermissionRepo.FindMenuParentByPkRole(role.Pk, valMod.Pk); menu != nil {
						for _, valMenu := range *menu {
							var shardMenu []role_permission_resp.ShardSubMenu

							var parentPerm []role_permission_resp.JoinedUserMenuPermissionMatrix
							if permission, _ := r.rolePermissionRepo.FindRoleMatrixByPkRoleAndPkMenu(role.Pk, valMenu.Pk); permission != nil {
								for _, val3 := range *permission {
									parentPerm = append(parentPerm, role_permission_resp.JoinedUserMenuPermissionMatrix{
										Pk:       val3.Pk,
										PkMenu:   val3.PkMenu,
										Name:     val3.Name,
										PermCode: val3.PermCode,
									})
								}
							}
							if subMenu, _ := r.rolePermissionRepo.FindMenuChildByParentId(role.Pk, valMenu.MenuId); subMenu != nil {
								for _, val2 := range *subMenu {
									var shardPermission []role_permission_resp.JoinedUserMenuPermissionMatrix
									if permission, _ := r.rolePermissionRepo.FindRoleMatrixByPkRoleAndPkMenu(role.Pk, val2.Pk); permission != nil {
										for _, val3 := range *permission {
											shardPermission = append(shardPermission, role_permission_resp.JoinedUserMenuPermissionMatrix{
												Pk:       val3.Pk,
												PkMenu:   val3.PkMenu,
												Name:     val3.Name,
												PermCode: val3.PermCode,
											})
										}
									}
									shardMenu = append(shardMenu, role_permission_resp.ShardSubMenu{
										ParentId:   val2.ParentId,
										MenuId:     val2.MenuId,
										Type:       val2.Type,
										MenuName:   val2.Name,
										Path:       val2.Path,
										Icon:       val2.MenuIcon,
										Permission: shardPermission,
									})
								}
							}

							parentMenu = append(parentMenu, role_permission_resp.ShardMenu{
								ParentId:   valMenu.ParentId,
								MenuId:     valMenu.MenuId,
								Type:       valMenu.Type,
								Path:       valMenu.Path,
								MenuName:   valMenu.Name,
								Icon:       valMenu.MenuIcon,
								Permission: parentPerm,
								SubMenu:    shardMenu,
							})
						}
					}
					perm = append(perm, role_permission_resp.ShardPermission{
						ModuleId:        valMod.ModuleId,
						Name:            valMod.Name,
						MaintenanceMode: valMod.MaintenanceMode,
						MtcStart:        valMod.MtcStart,
						MtcEnd:          valMod.MtcEnd,
						Menu:            parentMenu,
					})
				}
			}

			res.Name = role.Name
			res.RoleId = role.RoleId
			res.PermissionMenu = perm

			return util.CreateGlobalResponse(constanta.Success, res)
		}
	}
	return util.CreateGlobalResponse(constanta.ErrorPermissionDataNotFound, nil)
}

func (r roleService) GetPermissionByUserId(userId int64, roleId string) (resp response.GlobalResponse) {
	//TODO implement me
	var user *tbl_users.Users
	var err = errors.New("")

	if user, err = r.userRepo.FindById(userId); err == nil && user == nil {
		return util.CreateGlobalResponse(constanta.ErrorUserNotFound, nil)
	}

	if roleId != "" {
		if getRole, err := r.rolePermissionRepo.FindRoleById(roleId); err == nil && getRole != nil {
			user.FkRole = getRole.Pk
		}
	}

	if dt, err := r.rolePermissionRepo.FindRolesMenusByPkRole(user.FkRole); err == nil && dt != nil {
		var res role_permission_resp.RolePermissionResp

		if role, err := r.rolePermissionRepo.FindRoleByPkId(user.FkRole); err == nil && role != nil {
			//var menu []tbl_menus.Menus

			var perm []role_permission_resp.ShardPermission
			if modules, _ := r.rolePermissionRepo.FindModuleByPkRole(role.Pk); modules != nil {
				for _, valMod := range *modules {
					var parentMenu []role_permission_resp.ShardMenu
					if menu, _ := r.rolePermissionRepo.FindMenuParentByPkRole(role.Pk, valMod.Pk); menu != nil {
						for _, valMenu := range *menu {
							var shardMenu []role_permission_resp.ShardSubMenu

							var parentPerm []role_permission_resp.JoinedUserMenuPermissionMatrix
							if permission, _ := r.rolePermissionRepo.FindRoleMatrixByPkRoleAndPkMenu(role.Pk, valMenu.Pk); permission != nil {
								for _, val3 := range *permission {
									parentPerm = append(parentPerm, role_permission_resp.JoinedUserMenuPermissionMatrix{
										Pk:       val3.Pk,
										PkMenu:   val3.PkMenu,
										Name:     val3.Name,
										PermCode: val3.PermCode,
									})
								}
							}

							if subMenu, _ := r.rolePermissionRepo.FindMenuChildByParentId(role.Pk, valMenu.MenuId); subMenu != nil {
								for _, val2 := range *subMenu {
									var shardPermission []role_permission_resp.JoinedUserMenuPermissionMatrix
									if permission, _ := r.rolePermissionRepo.FindRoleMatrixByPkRoleAndPkMenu(role.Pk, val2.Pk); permission != nil {
										for _, val3 := range *permission {
											shardPermission = append(shardPermission, role_permission_resp.JoinedUserMenuPermissionMatrix{
												Pk:       val3.Pk,
												PkMenu:   val3.PkMenu,
												Name:     val3.Name,
												PermCode: val3.PermCode,
											})
										}
									}
									shardMenu = append(shardMenu, role_permission_resp.ShardSubMenu{
										ParentId:   val2.ParentId,
										MenuId:     val2.MenuId,
										Type:       val2.Type,
										MenuName:   val2.Name,
										Path:       val2.Path,
										Icon:       val2.MenuIcon,
										Permission: shardPermission,
									})
								}
							}

							parentMenu = append(parentMenu, role_permission_resp.ShardMenu{
								ParentId:   valMenu.ParentId,
								MenuId:     valMenu.MenuId,
								Type:       valMenu.Type,
								Path:       valMenu.Path,
								MenuName:   valMenu.Name,
								Icon:       valMenu.MenuIcon,
								Permission: parentPerm,
								SubMenu:    shardMenu,
							})
						}
					}
					perm = append(perm, role_permission_resp.ShardPermission{
						ModuleId:        valMod.ModuleId,
						Name:            valMod.Name,
						MaintenanceMode: valMod.MaintenanceMode,
						MtcStart:        valMod.MtcStart,
						MtcEnd:          valMod.MtcEnd,
						Menu:            parentMenu,
					})
				}
			}

			var addMenu []role_permission_resp.ShardSubMenu

			if additionalPriv, _ := r.rolePermissionRepo.FindMenuPrivByUserId(userId); additionalPriv != nil {
				for _, val := range *additionalPriv {
					var shardPermission []role_permission_resp.JoinedUserMenuPermissionMatrix
					if permission, _ := r.rolePermissionRepo.FindUserMatrixByPkUserAndPkMenu(user.Pk, val.Pk); permission != nil {
						for _, val3 := range *permission {
							shardPermission = append(shardPermission, role_permission_resp.JoinedUserMenuPermissionMatrix{
								Pk:       val3.Pk,
								PkMenu:   val3.PkMenu,
								Name:     val3.Name,
								PermCode: val3.PermCode,
							})
						}
					}
					addMenu = append(addMenu, role_permission_resp.ShardSubMenu{
						ParentId:   val.ParentId,
						MenuId:     val.MenuId,
						Type:       val.Type,
						MenuName:   val.Name,
						Path:       val.Path,
						Permission: shardPermission,
					})
				}
			}

			res.Name = role.Name
			res.RoleId = role.RoleId
			res.PermissionMenu = perm
			res.AdditionalPrivilege = addMenu

			return util.CreateGlobalResponse(constanta.Success, res)
		}
	}
	return util.CreateGlobalResponse(constanta.ErrorPermissionDataNotFound, nil)
}

func (r roleService) AddAdditionalMenuToUser(req []role_permission_req.ShardAdditionalPrivilege) (resp response.GlobalResponse) {
	//TODO implement me
	var temp []tbl_additional_privileges.AdditionalPrivileges

	for _, val := range req {
		temp = append(temp, tbl_additional_privileges.AdditionalPrivileges{
			PkUserId: val.UserId,
			PkMenuId: val.MenuId,
		})
	}

	if dt, err := r.rolePermissionRepo.SaveAddPriv(&temp); err == nil {
		return util.CreateGlobalResponse(constanta.Success, dt)
	}
	return util.CreateGlobalResponse(constanta.ErrorAddPrivFailedToAdd, nil)
}

func (r roleService) RemoveAdditionalMenuFromUser(req []role_permission_req.ShardAdditionalPrivilege) (resp response.GlobalResponse) {
	//TODO implement me
	var temp []tbl_additional_privileges.AdditionalPrivileges

	for _, val := range req {
		temp = append(temp, tbl_additional_privileges.AdditionalPrivileges{
			PkUserId: val.UserId,
			PkMenuId: val.MenuId,
		})
	}

	if err := r.rolePermissionRepo.DeleteAddPriv(&temp); err == nil {
		return util.CreateGlobalResponse(constanta.Success, nil)
	}
	return util.CreateGlobalResponse(constanta.ErrorAddPrivFailedToRemove, nil)
}

//func checkDuplicateData()
