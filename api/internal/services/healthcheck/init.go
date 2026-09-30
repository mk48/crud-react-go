package healthcheck

import (
	"kfamily/internal/util"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
)

func Init(echo *echo.Echo, envs *util.AppENV, db *sqlx.DB) {
	handler := NewHandler(envs, db)

	//setup routes
	echo.GET("/healthcheck", handler.HealthCheckHandler)
	echo.GET("/healthcheck/ready", handler.ReadinessHandler)
}
