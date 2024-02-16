package handler

import (
	"github.com/labstack/echo/v4"
	"gopkg.in/guregu/null.v4"
	"net/http"
	"thor/src/constanta"
	"thor/src/payload/lov_req"
	"thor/src/service/lov_service"
	"thor/src/util"
)

type lovHandler struct {
	service lov_service.ILovService
}

func NewLovHandler(lovService lov_service.ILovService) *lovHandler {

	return &lovHandler{service: lovService}
}

func (h lovHandler) CreateNewLov(ctx echo.Context) (err error) {
	dt := lov_req.CreateLovReq{}

	if err = ctx.Bind(&dt); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorBindingRequest, err.Error()))
	}
	if err = ctx.Validate(dt); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, err.Error()))
	}

	for _, val := range dt.LovDetails {
		dtD := val
		//if err = ctx.Bind(&dtD); err != nil {
		//	return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorBindingRequest, err.Error()))
		//}
		if err = ctx.Validate(dtD); err != nil {
			return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, err.Error()))
		}
	}

	resp := h.service.CreateNewLov(dt)

	return ctx.JSON(http.StatusOK, resp)
}

func (h lovHandler) CreateNewLovDetail(ctx echo.Context) (err error) {
	par := ctx.Param("headerId")

	if par == "" {
		return ctx.JSON(http.StatusOK, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, nil))
	}

	var dt []lov_req.LovDetailReq

	if err = ctx.Bind(&dt); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorBindingRequest, err.Error()))
	}
	for _, val := range dt {
		if err = ctx.Validate(val); err != nil {
			return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, err.Error()))
		}
	}

	resp := h.service.CreateNewLovDetail(par, dt)

	return ctx.JSON(http.StatusOK, resp)
}

func (h lovHandler) GetAllLovHeader(ctx echo.Context) (err error) {
	resp := h.service.GetAllLovHeader()

	return ctx.JSON(http.StatusOK, resp)
}

func (h lovHandler) GetAllLovByHeaderIdOrName(ctx echo.Context) (err error) {
	par := ctx.Param("headerId")

	if par == "" {
		return ctx.JSON(http.StatusOK, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, nil))
	}

	resp := h.service.GetLovByHeaderName(par)

	return ctx.JSON(http.StatusOK, resp)
}

func (h lovHandler) GetAllLovByHeaderIdOrNameTest(ctx echo.Context) (err error) {
	par := ctx.QueryParam("header")

	if par == "" {
		return ctx.JSON(http.StatusOK, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, nil))
	}

	resp := h.service.GetLovByHeaderName(par)

	return ctx.JSON(http.StatusOK, resp)
}

func (h lovHandler) UpdateLovHeader(ctx echo.Context) (err error) {
	var dt lov_req.CreateLovReq

	par := ctx.Param("headerId")

	if par == "" {
		return ctx.JSON(http.StatusOK, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, nil))
	}
	if err = ctx.Bind(&dt); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorBindingRequest, err.Error()))
	}
	for _, val := range dt.LovDetails {
		if err = ctx.Validate(val); err != nil {
			return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, err.Error()))
		}
	}

	resp := h.service.UpdateLov(par, dt)

	return ctx.JSON(http.StatusOK, resp)
}

func (h lovHandler) RemoveLovHeader(ctx echo.Context) (err error) {
	par := ctx.Param("headerId")

	if par == "" {
		return ctx.JSON(http.StatusOK, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, nil))
	}

	resp := h.service.RemoveLovHeader(par)

	return ctx.JSON(http.StatusOK, resp)
}

func (h lovHandler) GetLovDetailByHeaderId(ctx echo.Context) (err error) {
	panic("not implemented")
}

func (h lovHandler) UpdateLovDetail(ctx echo.Context) (err error) {
	var dt []lov_req.LovDetailReq

	if err = ctx.Bind(&dt); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorBindingRequest, err.Error()))
	}
	for _, val := range dt {
		if val.LovDetailId == "" || !null.StringFrom(val.LovDetailId).Valid {
			return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, nil))
		}
		if err = ctx.Validate(val); err != nil {
			return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, err.Error()))
		}
	}

	resp := h.service.UpdateLovDetailByDetailId(dt)

	return ctx.JSON(http.StatusOK, resp)
}

func (h lovHandler) RemoveLovDetail(ctx echo.Context) (err error) {
	par := ctx.Param("detailId")

	if par == "" {
		return ctx.JSON(http.StatusOK, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, nil))
	}

	resp := h.service.RemoveLovDetail(par)

	return ctx.JSON(http.StatusOK, resp)
}

func (h lovHandler) CreateNewListOfPermission(ctx echo.Context) (err error) {
	dt := lov_req.CreatePermissionReq{}

	if err = ctx.Bind(&dt); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorBindingRequest, err.Error()))
	}
	if err = ctx.Validate(dt); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, err.Error()))
	}

	resp := h.service.CreateNewPermission(dt)

	return ctx.JSON(http.StatusOK, resp)
}

func (h lovHandler) UpdateListOfPermission(ctx echo.Context) (err error) {
	var dt lov_req.UpdatePermissionReq

	if err = ctx.Bind(&dt); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorBindingRequest, err.Error()))
	}
	if err = ctx.Validate(dt); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, err.Error()))
	}

	resp := h.service.UpdatePermission(dt)

	return ctx.JSON(http.StatusOK, resp)
}

func (h lovHandler) RemoveListOfPermission(ctx echo.Context) (err error) {
	par := ctx.Param("code")

	if par == "" {
		return ctx.JSON(http.StatusOK, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, nil))
	}

	resp := h.service.RemovePermission(par)

	return ctx.JSON(http.StatusOK, resp)
}

func (h lovHandler) GetAllListOfPermission(ctx echo.Context) (err error) {
	resp := h.service.GetAllListPermission()

	return ctx.JSON(http.StatusOK, resp)
}
