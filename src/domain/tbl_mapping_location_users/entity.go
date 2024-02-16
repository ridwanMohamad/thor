package tbl_mapping_location_users

import (
	"time"

	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type MappingLocationUsers struct {
	Id         *int64         `gorm:"primary_key;autoIncrement;type:integer"`
	UserId     int64          `gorm:"type:integer"`
	LocationId string         `gorm:"column:location_id;type:varchar(100)"`
	IsDefault  bool           `gorm:"column:is_default;type:boolean";default:"false"`
	CreatedAt  time.Time      `gorm:"type:timestamp"`
	UpdatedAt  null.Time      `gorm:"type:timestamp" json:"-"`
	DeletedAt  gorm.DeletedAt `gorm:"type:timestamp" json:"-"`
}

func (t *MappingLocationUsers) TableName() string {
	return "mapping_user_locations"
}
