package tbl_menus

import (
	"gopkg.in/guregu/null.v4"
	"thor/src/constanta/enum"
	"time"
)

type Menus struct {
	Pk         int64             `gorm:"primary_key;type:integer AUTO_INCREMENT" json:"-"`
	PkModuleId int64             `gorm:"type:integer" json:"-"`
	MenuId     null.String       `gorm:"type:varchar(36)"`
	ParentId   null.String       `gorm:"type:varchar(36)"`
	Name       string            `gorm:"type:varchar(50)"`
	Path       string            `gorm:"type:varchar(50)"`
	Type       enum.MenuTypeEnum `gorm:"type:varchar(6)"`
	MenuCode   string            `gorm:"type:varchar(50)"`
	MenuIcon   string            `gorm:"type:varchar(255)"`
	CreatedAt  time.Time         `gorm:"type:timestamp"`
	UpdatedAt  null.Time         `gorm:"type:timestamp"`
}

func (t *Menus) TableName() string {
	return "menus"
}
