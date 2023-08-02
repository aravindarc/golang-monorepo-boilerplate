package core

import "github.com/labstack/echo/v4"

func newRouter() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	return e
}
