package users

import (
	"thor/src/domain/tbl_mapping_location_users"
	"thor/src/domain/tbl_mapping_user_roles"
	"thor/src/domain/tbl_session"
	"thor/src/domain/tbl_user_locked"
	"thor/src/domain/tbl_users"
	"thor/src/payload/user_resp"
)

type IUsersRepository interface {
	Save(data *tbl_users.Users) (resp *tbl_users.Users, err error)
	FindById(id int64) (resp *tbl_users.Users, err error)
	FindByEmailAndPassword(username string, email string) (resp *tbl_users.Users, err error)
	FindByEmail(email string) (resp *tbl_users.Users, err error)
	FindByUsername(username string) (resp *tbl_users.Users, err error)
	FindAll() (resp []tbl_users.Users, err error)
	Update(data *tbl_users.Users) (err error)
	UpdateLoginData(data *tbl_users.Users) (err error)
	UpdateLastLoginAt(data *tbl_session.Session) (err error)
	FindUserAndRoleByUsername(username string) (resp *user_resp.UserJoinRoleResp, err error)

	SaveUserLocked(userId int64) (resp *tbl_user_locked.UserLocked, err error)
	CountUserLocked(userId int64) (resp int64, err error)
	RemoveUserLocked(userId int64) (err error)

	InActiveUser(userId int64) (err error)

	FindMappingLocationById(Id int64) (resp *user_resp.UserJoinMappingLocation, err error)
	FindMappingLocationDefaultByUserId(userId int64) (resp *user_resp.UserJoinMappingLocation, err error)
	FindMappingLocationByUserId(username string) (resp *[]user_resp.UserJoinMappingLocation, err error)
	FindMappingLocationByUserIdandLocation(userId int64, location string) (resp *user_resp.UserJoinMappingLocation, err error)
	UpdateInsertMappingLocationUsers(data []tbl_mapping_location_users.MappingLocationUsers) (err error)
	RemoveMappingLocationUser(Id int64) (err error)
	RemoveDefaultLocationUser(userId int64) (err error)

	FindMappingRoleById(Id int64) (resp *user_resp.UserJoinMappingRole, err error)
	FindMappingRoleDefaultByUserId(userId int64) (resp *user_resp.UserJoinMappingRole, err error)
	FindMappingRoleByUserId(username string) (resp *[]user_resp.UserJoinMappingRole, err error)
	FindMappingRoleByUserIdandRole(userId int64, roleId string) (resp *user_resp.UserJoinMappingRole, err error)
	UpdateInsertMappingRoleUsers(data []tbl_mapping_user_roles.MappingRoleUsers) (err error)
	RemoveMappingRoleUser(Id int64) (err error)
	RemoveDefaultRoleUser(userId int64) (err error)
}
