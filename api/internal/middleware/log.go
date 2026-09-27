package middleware

import (
	"fmt"

	"github.com/labstack/echo/v5"
)

func (mw *Middleware) LogMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		fmt.Sprintln("log middleware")
		// Call the next handler in the chain
		return next(c)
	}
}
