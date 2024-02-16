package permission

import (
	"thor/src/domain/tbl_additional_privileges"
)

type IPermissionRepository interface {
	SaveAddPriv(data tbl_additional_privileges.AdditionalPrivileges) (resp *tbl_additional_privileges.AdditionalPrivileges, err error)
	FindAddPrivByUserId(userId int64) (resp *tbl_additional_privileges.AdditionalPrivileges, err error)
	DeleteAddPriv(data tbl_additional_privileges.AdditionalPrivileges) (err error)
}
