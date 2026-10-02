package operation

import (
	"database/sql"
	"encoding/json"
	"errors"
	"kfamily/internal/dto"
	"kfamily/internal/util"
	"net/http"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

//------------------------------------------------------------------------------

// GetOne godoc
// @Summary      Get one operation
// @Tags         Operations
// @Security     ApiKeyAuth
// @Produce      json
// @Description  An operation with the record changes it caused (at most 1000; compare with changeCount), in the order they were made.
// @Param        id   path  string  true  "operation id (UUID)"
// @Success      200  {object}  util.HttpResult{result=operation.DetailDto}
// @Failure      400,401,404,500  {object}  util.HttpResult
// @Router       /api/v1/operations/{id} [get]
func (h *Handler) GetOne(c *echo.Context) error {
	id, err := util.ParseUUIDParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, util.HttpError(err, "Operation id is not a valid UUID"))
	}

	result, err := h.service.GetOne(c.Request().Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, util.HttpErrorMessage("Operation not found"))
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, util.HttpError(err, "Error while retrieving operation"))
	}

	return c.JSON(http.StatusOK, util.HttpData(result))
}

// List godoc
// @Summary      List operations
// @Tags         Operations
// @Security     ApiKeyAuth
// @Produce      json
// @Description  One page of operations (user actions that wrote data), optionally narrowed by a free-text search over kind, target table and client (app).
// @Param        pageIndex       query  int     false  "Page index, from 0"  default(0)
// @Param        recordsPerPage  query  int     false  "Rows per page (1-150)"  default(10)
// @Param        sortBy          query  string  false  "Sort as <column>:asc|desc, e.g. createdAt:desc"
// @Param        searchText      query  string  false  "Free-text search"
// @Success      200  {object}  util.HttpResult{result=dto.PaginationResponse[operation.Dto]}
// @Failure      400,401,500  {object}  util.HttpResult
// @Router       /api/v1/operations [get]
func (h *Handler) List(c *echo.Context) error {
	var filters dto.Filters = util.FiltersFromContext(c)

	result, err := h.service.List(c.Request().Context(), filters)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.HttpError(err, "Error in getting operation list"))
	}

	return c.JSON(http.StatusOK, util.HttpData(result))
}

// Metadata godoc
// @Summary      Operations table columns
// @Tags         Operations
// @Security     ApiKeyAuth
// @Produce      json
// @Description  Column names and Postgres data types - valid sort/filter columns for the query builder.
// @Success      200  {object}  util.HttpResult{result=[]dto.ColumnMeta}
// @Failure      401,500  {object}  util.HttpResult
// @Router       /api/v1/operations/meta [get]
func (h *Handler) Metadata(c *echo.Context) error {
	columns, err := h.service.Metadata(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.HttpError(err, "Error while retrieving operation metadata"))
	}

	return c.JSON(http.StatusOK, util.HttpData(columns))
}

type QueryInputDto struct {
	WhereCondition               string `query:"whereCondition"`
	WhereConditionParametersJson string `query:"whereConditionParametersJson"`
}

// Query godoc
// @Summary      Advanced query for operations
// @Tags         Operations
// @Security     ApiKeyAuth
// @Produce      json
// @Description  Operations matching a parameterized SQL WHERE fragment from the query builder: quoted column names and :p1, :p2... placeholders, validated server-side.
// @Param        pageIndex       query  int     false  "Page index, from 0"  default(0)
// @Param        recordsPerPage  query  int     false  "Rows per page (1-150)"  default(10)
// @Param        sortBy          query  string  false  "Sort as <column>:asc|desc, e.g. createdAt:desc"
// @Param        whereCondition                query  string  true   "WHERE fragment with double-quoted column names and :pN placeholders"
// @Param        whereConditionParametersJson  query  string  false  "JSON array of the :p1, :p2... values, in order"
// @Success      200  {object}  util.HttpResult{result=dto.PaginationResponse[operation.Dto]}
// @Failure      400,401,500  {object}  util.HttpResult
// @Router       /api/v1/operations/query [get]
func (h *Handler) Query(c *echo.Context) error {
	var input QueryInputDto
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, util.HttpError(err, "Unable to parse input values"))
	}

	if input.WhereCondition == "" {
		return c.JSON(http.StatusBadRequest, util.HttpErrorMessage("whereCondition is required"))
	}

	var whereConditionParams []any
	if input.WhereConditionParametersJson != "" {
		if err := json.Unmarshal([]byte(input.WhereConditionParametersJson), &whereConditionParams); err != nil {
			return c.JSON(http.StatusBadRequest, util.HttpError(err, "whereConditionParametersJson is not valid JSON"))
		}
	}

	var filters dto.Filters = util.FiltersFromContext(c)

	result, err := h.service.Query(c.Request().Context(), filters, input.WhereCondition, whereConditionParams)
	if err != nil {
		return c.JSON(http.StatusBadRequest, util.HttpError(err, "Error in running operation query"))
	}

	return c.JSON(http.StatusOK, util.HttpData(result))
}
