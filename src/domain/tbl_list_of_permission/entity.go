package tbl_list_of_permission

import (
	"gopkg.in/guregu/null.v4"
	"time"
)

type ListOfPermission struct {
	Pk        int64     `gorm:"primary_key;autoIncrement;type:integer"`
	Name      string    `gorm:"type:varchar(50)"`
	Code      string    `gorm:"type:varchar(50)"`
	CreatedAt time.Time `gorm:"type:timestamp;column:created_at"`
	UpdatedAt null.Time `gorm:"type:timestamp;column:updated_at"`
	DeletedAt null.Time `gorm:"type:timestamp;column:deleted_at"`
}

func (t *ListOfPermission) TableName() string {
	return "list_of_permission"
}
