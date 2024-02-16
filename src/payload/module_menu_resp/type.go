package module_menu_resp

import (
	"gopkg.in/guregu/null.v4"
	"thor/src/constanta/enum"
	"thor/src/domain/tbl_menu_permission"
	"time"
)

type MenuResponse struct {
	MenuId     null.String                          `json:"MenuId"`
	ParentId   null.String                          `json:"ParentId"`
	Name       string                               `json:"Name"`
	Path       string                               `json:"Path"`
	Type       enum.MenuTypeEnum                    `json:"Type"`
	MenuCode   string                               `json:"MenuCode"`
	MenuIcon   string                               `json:"MenuIcon"`
	CreatedAt  time.Time                            `json:"CreatedAt"`
	UpdatedAt  null.Time                            `json:"UpdatedAt"`
	Permission []tbl_menu_permission.MenuPermission `json:"MenuPermission"`
}
