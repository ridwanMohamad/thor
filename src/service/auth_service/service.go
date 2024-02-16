package auth_service

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"thor/src/constanta"
	"thor/src/constanta/enum"
	"thor/src/domain/tbl_session"
	"thor/src/domain/tbl_users"
	"thor/src/payload/auth_req"
	"thor/src/payload/auth_resp"
	"thor/src/payload/response"
	"thor/src/payload/role_permission_resp"
	"thor/src/payload/user_resp"
	"thor/src/properties"
	"thor/src/repository/lov"
	"thor/src/repository/modules_menus"
	"thor/src/repository/roles_permissions"
	"thor/src/repository/sessions"
	"thor/src/repository/users"
	"thor/src/server/config"
	"thor/src/server/mail"
	"thor/src/util"
	"thor/src/util/crypto_util"
	"thor/src/util/password_generator_util"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type auth struct {
	mailer      *mail.GoMailType
	sessionRepo sessions.ISessionRepository
	userRepo    users.IUsersRepository
	menuRepo    modules_menus.IModulesMenusRepository
	roleRepo    roles_permissions.IRolesPermissionRepository
	lovRepo     lov.ILovRepository
	authConfig  config.AuthConfig
}

func NewAuthService(authRepo sessions.ISessionRepository,
	userRepo users.IUsersRepository,
	mail *mail.GoMailType,
	menuRepo modules_menus.IModulesMenusRepository,
	roleRepo roles_permissions.IRolesPermissionRepository,
	lovRepo lov.ILovRepository,
	authConfig config.AuthConfig) IAuthService {

	return &auth{
		sessionRepo: authRepo,
		userRepo:    userRepo,
		mailer:      mail,
		menuRepo:    menuRepo,
		roleRepo:    roleRepo,
		lovRepo:     lovRepo,
		authConfig:  authConfig,
	}
}

func (a auth) Login(req auth_req.LoginReq) (resp response.GlobalResponse) {
	user := &user_resp.UserJoinRoleResp{}
	err := errors.New("")
	if user, err = a.userRepo.FindUserAndRoleByUsername(req.Username); err == nil && user == nil {
		return util.CreateGlobalResponse(constanta.ErrorFailedAuthenticate, nil)
	}

	if err != nil {
		return util.CreateGlobalResponse(constanta.ErrorFailedAuthenticate, nil)
	}

	if user != nil {
		if user.EffectiveAt.After(time.Now()) {
			return util.CreateGlobalResponse(constanta.ErrorAccountNotEffective, nil)
		}
		if user.ExpiredAt.Before(time.Now()) {
			return util.CreateGlobalResponse(constanta.ErrorAccountExpired, nil)
		}
		if user.IsLocked {
			return util.CreateGlobalResponse(constanta.ErrorAccountLocked, nil)
		}

		if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			if lockedCount, err := a.userRepo.CountUserLocked(user.UserPk); err == nil && lockedCount >= a.authConfig.MaximumLoginFailed {
				a.userRepo.UpdateLoginData(&tbl_users.Users{
					Pk:       user.UserPk,
					IsLocked: true,
					LockedAt: null.TimeFrom(time.Now()),
				})

				return util.CreateGlobalResponse(constanta.ErrorAccountLocked, map[string]interface{}{
					"failedLogin":  lockedCount,
					"maximumLogin": a.authConfig.MaximumLoginFailed,
				})
			}
			a.userRepo.SaveUserLocked(user.UserPk)
			return util.CreateGlobalResponse(constanta.ErrorFailedAuthenticate, nil)
		}

		session := &tbl_session.Session{}
		if session, err = a.sessionRepo.FindByUserId(user.UserPk); err != nil {
			if errors.Is(err, constanta.DbFailedToExecuteQuery) {
				return util.CreateGlobalResponse(constanta.ErrorFailedExecuteQuery, nil)
			}
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return util.CreateGlobalResponse(constanta.ErrorFailedAuthenticate, nil)
			}
		}

		var randomId uuid.UUID

		if randomId, err = uuid.NewRandom(); err != nil {
			return util.CreateGlobalResponse(constanta.ErrorGeneral, err)
		}

		pid := strconv.Itoa(os.Getpid())
		AccessToken := crypto_util.ComputeHmac(randomId.String(), pid)
		fmt.Printf("%s:%s", "uuid", randomId)
		fmt.Printf("%s:%s", "pid", pid)
		fmt.Printf("%s:%s", "access token", AccessToken)

		var duration time.Duration
		if duration, err = time.ParseDuration(a.authConfig.SessionMin.String()); err != nil {
			return util.CreateGlobalResponse(constanta.ErrorGeneral, err)
		}
		//if session already exists
		if session != nil {
			//dt.SessionId = session.SessionId
			//dt.UserId = session.UserId
			session.AccessToken = AccessToken
			//dt.ExpiredAt = session.ExpiredAt
			//dt.CreatedAt = session.CreatedAt
			session.UpdatedAt = time.Now()
			expiredAt := null.NewTime(time.Now().Add(duration), true)
			//check expired is valid
			if session.ExpiredAt.Valid {
				if time.Now().After(session.ExpiredAt.Time) {
					session.ExpiredAt = expiredAt
				}
			} else {
				session.ExpiredAt = expiredAt
			}

			if err = a.sessionRepo.Update(session); err == constanta.DbFailedToUpdateData {
				return util.CreateGlobalResponse(constanta.ErrorGeneral, err)
			}

			if err = a.userRepo.UpdateLastLoginAt(session); err == constanta.DbFailedToUpdateData {
				return util.CreateGlobalResponse(constanta.ErrorGeneral, err)
			}

			return util.CreateGlobalResponse(constanta.Success, auth_resp.LoginResponse{
				AccessToken: session.AccessToken,
				ExpiredAt:   session.ExpiredAt,
				UserDetail: auth_resp.UserDetail{
					UserId:   user.UserPk,
					UserName: user.UserName,
					FullName: user.FullName,
					RolePk:   user.RolePk,
					RoleId:   user.RoleId,
					RoleName: user.RoleName,
				},
			})
		}

		//if session not exists
		dt := tbl_session.Session{}

		dt.FkUser = user.UserPk
		dt.AccessToken = AccessToken
		dt.ExpiredAt = null.NewTime(time.Now().Add(duration), true)

		if _, err = a.sessionRepo.Save(&dt); err != nil {
			return util.CreateGlobalResponse(constanta.ErrorGeneral, err)
		}

		return util.CreateGlobalResponse(constanta.Success, auth_resp.LoginResponse{
			AccessToken: dt.AccessToken,
			ExpiredAt:   dt.ExpiredAt,
			UserDetail: auth_resp.UserDetail{
				UserId:   user.UserPk,
				UserName: user.UserName,
				FullName: user.FullName,
				RolePk:   user.RolePk,
				RoleId:   user.RoleId,
				RoleName: user.RoleName,
			},
		})
	}
	return
}

