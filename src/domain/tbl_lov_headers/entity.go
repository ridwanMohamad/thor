package tbl_lov_headers

import (
	"gopkg.in/guregu/null.v4"
	"time"
)

type LovHeaders struct {
	Pk          int64     `gorm:"primary_key;type:integer AUTO_INCREMENT"`
	LovHeaderId string    `gorm:"type:varchar(36)"`
	Name        string    `gorm:"type:varchar(50)"`
	Description string    `gorm:"type:varchar(100)"`
	CreatedAt   time.Time `gorm:"type:timestamp"`
	UpdatedAt   null.Time `gorm:"type:timestamp"`
	DeletedAt   null.Time `gorm:"type:timestamp" json:"-"`
}

func (t *LovHeaders) TableName() string {
	return "lov_headers"
}
