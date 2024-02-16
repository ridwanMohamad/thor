package properties

import "github.com/labstack/echo/v4"

type CustomContext struct {
	echo.Context
}

func (c *CustomContext) Foo() {
	c.Logger().Info("Foo")
}
