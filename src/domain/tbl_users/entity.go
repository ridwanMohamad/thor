package tbl_users

import (
	"thor/src/constanta/enum"
	"time"

	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type Users struct {
	Pk            int64           `gorm:"primary_key;autoIncrement;type:integer"`
	Username      string          `gorm:"type:varchar(50);unique"`
	Email         string          `gorm:"type:varchar(100)"`
	MobilePhone   string          `gorm:"type:varchar(20)"`
	Password      string          `gorm:"type:varchar(100)" json:"-"`
	FullName      string          `gorm:"column:full_name;type:varchar(100)"`
	Department    string          `gorm:"column:department;type:varchar(100)"`
	Location      string          `gorm:"column:location;type:varchar(100)"`
	EffectiveAt   null.Time       `gorm:"type:date"`
	ExpiredAt     null.Time       `gorm:"type:date"`
	FkRole        int64           `gorm:"type:integer" json:"-"`
	Status        enum.StatusEnum `gorm:"type:varchar(8)"`
	CreatedAt     time.Time       `gorm:"type:timestamp"`
	UpdatedAt     null.Time       `gorm:"type:timestamp" json:"-"`
	DeletedAt     gorm.DeletedAt  `gorm:"type:timestamp" json:"-"`
	LastLoginAt   null.Time       `gorm:"column:last_login_at;type:timestamp"`
	IsLocked      bool            `gorm:"column:is_locked;type:boolean"`
	LockedAt      null.Time       `gorm:"column:locked_at;type:timestamp"`
	RememberMe    bool            `gorm:"column:remember_me;type:boolean"`
	RememberUntil null.Time       `gorm:"column:remember_until;type:timestamp"`
	LoginIpAddr   string          `gorm:"column:login_ip_addr;type:varchar(20)"`
	IsFirstLogin  bool            `gorm:"column:is_first_login;type:boolean;default:true"`
	EmployeeId    string          `gorm:"column:employee_id;type:varchar(10);unique"`
	ProfilePict   string          `gorm:"column:profile_pict;type:varchar(255);"`
}

func (t *Users) TableName() string {
	return "users"
}