func (a auth) LoginV2(req auth_req.LoginReq) (resp response.GlobalResponse) {
	user := &user_resp.UserJoinRoleResp{}
	var permission []role_permission_resp.PermittedPermissionRes
	var userDetail auth_resp.UserDetailV2
	err := errors.New("")
	if user, err = a.userRepo.FindUserAndRoleByUsername(req.Username); err == nil && user == nil {
		return util.CreateGlobalResponse(constanta.ErrorFailedAuthenticate, nil)
	}

	if err != nil {
		return util.CreateGlobalResponse(constanta.ErrorFailedAuthenticate, nil)
	}

	if user != nil {
		if user.EffectiveAt.After(time.Now()) {
			return util.CreateGlobalResponse(constanta.ErrorAccountNotEffective, nil)
		}
		if user.ExpiredAt.Before(time.Now()) {
			return util.CreateGlobalResponse(constanta.ErrorAccountExpired, nil)
		}
		if user.IsLocked {
			return util.CreateGlobalResponse(constanta.ErrorAccountLocked, nil)
		}

		if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			if lockedCount, err := a.userRepo.CountUserLocked(user.UserPk); err == nil && lockedCount >= a.authConfig.MaximumLoginFailed {
				a.userRepo.UpdateLoginData(&tbl_users.Users{
					Pk:       user.UserPk,
					IsLocked: true,
					LockedAt: null.TimeFrom(time.Now()),
				})

				return util.CreateGlobalResponse(constanta.ErrorAccountLocked, map[string]interface{}{
					"failedLogin":  lockedCount,
					"maximumLogin": a.authConfig.MaximumLoginFailed,
				})
			}
			a.userRepo.SaveUserLocked(user.UserPk)
			return util.CreateGlobalResponse(constanta.ErrorFailedAuthenticate, nil)
		}

		session := &tbl_session.Session{}
		if session, err = a.sessionRepo.FindByUserId(user.UserPk); err != nil {
			if errors.Is(err, constanta.DbFailedToExecuteQuery) {
				return util.CreateGlobalResponse(constanta.ErrorFailedExecuteQuery, nil)
			}
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return util.CreateGlobalResponse(constanta.ErrorFailedAuthenticate, nil)
			}
		}

		var randomId uuid.UUID

		if randomId, err = uuid.NewRandom(); err != nil {
			return util.CreateGlobalResponse(constanta.ErrorGeneral, err)
		}

		pid := strconv.Itoa(os.Getpid())
		AccessToken := crypto_util.ComputeHmac(randomId.String(), pid)
		fmt.Printf("%s:%s", "uuid", randomId)
		fmt.Printf("%s:%s", "pid", pid)
		fmt.Printf("%s:%s", "access token", AccessToken)

		var duration time.Duration
		if duration, err = time.ParseDuration(a.authConfig.SessionMin.String()); err != nil {
			return util.CreateGlobalResponse(constanta.ErrorGeneral, err)
		}

		permission = a.generatePermission(user.UserPk)
		location := a.generateLocation(user.Location)

		userDetail.Profile = auth_resp.UserProfile{
			UserId:      user.UserPk,
			UserName:    user.UserName,
			FullName:    user.UserName,
			ProfilePict: user.ProfilePict,
		}
		userDetail.Role = auth_resp.UserRole{
			Pk:       user.RolePk,
			RoleId:   user.RoleId,
			RoleName: user.RoleName,
		}

		if len(location) > 0 {
			userDetail.Location = location
		}

		if len(permission) > 0 {
			userDetail.Permission = permission
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
			"profile":    userDetail.Profile,
			"role":       userDetail.Role,
			"location":   userDetail.Location,
			"permission": userDetail.Permission,
		})

		//if session already exists
		if session != nil {
			session.AccessToken = AccessToken
			session.UpdatedAt = time.Now()
			expiredAt := null.NewTime(time.Now().Add(duration), true)
			//check expired is valid
			if session.ExpiredAt.Valid {
				if time.Now().After(session.ExpiredAt.Time) {
					session.ExpiredAt = expiredAt
				}
			} else {
				session.ExpiredAt = expiredAt
			}

			if err = a.sessionRepo.Update(session); err == constanta.DbFailedToUpdateData {
				return util.CreateGlobalResponse(constanta.ErrorGeneral, err)
			}

			if err = a.userRepo.UpdateLastLoginAt(session); err == constanta.DbFailedToUpdateData {
				return util.CreateGlobalResponse(constanta.ErrorGeneral, err)
			}

			tokenString, err := token.SignedString([]byte(session.AccessToken))

			if err != nil {
				fmt.Println(err)
			}

			return util.CreateGlobalResponse(constanta.Success, auth_resp.LoginResponseV2{
				SessionToken: session.AccessToken,
				ExpiredAt:    session.ExpiredAt,
				UserDetail:   tokenString,
				IsFirstLogin: user.IsFirstLogin,
			})
		}

		//if session not exists
		dt := tbl_session.Session{}

		dt.FkUser = user.UserPk
		dt.AccessToken = AccessToken
		dt.ExpiredAt = null.NewTime(time.Now().Add(duration), true)

		if _, err = a.sessionRepo.Save(&dt); err != nil {
			return util.CreateGlobalResponse(constanta.ErrorGeneral, err)
		}

		tokenString, err := token.SignedString([]byte(AccessToken))

		if err != nil {
			fmt.Println(err)
		}

		return util.CreateGlobalResponse(constanta.Success, auth_resp.LoginResponseV2{
			SessionToken: dt.AccessToken,
			ExpiredAt:    dt.ExpiredAt,
			UserDetail:   tokenString,
			IsFirstLogin: user.IsFirstLogin,
		})
	}
	return
}

