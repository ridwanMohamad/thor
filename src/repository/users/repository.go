package users

import (
	"errors"
	"strings"
	"thor/src/constanta"
	"thor/src/domain/tbl_mapping_location_users"
	"thor/src/domain/tbl_mapping_user_roles"
	"thor/src/domain/tbl_session"
	"thor/src/domain/tbl_user_locked"
	"thor/src/domain/tbl_users"
	"thor/src/payload/user_resp"
	"thor/src/server/database"
	"time"

	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type users struct {
	db *database.Database
}

func NewUsersRepository(Db *database.Database) IUsersRepository {
	if Db == nil {
		panic(constanta.SysDatabaseFailedInit)
	}
	return &users{db: Db}
}

func (u users) Save(data *tbl_users.Users) (resp *tbl_users.Users, err error) {
	//TODO implement me
	tx := u.db

	if err = tx.Create(&data).Error; err != nil {
		return nil, constanta.DbFailedToInsertData
	}

	if err = tx.Where("lower(username) = lower(?)", data.Username).Find(&resp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, constanta.DbFailedToExecuteQuery
	}
	//tx.Commit()
	return
}

func (u users) FindById(id int64) (resp *tbl_users.Users, err error) {
	//TODO implement me
	tx := u.db

	if err = tx.First(&resp, "pk = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (u users) FindByEmailAndPassword(username string, email string) (resp *tbl_users.Users, err error) {
	tx := u.db

	if err = tx.Unscoped().
		Where("lower(email) = ?", strings.ToLower(email)).
		Or("lower(username) = ?", strings.ToLower(username)).
		First(&resp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, constanta.DbFailedToExecuteQuery
	}

	return
}

func (u users) FindByEmail(email string) (resp *tbl_users.Users, err error) {
	tx := u.db

	if err = tx.Unscoped().Where("lower(email) = ?", strings.ToLower(email)).First(&resp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, constanta.DbFailedToExecuteQuery
	}

	return
}

func (u users) FindByUsername(username string) (resp *tbl_users.Users, err error) {
	//TODO implement me
	tx := u.db

	if err = tx.Unscoped().Where("lower(username) = ?", strings.ToLower(username)).First(&resp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, constanta.DbFailedToExecuteQuery
	}

	return
}

func (u users) FindAll() (resp []tbl_users.Users, err error) {
	//TODO implement me
	tx := u.db

	if err = tx.Find(&resp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (u users) Update(data *tbl_users.Users) (err error) {
	//TODO implement me
	tx := u.db.Begin()

	if err = tx.Model(&tbl_users.Users{}).
		Where("lower(username) = lower(?)", data.Username).
		Updates(map[string]interface{}{
			"email":          data.Email,
			"mobile_phone":   data.MobilePhone,
			"password":       data.Password,
			"full_name":      data.FullName,
			"location":       data.Location,
			"department":     data.Department,
			"effective_at":   data.EffectiveAt,
			"expired_at":     data.ExpiredAt,
			"fk_role":        data.FkRole,
			"status":         data.Status,
			"is_locked":      data.IsLocked,
			"locked_at":      data.LockedAt,
			"updated_at":     null.TimeFrom(time.Now()),
			"is_first_login": data.IsFirstLogin,
			"employee_id":    data.EmployeeId,
			"profile_pict":   data.ProfilePict,
		}).
		Error; err != nil {

		tx.Rollback()
		return constanta.DbFailedToUpdateData
	}
	tx.Commit()
	return
}

func (u users) UpdateLoginData(data *tbl_users.Users) (err error) {
	//TODO implement me
	tx := u.db.Begin()
	if err = tx.Model(&tbl_users.Users{}).
		Where("Pk = ?", data.Pk).
		Updates(map[string]interface{}{
			"is_locked":      data.IsLocked,
			"locked_at":      data.LockedAt,
			"remember_me":    data.RememberMe,
			"remember_until": data.RememberUntil,
			"login_ip_addr":  data.LoginIpAddr,
		}).
		Error; err != nil {

		tx.Rollback()
		return constanta.DbFailedToUpdateData
	}
	tx.Commit()
	return
}

func (u users) UpdateLastLoginAt(data *tbl_session.Session) (err error) {
	tx := u.db.Begin()
	if err = tx.Model(&tbl_users.Users{}).
		Where("Pk = ?", data.FkUser).
		Updates(tbl_users.Users{
			LastLoginAt: null.TimeFrom(time.Now()),
		}).
		Error; err != nil {

		tx.Rollback()
		return constanta.DbFailedToUpdateData
	}
	tx.Commit()
	return
}

func (u users) FindUserAndRoleByUsername(username string) (resp *user_resp.UserJoinRoleResp, err error) {
	//TODO implement me
	tx := u.db

	if err = tx.Model(&tbl_users.Users{}).
		Select("users.pk as user_pk, users.username as user_name, users.full_name, users.password, users.is_first_login, "+
			"roles.pk as role_pk, roles.role_id, roles.name as role_name, users.location ,users.effective_at, users.expired_at, users.is_locked, users.profile_pict").
		Joins("inner join roles on roles.pk = users.fk_role").
		Where("users.username = ?", username).Where("users.deleted_at IS NULL").
		Scan(&resp).Error; err != nil {
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (u users) SaveUserLocked(userId int64) (resp *tbl_user_locked.UserLocked, err error) {
	//TODO implement me
	tx := u.db

	if err = tx.Save(&tbl_user_locked.UserLocked{FkUser: userId, LoginAttempt: 1, CreatedAt: time.Now()}).Error; err != nil {
		return nil, constanta.DbFailedToInsertData
	}

	if err = tx.First(&resp, "fk_user = ?", userId).Error; err != nil {
		return nil, constanta.DbFailedToInsertData
	}
	return
}

func (u users) CountUserLocked(userId int64) (resp int64, err error) {
	//TODO implement me
	tx := u.db

	if err = tx.Model(&tbl_user_locked.UserLocked{}).Where("fk_user = ?", userId).Count(&resp).Error; err != nil {
		return 0, constanta.DbFailedToExecuteQuery
	}
	return
}

func (u users) RemoveUserLocked(userId int64) (err error) {
	//TODO implement me
	tx := u.db

	if err = tx.Where("fk_user = ?", userId).Delete(&tbl_user_locked.UserLocked{}).Error; err != nil {
		return constanta.DbFailedToDeleteData
	}
	return
}

func (u users) InActiveUser(userId int64) (err error) {
	tx := u.db.Begin()
	if err = tx.
		Where("Pk = ?", userId).Delete(&tbl_users.Users{}).
		Error; err != nil {

		tx.Rollback()
		return constanta.DbFailedToUpdateData
	}
	tx.Commit()
	return
}

func (u users) FindMappingLocationById(Id int64) (resp *user_resp.UserJoinMappingLocation, err error) {
	//TODO implement me
	tx := u.db

	if err = tx.Model(&tbl_mapping_location_users.MappingLocationUsers{}).
		Select("mapping_user_locations.id, mapping_user_locations.user_id, mapping_user_locations.location_id, mapping_user_locations.is_default").
		Where("mapping_user_locations.id = ?", Id).Where("mapping_user_locations.deleted_at IS NULL").
		Scan(&resp).Error; err != nil {
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}
func (u users) FindMappingLocationDefaultByUserId(userId int64) (resp *user_resp.UserJoinMappingLocation, err error) {
	//TODO implement me
	tx := u.db

	if err = tx.Model(&tbl_mapping_location_users.MappingLocationUsers{}).
		Select("mapping_user_locations.id, mapping_user_locations.user_id, mapping_user_locations.location_id, mapping_user_locations.is_default").
		Where("mapping_user_locations.user_id = ?", userId).Where("is_default is true").Where("mapping_user_locations.deleted_at IS NULL").
		Scan(&resp).Error; err != nil {
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (u users) FindMappingLocationByUserId(username string) (resp *[]user_resp.UserJoinMappingLocation, err error) {
	//TODO implement me
	tx := u.db

	if err = tx.Model(&tbl_mapping_location_users.MappingLocationUsers{}).
		Select("mapping_user_locations.id, mapping_user_locations.user_id, mapping_user_locations.location_id, lov_details.name as location_name,  mapping_user_locations.is_default, lov_details.value_str_1 as location_value1, lov_details.value_str_2 as location_value2, lov_details.value_str_3 as location_value3").
		Joins("inner join users on users.pk = mapping_user_locations.user_id").
		Joins("inner join lov_details on mapping_user_locations.location_id = lov_details.lov_detail_id").
		Where("users.username = ?", username).Where("mapping_user_locations.deleted_at IS NULL").
		Order("is_default desc").Order("id asc").
		Scan(&resp).Error; err != nil {
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}
func (u users) FindMappingLocationByUserIdandLocation(userId int64, location string) (resp *user_resp.UserJoinMappingLocation, err error) {
	//TODO implement me
	tx := u.db

	if err = tx.Model(&tbl_mapping_location_users.MappingLocationUsers{}).
		Select("mapping_user_locations.id, mapping_user_locations.user_id, mapping_user_locations.location_id,  mapping_user_locations.is_default").
		Where("mapping_user_locations.user_id = ?", userId).Where("mapping_user_locations.location_id = ?", location).Where("mapping_user_locations.deleted_at IS NULL").
		Order("is_default desc").Order("id asc").
		Scan(&resp).Error; err != nil {
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (u users) UpdateInsertMappingLocationUsers(data []tbl_mapping_location_users.MappingLocationUsers) (err error) {
	//TODO implement me
	tx := u.db.Begin()
	for _, v := range data {
		if v.Id != nil {
			if err = tx.Model(&tbl_mapping_location_users.MappingLocationUsers{}).
				Where("id", v.Id).
				Updates(map[string]interface{}{
					"user_id":     v.UserId,
					"location_id": v.LocationId,
					"updated_at":  null.TimeFrom(time.Now()),
					"is_default":  v.IsDefault,
				}).
				Error; err != nil {

				tx.Rollback()
				return constanta.DbFailedToUpdateData
			}
		} else {
			if err = tx.Create(&v).Error; err != nil {
				return constanta.DbFailedToInsertData
			}
		}
	}

	tx.Commit()
	return
}

func (u users) RemoveMappingLocationUser(Id int64) (err error) {
	//TODO implement me
	tx := u.db.Begin()
	if err = tx.Model(&tbl_mapping_location_users.MappingLocationUsers{}).
		Where("id", Id).
		Updates(map[string]interface{}{
			"deleted_at": null.TimeFrom(time.Now()),
		}).
		Error; err != nil {

		tx.Rollback()
		return constanta.DbFailedToUpdateData
	}

	tx.Commit()
	return
}
func (u users) RemoveDefaultLocationUser(userId int64) (err error) {
	//TODO implement me
	tx := u.db.Begin()
	if err = tx.Model(&tbl_mapping_location_users.MappingLocationUsers{}).
		Where("user_id", userId).
		Where("is_default is true").
		Updates(map[string]interface{}{
			"deleted_at": null.TimeFrom(time.Now()),
		}).
		Error; err != nil {

		tx.Rollback()
		return constanta.DbFailedToUpdateData
	}

	tx.Commit()
	return
}

func (u users) FindMappingRoleById(Id int64) (resp *user_resp.UserJoinMappingRole, err error) {
	//TODO implement me
	tx := u.db

	if err = tx.Model(&tbl_mapping_user_roles.MappingRoleUsers{}).
		Select("mapping_user_roles.id, mapping_user_roles.user_id, mapping_user_roles.role_id, mapping_user_roles.is_default").
		Where("mapping_user_roles.id = ?", Id).Where("mapping_user_roles.deleted_at IS NULL").
		Scan(&resp).Error; err != nil {
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (u users) FindMappingRoleByUserId(username string) (resp *[]user_resp.UserJoinMappingRole, err error) {
	//TODO implement me
	tx := u.db

	if err = tx.Model(&tbl_mapping_user_roles.MappingRoleUsers{}).
		Select("mapping_user_roles.id, mapping_user_roles.user_id, mapping_user_roles.role_id, roles.name as role_name,  mapping_user_roles.is_default").
		Joins("inner join users on users.pk = mapping_user_roles.user_id").
		Joins("inner join roles on mapping_user_roles.role_id = roles.role_id").
		Where("users.username = ?", username).Where("mapping_user_roles.deleted_at IS NULL").
		Order("is_default desc").Order("id asc").
		Scan(&resp).Error; err != nil {
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}
func (u users) FindMappingRoleByUserIdandRole(userId int64, roleId string) (resp *user_resp.UserJoinMappingRole, err error) {
	//TODO implement me
	tx := u.db

	if err = tx.Model(&tbl_mapping_user_roles.MappingRoleUsers{}).
		Select("mapping_user_roles.id, mapping_user_roles.user_id, mapping_user_roles.role_id,  mapping_user_roles.is_default").
		Where("mapping_user_roles.user_id = ?", userId).
		Where("mapping_user_roles.role_id = ?", roleId).
		Where("mapping_user_roles.deleted_at IS NULL").
		Order("is_default desc").Order("id asc").
		Scan(&resp).Error; err != nil {
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (u users) FindMappingRoleDefaultByUserId(userId int64) (resp *user_resp.UserJoinMappingRole, err error) {
	//TODO implement me
	tx := u.db

	if err = tx.Model(&tbl_mapping_user_roles.MappingRoleUsers{}).
		Select("mapping_user_roles.id, mapping_user_roles.user_id, mapping_user_roles.role_id, mapping_user_roles.is_default").
		Where("mapping_user_roles.user_id = ?", userId).Where("mapping_user_roles.is_default is true").
		Where("mapping_user_roles.deleted_at IS NULL").
		Scan(&resp).Error; err != nil {
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (u users) UpdateInsertMappingRoleUsers(data []tbl_mapping_user_roles.MappingRoleUsers) (err error) {
	//TODO implement me
	tx := u.db.Begin()
	for _, v := range data {
		if v.Id != nil {
			if err = tx.Model(&tbl_mapping_user_roles.MappingRoleUsers{}).
				Where("id", v.Id).
				Updates(map[string]interface{}{
					"user_id":    v.UserId,
					"role_id":    v.RoleId,
					"updated_at": null.TimeFrom(time.Now()),
					"is_default": v.IsDefault,
				}).
				Error; err != nil {

				tx.Rollback()
				return constanta.DbFailedToUpdateData
			}
		} else {
			if err = tx.Create(&v).Error; err != nil {
				return constanta.DbFailedToInsertData
			}
		}
	}

	tx.Commit()
	return
}

func (u users) RemoveMappingRoleUser(Id int64) (err error) {
	//TODO implement me
	tx := u.db.Begin()
	if err = tx.Model(&tbl_mapping_user_roles.MappingRoleUsers{}).
		Where("id", Id).
		Updates(map[string]interface{}{
			"deleted_at": null.TimeFrom(time.Now()),
		}).
		Error; err != nil {

		tx.Rollback()
		return constanta.DbFailedToUpdateData
	}

	tx.Commit()
	return
}
func (u users) RemoveDefaultRoleUser(userId int64) (err error) {
	//TODO implement me
	tx := u.db.Begin()
	if err = tx.Model(&tbl_mapping_user_roles.MappingRoleUsers{}).
		Where("user_id", userId).
		Where("is_default is true").
		Where("deleted_at is null").
		Updates(map[string]interface{}{
			"deleted_at": null.TimeFrom(time.Now()),
		}).
		Error; err != nil {

		tx.Rollback()
		return constanta.DbFailedToUpdateData
	}

	tx.Commit()
	return
}
