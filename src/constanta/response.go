package constanta

import "errors"

var (
	ErrorUserNotFound         = ErrMsg{Code: "THOR_001", Msg: "User not found"}
	ErrorDuplicateUser        = ErrMsg{Code: "THOR_002", Msg: "Username Already Taken"}
	ErrorDuplicateEmailUser   = ErrMsg{Code: "THOR_003", Msg: "Email Already Taken"}
	ErrorHashingPassword      = ErrMsg{Code: "THOR_004", Msg: "Password Hashing Failed"}
	ErrorInvalidFormatDate    = ErrMsg{Code: "THOR_005", Msg: "Format Date Must Be YYYY-MM-DD"}
	ErrorFailedToSaveUser     = ErrMsg{Code: "THOR_006", Msg: "Failed to Create User"}
	ErrorFailedToUpdateUser   = ErrMsg{Code: "THOR_007", Msg: "Failed to Update User"}
	ErrorFailedAuthenticate   = ErrMsg{Code: "THOR_008", Msg: "Failed to Authenticate User, Check Username / Password"}
	ErrorInvalidAccessToken   = ErrMsg{Code: "THOR_009", Msg: "Invalid Access Token"}
	ErrorExpiredAccessToken   = ErrMsg{Code: "THOR_010", Msg: "Expired Access Token"}
	ErrorChangePasswordFailed = ErrMsg{Code: "THOR_011", Msg: "Failed to change Password"}

	ErrorLovHeaderFailedToRemove  = ErrMsg{Code: "THOR_012", Msg: "Failed to Remove LOV Header"}
	ErrorLovDetailFailedToRemove  = ErrMsg{Code: "THOR_013", Msg: "Failed to Remove LOV Detail"}
	ErrorLovHeaderFailedToCreate  = ErrMsg{Code: "THOR_014", Msg: "Failed to Create LOV Header"}
	ErrorLovDetailFailedToCreate  = ErrMsg{Code: "THOR_015", Msg: "Failed to Create LOV Details"}
	ErrorLovHeaderIdNotFound      = ErrMsg{Code: "THOR_016", Msg: "LOV Header Id Not Found"}
	ErrorLovHeaderNotFound        = ErrMsg{Code: "THOR_017", Msg: "Lov Header Is Empty"}
	ErrorLovDetailRequestNotFound = ErrMsg{Code: "THOR_018", Msg: "Lov Detail Request Is Empty"}
	ErrorLovDetailNotFound        = ErrMsg{Code: "THOR_019", Msg: "Lov Detail Data Not Found"}
	ErrorLovDetailFailedUpdate    = ErrMsg{Code: "THOR_020", Msg: "Failed to Update Lov Details"}
	ErrorLovHeaderFailedUpdate    = ErrMsg{Code: "THOR_021", Msg: "Failed to Update Lov Details"}
	ErrorLovDuplicate             = ErrMsg{Code: "THOR_022", Msg: "Lookup Already Exist"}

	ErrorRoleDataNotFound       = ErrMsg{Code: "THOR_025", Msg: "Role Data Not Found"}
	ErrorFailedToCreateRole     = ErrMsg{Code: "THOR_026", Msg: "Failed to Create Role"}
	ErrorFailedToUpdateRole     = ErrMsg{Code: "THOR_027", Msg: "Failed to Update Role"}
	ErrorPermissionDataNotFound = ErrMsg{Code: "THOR_028", Msg: "Permission Data Not Found"}
	ErrorAddPrivFailedToAdd     = ErrMsg{Code: "THOR_029", Msg: "Failed to Mapping Additional Privilege to User"}
	ErrorAddPrivFailedToRemove  = ErrMsg{Code: "THOR_030", Msg: "Failed to Remove Additional Privilege From User"}

	ErrorModuleDataNotFound       = ErrMsg{Code: "THOR_031", Msg: "Module Data Not Found"}
	ErrorMenuDataNotFound         = ErrMsg{Code: "THOR_032", Msg: "Menu Data Not Found"}
	ErrorFailedToCreateModule     = ErrMsg{Code: "THOR_033", Msg: "Failed to Create Module"}
	ErrorFailedToCreateMenu       = ErrMsg{Code: "THOR_034", Msg: "Failed to Create Menu"}
	ErrorFailedToUpdateModule     = ErrMsg{Code: "THOR_035", Msg: "Failed to Update Module"}
	ErrorFailedToUpdateMenu       = ErrMsg{Code: "THOR_036", Msg: "Failed to Update Menu"}
	ErrorMenuPartialDataNotFound  = ErrMsg{Code: "THOR_037", Msg: "Some Menu Data Not Found"}
	ErrorMenuParentIdBelongTo     = ErrMsg{Code: "THOR_038", Msg: "Parent belong to others module"}
	ErrorRoleMenuStateNotValid    = ErrMsg{Code: "THOR_039", Msg: "State not valid for creating new role"}
	ErrorMenuIdDuplicate          = ErrMsg{Code: "THOR_040", Msg: "Menu id duplicated on request"}
	ErrorMenuNotFoundAddPrivilege = ErrMsg{Code: "THOR_041", Msg: "Some menu not found for additional privileges"}
	ErrorPermissionNotFound       = ErrMsg{Code: "THOR_042", Msg: "Permission not found"}
	ErrorFailedToRemovePermission = ErrMsg{Code: "THOR_043", Msg: "Failed to remove permission"}
	ErrorFailedToAddPermission    = ErrMsg{Code: "THOR_044", Msg: "Failed to add permission"}
	ErrorFailedToModuleExists     = ErrMsg{Code: "THOR_045", Msg: "Module name already exists"}
	ErrorFailedToMenuExists       = ErrMsg{Code: "THOR_046", Msg: "Menu name already exists"}
	ErrorMenuPermissionNotExists  = ErrMsg{Code: "THOR_047", Msg: "Menu does not have permission"}
	ErrorAddMenuPermission        = ErrMsg{Code: "THOR_048", Msg: "Failed to add menu permission"}
	ErrorRemovedMenuPermission    = ErrMsg{Code: "THOR_049", Msg: "Failed to remove menu permission"}
	ErrorNotFoundLoPermission     = ErrMsg{Code: "THOR_051", Msg: "List of permission not found"}
	ErrorCreateLoPermission       = ErrMsg{Code: "THOR_051", Msg: "Failed to create new list of permission"}
	ErrorUpdateLoPermission       = ErrMsg{Code: "THOR_052", Msg: "Failed to update list of permission"}
	ErrorRemoveLoPermission       = ErrMsg{Code: "THOR_053", Msg: "Failed to remove list of permission"}

	ErrorFailedExecuteQuery = ErrMsg{Code: "THOR_080", Msg: "Failed to Execute Query"}

	ErrorBindingRequest           = ErrMsg{Code: "THOR_090", Msg: "Invalid Binding"}
	ErrorInvalidRequest           = ErrMsg{Code: "THOR_091", Msg: "Invalid Request"}
	ErrorInvalidDate              = ErrMsg{Code: "THOR_092", Msg: "Invalid Date"}
	ErrorInvalidDateFormat        = ErrMsg{Code: "THOR_093", Msg: "Format Date Is Not ISO-8601"}
	ErrorComparisonBetween2Date   = ErrMsg{Code: "THOR_094", Msg: "Invalid Compare Between 2 Dates"}
	ErrorInvalidDefaultStatus     = ErrMsg{Code: "THOR_095", Msg: "Invalid Default Status (active / inactive)"}
	ErrorNoDataUpdated            = ErrMsg{Code: "THOR_096", Msg: "Similar request with data in database, data will not updated"}
	ErrorNewAndOldSamePassword    = ErrMsg{Code: "THOR_100", Msg: "Old Password and New Password is same, please change new password"}
	ErrorNewAndReconfirmDifferent = ErrMsg{Code: "THOR_101", Msg: "New Password And Reconfirm Is different, please check your password"}
	ErrorInvalidOldPassword       = ErrMsg{Code: "THOR_102", Msg: "Invalid old password"}
	ErrorAccountLocked            = ErrMsg{Code: "THOR_103", Msg: "Your account has been locked"}
	ErrorAccountNotEffective      = ErrMsg{Code: "THOR_103", Msg: "Your account cannot login before effective date"}
	ErrorAccountExpired           = ErrMsg{Code: "THOR_103", Msg: "Your account has expired"}
	ErrorFailedToUpdateSession    = ErrMsg{Code: "THOR_120", Msg: "Failed to update session data"}
	ErrorFailedToRemoveLockedData = ErrMsg{Code: "THOR_121", Msg: "Failed to locked data"}
	ErrorFailedInActiveUser       = ErrMsg{Code: "THOR_122", Msg: "Failed to inactive user"}
	ErrorFileNotFound             = ErrMsg{Code: "THOR_123", Msg: "File not found"}
	ErrorLocationInUsed           = ErrMsg{Code: "THOR_124", Msg: "Location is in used"}

	ErrorGeneral = ErrMsg{Code: "VY3_999", Msg: "General Error"}
)

var (
	Success               = ErrMsg{Code: "000", Msg: "Success"}
	SuccessForgotPassword = ErrMsg{Code: "000", Msg: "Your password will be send to your email"}
)

var (
	DbFailedToExecuteQuery = errors.New("failed to execute query from db")
	DbFailedToInsertData   = errors.New("failed to insert data to db")
	DbFailedToUpdateData   = errors.New("failed to update data to db")
	DbFailedToDeleteData   = errors.New("failed to delete data from db")
)

var (
	DateFormatInvalid = errors.New("format date invalid")
	DateEmpty         = errors.New("date field empty")
	Date1AfterDate2   = errors.New("date1 after date2")
	Date2BeforeDate1  = errors.New("date2 before date1")
)

var (
	SysDatabaseFailedInit   = "failed to init database"
	SysConfigFailedRead     = "failed to read config"
	SysConfigUnmarshall     = "failed to unmarshall config"
	SysRepositoryFailedInit = "failed to init repository"
	SysServiceFailedInit    = "failed to init service"
)

type ErrMsg struct {
	Code string
	Msg  string
}
