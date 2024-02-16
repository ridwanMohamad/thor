package user_service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"thor/src/constanta"
	"thor/src/constanta/enum"
	"thor/src/domain/tbl_additional_privileges"
	"thor/src/domain/tbl_mapping_location_users"
	"thor/src/domain/tbl_mapping_user_roles"
	"thor/src/domain/tbl_menus"
	"thor/src/domain/tbl_roles"
	"thor/src/domain/tbl_session"
	"thor/src/domain/tbl_users"
	"thor/src/payload/response"
	"thor/src/payload/role_permission_resp"
	"thor/src/payload/user_req"
	"thor/src/payload/user_resp"
	"thor/src/properties"
	"thor/src/repository/lov"
	"thor/src/repository/modules_menus"
	"thor/src/repository/roles_permissions"
	"thor/src/repository/sessions"
	"thor/src/repository/users"
	"thor/src/server/mail"
	"thor/src/server/service"
	"thor/src/util"
	"thor/src/util/date_util"
	"thor/src/util/password_generator_util"

	"gopkg.in/gomail.v2"
	"gopkg.in/guregu/null.v4"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo    users.IUsersRepository
	menuRepo    modules_menus.IModulesMenusRepository
	roleRepo    roles_permissions.IRolesPermissionRepository
	lovRepo     lov.ILovRepository
	sessionRepo sessions.ISessionRepository
	mailer      *mail.GoMailType
}

func NewUserService(iUser users.IUsersRepository, menuRepo modules_menus.IModulesMenusRepository, roleRepo roles_permissions.IRolesPermissionRepository, lovRepo lov.ILovRepository, sessionRepo sessions.ISessionRepository, mail *mail.GoMailType) IUserService {

	return userService{userRepo: iUser, roleRepo: roleRepo, menuRepo: menuRepo, lovRepo: lovRepo, sessionRepo: sessionRepo, mailer: mail}
}

func (u userService) CreateNewUser(prop *properties.CustomContext, req user_req.UserDTO) (resp response.GlobalResponse) {
	//TODO implement me
	dt := &tbl_users.Users{}

	var err = errors.New("")
	pid := os.Getpid()
	fmt.Printf("pid : %s", pid)

	if dt, err = u.userRepo.FindByEmailAndPassword(req.Username, req.Email); err != nil || dt != nil {
		if dt != nil {
			if dt.Username == req.Username {
				return util.CreateGlobalResponse(constanta.ErrorDuplicateUser, "")
			}
			if dt.Email == req.Email {
				return util.CreateGlobalResponse(constanta.ErrorDuplicateEmailUser, "")
			}
		}
		if errors.Is(err, constanta.DbFailedToExecuteQuery) {
			return util.CreateGlobalResponse(constanta.ErrorFailedExecuteQuery, "")
		}
	}

	var dtRole *tbl_roles.Roles

	if dtRole, err = u.roleRepo.FindRoleById(req.Role); err != nil || dtRole == nil {
		if errors.Is(err, constanta.DbFailedToExecuteQuery) {
			return util.CreateGlobalResponse(constanta.ErrorFailedExecuteQuery, "")
		}
		if err == nil && dtRole == nil {
			return util.CreateGlobalResponse(constanta.ErrorRoleDataNotFound, "")
		}
	}

	var dtMenu *[]tbl_menus.Menus

	if dtMenu, err = u.menuRepo.FindAllMenus(); err != nil || dtMenu == nil {

		if errors.Is(err, constanta.DbFailedToExecuteQuery) {
			return util.CreateGlobalResponse(constanta.ErrorFailedExecuteQuery, "")
		}
		if err == nil && dtMenu == nil {
			return util.CreateGlobalResponse(constanta.ErrorMenuDataNotFound, "")
		}
	}

	var temp *[]user_req.AdditionalPrivilegeShard
	flag := true

	if flag, temp = removeDuplicateAndValidate(req.AdditionalPrivilege, dtMenu); flag && temp == nil {
		return util.CreateGlobalResponse(constanta.ErrorMenuNotFoundAddPrivilege, "")
	}

	if dt == nil {
		dt = &tbl_users.Users{}
	}

	// password := []byte(req.Password)

	// hashedPassword, err := bcrypt.GenerateFromPassword(password, 15)
	passwd := password_generator_util.GeneratePassword(8, 2, 2, 2)
	hash, err := bcrypt.GenerateFromPassword([]byte(passwd), 15)

	if err != nil {
		return util.CreateGlobalResponse(constanta.ErrorHashingPassword, "")
	}

	// if err != nil {
	// 	return util.CreateGlobalResponse(constanta.ErrorHashingPassword, "")
	// }

	//validate between 2 Dates
	effectiveDate, expiredAt, err := date_util.T1BelowT2(req.EffectiveAt, req.ExpiredAt)

	if err != nil {
		return util.CreateGlobalResponse(constanta.ErrorComparisonBetween2Date, "")
	}

	dt.Username = req.Username
	dt.Email = req.Email
	dt.MobilePhone = req.MobilePhone
	dt.FullName = req.FullName
	// dt.Password = string(hashedPassword)
	dt.Password = string(hash)
	dt.Department = req.Department
	dt.Location = req.Location
	dt.EffectiveAt = effectiveDate
	dt.ExpiredAt = expiredAt
	dt.Status = enum.Active
	dt.FkRole = dtRole.Pk
	dt.EmployeeId = req.EmployeeId

	final, err := u.userRepo.Save(dt)

	if err != nil {
		return util.CreateGlobalResponse(constanta.ErrorFailedToSaveUser, err)
	}
	var insertDefaultMappingLocation []tbl_mapping_location_users.MappingLocationUsers
	insertDefaultMappingLocation = append(insertDefaultMappingLocation, tbl_mapping_location_users.MappingLocationUsers{
		LocationId: final.Location,
		UserId:     final.Pk,
		IsDefault:  true,
	})
	if len(insertDefaultMappingLocation) > 0 {
		err := u.userRepo.UpdateInsertMappingLocationUsers(insertDefaultMappingLocation)
		if err != nil {
			return util.CreateGlobalResponse(constanta.ErrorFailedToSaveUser, err)
		}
	}
	var insertDefaultMappingRoles []tbl_mapping_user_roles.MappingRoleUsers
	insertDefaultMappingRoles = append(insertDefaultMappingRoles, tbl_mapping_user_roles.MappingRoleUsers{
		RoleId:    dtRole.RoleId,
		UserId:    final.Pk,
		IsDefault: true,
	})
	if len(insertDefaultMappingRoles) > 0 {
		err := u.userRepo.UpdateInsertMappingRoleUsers(insertDefaultMappingRoles)
		if err != nil {
			return util.CreateGlobalResponse(constanta.ErrorFailedToSaveUser, err)
		}
	}

	if len(*temp) > 0 {
		var additionalPrivilege []tbl_additional_privileges.AdditionalPrivileges

		for _, v := range *temp {
			for _, v2 := range v.Permission {
				additionalPrivilege = append(additionalPrivilege, tbl_additional_privileges.AdditionalPrivileges{
					PkUserId:     final.Pk,
					PkMenuId:     v.PkMenu,
					PkMenuPermId: v2.PermissionId,
				})
			}
		}

		if _, err = u.roleRepo.SaveAddPriv(&additionalPrivilege); err != nil {
			return util.CreateGlobalResponse(constanta.ErrorAddPrivFailedToAdd, nil)
		}
	}

	mailContent := u.composeRegisterEmail(final.Email, final.FullName, final.Username, string(passwd))

	if errm := u.mailer.GMail.DialAndSend(mailContent); errm != nil {
		prop.Context.Logger().Error(errm)
	}
	prop.Context.Logger().Info("send email success")

	return util.CreateGlobalResponse(constanta.Success, final)
}

