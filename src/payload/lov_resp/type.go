package lov_resp

import (
	"gopkg.in/guregu/null.v4"
)

type LovAfterCreateResp struct {
	LovHeaderId string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Detail      []ShardLovDetail `json:"detail"`
}

type ShardLovHeaders struct {
	HeaderId    string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

type ShardLovDetail struct {
	LovDetailId string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ValueStr1   string    `json:"valueString1"`
	ValueStr2   string    `json:"valueString2"`
	ValueStr3   string    `json:"valueString3"`
	ValueNum1   int64     `json:"valueNumber1"`
	ValueNum2   int64     `json:"valueNumber2"`
	ValueNum3   int64     `json:"valueNumber3"`
	ValueDate1  null.Time `json:"valueDate1"`
	ValueDate2  null.Time `json:"valueDate2"`
	ValueDate3  null.Time `json:"valueDate3"`
}
