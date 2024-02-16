package handler

import (
	"net/http"
	"strconv"
	"thor/src/constanta"
	"thor/src/payload/user_req"
	"thor/src/properties"
	"thor/src/service/user_service"
	"thor/src/util"

	"github.com/labstack/echo/v4"
)

type userHandler struct {
	userService user_service.IUserService
}

func NewUserHandler(usrService user_service.IUserService) *userHandler {

	return &userHandler{userService: usrService}
}

func (h *userHandler) CreateNewUser(ctx echo.Context) (err error) {
	c := ctx.(*properties.CustomContext)
	dt := user_req.UserDTO{}

	if err = ctx.Bind(&dt); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorBindingRequest, err.Error()))
	}
	if err = ctx.Validate(dt); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, err.Error()))
	}

	resp := h.userService.CreateNewUser(c, dt)

	return ctx.JSON(http.StatusOK, resp)
}

func (h *userHandler) UpdateUser(ctx echo.Context) (err error) {

	dt := user_req.UserUpdateDTO{}

	if err = ctx.Bind(&dt); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorBindingRequest, err.Error()))
	}
	if err = ctx.Validate(dt); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, err.Error()))
	}

	resp := h.userService.UpdateUser(dt)

	return ctx.JSON(http.StatusOK, resp)
}

func (h *userHandler) UpdateUserProfile(ctx echo.Context) (err error) {

	dt := user_req.UpdateUserProfile{}

	if err = ctx.Bind(&dt); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorBindingRequest, err.Error()))
	}
	if err = ctx.Validate(dt); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, err.Error()))
	}

	resp := h.userService.UpdateUserProfile(dt)

	return ctx.JSON(http.StatusOK, resp)
}

func (h *userHandler) GetAllUsers(ctx echo.Context) (err error) {
	resp := h.userService.GetAllUser()

	return ctx.JSON(http.StatusOK, resp)
}

func (h *userHandler) GetUserByUserId(ctx echo.Context) (err error) {
	par := ctx.Param("userId")

	if par == "" {
		return ctx.JSON(http.StatusOK, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, nil))
	}
	id := -1

	if id, err = strconv.Atoi(par); err != nil {
		return ctx.JSON(http.StatusOK, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, nil))
	}
	resp := h.userService.GetUserById(int64(id))

	return ctx.JSON(http.StatusOK, resp)
}

func (h *userHandler) ChangeUserPassword(ctx echo.Context) (err error) {
	dt := user_req.UserChangePasswordReq{}

	if err = ctx.Bind(&dt); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorBindingRequest, err.Error()))
	}
	if err = ctx.Validate(dt); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, err.Error()))
	}

	resp := h.userService.ChangePassword(dt)

	return ctx.JSON(http.StatusOK, resp)
}

func (h *userHandler) InActiveUser(ctx echo.Context) (err error) {
	par := ctx.Param("userId")

	if par == "" {
		return ctx.JSON(http.StatusOK, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, nil))
	}

	id := -1

	if id, err = strconv.Atoi(par); err != nil {
		return ctx.JSON(http.StatusOK, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, nil))
	}

	resp := h.userService.InActiveUser(int64(id))

	return ctx.JSON(http.StatusOK, resp)
}

func (h *userHandler) GetMappingLocationUser(ctx echo.Context) (err error) {
	par := ctx.Param("username")

	if par == "" {
		return ctx.JSON(http.StatusOK, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, nil))
	}
	resp := h.userService.GetMappingLocationByUsername(par)

	return ctx.JSON(http.StatusOK, resp)
}

func (h *userHandler) UpdateMappingLocationUser(ctx echo.Context) (err error) {

	dt := user_req.MappingUserLocationReq{}

	if err = ctx.Bind(&dt); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorBindingRequest, err.Error()))
	}
	if err = ctx.Validate(dt); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, err.Error()))
	}

	resp := h.userService.UpdateMappingLocationUser(dt)

	return ctx.JSON(http.StatusOK, resp)
}
func (h *userHandler) RemoveMappingLocationUser(ctx echo.Context) (err error) {
	par := ctx.Param("id")

	if par == "" {
		return ctx.JSON(http.StatusOK, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, nil))
	}
	id := -1

	if id, err = strconv.Atoi(par); err != nil {
		return ctx.JSON(http.StatusOK, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, nil))
	}
	resp := h.userService.RemoveMappingLocationUser(int64(id))

	return ctx.JSON(http.StatusOK, resp)
}

func (h *userHandler) GetMappingRoleUser(ctx echo.Context) (err error) {
	par := ctx.Param("username")

	if par == "" {
		return ctx.JSON(http.StatusOK, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, nil))
	}
	resp := h.userService.GetMappingRoleByUsername(par)

	return ctx.JSON(http.StatusOK, resp)
}

func (h *userHandler) UpdateMappingRoleUser(ctx echo.Context) (err error) {

	dt := user_req.MappingUserRoleReq{}

	if err = ctx.Bind(&dt); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorBindingRequest, err.Error()))
	}
	if err = ctx.Validate(dt); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, err.Error()))
	}

	resp := h.userService.UpdateMappingRoleUser(dt)

	return ctx.JSON(http.StatusOK, resp)
}
func (h *userHandler) RemoveMappingRoleUser(ctx echo.Context) (err error) {
	par := ctx.Param("id")

	if par == "" {
		return ctx.JSON(http.StatusOK, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, nil))
	}
	id := -1

	if id, err = strconv.Atoi(par); err != nil {
		return ctx.JSON(http.StatusOK, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, nil))
	}
	resp := h.userService.RemoveMappingRoleUser(int64(id))

	return ctx.JSON(http.StatusOK, resp)
}
