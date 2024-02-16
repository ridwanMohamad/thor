package lov_req

type CreateLovReq struct {
	Name        string         `json:"name" validate:"required"`
	Description string         `json:"description" validate:"required"`
	Enabled     bool           `json:"enabled"`
	LovDetails  []LovDetailReq `json:"lovDetails" validate:"required"`
}

type LovHeaderReq struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Enabled     bool   `json:"enabled"`
}

type LovDetailReq struct {
	LovDetailId string `json:"lovDetailId"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	ValueStr1   string `json:"valueString1"`
	ValueStr2   string `json:"valueString2"`
	ValueStr3   string `json:"valueString3"`
	ValueNum1   int64  `json:"valueNumber1"`
	ValueNum2   int64  `json:"valueNumber2"`
	ValueNum3   int64  `json:"valueNumber3"`
	ValueDate1  string `json:"valueDate1" validate:"ISO8601date"`
	ValueDate2  string `json:"valueDate2" validate:"ISO8601date"`
	ValueDate3  string `json:"valueDate3" validate:"ISO8601date"`
}

type CreatePermissionReq struct {
	Name string `json:"name" validate:"required"`
}
type UpdatePermissionReq struct {
	OldName string `json:"oldName" validate:"required"`
	NewName string `json:"newName"`
	//Status  enum.AddRemoveEnum `json:"status"`
}
