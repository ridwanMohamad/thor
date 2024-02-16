package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"thor/src/constanta"
	"thor/src/payload/auth_req"
	"thor/src/properties"
	"thor/src/service/auth_service"
	"thor/src/util"
)

type authHandler struct {
	authService auth_service.IAuthService
}

func NewAuthHandler(service auth_service.IAuthService) *authHandler {
	return &authHandler{authService: service}
}

func (h authHandler) Login(ctx echo.Context) (err error) {
	req := auth_req.LoginReq{}

	if err = ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorBindingRequest, err))
	}

	if err = ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, err))
	}

	resp := h.authService.Login(req)

	return ctx.JSON(http.StatusOK, resp)
}

func (h authHandler) LoginV2(ctx echo.Context) (err error) {
	req := auth_req.LoginReq{}

	if err = ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorBindingRequest, err))
	}

	if err = ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, err))
	}

	resp := h.authService.LoginV2(req)

	return ctx.JSON(http.StatusOK, resp)
}

func (h authHandler) Logout(ctx echo.Context) (err error) {

	accessToken := ctx.Request().Header.Get("accessToken")

	if accessToken == "" {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, err))
	}

	resp := h.authService.Logout(accessToken)

	return ctx.JSON(http.StatusOK, resp)
}

func (h authHandler) ForgotPassword(ctx echo.Context) (err error) {
	c := ctx.(*properties.CustomContext)
	email := ctx.Request().Header.Get("email")

	if email == "" {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, err))
	}

	resp := h.authService.ForgotPassword(c, email)

	return ctx.JSON(http.StatusOK, resp)
}

func (h authHandler) CheckToken(ctx echo.Context) (err error) {
	accessToken := ctx.Request().Header.Get("accessToken")

	if accessToken == "" {
		return ctx.JSON(http.StatusBadRequest, util.CreateGlobalResponse(constanta.ErrorInvalidRequest, err))
	}

	resp := h.authService.CheckToken(accessToken)

	httpCode := http.StatusOK

	if resp.Status == constanta.ErrorAccountLocked.Code || resp.Status == constanta.ErrorInvalidAccessToken.Code {
		httpCode = http.StatusUnauthorized
	}

	return ctx.JSON(httpCode, resp)
}
