package tbl_session

import (
	"gopkg.in/guregu/null.v4"
	"time"
)

type Session struct {
	Pk          int64     `gorm:"primary_key;type:integer AUTO_INCREMENT"`
	FkUser      int64     `gorm:"type:integer"`
	AccessToken string    `gorm:"type:varchar(1000)"`
	ExpiredAt   null.Time `gorm:"type:timestamp"`
	CreatedAt   time.Time `sql:"DEFAULT:now()" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `sql:"DEFAULT:now()" gorm:"autoCreateTime"`
}

func (t *Session) TableName() string {
	return "sessions"
}
