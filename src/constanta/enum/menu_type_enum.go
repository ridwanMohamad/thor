package enum

type MenuTypeEnum string

func (a MenuTypeEnum) string() string {
	return string(a)
}

var (
	Parent = MenuTypeEnum("parent")
	Child  = MenuTypeEnum("child")
)
