package tbl_mapping_user_roles

import (
	"time"

	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type MappingRoleUsers struct {
	Id        *int64         `gorm:"primary_key;autoIncrement;type:integer"`
	UserId    int64          `gorm:"type:integer"`
	RoleId    string         `gorm:"column:role_id;type:varchar(100)"`
	IsDefault bool           `gorm:"column:is_default;type:boolean";default:"false"`
	CreatedAt time.Time      `gorm:"type:timestamp"`
	UpdatedAt null.Time      `gorm:"type:timestamp" json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamp" json:"-"`
}

func (t *MappingRoleUsers) TableName() string {
	return "mapping_user_roles"
}
