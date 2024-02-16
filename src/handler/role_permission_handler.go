package handler

import (
	"net/http"
	"strconv"
	"thor/src/constanta"
	"thor/src/payload/role_permission_req"
	"thor/src/service/role_permission_service"
	"thor/src/util"

	"github.com/labstack/echo/v4"
)

type rolePermissionHandler struct {
	service role_permission_service.IRoleService
}

func NewRolePermissionHandler(service role_permission_service.IRoleService) *rolePermissionHandler {

	return &rolePermissionHandler{service: service}
}

func (h rolePermissionHandler) CreateNewRoleWithPermission(ctx echo.Context) (err error) {
	dt := role_permission_req.RoleWithPermissionReq{}

	if err = ctx.Bind(&dt); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorBindingRequest, err.Error()))
	}
	if err = ctx.Validate(dt); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, err.Error()))
	}

	resp := h.service.CreateNewRole(dt)

	return ctx.JSON(http.StatusOK, resp)
}

func (h rolePermissionHandler) UpdateRoleData(ctx echo.Context) (err error) {
	par := ctx.Param("roleId")

	if par == "" {
		return ctx.JSON(http.StatusOK, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, nil))
	}

	dt := role_permission_req.RoleWithPermissionReq{}

	if err = ctx.Bind(&dt); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorBindingRequest, err.Error()))
	}
	if err = ctx.Validate(dt); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, err.Error()))
	}

	resp := h.service.UpdateRoleAndPermission(par, dt)

	return ctx.JSON(http.StatusOK, resp)
}

func (h rolePermissionHandler) DeleteRole(ctx echo.Context) (err error) {
	par := ctx.Param("roleId")

	if par == "" {
		return ctx.JSON(http.StatusOK, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, nil))
	}

	resp := h.service.DeleteRole(par)

	return ctx.JSON(http.StatusOK, resp)
}

func (h rolePermissionHandler) GetAllRole(ctx echo.Context) (err error) {
	resp := h.service.GetAllRole()

	return ctx.JSON(http.StatusOK, resp)
}

func (h rolePermissionHandler) GetPermissionByUserId(ctx echo.Context) (err error) {
	par := ctx.Param("userId")
	roleId := ctx.QueryParam("role")

	if par == "" {
		return ctx.JSON(http.StatusOK, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, nil))
	}

	var id int

	if id, _ = strconv.Atoi(par); err != nil {
		return ctx.JSON(http.StatusOK, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, nil))
	}

	resp := h.service.GetPermissionByUserId(int64(id), roleId)

	return ctx.JSON(http.StatusOK, resp)
}

func (h rolePermissionHandler) GetPermissionByRoleId(ctx echo.Context) (err error) {
	par := ctx.Param("roleId")

	if par == "" {
		return ctx.JSON(http.StatusOK, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, nil))
	}

	resp := h.service.GetPermissionByRoleId(par)

	return ctx.JSON(http.StatusOK, resp)
}

func (h rolePermissionHandler) CreateNewAdditionalPrivilege(ctx echo.Context) (err error) {
	var dt []role_permission_req.ShardAdditionalPrivilege

	if err = ctx.Bind(&dt); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorBindingRequest, err.Error()))
	}
	if err = ctx.Validate(dt); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, err.Error()))
	}

	resp := h.service.AddAdditionalMenuToUser(dt)

	return ctx.JSON(http.StatusOK, resp)
}

func (h rolePermissionHandler) RemoveAdditionalPrivilege(ctx echo.Context) (err error) {
	var dt []role_permission_req.ShardAdditionalPrivilege

	if err = ctx.Bind(&dt); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorBindingRequest, err.Error()))
	}
	if err = ctx.Validate(dt); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, err.Error()))
	}

	resp := h.service.RemoveAdditionalMenuFromUser(dt)

	return ctx.JSON(http.StatusOK, resp)
}
