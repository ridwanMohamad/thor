package lov

import (
	"thor/src/domain/tbl_list_of_permission"
	"thor/src/domain/tbl_lov_details"
	"thor/src/domain/tbl_lov_headers"
)

type ILovRepository interface {
	SaveLovHeader(data *tbl_lov_headers.LovHeaders) (resp *tbl_lov_headers.LovHeaders, err error)
	UpdateLovHeader(data *tbl_lov_headers.LovHeaders) (err error)
	DeleteLovHeader(headerId string) (err error)
	FindAllLovHeaders() (resp *[]tbl_lov_headers.LovHeaders, err error)
	FindLovHeaderById(headerId string) (resp *tbl_lov_headers.LovHeaders, err error)
	FindLovHeaderByName(headerName string) (resp *tbl_lov_headers.LovHeaders, err error)

	SaveLovDetail(data *[]tbl_lov_details.LovDetails) (resp *[]tbl_lov_details.LovDetails, err error)
	UpdateLovDetail(data *[]tbl_lov_details.LovDetails) (err error)
	DeleteLovDetail(detailId string) (err error)
	FindAllLovDetails() (resp *[]tbl_lov_details.LovDetails, err error)
	FindLovDetailByDetailId(detailId string) (resp *tbl_lov_details.LovDetails, err error)
	FindAllLovDetailByHeaderId(headerId int64) (resp *[]tbl_lov_details.LovDetails, err error)
	FindLovDetailByHeaderId(headerId int64) (resp *tbl_lov_details.LovDetails, err error)

	SavePermission(data *tbl_list_of_permission.ListOfPermission) (resp *[]tbl_list_of_permission.ListOfPermission, err error)
	UpdatePermission(oldCode string, newCode string, newName string) (err error)
	DeletePermission(code string) (err error)
	FindAllPermission() (resp *[]tbl_list_of_permission.ListOfPermission, err error)
}