func (a auth) Logout(accessToken string) (resp response.GlobalResponse) {
	//TODO implement me
	session := &tbl_session.Session{}
	err := errors.New("")
	if session, err = a.sessionRepo.FindByAccessToken(accessToken); err == nil && session == nil {
		return util.CreateGlobalResponse(constanta.ErrorInvalidAccessToken, nil)
	}
	if session != nil {

		session.AccessToken = ""
		session.ExpiredAt = null.TimeFromPtr(nil)
		session.UpdatedAt = time.Now()

		if err = a.sessionRepo.Update(session); err == constanta.DbFailedToUpdateData {
			return util.CreateGlobalResponse(constanta.ErrorInvalidAccessToken, nil)
		}
	}
	return util.CreateGlobalResponse(constanta.Success, nil)
}

func (a auth) ForgotPassword(prop *properties.CustomContext, email string) (resp response.GlobalResponse) {
	prop.Context.Logger().Info("fooo")
	if dt, err := a.userRepo.FindByEmail(email); dt != nil && err == nil {
		passwd := password_generator_util.GeneratePassword(8, 2, 2, 2)
		hash, err := bcrypt.GenerateFromPassword([]byte(passwd), 15)

		if err != nil {
			return util.CreateGlobalResponse(constanta.ErrorHashingPassword, "")
		}

		if sess, err := a.sessionRepo.FindByUserId(dt.Pk); sess != nil && err == nil {
			sess.AccessToken = ""
			sess.UpdatedAt = time.Now()

			if err = a.sessionRepo.Update(sess); err != nil {
				return util.CreateGlobalResponse(constanta.ErrorFailedToUpdateSession, nil)
			}
		}

		dt.Password = string(hash)
		dt.UpdatedAt = null.TimeFrom(time.Now())
		if dt.IsLocked {
			dt.IsLocked = false
			dt.LockedAt = null.TimeFromPtr(nil)
		}

		if err = a.userRepo.Update(dt); err != nil {
			return util.CreateGlobalResponse(constanta.ErrorFailedToUpdateUser, nil)
		}

		if err = a.userRepo.RemoveUserLocked(dt.Pk); err != nil {
			return util.CreateGlobalResponse(constanta.ErrorFailedToRemoveLockedData, nil)
		}

		mailContent := a.composeForgotEmail(dt.Email, dt.FullName, dt.Username, passwd)

		if errm := a.mailer.GMail.DialAndSend(mailContent); errm != nil {
			prop.Context.Logger().Error(errm)
		}
		prop.Context.Logger().Info("send email success")
	}
	return util.CreateGlobalResponse(constanta.SuccessForgotPassword, nil)
}

