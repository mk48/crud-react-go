package auth

import (
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type Product struct {
	Controller *Handler
}

// Init wires up the public (unauthenticated) sign-in endpoint - this is how
// a client gets a token in the first place, so it can't sit behind
// AuthMiddleware. It's rate limited per client IP instead: a real sign-in
// is one request, so a burst of 5 then 1 every 6s is plenty for people and
// stops anyone hammering the Casdoor code exchange through it.
func Init(echo *echo.Group) *Product {
	handler := NewHandler()

	signinLimiter := middleware.RateLimiter(middleware.NewRateLimiterMemoryStoreWithConfig(
		middleware.RateLimiterMemoryStoreConfig{
			Rate:      1.0 / 6, // requests per second
			Burst:     5,
			ExpiresIn: 10 * time.Minute,
		},
	))

	apiGroup := echo.Group("/v1/auth")
	apiGroup.POST("/signin", handler.Signin, signinLimiter)

	return &Product{
		Controller: handler,
	}
}
