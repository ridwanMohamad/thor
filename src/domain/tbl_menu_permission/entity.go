package tbl_menu_permission

import (
	"gopkg.in/guregu/null.v4"
	"time"
)

type MenuPermission struct {
	Pk        int64     `gorm:"primary_key;type:integer AUTO_INCREMENT;column:pk" json:"-"`
	PkMenu    int64     `gorm:"type:integer;column:pk_menu" json:"-"`
	Name      string    `gorm:"type:varchar(50);column:name"`
	PermCode  string    `gorm:"type:varchar(50);column:perm_code"`
	CreatedAt time.Time `gorm:"type:timestamp;column:created_at"`
	UpdatedAt null.Time `gorm:"type:timestamp;column:updated_at"`
	DeletedAt null.Time `gorm:"type:timestamp;column:deleted_at"`
}

func (T *MenuPermission) TableName() string {
	return "menu_permission"
}
