package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupMiddleware(e *echo.Echo) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		//AllowHeaders: []string{"Origin", "Authorization", "Access-Control-Allow-Origin", echo.HeaderContentType, "Accept", "Content-Length", "Accept-Encoding", "X-CSRF-Token"},
	}))
}
