package util

import (
	"fmt"

	"kfamily/internal/dto"
	"kfamily/internal/model"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

// ParseUUIDParam parses c's named path param as a UUID, so every handler
// reports the same "not a valid UUID" message instead of each writing its own.
func ParseUUIDParam(c *echo.Context, name string) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Param(name))
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%s is not a valid UUID: %s", name, c.Param(name))
	}
	return id, nil
}

// CurrentUser returns the user AuthMiddleware attached to c.
func CurrentUser(c *echo.Context) *model.User {
	return c.Get("user").(*model.User)
}

// FiltersFromContext returns the filters middleware.PaginationFilter attached to c.
func FiltersFromContext(c *echo.Context) dto.Filters {
	return c.Get("filters").(dto.Filters)
}
