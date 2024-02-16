package enum

type StatusEnum string

func (e StatusEnum) string() string {
	return string(e)
}

var (
	Active   = StatusEnum("active")
	InActive = StatusEnum("inactive")
)
