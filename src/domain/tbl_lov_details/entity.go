package tbl_lov_details

import (
	"gopkg.in/guregu/null.v4"
	"time"
)

type LovDetails struct {
	Pk          int64     `gorm:"primary_key;type:integer AUTO_INCREMENT" json:"-"`
	PkLovHeader int64     `gorm:"type:integer" json:"-"`
	LovDetailId string    `gorm:"type:varchar(36)"`
	Name        string    `gorm:"type:varchar(50)"`
	Description string    `gorm:"type:varchar(100)"`
	ValueStr1   string    `gorm:"type:varchar(100);column:value_str_1" json:"valueString1"`
	ValueStr2   string    `gorm:"type:varchar(100);column:value_str_2" json:"valueString2"`
	ValueStr3   string    `gorm:"type:varchar(100);column:value_str_3" json:"valueString3"`
	ValueNum1   int64     `gorm:"type:integer;column:value_num_1" json:"valueNumber1"`
	ValueNum2   int64     `gorm:"type:integer;column:value_num_2" json:"valueNumber2"`
	ValueNum3   int64     `gorm:"type:integer;column:value_num_3" json:"valueNumber3"`
	ValueDate1  null.Time `gorm:"type:date;column:value_date_1" json:"valueDate1"`
	ValueDate2  null.Time `gorm:"type:date;column:value_date_2" json:"valueDate2"`
	ValueDate3  null.Time `gorm:"type:date;column:value_date_3" json:"valueDate3"`
	CreatedAt   time.Time `gorm:"type:timestamp"`
	UpdatedAt   null.Time `gorm:"type:timestamp" json:"-"`
}

func (t *LovDetails) TableName() string {
	return "lov_details"
}
