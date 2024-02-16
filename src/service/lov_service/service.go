package lov_service

import (
	"errors"
	"strings"
	"thor/src/constanta"
	"thor/src/domain/tbl_list_of_permission"
	"thor/src/domain/tbl_lov_details"
	"thor/src/domain/tbl_lov_headers"
	"thor/src/payload/lov_req"
	"thor/src/payload/lov_resp"
	"thor/src/payload/response"
	"thor/src/repository/lov"
	"thor/src/util"
	"thor/src/util/date_util"
	"thor/src/util/string_util"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

type lovService struct {
	lovRepo lov.ILovRepository
}

func NewLovService(repo lov.ILovRepository) ILovService {

	return lovService{lovRepo: repo}
}

func (l lovService) CreateNewLov(req lov_req.CreateLovReq) (resp response.GlobalResponse) {
	//TODO implement me
	var headerDt = &tbl_lov_headers.LovHeaders{}
	headerId, err := uuid.NewRandom()

	if err != nil {
		return util.CreateGlobalResponse(constanta.ErrorGeneral, nil)
	}

	if headerDt, err = l.lovRepo.FindLovHeaderByName(req.Name); err != nil || headerDt != nil {
		if headerDt != nil {
			if headerDt.Name == req.Name {
				return util.CreateGlobalResponse(constanta.ErrorLovDuplicate, "")
			}
		}
		if errors.Is(err, constanta.DbFailedToExecuteQuery) {
			return util.CreateGlobalResponse(constanta.ErrorFailedExecuteQuery, "")
		}
	}

	if req.LovDetails == nil || len(req.LovDetails) < 0 {
		return util.CreateGlobalResponse(constanta.ErrorLovDetailRequestNotFound, nil)
	}

	if headerDt, err = l.lovRepo.SaveLovHeader(&tbl_lov_headers.LovHeaders{
		LovHeaderId: headerId.String(),
		Name:        req.Name,
		Description: req.Description,
	}); err != nil {
		return util.CreateGlobalResponse(constanta.ErrorLovHeaderFailedToCreate, nil)
	}

	var detailDt []tbl_lov_details.LovDetails
	var lovDetail []lov_resp.ShardLovDetail

	for _, val := range req.LovDetails {
		detailId, err := uuid.NewRandom()

		if err != nil {
			return util.CreateGlobalResponse(constanta.ErrorGeneral, nil)
		}

		detailDt = append(detailDt, tbl_lov_details.LovDetails{
			PkLovHeader: headerDt.Pk,
			LovDetailId: detailId.String(),
			Name:        val.Name,
			Description: val.Description,
			ValueStr1:   val.ValueStr1,
			ValueStr2:   val.ValueStr2,
			ValueStr3:   val.ValueStr3,
			ValueNum1:   val.ValueNum1,
			ValueNum2:   val.ValueNum2,
			ValueNum3:   val.ValueNum3,
			ValueDate1:  date_util.StringToNilTime(val.ValueDate1),
			ValueDate2:  date_util.StringToNilTime(val.ValueDate2),
			ValueDate3:  date_util.StringToNilTime(val.ValueDate3),
		})

		lovDetail = append(lovDetail, lov_resp.ShardLovDetail{
			LovDetailId: detailId.String(),
			Name:        val.Name,
			Description: val.Description,
			ValueStr1:   val.ValueStr1,
			ValueStr2:   val.ValueStr2,
			ValueStr3:   val.ValueStr3,
			ValueNum1:   val.ValueNum1,
			ValueNum2:   val.ValueNum2,
			ValueNum3:   val.ValueNum3,
			ValueDate1:  date_util.StringToNilTime(val.ValueDate1),
			ValueDate2:  date_util.StringToNilTime(val.ValueDate2),
			ValueDate3:  date_util.StringToNilTime(val.ValueDate3),
		})
	}

	if _, err = l.lovRepo.SaveLovDetail(&detailDt); err != nil {
		return util.CreateGlobalResponse(constanta.ErrorLovDetailFailedToCreate, nil)
	}

	return util.CreateGlobalResponse(constanta.Success, lov_resp.LovAfterCreateResp{
		LovHeaderId: headerDt.LovHeaderId,
		Name:        headerDt.Name,
		Description: headerDt.Description,
		Detail:      lovDetail,
	})
}

func (l lovService) UpdateLov(headerId string, req lov_req.CreateLovReq) (resp response.GlobalResponse) {
	if dt, err := l.lovRepo.FindLovHeaderById(headerId); dt != nil && err == nil {

		if req.Enabled {
			dt.DeletedAt = null.TimeFromPtr(nil)
		}

		if err = l.lovRepo.UpdateLovHeader(dt); err != nil {
			return util.CreateGlobalResponse(constanta.ErrorLovHeaderFailedUpdate, nil)
		}

		var tmpI = []tbl_lov_details.LovDetails{}
		var tmpU = []tbl_lov_details.LovDetails{}

		for _, val := range req.LovDetails {
			if val.LovDetailId == "" {
				detailId, _ := uuid.NewRandom()

				tmpI = append(tmpI, tbl_lov_details.LovDetails{
					PkLovHeader: dt.Pk,
					LovDetailId: detailId.String(),
					Name:        val.Name,
					Description: val.Description,
					ValueStr1:   val.ValueStr1,
					ValueStr2:   val.ValueStr2,
					ValueStr3:   val.ValueStr3,
					ValueNum1:   val.ValueNum1,
					ValueNum2:   val.ValueNum2,
					ValueNum3:   val.ValueNum3,
					ValueDate1:  date_util.StringToNilTime(val.ValueDate1),
					ValueDate2:  date_util.StringToNilTime(val.ValueDate2),
					ValueDate3:  date_util.StringToNilTime(val.ValueDate3),
				})
			} else {
				if lod, err := l.lovRepo.FindLovDetailByDetailId(val.LovDetailId); lod != nil && err == nil {
					tmpU = append(tmpU, tbl_lov_details.LovDetails{
						PkLovHeader: lod.PkLovHeader,
						LovDetailId: lod.LovDetailId,
						Name:        val.Name,
						Description: val.Description,
						ValueStr1:   val.ValueStr1,
						ValueStr2:   val.ValueStr2,
						ValueStr3:   val.ValueStr3,
						ValueNum1:   val.ValueNum1,
						ValueNum2:   val.ValueNum2,
						ValueNum3:   val.ValueNum3,
						ValueDate1:  date_util.StringToNilTime(val.ValueDate1),
						ValueDate2:  date_util.StringToNilTime(val.ValueDate2),
						ValueDate3:  date_util.StringToNilTime(val.ValueDate3),
					})
				}
			}
		}
		if len(tmpI) > 0 {
			if dta, err := l.lovRepo.SaveLovDetail(&tmpI); dta != nil && err == nil {

			}
		}
		if len(tmpU) > 0 {
			if err = l.lovRepo.UpdateLovDetail(&tmpU); err == nil {

			}
		}
		return util.CreateGlobalResponse(constanta.Success, nil)
	}
	return util.CreateGlobalResponse(constanta.ErrorLovHeaderNotFound, nil)
}

func (l lovService) RemoveLovHeader(headerId string) (resp response.GlobalResponse) {
	if dt, err := l.lovRepo.FindLovHeaderById(headerId); dt != nil && err == nil {
		if err = l.lovRepo.DeleteLovHeader(dt.LovHeaderId); err == nil {
			return util.CreateGlobalResponse(constanta.Success, nil)
		}
		return util.CreateGlobalResponse(constanta.ErrorLovHeaderFailedToRemove, nil)
	}
	return util.CreateGlobalResponse(constanta.ErrorLovHeaderNotFound, nil)
}

func (l lovService) CreateNewLovDetail(headerId string, req []lov_req.LovDetailReq) (resp response.GlobalResponse) {
	//TODO implement me
	var lovHeader *tbl_lov_headers.LovHeaders
	err := errors.New("")

	if lovHeader, err = l.lovRepo.FindLovHeaderById(headerId); err != nil {
		return util.CreateGlobalResponse(constanta.ErrorLovHeaderIdNotFound, nil)
	}

	var detailId uuid.UUID

	var detailDt []tbl_lov_details.LovDetails
	var lovDetail []lov_resp.ShardLovDetail

	for _, val := range req {
		if detailId, err = uuid.NewRandom(); err != nil {
			return util.CreateGlobalResponse(constanta.ErrorGeneral, nil)
		}

		detailDt = append(detailDt, tbl_lov_details.LovDetails{
			PkLovHeader: lovHeader.Pk,
			LovDetailId: detailId.String(),
			Name:        val.Name,
			Description: val.Description,
			ValueStr1:   val.ValueStr1,
			ValueStr2:   val.ValueStr2,
			ValueStr3:   val.ValueStr3,
			ValueNum1:   val.ValueNum1,
			ValueNum2:   val.ValueNum2,
			ValueNum3:   val.ValueNum3,
			ValueDate1:  date_util.StringToNilTime(val.ValueDate1),
			ValueDate2:  date_util.StringToNilTime(val.ValueDate2),
			ValueDate3:  date_util.StringToNilTime(val.ValueDate3),
		})

		lovDetail = append(lovDetail, lov_resp.ShardLovDetail{
			LovDetailId: detailId.String(),
			Name:        val.Name,
			Description: val.Description,
			ValueStr1:   val.ValueStr1,
			ValueStr2:   val.ValueStr2,
			ValueStr3:   val.ValueStr3,
			ValueNum1:   val.ValueNum1,
			ValueNum2:   val.ValueNum2,
			ValueNum3:   val.ValueNum3,
			ValueDate1:  date_util.StringToNilTime(val.ValueDate1),
			ValueDate2:  date_util.StringToNilTime(val.ValueDate2),
			ValueDate3:  date_util.StringToNilTime(val.ValueDate3),
		})
	}

	if _, err = l.lovRepo.SaveLovDetail(&detailDt); err != nil {
		return util.CreateGlobalResponse(constanta.ErrorLovDetailFailedUpdate, nil)
	}
	return util.CreateGlobalResponse(constanta.Success, lovDetail)
}

func (l lovService) GetAllLovDetailByHeaderId(id int64) (resp response.GlobalResponse) {
	//TODO implement me
	panic("implement me")
}

func (l lovService) GetAllLovHeader() (resp response.GlobalResponse) {
	//TODO implement me
	if val, err := l.lovRepo.FindAllLovHeaders(); err == nil {
		if val != nil {
			var dt []lov_resp.ShardLovHeaders

			for _, d := range *val {
				var enabled = !d.DeletedAt.Valid

				dt = append(dt, lov_resp.ShardLovHeaders{
					HeaderId:    d.LovHeaderId,
					Name:        d.Name,
					Description: d.Description,
					Enabled:     enabled,
				})
			}
			return util.CreateGlobalResponse(constanta.Success, dt)
		}
	}
	return util.CreateGlobalResponse(constanta.ErrorLovHeaderNotFound, nil)
}

func (l lovService) GetLovByHeaderId(headerId string) (resp response.GlobalResponse) {
	//TODO implement me
	var lovHeader *tbl_lov_headers.LovHeaders
	err := errors.New("")

	if lovHeader, err = l.lovRepo.FindLovHeaderByName(headerId); err != nil {
		return util.CreateGlobalResponse(constanta.ErrorLovHeaderIdNotFound, nil)
	}

	if val, err := l.lovRepo.FindAllLovDetailByHeaderId(lovHeader.Pk); err == nil {
		return util.CreateGlobalResponse(constanta.Success, val)
	}
	return util.CreateGlobalResponse(constanta.ErrorLovHeaderNotFound, nil)
}

func (l lovService) GetLovByHeaderName(headerName string) (resp response.GlobalResponse) {
	var lovHeader *tbl_lov_headers.LovHeaders
	err := errors.New("")

	if lovHeader, err = l.lovRepo.FindLovHeaderByName(headerName); err != nil {
		return util.CreateGlobalResponse(constanta.ErrorLovHeaderIdNotFound, nil)
	}

	if val, err := l.lovRepo.FindAllLovDetailByHeaderId(lovHeader.Pk); err == nil {
		return util.CreateGlobalResponse(constanta.Success, val)
	}
	return util.CreateGlobalResponse(constanta.ErrorLovHeaderNotFound, nil)
}

func (l lovService) GetLovDetailByDetailId(detailId string) (resp response.GlobalResponse) {
	if dt, err := l.lovRepo.FindLovDetailByDetailId(detailId); dt != nil && err == nil {
		return util.CreateGlobalResponse(constanta.Success, dt)
	}
	return util.CreateGlobalResponse(constanta.ErrorLovDetailNotFound, nil)
}

func (l lovService) UpdateLovDetailByDetailId(req []lov_req.LovDetailReq) (resp response.GlobalResponse) {
	var dt []tbl_lov_details.LovDetails
	for _, val := range req {
		if tmp, err := l.lovRepo.FindLovDetailByDetailId(val.LovDetailId); tmp == nil && err == nil {
			return util.CreateGlobalResponse(constanta.ErrorLovDetailNotFound, nil)
		}

		dt = append(dt, tbl_lov_details.LovDetails{
			LovDetailId: val.LovDetailId,
			Name:        val.Name,
			Description: val.Description,
			ValueStr1:   val.ValueStr1,
			ValueStr2:   val.ValueStr2,
			ValueStr3:   val.ValueStr3,
			ValueNum1:   val.ValueNum1,
			ValueNum2:   val.ValueNum2,
			ValueNum3:   val.ValueNum3,
			ValueDate1:  date_util.StringToNilTime(val.ValueDate1),
			ValueDate2:  date_util.StringToNilTime(val.ValueDate2),
			ValueDate3:  date_util.StringToNilTime(val.ValueDate3),
		})
	}

	if err := l.lovRepo.UpdateLovDetail(&dt); errors.Is(err, constanta.DbFailedToUpdateData) {
		return util.CreateGlobalResponse(constanta.ErrorLovDetailFailedUpdate, nil)
	}
	return util.CreateGlobalResponse(constanta.Success, nil)
}

func (l lovService) RemoveLovDetail(detailId string) (resp response.GlobalResponse) {
	if dt, err := l.lovRepo.FindLovDetailByDetailId(detailId); dt != nil && err == nil {
		if err = l.lovRepo.DeleteLovDetail(dt.LovDetailId); err == nil {
			return util.CreateGlobalResponse(constanta.Success, nil)
		}
		return util.CreateGlobalResponse(constanta.ErrorLovDetailFailedToRemove, nil)
	}
	return util.CreateGlobalResponse(constanta.ErrorLovDetailNotFound, nil)
}

func (l lovService) CreateNewPermission(req lov_req.CreatePermissionReq) (resp response.GlobalResponse) {
	if dt, err := l.lovRepo.SavePermission(&tbl_list_of_permission.ListOfPermission{
		Name: req.Name,
		Code: strings.ToUpper(string_util.TrimSpace(req.Name)),
	}); err == nil && dt != nil {
		return util.CreateGlobalResponse(constanta.Success, dt)
	}
	return util.CreateGlobalResponse(constanta.ErrorCreateLoPermission, nil)
}

func (l lovService) UpdatePermission(req lov_req.UpdatePermissionReq) (resp response.GlobalResponse) {
	oldCode := strings.ToUpper(string_util.TrimSpace(req.OldName))
	newCode := strings.ToUpper(string_util.TrimSpace(req.NewName))

	if err := l.lovRepo.UpdatePermission(oldCode, newCode, req.NewName); err == nil {
		return util.CreateGlobalResponse(constanta.Success, nil)
	}
	return util.CreateGlobalResponse(constanta.ErrorUpdateLoPermission, nil)
}

func (l lovService) RemovePermission(code string) (resp response.GlobalResponse) {
	temp := strings.ToUpper(string_util.TrimSpace(code))
	if err := l.lovRepo.DeletePermission(temp); err == nil {
		return util.CreateGlobalResponse(constanta.Success, nil)
	}
	return util.CreateGlobalResponse(constanta.ErrorRemoveLoPermission, nil)
}

func (l lovService) GetAllListPermission() (resp response.GlobalResponse) {
	if dt, err := l.lovRepo.FindAllPermission(); dt != nil && err == nil {
		return util.CreateGlobalResponse(constanta.Success, dt)
	}
	return util.CreateGlobalResponse(constanta.ErrorNotFoundLoPermission, nil)
}
