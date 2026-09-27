package healthcheck

import (
	"kfamily/internal/util"

	"github.com/labstack/echo/v5"
)

func Init(echo *echo.Echo, envs *util.AppENV) {

	handler := NewHandler(envs)

	//setup routes
	echo.GET("/healthcheck", handler.HealthCheckHandler)
}
