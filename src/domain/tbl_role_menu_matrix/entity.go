package tbl_role_menu_matrix

type RoleMenuMatrix struct {
	PkRole     int64 `gorm:"type:bigint;column:pk_role"`
	PkMenu     int64 `gorm:"type:bigint;column:pk_menu"`
	PkMenuPerm int64 `gorm:"type:bigint;column:pk_menu_perm"`
}

func (T *RoleMenuMatrix) TableName() string {
	return "role_menu_matrix"
}
