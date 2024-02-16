package tbl_role_menu

type RoleMenu struct {
	RoleMenuId int64 `gorm:"primary_key;type:integer;autoIncrement"`
	PkRole     int64 `gorm:"unique;type:integer"`
	PkMenu     int64 `gorm:"type:integer"`
}

func (t *RoleMenu) TableName() string {
	return "role_menu"
}
