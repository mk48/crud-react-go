package auth

import (
	"github.com/labstack/echo/v5"
)

type Product struct {
	Controller *Handler
}

// Init wires up the public (unauthenticated) sign-in endpoint - this is how
// a client gets a token in the first place, so it can't sit behind
// AuthMiddleware.
func Init(echo *echo.Group) *Product {
	handler := NewHandler()

	apiGroup := echo.Group("/v1/auth")
	apiGroup.POST("/signin", handler.Signin)

	return &Product{
		Controller: handler,
	}
}
