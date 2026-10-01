package operation

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

// Init registers the read-only operations endpoints. Operations are written
// only by util.RunOperation, as part of the action they record.
func Init(ctx context.Context, log *slog.Logger, db *sqlx.DB, echo *echo.Group, middleware *middleware.Middleware) *Product {

	service := NewService(db)
	handler := NewHandler(service)

	//setup routes
	apiGroup := echo.Group("/v1/operations")
	apiGroup.Use(middleware.AuthMiddleware)

	sortableColumns := util.BuildSortableColumns(ctx, log, db, tableName, nil)

	apiGroup.GET("", handler.List, middleware.PaginationFilter(sortableColumns))
	apiGroup.GET("/query", handler.Query, middleware.PaginationFilter(sortableColumns))
	apiGroup.GET("/meta", handler.Metadata)
	apiGroup.GET("/:id", handler.GetOne)

	return &Product{
		Controller: handler,
		Service:    service,
	}
}
