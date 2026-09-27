package audithistory

import (
	"kfamily/internal/middleware"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
)

type (
	Product struct {
		Controller *Handler
		Service    *Service
	}
)

func Init(db *sqlx.DB, echo *echo.Group, middleware *middleware.Middleware) *Product {

	service := NewService(db)
	handler := NewHandler(service)

	//setup routes
	apiGroup := echo.Group("/v1/audit-history")
	apiGroup.Use(middleware.AuthMiddleware)

	apiGroup.GET("/:id", handler.List)

	return &Product{
		Controller: handler,
		Service:    service,
	}
}
