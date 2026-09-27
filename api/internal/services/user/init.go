package user

import (
	"context"
	"kfamily/internal/middleware"
	"kfamily/internal/util"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
)

type (
	Product struct {
		Controller *Handler
		Service    *Service
	}
)

func Init(ctx context.Context, log *slog.Logger, db *sqlx.DB, echo *echo.Group, middleware *middleware.Middleware) *Product {

	service := NewService(db)
	handler := NewHandler(service)

	//setup routes
	apiGroup := echo.Group("/v1/users")
	apiGroup.Use(middleware.AuthMiddleware)

	// "actionAt"/"actionBy" are a composite audit column in the UI (most
	// recent of created/updated) with no single matching DB column - sort by
	// the updated column as a representative approximation.
	sortableColumns := util.BuildSortableColumns(ctx, log, db, tableName, map[string]string{
		"actionAt": "updated_at",
		"actionBy": "updated_by",
	})

	// crud - no Create: users are provisioned automatically by
	// AuthMiddleware the first time they sign in.
	apiGroup.GET("", handler.List, middleware.PaginationFilter(sortableColumns))
	apiGroup.GET("/query", handler.Query, middleware.PaginationFilter(sortableColumns))
	apiGroup.GET("/meta", handler.Metadata)
	// "/me" must be registered before "/:id" so it isn't shadowed by the
	// param route.
	apiGroup.GET("/me", handler.Me)
	apiGroup.GET("/:id", handler.GetOne)
	apiGroup.PUT("/:id", handler.Update, middleware.AdminMiddleware)
	apiGroup.DELETE("/:id", handler.Delete, middleware.AdminMiddleware)

	return &Product{
		Controller: handler,
		Service:    service,
	}
}
