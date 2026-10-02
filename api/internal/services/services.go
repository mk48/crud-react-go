package services

import (
	"context"
	"kfamily/internal/middleware"
	"kfamily/internal/services/audithistory"
	"kfamily/internal/services/auth"
	"kfamily/internal/services/healthcheck"
	"kfamily/internal/services/operation"
	"kfamily/internal/services/sample"
	"kfamily/internal/services/samplechild"
	"kfamily/internal/services/user"
	"kfamily/internal/util"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
)

func CreateServices(ctx context.Context, db *sqlx.DB, echo *echo.Echo, envs *util.AppENV, clients *util.ClientRegistry) {

	m := middleware.NewMiddleware(db, envs, clients)

	healthcheck.Init(echo, envs, db)

	apiGroup := echo.Group("/api")
	auth.Init(apiGroup)
	user.Init(ctx, echo.Logger, db, apiGroup, m)
	sample.Init(ctx, echo.Logger, db, apiGroup, m)
	samplechild.Init(ctx, echo.Logger, db, apiGroup, m)
	audithistory.Init(db, apiGroup, m)
	operation.Init(ctx, echo.Logger, db, apiGroup, m)

}