func (a auth) RevokeAccess(username string) (resp response.GlobalResponse) {
	//TODO implement me
	panic("implement me")
}

func (a auth) CheckToken(accessToken string) (resp response.GlobalResponse) {
	//TODO implement me
	session := &tbl_session.Session{}
	err := errors.New("")
	if session, err = a.sessionRepo.FindByAccessToken(accessToken); err == nil && session == nil {
		return util.CreateGlobalResponse(constanta.ErrorInvalidAccessToken, nil)
	}

	if lockedCount, err := a.userRepo.CountUserLocked(session.FkUser); err == nil && lockedCount >= a.authConfig.MaximumLoginFailed {
		return util.CreateGlobalResponse(constanta.ErrorAccountLocked, nil)
	}

	if session != nil {
		if (session.AccessToken != "") || (session.ExpiredAt.Valid == true) {
			if time.Now().Before(session.ExpiredAt.Time) {
				return util.CreateGlobalResponse(constanta.Success, auth_resp.LoginResponse{
					AccessToken: session.AccessToken,
					ExpiredAt:   session.ExpiredAt,
				})
			}
		}
	}
	return util.CreateGlobalResponse(constanta.ErrorExpiredAccessToken, nil)
}

func (a auth) composeForgotEmail(mailTo string, recipientName string, username string, content string) *gomail.Message {
	r, err := ioutil.ReadFile("./resources/mail_content/mail.html")

	if err != nil {
		return nil
	}
	cont := string(r)
	cont = strings.Replace(cont, "{user-fullname}", recipientName, -1)
	cont = strings.Replace(cont, "{username}", username, -1)
	cont = strings.Replace(cont, "{password}", content, -1)

	//fmt.Println(r)
	m := gomail.NewMessage()
	m.SetHeader("To", mailTo)
	m.SetHeader("From", *a.mailer.MailFrom)
	m.SetHeader("Subject", "Lupa Password")
	m.SetBody("text/html", cont)

	return m
}

func (a auth) generateLocation(locationId string) (resp []auth_resp.UserLocation) {
	var res []auth_resp.UserLocation
	if dt, err := a.lovRepo.FindLovDetailByDetailId(locationId); err == nil && dt != nil {
		res = append(res, auth_resp.UserLocation{
			LocationId:     dt.LovDetailId,
			LocationName:   dt.Name,
			LocationValue1: dt.ValueStr1,
			LocationValue2: dt.ValueStr2,
			LocationValue3: dt.ValueStr3,
		})

		return res
	}
	return
}

func (a auth) generatePermission(userPk int64) (resp []role_permission_resp.PermittedPermissionRes) {
	dtMenu, _ := a.menuRepo.FindALlPermittedMenu(userPk)
	dtPermission, _ := a.roleRepo.FindPermittedPermission(userPk)

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
				if dtPermission != nil {
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
				}

				//get menu by parent id && child type
				linq.From(*dtMenu).Where(func(i interface{}) bool {
					t := i.(role_permission_resp.JoinedPermittedMenu)

					return t.ParentId == val2.MenuId && t.Type == string(enum.Child)
				}).ForEach(func(x interface{}) {
					t1 := x.(role_permission_resp.JoinedPermittedMenu)
					var childPerm []role_permission_resp.JoinedRoleMenuPermissionMatrix

					//get permission for child menu by pk menu
					if dtPermission != nil {
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
					}

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
		return result
	}
	return
}
