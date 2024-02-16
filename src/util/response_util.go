package util

import (
	"thor/src/constanta"
	"thor/src/payload/response"
)

func CreateGlobalResponse(dt constanta.ErrMsg, data interface{}) (resp response.GlobalResponse) {
	resp = response.GlobalResponse{
		Status:  dt.Code,
		Message: dt.Msg,
		Data:    data,
	}
	return
}
