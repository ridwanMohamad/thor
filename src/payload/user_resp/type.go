package user_resp

import (
	"thor/src/constanta/enum"
	"thor/src/payload/role_permission_resp"
	"time"

	"gopkg.in/guregu/null.v4"
)

type UserRegisterResp struct {
	UserId      int64           `json:"userId"`
	EmployeeId  string          `json:"employeeId"`
	Username    string          `json:"username"`
	Email       string          `json:"email"`
	MobilePhone string          `json:"mobilePhone"`
	FullName    string          `json:"fullName"`
	Status      enum.StatusEnum `json:"status"`
	EffectiveAt time.Time       `json:"effectiveAt"`
	ExpiredAt   time.Time       `json:"expiredAt"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
	DeletedAt   time.Time       `json:"deletedAt"`
}

type UserDetailResp struct {
	UserId              int64            `json:"userId"`
	IsFirstLogin        bool             `json:"is_first_login"`
	Username            string           `json:"username"`
	Email               string           `json:"email"`
	MobilePhone         string           `json:"mobilePhone"`
	FullName            string           `json:"fullName"`
	Department          string           `json:"department"`
	DepartmentName      string           `json:"departmentName"`
	Location            string           `json:"location"`
	LocationName        string           `json:"locationName"`
	Status              enum.StatusEnum  `json:"status"`
	EffectiveAt         time.Time        `json:"effectiveAt"`
	ExpiredAt           time.Time        `json:"expiredAt"`
	CreatedAt           time.Time        `json:"createdAt"`
	UpdatedAt           null.Time        `json:"updatedAt"`
	LastLoginAt         null.Time        `json:"lastLoginAt"`
	IsLocked            bool             `json:"isLocked"`
	LockedAt            null.Time        `json:"lockedAt"`
	Role                CommonPkIdName   `json:"role"`
	EmployeeId          string           `json:"employeeId"`
	ProfilePict         string           `json:"profilePict"`
	AdditionalPrivilege []CommonPkIdName `json:"additionalPrivilege"`
}

type CommonPkIdName struct {
	Pk   int64  `json:"pk"`
	Id   string `json:"id"`
	Name string `json:"name"`
}

type AdditionalMenu struct {
	ParentId         string                                                `json:"parent_id"`
	MenuId           string                                                `json:"menu_id"`
	MenuType         string                                                `json:"menu_type"`
	Name             string                                                `json:"name"`
	Icon             string                                                `json:"icon"`
	RegisteredStatus string                                                `json:"registered_status"`
	Permission       []role_permission_resp.JoinedUserMenuPermissionMatrix `json:"permission"`
}

type UserJoinRoleResp struct {
	UserPk       int64
	UserName     string
	FullName     string
	Password     string
	RolePk       int64
	RoleId       string
	RoleName     string
	Location     string
	ProfilePict  string
	EffectiveAt  time.Time
	ExpiredAt    time.Time
	IsLocked     bool
	IsFirstLogin bool
}

type MappingUserLocationResp struct {
	Username string         `json:"username"`
	Email    string         `json:"email"`
	FullName string         `json:"fullName"`
	Location []UserLocation `json:"location"`
}
type UserLocation struct {
	Id             int64  `json:"id"`
	LocationId     string `json:"locationId"`
	LocationName   string `json:"locationName"`
	LocationValue1 string `json:"locationValue1"`
	LocationValue2 string `json:"locationValue2"`
	LocationValue3 string `json:"locationValue3"`
	IsDefault      bool   `json:"isDefault"`
}
type UserJoinMappingLocation struct {
	Id             int64
	UserId         int64
	IsDefault      bool
	LocationId     string
	LocationName   string
	LocationValue1 string
	LocationValue2 string
	LocationValue3 string
}

type MappingUserRoleResp struct {
	Username string     `json:"username"`
	Email    string     `json:"email"`
	FullName string     `json:"fullName"`
	Roles    []UserRole `json:"roles"`
}
type UserRole struct {
	Id        int64  `json:"id"`
	RoleId    string `json:roleId`
	RoleName  string `json:roleName`
	IsDefault bool   `json:isDefault`
}
type UserJoinMappingRole struct {
	Id        int64
	UserId    int64
	IsDefault bool
	RoleId    string
	RoleName  string
}
