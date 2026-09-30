package services

import (
	"context"
	"kfamily/internal/crud"
	"kfamily/internal/middleware"
	"kfamily/internal/services/audithistory"
	"kfamily/internal/services/auth"
	"kfamily/internal/services/healthcheck"
	"kfamily/internal/services/sample"
	"kfamily/internal/services/samplechild"
	"kfamily/internal/services/user"
	"kfamily/internal/util"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
)

func CreateServices(ctx context.Context, db *sqlx.DB, echo *echo.Echo, envs *util.AppENV) {

	m := middleware.NewMiddleware(db, envs)

	healthcheck.Init(echo, envs, db)

	apiGroup := echo.Group("/api")
	auth.Init(apiGroup)

	// Resources on the generic crud endpoints (see crud.Register).
	deps := crud.Deps{DB: db, Log: echo.Logger, MW: m}
	user.Init(ctx, deps, apiGroup)
	sample.Init(ctx, deps, apiGroup)
	samplechild.Init(ctx, deps, apiGroup)

	audithistory.Init(db, apiGroup, m)

}
