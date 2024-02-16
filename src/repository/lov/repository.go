package lov

import (
	"errors"
	"thor/src/constanta"
	"thor/src/domain/tbl_list_of_permission"
	"thor/src/domain/tbl_lov_details"
	"thor/src/domain/tbl_lov_headers"
	"thor/src/server/database"
	"time"

	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type lov struct {
	db *database.Database
}

func NewLovRepository(DB *database.Database) ILovRepository {

	return &lov{db: DB}
}

func (l lov) SaveLovHeader(data *tbl_lov_headers.LovHeaders) (resp *tbl_lov_headers.LovHeaders, err error) {
	//TODO implement me
	tx := l.db

	if err = tx.Create(data).Error; err != nil {
		return nil, constanta.DbFailedToExecuteQuery
	}

	if err = tx.Where("lov_header_id = ?", data.LovHeaderId).Find(&resp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, constanta.DbFailedToExecuteQuery
	}

	return
}

func (l lov) UpdateLovHeader(data *tbl_lov_headers.LovHeaders) (err error) {
	//TODO implement me
	tx := l.db.Begin()

	if err = tx.Unscoped().Model(&tbl_lov_headers.LovHeaders{}).
		Where("lov_header_id = ? ", data.LovHeaderId).
		Updates(tbl_lov_headers.LovHeaders{
			Name:        data.Name,
			Description: data.Description,
			UpdatedAt:   null.TimeFrom(time.Now()),
			DeletedAt:   data.DeletedAt,
		}).
		Error; err != nil {

		tx.Rollback()
		return constanta.DbFailedToUpdateData
	}
	tx.Commit()
	return
}

func (l lov) DeleteLovHeader(headerId string) (err error) {
	tx := l.db.Begin()

	if err = tx.Where("lov_header_id = ?", headerId).Delete(&tbl_lov_headers.LovHeaders{}).Error; err != nil {
		tx.Rollback()
		return constanta.DbFailedToExecuteQuery
	}
	tx.Commit()
	return
}

func (l lov) FindAllLovHeaders() (resp *[]tbl_lov_headers.LovHeaders, err error) {
	//TODO implement me
	tx := l.db

	if err = tx.Unscoped().Find(&resp).Error; err != nil {
		return nil, gorm.ErrRecordNotFound
	}

	return
}

func (l lov) FindLovHeaderById(headerId string) (resp *tbl_lov_headers.LovHeaders, err error) {
	//TODO implement me
	tx := l.db

	if err = tx.Unscoped().Where("lov_header_id = ?", headerId).First(&resp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (l lov) FindLovHeaderByName(headerName string) (resp *tbl_lov_headers.LovHeaders, err error) {
	//TODO implement me
	tx := l.db

	if err = tx.Unscoped().
		Where("upper(name) = upper(?)", headerName).
		First(&resp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (l lov) SaveLovDetail(data *[]tbl_lov_details.LovDetails) (resp *[]tbl_lov_details.LovDetails, err error) {
	//TODO implement me
	tx := l.db
	var headerId int64

	for _, val := range *data {
		headerId = val.PkLovHeader
		if err = tx.Create(&val).Error; err != nil {
			return nil, constanta.DbFailedToExecuteQuery
		}
	}

	if err = tx.Model(&resp).Where("pk_header = ", headerId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (l lov) UpdateLovDetail(data *[]tbl_lov_details.LovDetails) (err error) {
	//TODO implement me
	tx := l.db.Begin()

	for _, val := range *data {
		if err = tx.Model(&tbl_lov_details.LovDetails{}).
			Where("lov_detail_id = ? ", val.LovDetailId).
			//Select("name, description, value_str_1, value_str_2, value_str_3, value_num_1, value_num_2, value_num_3, value_date_1, value_date_2, value_date_3, updated_at").
			Updates(tbl_lov_details.LovDetails{
				Name:        val.Name,
				Description: val.Description,
				ValueStr1:   val.ValueStr1,
				ValueStr2:   val.ValueStr2,
				ValueStr3:   val.ValueStr3,
				ValueNum1:   val.ValueNum1,
				ValueNum2:   val.ValueNum2,
				ValueNum3:   val.ValueNum3,
				ValueDate1:  val.ValueDate1,
				ValueDate2:  val.ValueDate2,
				ValueDate3:  val.ValueDate3,
				UpdatedAt:   null.TimeFrom(time.Now()),
			}).
			Error; err != nil {

			tx.Rollback()
			return constanta.DbFailedToUpdateData
		}
	}
	tx.Commit()

	return
}

func (l lov) DeleteLovDetail(detailId string) (err error) {
	//TODO implement me
	tx := l.db.Begin()

	if err = tx.Where("lov_detail_id = ?", detailId).Delete(&tbl_lov_details.LovDetails{}).Error; err != nil {
		tx.Rollback()
		return constanta.DbFailedToExecuteQuery
	}
	tx.Commit()
	return
}

func (l lov) FindAllLovDetails() (resp *[]tbl_lov_details.LovDetails, err error) {
	//TODO implement me
	panic("implement me")
}

func (l lov) FindLovDetailByDetailId(detailId string) (resp *tbl_lov_details.LovDetails, err error) {
	//TODO implement me
	tx := l.db

	if err = tx.First(&resp, "lov_detail_id = ?", detailId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (l lov) FindAllLovDetailByHeaderId(headerId int64) (resp *[]tbl_lov_details.LovDetails, err error) {
	//TODO implement me
	tx := l.db

	if err = tx.Where("pk_lov_header = ?", headerId).Order("name asc").Find(&resp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (l lov) FindLovDetailByHeaderId(headerId int64) (resp *tbl_lov_details.LovDetails, err error) {
	//TODO implement me
	panic("implement me")
}

func (l lov) SavePermission(data *tbl_list_of_permission.ListOfPermission) (resp *[]tbl_list_of_permission.ListOfPermission, err error) {
	tx := l.db

	//for _, val := range *data {
	var temp *tbl_list_of_permission.ListOfPermission
	if err = tx.Find(&temp, "code = ? and deleted_at <> null", data.Code).Error; err == nil && temp.Code != "" {
		if err = tx.Unscoped().Model(&tbl_list_of_permission.ListOfPermission{}).
			Where("code = ? ", data.Code).
			Updates(tbl_list_of_permission.ListOfPermission{
				DeletedAt: null.TimeFromPtr(nil),
			}).
			Error; err != nil {
			return nil, constanta.DbFailedToExecuteQuery
		}
	} else {
		if err = tx.Create(&data).Error; err != nil {
			return nil, constanta.DbFailedToExecuteQuery
		}
	}
	//}

	if err = tx.Find(&resp, "code = ? and deleted_at IS NULL", data.Code).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}
func (l lov) UpdatePermission(oldCode string, newCode string, newName string) (err error) {
	tx := l.db.Begin()

	if err = tx.Unscoped().Model(&tbl_list_of_permission.ListOfPermission{}).
		Where("code = ? ", oldCode).
		Updates(tbl_list_of_permission.ListOfPermission{
			Name:      newName,
			Code:      newCode,
			UpdatedAt: null.TimeFrom(time.Now()),
			DeletedAt: null.TimeFromPtr(nil),
		}).
		Error; err != nil {

		tx.Rollback()
		return constanta.DbFailedToUpdateData
	}
	tx.Commit()
	return
}
func (l lov) DeletePermission(code string) (err error) {
	tx := l.db.Begin()

	if err = tx.Unscoped().Model(&tbl_list_of_permission.ListOfPermission{}).
		Where("code = ? ", code).
		Updates(tbl_list_of_permission.ListOfPermission{
			DeletedAt: null.TimeFrom(time.Now()),
		}).
		Error; err != nil {

		tx.Rollback()
		return constanta.DbFailedToUpdateData
	}
	tx.Commit()
	return
}
func (l lov) FindAllPermission() (resp *[]tbl_list_of_permission.ListOfPermission, err error) {
	tx := l.db

	if err = tx.Find(&resp, "deleted_at IS NULL").Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}
