package tbl_modules

import (
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
	"thor/src/constanta/enum"
	"time"
)

type Modules struct {
	Pk              int64           `gorm:"primary_key;type:integer AUTO_INCREMENT" json:"-"`
	ModuleId        string          `gorm:"type:varchar(36)"`
	Name            string          `gorm:"type:varchar(50)"`
	Description     string          `gorm:"type:varchar(100)"`
	Status          enum.StatusEnum `gorm:"type:varchar(8)"`
	MaintenanceMode bool            `gorm:"column:maintenance_mode;type:boolean"`
	MtcStart        null.Time       `gorm:"type:date"`
	MtcEnd          null.Time       `gorm:"type:date"`
	ModuleCode      string          `gorm:"type:varchar(50)"`
	CreatedAt       time.Time       `gorm:"type:timestamp"`
	UpdatedAt       null.Time       `gorm:"type:timestamp"`
	DeletedAt       gorm.DeletedAt  `gorm:"type:timestamp" json:"-"`
}

func (t *Modules) TableName() string {
	return "modules"
}
