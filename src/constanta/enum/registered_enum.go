package enum

type RegisteredEnum string

func (e RegisteredEnum) string() string {
	return string(e)
}

var (
	Registered   = RegisteredEnum("registered")
	UnRegistered = RegisteredEnum("unregistered")
)
