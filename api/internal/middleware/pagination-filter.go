package middleware

import (
	"errors"
	"fmt"
	"kfamily/internal/dto"
	"kfamily/internal/util"
	"net/http"
	"sort"
	"strings"

	"github.com/labstack/echo/v5"
)

// PaginationFilter validates and applies sort/pagination/search query params.
// sortableColumns maps the sort key a client sends (matching the frontend
// table's column id, e.g. "createdAt") to the real column to order by in SQL
// (e.g. "created_at"). Keeping these separate lets UI-facing camelCase sort
// keys differ from the underlying snake_case Postgres columns, and lets a
// key map to a column that doesn't share its name at all (e.g. a composite
// display column mapped to one representative column).
func (mw *Middleware) PaginationFilter(sortableColumns map[string]string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			//parse input
			var filtersInput dto.FiltersRequest
			var filters dto.Filters

			err := c.Bind(&filtersInput)
			if err != nil {
				return c.JSON(http.StatusBadRequest, util.HttpError(err, "Unable to parse input filter values"))
			}

			//page query validation
			if err = filtersInput.Validate(); err != nil {
				return c.JSON(http.StatusBadRequest, util.HttpError(err, "Filter query values not correct"))
			}

			filters.PageIndex = filtersInput.PageIndex
			filters.ResultsPerPage = filtersInput.ResultsPerPage
			filters.Search = filtersInput.Search
			filters.IncludeDeletedRecords = filtersInput.IncludeDeleted

			// sort
			if filtersInput.Sort != "" {
				var sortColumnAndDirection = strings.Split(filtersInput.Sort, ":")
				if len(sortColumnAndDirection) != 2 {
					return c.JSON(http.StatusBadRequest, util.HttpError(err, "sort column should be formatted with <colname>:asc|desc"))
				}

				var sortColumn = sortColumnAndDirection[0]
				var sortDirection = sortColumnAndDirection[1]
				// validate sort column and direction
				dbColumn, ok := sortableColumns[sortColumn]
				if !ok {
					allowed := make([]string, 0, len(sortableColumns))
					for k := range sortableColumns {
						allowed = append(allowed, k)
					}
					sort.Strings(allowed)
					var errMsg = fmt.Sprintf("sort column must be one of the following. %#v", allowed)
					return c.JSON(http.StatusBadRequest, util.HttpError(errors.New(errMsg), errMsg))
				}

				if sortDirection != "asc" && sortDirection != "desc" {
					return c.JSON(http.StatusBadRequest, util.HttpError(err, "sort direction must be asc or desc"))
				}

				//assign the column and direction into filter
				filters.SortColumn = dbColumn
				filters.SortDirection = sortDirection
			}

			// set user in context
			c.Set("filters", filters)

			return next(c)
		}
	}
}
