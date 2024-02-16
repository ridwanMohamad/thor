package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"thor/src/constanta"
	"thor/src/payload/module_menu_req"
	"thor/src/service/module_menu_service"
	"thor/src/util"
	"time"
)

type moduleMenuHandler struct {
	service module_menu_service.IModuleMenuService
}

func NewModuleMenuHandler(svc module_menu_service.IModuleMenuService) *moduleMenuHandler {

	return &moduleMenuHandler{service: svc}
}

func (h *moduleMenuHandler) CreateNewModule(ctx echo.Context) (err error) {
	dt := module_menu_req.CreateModuleReq{}

	if err = ctx.Bind(&dt); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorBindingRequest, err.Error()))
	}
	if err = ctx.Validate(dt); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, err.Error()))
	}

	if dt.Maintenance {
		if _, err := time.Parse("2006-01-02", dt.MtcStart); err != nil {
			return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorInvalidFormatDate, err.Error()))

		}
		if _, err := time.Parse("2006-01-02", dt.MtcEnd); err != nil {
			return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorInvalidFormatDate, err.Error()))
		}
	}

	resp := h.service.CreateModule(dt)

	return ctx.JSON(http.StatusOK, resp)
}

func (h *moduleMenuHandler) UpdateModuleData(ctx echo.Context) (err error) {
	par := ctx.Param("moduleId")

	if par == "" {
		return ctx.JSON(http.StatusOK, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, nil))
	}

	dt := module_menu_req.CreateModuleReq{}

	if err = ctx.Bind(&dt); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorBindingRequest, err.Error()))
	}
	if err = ctx.Validate(dt); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, err.Error()))
	}

	if dt.Maintenance {
		if _, err := time.Parse("2006-01-02", dt.MtcStart); err != nil {
			return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorInvalidFormatDate, err.Error()))

		}
		if _, err := time.Parse("2006-01-02", dt.MtcEnd); err != nil {
			return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorInvalidFormatDate, err.Error()))
		}
	}

	resp := h.service.UpdateModule(par, dt)

	return ctx.JSON(http.StatusOK, resp)
}

func (h *moduleMenuHandler) GetAllModule(ctx echo.Context) (err error) {
	resp := h.service.GetAllModule()

	return ctx.JSON(http.StatusOK, resp)
}

func (h *moduleMenuHandler) GetModuleByIdModule(ctx echo.Context) (err error) {

	par := ctx.Param("moduleId")

	if par == "" {
		return ctx.JSON(http.StatusOK, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, nil))
	}
	resp := h.service.GetModuleById(par)

	return ctx.JSON(http.StatusOK, resp)
}

func (h *moduleMenuHandler) CreateNewMenu(ctx echo.Context) (err error) {
	dt := module_menu_req.CreateMenuReq{}

	if err = ctx.Bind(&dt); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorBindingRequest, err.Error()))
	}
	if err = ctx.Validate(dt); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, err.Error()))
	}

	resp := h.service.CreateMenu(dt)

	return ctx.JSON(http.StatusOK, resp)
}

func (h *moduleMenuHandler) UpdateMenuData(ctx echo.Context) (err error) {
	par := ctx.Param("menuId")

	if par == "" {
		return ctx.JSON(http.StatusOK, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, nil))
	}

	dt := module_menu_req.CreateMenuReq{}

	if err = ctx.Bind(&dt); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorBindingRequest, err.Error()))
	}
	if err = ctx.Validate(dt); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, err.Error()))
	}

	resp := h.service.UpdateMenu(par, dt)

	return ctx.JSON(http.StatusOK, resp)
}

func (h *moduleMenuHandler) GetAllMenu(ctx echo.Context) (err error) {
	resp := h.service.GetAllMenu()

	return ctx.JSON(http.StatusOK, resp)
}

func (h *moduleMenuHandler) GetAllMenuByModuleId(ctx echo.Context) (err error) {
	par := ctx.Param("moduleId")

	if par == "" {
		return ctx.JSON(http.StatusOK, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, nil))
	}

	resp := h.service.GetAllMenuByModuleId(par)

	return ctx.JSON(http.StatusOK, resp)
}

func (h *moduleMenuHandler) GetDetailMenuByMenuId(ctx echo.Context) (err error) {
	par := ctx.Param("menuId")

	if par == "" {
		return ctx.JSON(http.StatusOK, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, nil))
	}

	resp := h.service.GetMenuById(par)

	return ctx.JSON(http.StatusOK, resp)
}

func (h *moduleMenuHandler) GetMenuNotExistsInRoleMenu(ctx echo.Context) (err error) {
	par := ctx.Param("roleId")

	if par == "" {
		return ctx.JSON(http.StatusOK, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, nil))
	}

	resp := h.service.GetMenuByPermission(par)

	return ctx.JSON(http.StatusOK, resp)
}

func (h *moduleMenuHandler) GetAllMenuAndPermissionByUserId(ctx echo.Context) (err error) {
	par := ctx.Param("userPk")
	var pk int
	if par == "" {
		return ctx.JSON(http.StatusOK, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, nil))
	}

	if pk, err = strconv.Atoi(par); err != nil {
		return ctx.JSON(http.StatusOK, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, nil))
	}

	resp := h.service.GetUnregisteredMenuByUserId(int64(pk))

	return ctx.JSON(http.StatusOK, resp)
}

func (h *moduleMenuHandler) GetTestAllMenu(ctx echo.Context) (err error) {
	resp := h.service.TestGetAllMenu()

	return ctx.JSON(http.StatusOK, resp)
}
