package tbl_additional_privileges

type AdditionalPrivileges struct {
	AddPrivId    int64 `gorm:"primary_key;autoIncrement;type:integer"`
	PkUserId     int64 `gorm:"type:integer"`
	PkMenuId     int64 `gorm:"type:integer"`
	PkMenuPermId int64 `gorm:"type:integer"`
}

func (t *AdditionalPrivileges) TableName() string {
	return "additional_privileges"
}