func (a userService) composeRegisterEmail(mailTo string, recipientName string, username string, content string) *gomail.Message {
	r, err := ioutil.ReadFile("./resources/mail_content/register-mail.html")

	if err != nil {
		return nil
	}
	cont := string(r)
	cont = strings.Replace(cont, "{user-fullname}", recipientName, -1)
	cont = strings.Replace(cont, "{username}", username, -1)
	cont = strings.Replace(cont, "{password}", content, -1)
	cont = strings.Replace(cont, "{link-app}", "https://apps.titipaja.id", -1)

	//fmt.Println(r)
	m := gomail.NewMessage()
	m.SetHeader("To", mailTo)
	m.SetHeader("From", *a.mailer.MailFrom)
	m.SetHeader("Subject", "Information New Account")
	m.SetBody("text/html", cont)

	return m
}

func (u userService) GetUserById(id int64) (resp response.GlobalResponse) {
	//TODO implement me
	var result user_resp.UserDetailResp
	if user, err := u.userRepo.FindById(id); user != nil && err == nil {
		result.UserId = user.Pk
		result.Username = user.Username
		result.Email = user.Email
		result.FullName = user.FullName
		result.Location = user.Location
		result.Department = user.Department
		result.MobilePhone = user.MobilePhone
		result.Status = user.Status
		result.EffectiveAt = user.EffectiveAt.Time
		result.ExpiredAt = user.ExpiredAt.Time
		result.CreatedAt = user.CreatedAt
		result.UpdatedAt = user.UpdatedAt
		result.LastLoginAt = user.LastLoginAt
		result.IsLocked = user.IsLocked
		result.LockedAt = user.LockedAt
		result.EmployeeId = user.EmployeeId
		result.ProfilePict = user.ProfilePict

		if dt, err := u.roleRepo.FindRoleByPkId(user.FkRole); dt != nil && err == nil {
			result.Role = user_resp.CommonPkIdName{
				Pk:   dt.Pk,
				Id:   dt.RoleId,
				Name: dt.Name,
			}
		}
		if dt, err := u.lovRepo.FindLovDetailByDetailId(user.Department); dt != nil && err == nil {
			result.DepartmentName = dt.Name
		}
		if dt, err := u.lovRepo.FindLovDetailByDetailId(user.Location); dt != nil && err == nil {
			result.LocationName = dt.Name
		}
		if dt, err := u.roleRepo.FindMenuPrivByUserId(user.Pk); dt != nil && err == nil {
			var addPrivilege []user_resp.CommonPkIdName
			for _, val := range *dt {
				addPrivilege = append(addPrivilege, user_resp.CommonPkIdName{
					Pk:   val.Pk,
					Id:   val.MenuId,
					Name: val.Name,
				})
			}
			result.AdditionalPrivilege = addPrivilege
		}

		return util.CreateGlobalResponse(constanta.Success, result)
	}
	return util.CreateGlobalResponse(constanta.ErrorUserNotFound, nil)
}

