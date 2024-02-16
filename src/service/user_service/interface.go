package user_service

import (
	"thor/src/payload/response"
	"thor/src/payload/user_req"
	"thor/src/properties"
)

type IUserService interface {
	CreateNewUser(prop *properties.CustomContext, req user_req.UserDTO) (resp response.GlobalResponse)
	GetUserById(id int64) (resp response.GlobalResponse)
	GetUserByUsername(username string) (resp response.GlobalResponse)
	GetAllUser() (resp response.GlobalResponse)
	UpdateUser(req user_req.UserUpdateDTO) (resp response.GlobalResponse)
	ChangePassword(req user_req.UserChangePasswordReq) (resp response.GlobalResponse)
	InActiveUser(id int64) (resp response.GlobalResponse)
	UpdateUserProfile(req user_req.UpdateUserProfile) (resp response.GlobalResponse)
	GetMappingLocationByUsername(username string) (resp response.GlobalResponse)
	UpdateMappingLocationUser(req user_req.MappingUserLocationReq) (resp response.GlobalResponse)
	RemoveMappingLocationUser(id int64) (resp response.GlobalResponse)
	RemoveDefaultLocationUser(id int64) (resp response.GlobalResponse)
	GetMappingRoleByUsername(username string) (resp response.GlobalResponse)
	UpdateMappingRoleUser(req user_req.MappingUserRoleReq) (resp response.GlobalResponse)
	RemoveMappingRoleUser(id int64) (resp response.GlobalResponse)
	RemoveDefaultRoleUser(userId int64) (resp response.GlobalResponse)
}
