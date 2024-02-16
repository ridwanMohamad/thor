package tbl_user_locked

import (
	"time"
)

type UserLocked struct {
	FkUser       int64     `gorm:"primary_key;type:integer AUTO_INCREMENT"`
	LoginAttempt int64     `gorm:"primary_key;type:integer"`
	CreatedAt    time.Time `gorm:"primary_key;column:created_at;type:timestamp"`
}

func (u UserLocked) TableName() string {
	return "user_locked"
}
