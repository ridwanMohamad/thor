package tbl_roles

import (
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
	"thor/src/constanta/enum"
	"time"
)

type Roles struct {
	Pk        int64           `gorm:"primary_key;type:integer;autoIncrement" json:"-"`
	RoleId    string          `gorm:"unique;type:varchar(36)"`
	Name      string          `gorm:"type:integer"`
	Status    enum.StatusEnum `gorm:"type:varchar(8)"`
	CreatedAt time.Time       `gorm:"type:timestamp"`
	UpdatedAt null.Time       `gorm:"type:timestamp"`
	DeletedAt gorm.DeletedAt  `gorm:"type:timestamp"`
}

func (t *Roles) TableName() string {
	return "roles"
}
