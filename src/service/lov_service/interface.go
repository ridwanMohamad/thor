package lov_service

import (
	"thor/src/payload/lov_req"
	"thor/src/payload/response"
)

type ILovService interface {
	CreateNewLov(req lov_req.CreateLovReq) (resp response.GlobalResponse)
	UpdateLov(headerId string, req lov_req.CreateLovReq) (resp response.GlobalResponse)
	RemoveLovHeader(headerId string) (resp response.GlobalResponse)
	CreateNewLovDetail(headerId string, req []lov_req.LovDetailReq) (resp response.GlobalResponse)
	GetAllLovDetailByHeaderId(id int64) (resp response.GlobalResponse)
	GetAllLovHeader() (resp response.GlobalResponse)
	GetLovByHeaderId(headerId string) (resp response.GlobalResponse)
	GetLovByHeaderName(headerName string) (resp response.GlobalResponse)
	GetLovDetailByDetailId(detailId string) (resp response.GlobalResponse)
	UpdateLovDetailByDetailId(req []lov_req.LovDetailReq) (resp response.GlobalResponse)
	RemoveLovDetail(detailId string) (resp response.GlobalResponse)
	//GetLovDetailById(id int64) (resp *response.GlobalResponse)
	CreateNewPermission(req lov_req.CreatePermissionReq) (resp response.GlobalResponse)
	UpdatePermission(req lov_req.UpdatePermissionReq) (resp response.GlobalResponse)
	RemovePermission(code string) (resp response.GlobalResponse)
	GetAllListPermission() (resp response.GlobalResponse)
}
