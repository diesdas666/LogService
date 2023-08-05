package internal

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

var AppVersion = VersionRest{
	Service: "rest-echo/http",
	Version: "0.1.0",
	Build:   "1",
}

func GetVersion() func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, AppVersion)
	}
}
