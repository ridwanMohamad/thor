package user_req

import "thor/src/constanta/enum"

type UserDTO struct {
	Username    string `json:"username" validate:"required"`
	Email       string `json:"email" validate:"email"`
	FullName    string `json:"fullName" validate:"required"`
	MobilePhone string `json:"mobilePhone" validate:"required"`
	// Password            string                     `json:"password" validate:"required"`
	EffectiveAt         string                     `json:"effectiveAt" validate:"required,ISO8601date"`
	ExpiredAt           string                     `json:"expiredAt" validate:"required,ISO8601date"`
	Location            string                     `json:"location" validate:"required"`
	Department          string                     `json:"department" validate:"required"`
	Role                string                     `json:"role" validate:"required"`
	EmployeeId          string                     `json:"employeeId"`
	AdditionalPrivilege []AdditionalPrivilegeShard `json:"additionalPrivileges"`
}

type AdditionalPrivilegeShard struct {
	PkMenu     int64                 `json:"-"`
	MenuId     string                `json:"menuId"`
	State      enum.AddRemoveEnum    `json:"state"`
	Permission []MenuPermissionShard `json:"permission"`
}

type MenuPermissionShard struct {
	PermissionId int64 `json:"permId"`
	State        enum.AddRemoveEnum
}

type UserUpdateDTO struct {
	Username            string                     `json:"username"  validate:"required"`
	Email               string                     `json:"email" validate:"email"`
	FullName            string                     `json:"fullName" validate:"required"`
	MobilePhone         string                     `json:"mobilePhone" validate:"required"`
	RoleId              string                     `json:"roleId"`
	Password            string                     `json:"password"`
	EffectiveAt         string                     `json:"effectiveAt" validate:"required,ISO8601date"`
	ExpiredAt           string                     `json:"expiredAt" validate:"required,ISO8601date"`
	Location            string                     `json:"location" validate:"required"`
	Department          string                     `json:"department" validate:"required"`
	Status              string                     `json:"status" validate:"required"`
	Unlock              bool                       `json:"unlockAccount"`
	EmployeeId          string                     `json:"employeeId"`
	AdditionalPrivilege []AdditionalPrivilegeShard `json:"additionalPrivileges"`
	//AddPrivilege        []AdditionalPrivilegeShard `json:"addPrivileges"`
	//RemovePrivilege     []AdditionalPrivilegeShard `json:"removePrivileges"`
}

type UserChangePasswordReq struct {
	SessionId            string `json:"sessionId" validate:"required"`
	UserPk               int64  `json:"userPk" validate:"required"`
	OldPassword          string `json:"oldPassword" validate:"required"`
	NewPassword          string `json:"newPassword" validate:"required"`
	ReConfirmNewPassword string `json:"reConfirmNewPassword" validate:"required"`
}

type DeactivateUserReq struct {
	UserId string `json:"userId" validate:"required"`
}

type UpdateUserProfile struct {
	Username    string `json:"username"  validate:"required"`
	Email       string `json:"email" validate:"email,required"`
	MobilePhone string `json:"mobilePhone" validate:"required"`
	EmployeeId  string `json:"employeeId"`
	// RoleId      string `json:"roleId"`
	// Location    string `json:"location"`
	ProfilePict string `json:"profilePict"`
}
type MappingUserLocationReq struct {
	Username string         `json:"username" validate:"required"`
	Email    string         `json:"email"`
	FullName string         `json:"fullName"`
	Location []UserLocation `json:"location"`
}
type UserLocation struct {
	Id           *int64 `json:"id"`
	UserId       int64  `json:userId`
	LocationId   string `json:locationId validate:"required"`
	LocationName string `json:locationName`
	IsDefault    bool   `json:isDefault default:"false"`
}

type MappingUserRoleReq struct {
	Username string     `json:"username" validate:"required"`
	Email    string     `json:"email"`
	FullName string     `json:"fullName"`
	Roles    []UserRole `json:"roles"`
}
type UserRole struct {
	Id        *int64 `json:"id"`
	UserId    int64  `json:userId`
	RoleId    string `json:roleId validate:"required"`
	RoleName  string `json:roleName`
	IsDefault bool   `json:isDefault default:"false"`
}
