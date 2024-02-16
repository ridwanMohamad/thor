package auth_service

import (
	"thor/src/payload/auth_req"
	"thor/src/payload/response"
	"thor/src/properties"
)

type IAuthService interface {
	Login(req auth_req.LoginReq) (resp response.GlobalResponse)
	LoginV2(req auth_req.LoginReq) (resp response.GlobalResponse)
	Logout(token string) (resp response.GlobalResponse)
	ForgotPassword(prop *properties.CustomContext, email string) (resp response.GlobalResponse)
	RevokeAccess(username string) (resp response.GlobalResponse)
	CheckToken(accessToken string) (resp response.GlobalResponse)
}
