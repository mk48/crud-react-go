package sample

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
	apiGroup := echo.Group("/v1/samples")
	apiGroup.Use(middleware.AuthMiddleware)

	sortableColumns := util.BuildSortableColumns(ctx, log, db, tableName, map[string]string{
		"actionAt": "updated_at",
		"actionBy": "updated_by",
	})

	//crud
	apiGroup.GET("", handler.List, middleware.PaginationFilter(sortableColumns))
	apiGroup.GET("/query", handler.Query, middleware.PaginationFilter(sortableColumns))
	apiGroup.GET("/meta", handler.Metadata)
	apiGroup.GET("/:id", handler.GetOne)
	apiGroup.PUT("/:id", handler.Update, middleware.AdminMiddleware)
	apiGroup.POST("", handler.Create, middleware.AdminMiddleware)
	apiGroup.DELETE("/:id", handler.Delete, middleware.AdminMiddleware)

	return &Product{
		Controller: handler,
		Service:    service,
	}
}
