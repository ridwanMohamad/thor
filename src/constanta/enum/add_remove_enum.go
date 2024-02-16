package enum

type AddRemoveEnum string

func (e AddRemoveEnum) string() string {
	return string(e)
}

var (
	Add      = AddRemoveEnum("add")
	Exist    = AddRemoveEnum("exist")
	Modified = AddRemoveEnum("modified")
	Removed  = AddRemoveEnum("removed")
)