func (u userService) GetUserByUsername(username string) (resp response.GlobalResponse) {
	//TODO implement me
	panic("implement me")
}

func (u userService) GetAllUser() (resp response.GlobalResponse) {
	//TODO implement me
	if dt, err := u.userRepo.FindAll(); err == nil || dt != nil {
		return util.CreateGlobalResponse(constanta.Success, dt)
	}
	return util.CreateGlobalResponse(constanta.ErrorUserNotFound, nil)
}

func (u userService) UpdateUser(req user_req.UserUpdateDTO) (resp response.GlobalResponse) {
	//TODO implement me
	dt := &tbl_users.Users{}
	var err = errors.New("")

	if dt, err = u.userRepo.FindByUsername(req.Username); err == nil && dt == nil {
		return util.CreateGlobalResponse(constanta.ErrorUserNotFound, "")
	}

	rl, _ := u.roleRepo.FindRoleByPkId(dt.FkRole)
	//var dtRole *tbl_roles.Roles

	//if dtRole, err = u.roleRepo.FindRoleById(req.Role); err != nil || dtRole == nil {
	//	if errors.Is(err, constanta.DbFailedToExecuteQuery) {
	//		return util.CreateGlobalResponse(constanta.ErrorFailedExecuteQuery, "")
	//	}
	//	if err == nil && dtRole == nil {
	//		return util.CreateGlobalResponse(constanta.ErrorRoleDataNotFound, "")
	//	}
	//}

	var dtMenu *[]tbl_menus.Menus

	if dtMenu, err = u.menuRepo.FindAllMenus(); err != nil || dtMenu == nil {

		if errors.Is(err, constanta.DbFailedToExecuteQuery) {
			return util.CreateGlobalResponse(constanta.ErrorFailedExecuteQuery, "")
		}
		if err == nil && dtMenu == nil {
			return util.CreateGlobalResponse(constanta.ErrorMenuDataNotFound, "")
		}
	}

	var tempAdditional *[]user_req.AdditionalPrivilegeShard
	//var tempAdd *[]user_req.AdditionalPrivilegeShard
	//var tempDel *[]user_req.AdditionalPrivilegeShard

	flag := true

	if flag, tempAdditional = removeDuplicateAndValidate(req.AdditionalPrivilege, dtMenu); flag && tempAdditional == nil {
		return util.CreateGlobalResponse(constanta.ErrorMenuNotFoundAddPrivilege, "")
	}

	//if flag, tempAdd = removeDuplicateAndValidate(req.AddPrivilege, dtMenu); flag && tempAdd == nil {
	//	return util.CreateGlobalResponse(constanta.ErrorMenuNotFoundAddPrivilege, "")
	//}
	//
	//if flag, tempDel = removeDuplicateAndValidate(req.RemovePrivilege, dtMenu); flag && tempDel == nil {
	//	return util.CreateGlobalResponse(constanta.ErrorMenuNotFoundAddPrivilege, "")
	//}

	if req.Password != "" {
		password := []byte(req.Password)

		hashedPassword, err := bcrypt.GenerateFromPassword(password, 15)

		if err != nil {
			return util.CreateGlobalResponse(constanta.ErrorHashingPassword, "")
		}

		//validate if password not equal with previous data then replace with new password
		if err = bcrypt.CompareHashAndPassword([]byte(dt.Password), password); err == nil {
			dt.Password = string(hashedPassword)
		}
	}
	//validate between 2 Dates
	effectiveDate, expiredAt, err := date_util.T1BelowT2(req.EffectiveAt, req.ExpiredAt)

	if err != nil {
		return util.CreateGlobalResponse(constanta.ErrorComparisonBetween2Date, "")
	}

	dt.EffectiveAt = effectiveDate
	dt.ExpiredAt = expiredAt

	//validate email with previous data
	if dt.Email != req.Email {
		dt.Email = req.Email
	}

	//validate mobile phone with previous data
	if dt.MobilePhone != req.MobilePhone {
		dt.MobilePhone = req.MobilePhone
	}

	//validate fullname with previous data
	if dt.FullName != req.FullName {
		dt.FullName = req.FullName
	}
	//validate status
	if string(dt.Status) != req.Status {
		if req.Status != string(enum.Active) && req.Status != string(enum.InActive) {
			return util.CreateGlobalResponse(constanta.ErrorInvalidDefaultStatus, "")
		}
		dt.Status = enum.StatusEnum(req.Status)
	}
	//validate department
	if dt.Department != req.Department {
		dt.Department = req.Department
	}
	//validate location
	if dt.Location != req.Location {
		dt.Location = req.Location
		err := u.userRepo.RemoveDefaultLocationUser(dt.Pk)
		if err != nil {
			return util.CreateGlobalResponse(constanta.ErrorFailedToUpdateUser, err)
		}
		var insertDefaultMappingLocation []tbl_mapping_location_users.MappingLocationUsers
		if dtMappingLocationUser, _ := u.userRepo.FindMappingLocationByUserIdandLocation(dt.Pk, req.Location); dtMappingLocationUser != nil {
			insertDefaultMappingLocation = append(insertDefaultMappingLocation, tbl_mapping_location_users.MappingLocationUsers{
				Id:         &dtMappingLocationUser.Id,
				UserId:     dtMappingLocationUser.UserId,
				LocationId: dtMappingLocationUser.LocationId,
				IsDefault:  true,
			})
		} else {
			insertDefaultMappingLocation = append(insertDefaultMappingLocation, tbl_mapping_location_users.MappingLocationUsers{
				LocationId: req.Location,
				UserId:     dt.Pk,
				IsDefault:  true,
			})
		}

		if len(insertDefaultMappingLocation) > 0 {
			err := u.userRepo.UpdateInsertMappingLocationUsers(insertDefaultMappingLocation)
			if err != nil {
				return util.CreateGlobalResponse(constanta.ErrorFailedToSaveUser, err)
			}
		}

	}
	//validate locked enable
	if req.Unlock && dt.IsLocked {
		dt.IsLocked = false
		dt.LockedAt = null.TimeFromPtr(nil)
	}
	//validate Employee Id
	if req.EmployeeId != dt.EmployeeId {
		dt.EmployeeId = req.EmployeeId
	}
	flagRoleUpdated := false
	//validate role is changed
	if req.RoleId != rl.RoleId {

		dtr, _ := u.roleRepo.FindRoleById(req.RoleId)
		dt.FkRole = dtr.Pk
		flagRoleUpdated = true

		err := u.userRepo.RemoveDefaultRoleUser(dt.Pk)
		if err != nil {
			return util.CreateGlobalResponse(constanta.ErrorFailedToUpdateUser, err)
		}
		var insertDefaultMappingRoles []tbl_mapping_user_roles.MappingRoleUsers

		if dtMappingRoleUser, _ := u.userRepo.FindMappingRoleByUserIdandRole(dt.Pk, req.RoleId); dtMappingRoleUser != nil {
			insertDefaultMappingRoles = append(insertDefaultMappingRoles, tbl_mapping_user_roles.MappingRoleUsers{
				Id:        &dtMappingRoleUser.Id,
				UserId:    dtMappingRoleUser.UserId,
				RoleId:    dtMappingRoleUser.RoleId,
				IsDefault: true,
			})
		} else {
			insertDefaultMappingRoles = append(insertDefaultMappingRoles, tbl_mapping_user_roles.MappingRoleUsers{
				RoleId:    req.RoleId,
				UserId:    dt.Pk,
				IsDefault: true,
			})
		}

		if len(insertDefaultMappingRoles) > 0 {
			err := u.userRepo.UpdateInsertMappingRoleUsers(insertDefaultMappingRoles)
			if err != nil {
				return util.CreateGlobalResponse(constanta.ErrorFailedToSaveUser, err)
			}
		}
	}
	//dt.FkRole = dtRole.Pk

	err = u.userRepo.Update(dt)

	if err != nil {
		return util.CreateGlobalResponse(constanta.ErrorFailedToSaveUser, err)
	}

	//userAddPrivileges, _ := u.roleRepo.FindMenuPrivByUserId(dt.Pk)

	//addPrivileges := addOrRemovePrivileges(dt.Pk, tempAdd, userAddPrivileges)
	//removePrivileges := addOrRemovePrivileges(dt.Pk, tempDel, userAddPrivileges)

	if req.Unlock {
		if err = u.userRepo.RemoveUserLocked(dt.Pk); err != nil {
			return util.CreateGlobalResponse(constanta.ErrorFailedToRemoveLockedData, nil)
		}
	}
	stateModifiedPermission := false
	stateInsert := false
	stateRemoved := false

	var addPermission []tbl_additional_privileges.AdditionalPrivileges
	var removedPermission []tbl_additional_privileges.AdditionalPrivileges

	if flagRoleUpdated == false {
		if len(*tempAdditional) > 0 {
			for _, v := range *tempAdditional {
				if v.State == enum.Modified {
					if userMatrix, err := u.roleRepo.FindUserMatrixByPkUserAndPkMenu(dt.Pk, v.PkMenu); userMatrix != nil && err == nil {
						for _, valA := range v.Permission {
							count := 0
							for _, valB := range *userMatrix {
								if valA.PermissionId == valB.Pk {
									count += 1
									break
								}
							}
							if valA.State == enum.Add && count == 0 {
								//add permission
								if tmp, err := u.menuRepo.FindAllMenuPermByMenuPkAndPermPk(v.PkMenu, valA.PermissionId); err == nil && tmp != nil {
									addPermission = append(addPermission, tbl_additional_privileges.AdditionalPrivileges{
										PkUserId:     dt.Pk,
										PkMenuId:     v.PkMenu,
										PkMenuPermId: valA.PermissionId,
									})
								} else {
									return util.CreateGlobalResponse(constanta.ErrorMenuPermissionNotExists, nil)
								}
							}
							if valA.State == enum.Removed && count == 1 {
								//remove permission
								removedPermission = append(removedPermission, tbl_additional_privileges.AdditionalPrivileges{
									PkUserId:     dt.Pk,
									PkMenuId:     v.PkMenu,
									PkMenuPermId: valA.PermissionId,
								})
							}
							if len(addPermission) > 0 || len(removedPermission) > 0 {
								stateModifiedPermission = true
							}
						}
					} else {
						for _, valA := range v.Permission {
							if tmp, err := u.menuRepo.FindAllMenuPermByMenuPkAndPermPk(v.PkMenu, valA.PermissionId); err == nil && tmp != nil {
								addPermission = append(addPermission, tbl_additional_privileges.AdditionalPrivileges{
									PkUserId:     dt.Pk,
									PkMenuId:     tmp.PkMenu,
									PkMenuPermId: tmp.Pk,
								})
							} else {
								return util.CreateGlobalResponse(constanta.ErrorMenuPermissionNotExists, nil)
							}
						}
						if len(addPermission) > 0 {
							stateModifiedPermission = true
						}
					}

				}
				if v.State == enum.Removed {
					if perm, err := u.roleRepo.FindUserMatrixByPkUserAndPkMenu(dt.Pk, v.PkMenu); perm != nil && err == nil {
						for _, val := range *perm {
							removedPermission = append(removedPermission, tbl_additional_privileges.AdditionalPrivileges{
								PkUserId:     dt.Pk,
								PkMenuId:     val.PkMenu,
								PkMenuPermId: val.Pk,
							})
						}
						stateRemoved = true
					} else {
						return util.CreateGlobalResponse(constanta.ErrorPermissionNotFound, nil)
					}
				}
				if v.State == enum.Add {
					if perm, err := u.roleRepo.FindUserMatrixByPkUserAndPkMenu(dt.Pk, v.PkMenu); perm != nil && err == nil {
						continue
					}

					if menu, err := u.menuRepo.FindMenuById(v.MenuId); menu != nil || err == nil {

						for _, valA := range v.Permission {
							if tmp, err := u.menuRepo.FindAllMenuPermByMenuPkAndPermPk(menu.Pk, valA.PermissionId); err == nil && tmp != nil {
								addPermission = append(addPermission, tbl_additional_privileges.AdditionalPrivileges{
									PkUserId:     dt.Pk,
									PkMenuId:     tmp.PkMenu,
									PkMenuPermId: tmp.Pk,
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
			}

			if stateRemoved {
				if len(removedPermission) > 0 && err == nil {
					if err = u.roleRepo.DeleteAddPriv(&removedPermission); err != nil {
						return util.CreateGlobalResponse(constanta.ErrorRemovedMenuPermission, nil)
					}
				}
			}
			var skip = false

			if stateInsert && stateModifiedPermission {
				skip = true
			}

			if stateInsert {
				if len(addPermission) > 0 && err == nil {
					if _, err = u.roleRepo.SaveAddPriv(&addPermission); err != nil {
						return util.CreateGlobalResponse(constanta.ErrorAddMenuPermission, nil)
					}
				}
			}

			if stateModifiedPermission {
				if len(removedPermission) > 0 {
					if err = u.roleRepo.DeleteAddPriv(&removedPermission); err != nil {
						return util.CreateGlobalResponse(constanta.ErrorRemovedMenuPermission, nil)
					}
				}
				if len(addPermission) > 0 && skip == false {
					if _, err = u.roleRepo.SaveAddPriv(&addPermission); err != nil {
						return util.CreateGlobalResponse(constanta.ErrorAddMenuPermission, nil)
					}
				}
			}
		}
		//if addPrivileges != nil {
		//	if _, err = u.roleRepo.SaveAddPriv(addPrivileges); err != nil {
		//		return util.CreateGlobalResponse(constanta.ErrorAddPrivFailedToAdd, nil)
		//	}
		//}
		//if removePrivileges != nil {
		//	if err = u.roleRepo.DeleteAddPriv(removePrivileges); err != nil {
		//		return util.CreateGlobalResponse(constanta.ErrorAddPrivFailedToRemove, nil)
		//	}
		//}
	}

	if flagRoleUpdated == true {
		if err = u.roleRepo.DeleteAddPrivByUserPk(dt.Pk); err != nil {
			return util.CreateGlobalResponse(constanta.ErrorAddPrivFailedToRemove, nil)
		}
	}

	return util.CreateGlobalResponse(constanta.Success, dt)
}

func (u userService) UpdateUserProfile(req user_req.UpdateUserProfile) (resp response.GlobalResponse) {
	ctx := context.Background()
	dt := &tbl_users.Users{}
	var err = errors.New("")
	if dt, err = u.userRepo.FindByUsername(req.Username); err == nil && dt == nil {
		return util.CreateGlobalResponse(constanta.ErrorUserNotFound, "")
	}

	// rl, _ := u.roleRepo.FindRoleByPkId(dt.FkRole)

	if dt.Email != req.Email {
		dt.Email = req.Email
	}

	if dt.MobilePhone != req.MobilePhone {
		dt.MobilePhone = req.MobilePhone
	}

	if dt.EmployeeId != req.EmployeeId {
		dt.EmployeeId = req.EmployeeId
	}

	// if req.RoleId != rl.RoleId {
	// 	dtr, _ := u.roleRepo.FindRoleById(req.RoleId)

	// 	if dtr == nil {
	// 		return util.CreateGlobalResponse(constanta.ErrorRoleDataNotFound, nil)
	// 	}
	// 	dt.FkRole = dtr.Pk
	// 	// flagRoleUpdated = true
	// }

	// if dt.Location != req.Location {
	// 	dt.Location = req.Location
	// }

	if dt.ProfilePict != req.ProfilePict {
		if req.ProfilePict != "" {
			profilePict := req.ProfilePict

			i := strings.Index(profilePict, ",")
			if i < 0 {
				return util.CreateGlobalResponse(constanta.ErrorFileNotFound, nil)
			}

			// pass reader to NewDecoder
			dec := base64.NewDecoder(base64.StdEncoding, strings.NewReader(profilePict[i+1:]))

			filename := os.Getenv("ENV") + "/user-profile/" + req.Username

			err = service.GcpClient.UploadBase64(ctx, filename, dec)
			if err != nil {
				return util.CreateGlobalResponse(constanta.ErrorFileNotFound, err.Error())
			}

			dt.ProfilePict = service.GcpClient.URL(filename)
		} else {
			dt.ProfilePict = ""
		}
	}

	err = u.userRepo.Update(dt)

	if err != nil {
		return util.CreateGlobalResponse(constanta.ErrorFailedToSaveUser, err)
	}

	return util.CreateGlobalResponse(constanta.Success, dt)

}

func (u userService) ChangePassword(req user_req.UserChangePasswordReq) (resp response.GlobalResponse) {
	var sess *tbl_session.Session
	err := errors.New("")

	if sess, err = u.sessionRepo.FindByAccessToken(req.SessionId); sess == nil && (err != nil || err == nil) {
		return util.CreateGlobalResponse(constanta.ErrorInvalidAccessToken, nil)
	}

	var user *tbl_users.Users

	if user, err = u.userRepo.FindById(req.UserPk); user == nil && (err != nil || err == nil) {
		return util.CreateGlobalResponse(constanta.ErrorUserNotFound, nil)
	}

	if user.IsLocked {
		return util.CreateGlobalResponse(constanta.ErrorAccountLocked, nil)
	}

	if user.Pk != sess.FkUser {
		return util.CreateGlobalResponse(constanta.ErrorUserNotFound, nil)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return util.CreateGlobalResponse(constanta.ErrorInvalidOldPassword, nil)
	}

	if req.NewPassword != req.ReConfirmNewPassword {
		return util.CreateGlobalResponse(constanta.ErrorNewAndReconfirmDifferent, nil)
	}

	if req.OldPassword == req.NewPassword {
		return util.CreateGlobalResponse(constanta.ErrorNewAndOldSamePassword, nil)
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 15)
	if err != nil {
		return util.CreateGlobalResponse(constanta.ErrorGeneral, nil)
	}

	user.Password = string(pass)
	user.IsFirstLogin = false

	if err = u.userRepo.Update(user); err != nil {
		return util.CreateGlobalResponse(constanta.ErrorChangePasswordFailed, nil)

	}

	return util.CreateGlobalResponse(constanta.Success, nil)
}

func addOrRemovePrivileges(userPk int64, req *[]user_req.AdditionalPrivilegeShard, data *[]role_permission_resp.JoinedPrivilegeMenu) (resp *[]tbl_additional_privileges.AdditionalPrivileges) {
	if len(*req) > 0 {
		if data != nil {
			var priv []tbl_additional_privileges.AdditionalPrivileges
			for _, val := range *req {
				var pkMenu int64
				exists := 0

				for _, val2 := range *data {
					if val.MenuId == val2.MenuId {
						pkMenu = val2.Pk
						exists += 1
						break
					}
				}

				if exists == 0 || exists >= 1 {
					if len(val.Permission) <= 0 {
						for _, val2 := range val.Permission {
							priv = append(priv, tbl_additional_privileges.AdditionalPrivileges{
								PkUserId:     userPk,
								PkMenuId:     pkMenu,
								PkMenuPermId: val2.PermissionId,
							})
						}
					} else {
						priv = append(priv, tbl_additional_privileges.AdditionalPrivileges{
							PkUserId: userPk,
							PkMenuId: pkMenu,
						})
					}
				}

				return &priv
			}
		} else {
			var priv []tbl_additional_privileges.AdditionalPrivileges
			for _, val := range *req {
				if len(val.Permission) <= 0 {
					for _, val2 := range val.Permission {
						priv = append(priv, tbl_additional_privileges.AdditionalPrivileges{
							PkUserId:     userPk,
							PkMenuId:     val.PkMenu,
							PkMenuPermId: val2.PermissionId,
						})
					}
				} else {
					priv = append(priv, tbl_additional_privileges.AdditionalPrivileges{
						PkUserId:     userPk,
						PkMenuId:     val.PkMenu,
						PkMenuPermId: 0,
					})
				}

			}
			return &priv
		}
	}
	return
}

func removeDuplicateAndValidate(req []user_req.AdditionalPrivilegeShard, data *[]tbl_menus.Menus) (flag bool, dt *[]user_req.AdditionalPrivilegeShard) {

	temp := map[string]bool{}
	res := []user_req.AdditionalPrivilegeShard{}

	//remove duplicate
	for v := range req {
		if temp[req[v].MenuId] != true {
			temp[req[v].MenuId] = true
			res = append(res, req[v])
		}
	}

	//validating
	result := []user_req.AdditionalPrivilegeShard{}
	for v := range res {
		count := 0
		for _, v2 := range *data {
			if res[v].MenuId == v2.MenuId.String {
				count += 1
				result = append(result, user_req.AdditionalPrivilegeShard{
					PkMenu:     v2.Pk,
					MenuId:     v2.MenuId.String,
					State:      res[v].State,
					Permission: res[v].Permission,
				})
				break
			}
		}
		if count <= 0 {
			return true, nil
		}
	}
	return false, &result
}

func (u userService) InActiveUser(id int64) (resp response.GlobalResponse) {
	dt := &tbl_users.Users{}
	var err = errors.New("")

	if dt, err = u.userRepo.FindById(id); err == nil && dt == nil {
		return util.CreateGlobalResponse(constanta.ErrorUserNotFound, "")
	}

	if err = u.userRepo.InActiveUser(id); err != nil {
		return util.CreateGlobalResponse(constanta.ErrorFailedInActiveUser, nil)
	}

	return util.CreateGlobalResponse(constanta.Success, dt)
}

func (u userService) GetMappingLocationByUsername(username string) (resp response.GlobalResponse) {
	//TODO implement me
	var result user_resp.MappingUserLocationResp
	if user, err := u.userRepo.FindByUsername(username); user != nil && err == nil {
		result.Username = user.Username
		result.Email = user.Email
		result.FullName = user.FullName
		if dt, err := u.userRepo.FindMappingLocationByUserId(user.Username); dt != nil && err == nil {
			var userLocation []user_resp.UserLocation
			for _, val := range *dt {
				userLocation = append(userLocation, user_resp.UserLocation{
					Id:             val.Id,
					LocationId:     val.LocationId,
					LocationName:   val.LocationName,
					IsDefault:      val.IsDefault,
					LocationValue1: val.LocationValue1,
					LocationValue2: val.LocationValue2,
					LocationValue3: val.LocationValue3,
				})
			}
			result.Location = userLocation
		}

		return util.CreateGlobalResponse(constanta.Success, result)
	}
	return util.CreateGlobalResponse(constanta.ErrorUserNotFound, nil)
}

func (u userService) UpdateMappingLocationUser(req user_req.MappingUserLocationReq) (resp response.GlobalResponse) {
	//TODO implement me
	dt := &tbl_users.Users{}
	var err = errors.New("")

	if dt, err = u.userRepo.FindByUsername(req.Username); err == nil && dt == nil {
		return util.CreateGlobalResponse(constanta.ErrorUserNotFound, "")
	}

	var dataMappingLocationUser []tbl_mapping_location_users.MappingLocationUsers
	for _, v := range req.Location {
		dataMappingLocationUser = append(dataMappingLocationUser, tbl_mapping_location_users.MappingLocationUsers{
			Id:         v.Id,
			UserId:     dt.Pk,
			LocationId: v.LocationId,
			IsDefault:  v.IsDefault,
		})
	}
	err = u.userRepo.UpdateInsertMappingLocationUsers(dataMappingLocationUser)

	if err != nil {
		return util.CreateGlobalResponse(constanta.ErrorFailedToSaveUser, err)
	}

	result := u.GetMappingLocationByUsername(dt.Username)

	return result
}
func (u userService) RemoveMappingLocationUser(Id int64) (resp response.GlobalResponse) {
	//TODO implement me
	dt, err := u.userRepo.FindMappingLocationById(Id)
	if err != nil {
		return util.CreateGlobalResponse(constanta.ErrorFailedToSaveUser, err)
	}

	if dtUser, errUser := u.userRepo.FindById(dt.UserId); errUser != nil || dtUser != nil {
		if dt != nil {
			if dtUser.Location == dt.LocationId {
				return util.CreateGlobalResponse(constanta.ErrorLocationInUsed, nil)
			}
		}
		if errors.Is(err, constanta.DbFailedToExecuteQuery) {
			return util.CreateGlobalResponse(constanta.ErrorFailedExecuteQuery, nil)
		}
	}

	err = u.userRepo.RemoveMappingLocationUser(Id)

	if err != nil {
		return util.CreateGlobalResponse(constanta.ErrorFailedToSaveUser, err)
	}

	return util.CreateGlobalResponse(constanta.Success, true)
}
func (u userService) RemoveDefaultLocationUser(userId int64) (resp response.GlobalResponse) {
	//TODO implement me
	dt, err := u.userRepo.FindMappingLocationDefaultByUserId(userId)
	if err != nil {
		return util.CreateGlobalResponse(constanta.ErrorFailedToSaveUser, err)
	}

	err = u.userRepo.RemoveDefaultLocationUser(dt.UserId)

	if err != nil {
		return util.CreateGlobalResponse(constanta.ErrorFailedToSaveUser, err)
	}

	return util.CreateGlobalResponse(constanta.Success, true)
}
func (u userService) GetMappingRoleByUsername(username string) (resp response.GlobalResponse) {
	//TODO implement me
	var result user_resp.MappingUserRoleResp
	if user, err := u.userRepo.FindByUsername(username); user != nil && err == nil {
		result.Username = user.Username
		result.Email = user.Email
		result.FullName = user.FullName
		if dt, err := u.userRepo.FindMappingRoleByUserId(user.Username); dt != nil && err == nil {
			var userRole []user_resp.UserRole
			for _, val := range *dt {
				userRole = append(userRole, user_resp.UserRole{
					Id:        val.Id,
					RoleId:    val.RoleId,
					RoleName:  val.RoleName,
					IsDefault: val.IsDefault,
				})
			}
			result.Roles = userRole
		}

		return util.CreateGlobalResponse(constanta.Success, result)
	}
	return util.CreateGlobalResponse(constanta.ErrorUserNotFound, nil)
}

func (u userService) UpdateMappingRoleUser(req user_req.MappingUserRoleReq) (resp response.GlobalResponse) {
	//TODO implement me
	dt := &tbl_users.Users{}
	var err = errors.New("")

	if dt, err = u.userRepo.FindByUsername(req.Username); err == nil && dt == nil {
		return util.CreateGlobalResponse(constanta.ErrorUserNotFound, "")
	}

	var dataMappingRoleUser []tbl_mapping_user_roles.MappingRoleUsers
	for _, v := range req.Roles {
		dataMappingRoleUser = append(dataMappingRoleUser, tbl_mapping_user_roles.MappingRoleUsers{
			Id:        v.Id,
			UserId:    dt.Pk,
			RoleId:    v.RoleId,
			IsDefault: v.IsDefault,
		})
	}
	err = u.userRepo.UpdateInsertMappingRoleUsers(dataMappingRoleUser)

	if err != nil {
		return util.CreateGlobalResponse(constanta.ErrorFailedToSaveUser, err)
	}

	result := u.GetMappingRoleByUsername(dt.Username)

	return result
}
func (u userService) RemoveMappingRoleUser(Id int64) (resp response.GlobalResponse) {
	//TODO implement me
	dt, err := u.userRepo.FindMappingRoleById(Id)
	if err != nil {
		return util.CreateGlobalResponse(constanta.ErrorFailedToSaveUser, err)
	}

	if dtUser, errUser := u.userRepo.FindById(dt.UserId); errUser != nil || dtUser != nil {
		if dtUser != nil {
			if dtRole, errUser := u.roleRepo.FindRoleByPkId(dtUser.FkRole); errUser != nil || dtUser != nil {
				if dtRole != nil {
					if dtRole.RoleId == dt.RoleId {
						return util.CreateGlobalResponse(constanta.ErrorLocationInUsed, nil)
					}
				}
			}
		}
		if errors.Is(err, constanta.DbFailedToExecuteQuery) {
			return util.CreateGlobalResponse(constanta.ErrorFailedExecuteQuery, nil)
		}
	}

	err = u.userRepo.RemoveMappingRoleUser(Id)

	if err != nil {
		return util.CreateGlobalResponse(constanta.ErrorFailedToSaveUser, err)
	}

	return util.CreateGlobalResponse(constanta.Success, true)
}
func (u userService) RemoveDefaultRoleUser(userId int64) (resp response.GlobalResponse) {
	//TODO implement me
	dt, err := u.userRepo.FindMappingRoleDefaultByUserId(userId)
	if err != nil {
		return util.CreateGlobalResponse(constanta.ErrorFailedToSaveUser, err)
	}

	err = u.userRepo.RemoveDefaultRoleUser(dt.UserId)

	if err != nil {
		return util.CreateGlobalResponse(constanta.ErrorFailedToSaveUser, err)
	}

	return util.CreateGlobalResponse(constanta.Success, true)
}
