package auth_resp

import (
	"thor/src/payload/role_permission_resp"

	"gopkg.in/guregu/null.v4"
)

type LoginResponse struct {
	AccessToken  string     `json:"accessToken,omitempty"`
	ExpiredAt    null.Time  `json:"expiredAt,omitempty"`
	LastLoginAt  null.Time  `json:"lastLoginAt,omitempty"`
	IsFirstLogin bool       `json:"isFirstLogin,omitempty"`
	UserDetail   UserDetail `json:"userDetail"`
}

type LoginResponseV2 struct {
	SessionToken string    `json:"sessionToken,omitempty"`
	ExpiredAt    null.Time `json:"expiredAt,omitempty"`
	LastLoginAt  null.Time `json:"lastLoginAt,omitempty"`
	UserDetail   string    `json:"userDetail"`
	IsFirstLogin bool      `json:"isFirstLogin"`
}

type UserDetail struct {
	UserId   int64  `json:"userId"`
	UserName string `json:"userName"`
	FullName string `json:"fullName"`
	RolePk   int64  `json:"rolePk"`
	RoleId   string `json:"roleId"`
	RoleName string `json:"roleName"`
}

type UserDetailV2 struct {
	Profile    UserProfile                                   `json:"profile"`
	Role       UserRole                                      `json:"role"`
	Location   []UserLocation                                `json:"location"`
	Permission []role_permission_resp.PermittedPermissionRes `json:"permission"`
}

type UserProfile struct {
	UserId      int64  `json:"userId"`
	UserName    string `json:"userName"`
	FullName    string `json:"fullName"`
	ProfilePict string `json:"profilePict"`
}

type UserRole struct {
	Pk       int64  `json:"pk"`
	RoleId   string `json:"roleId"`
	RoleName string `json:"roleName"`
}

type UserLocation struct {
	LocationId     string `json:"locationId"`
	LocationName   string `json:"locationName"`
	LocationValue1 string `json:"locationValue1,omitempty"`
	LocationValue2 string `json:"locationValue2,omitempty"`
	LocationValue3 string `json:"locationValue3,omitempty"`
}
