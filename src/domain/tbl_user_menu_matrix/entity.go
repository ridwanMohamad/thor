package tbl_user_menu_matrix

type UserMenuMatrix struct {
	PkUser     int64 `gorm:"type:bigint;column:pk_user"`
	PkMenu     int64 `gorm:"type:bigint;column:pk_menu"`
	PkMenuPerm int64 `gorm:"type:bigint;column:pk_menu_perm"`
}

func (T *UserMenuMatrix) TableName() string {
	return "user_menu_matrix"
}
