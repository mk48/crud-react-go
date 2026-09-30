package sample

import (
	"database/sql"
	"encoding/json"
	"errors"
	"kfamily/internal/dto"
	"kfamily/internal/util"
	"net/http"
	"strconv"

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
// @Summary      Get one sample item
// @Tags         Samples
// @Security     ApiKeyAuth
// @Produce      json
// @Param        id              path   string  true   "sample item id (UUID)"
// @Param        includeDeleted  query  bool    false  "Return it even if soft-deleted"
// @Success      200  {object}  util.HttpResult{result=sample.Dto}
// @Failure      400,401,404,500  {object}  util.HttpResult
// @Router       /api/v1/samples/{id} [get]
func (h *Handler) GetOne(c *echo.Context) error {
	id, err := util.ParseUUIDParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, util.HttpError(err, "Sample item id is not a valid UUID"))
	}

	includeDeleted, _ := strconv.ParseBool(c.QueryParam("includeDeleted"))

	row, err := h.service.GetOne(c.Request().Context(), id, includeDeleted)
	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, util.HttpErrorMessage("Sample item not found"))
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, util.HttpError(err, "Error while retrieving sample item"))
	}

	return c.JSON(http.StatusOK, util.HttpData(row))
}

// List godoc
// @Summary      List sample items
// @Tags         Samples
// @Security     ApiKeyAuth
// @Produce      json
// @Description  One page of sample items, optionally narrowed by a free-text search.
// @Param        pageIndex       query  int     false  "Page index, from 0"  default(0)
// @Param        recordsPerPage  query  int     false  "Rows per page (1-150)"  default(10)
// @Param        sortBy          query  string  false  "Sort as <column>:asc|desc, e.g. createdAt:desc"
// @Param        searchText      query  string  false  "Free-text search"
// @Param        includeDeleted  query  bool    false  "Include soft-deleted rows"
// @Success      200  {object}  util.HttpResult{result=dto.PaginationResponse[sample.Dto]}
// @Failure      400,401,500  {object}  util.HttpResult
// @Router       /api/v1/samples [get]
func (h *Handler) List(c *echo.Context) error {
	var filters dto.Filters = util.FiltersFromContext(c)

	result, err := h.service.List(c.Request().Context(), filters)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.HttpError(err, "Error in getting sample item list"))
	}

	return c.JSON(http.StatusOK, util.HttpData(result))
}

// Metadata godoc
// @Summary      Samples table columns
// @Tags         Samples
// @Security     ApiKeyAuth
// @Produce      json
// @Description  Column names and Postgres data types - valid sort/filter columns for the query builder.
// @Success      200  {object}  util.HttpResult{result=[]dto.ColumnMeta}
// @Failure      401,500  {object}  util.HttpResult
// @Router       /api/v1/samples/meta [get]
func (h *Handler) Metadata(c *echo.Context) error {
	columns, err := h.service.Metadata(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.HttpError(err, "Error while retrieving sample item metadata"))
	}

	return c.JSON(http.StatusOK, util.HttpData(columns))
}

type QueryInputDto struct {
	WhereCondition               string `query:"whereCondition"`
	WhereConditionParametersJson string `query:"whereConditionParametersJson"`
}

// Query godoc
// @Summary      Advanced query for sample items
// @Tags         Samples
// @Security     ApiKeyAuth
// @Produce      json
// @Description  Rows matching a parameterized SQL WHERE fragment from the query builder: quoted column names and :p1, :p2... placeholders, validated server-side. Includes soft-deleted rows (filter on deleted_at to exclude them).
// @Param        pageIndex       query  int     false  "Page index, from 0"  default(0)
// @Param        recordsPerPage  query  int     false  "Rows per page (1-150)"  default(10)
// @Param        sortBy          query  string  false  "Sort as <column>:asc|desc, e.g. createdAt:desc"
// @Param        whereCondition                query  string  true   "WHERE fragment with double-quoted column names and :pN placeholders"
// @Param        whereConditionParametersJson  query  string  false  "JSON array of the :p1, :p2... values, in order"
// @Success      200  {object}  util.HttpResult{result=dto.PaginationResponse[sample.Dto]}
// @Failure      400,401,500  {object}  util.HttpResult
// @Router       /api/v1/samples/query [get]
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
		return c.JSON(http.StatusBadRequest, util.HttpError(err, "Error in running sample item query"))
	}

	return c.JSON(http.StatusOK, util.HttpData(result))
}

// Create godoc
// @Summary      Create a sample item
// @Tags         Samples
// @Security     ApiKeyAuth
// @Produce      json
// @Description  Admin only.
// @Accept       json
// @Param        request  body  sample.CreateInputDto  true  "New sample item"
// @Success      201  {object}  util.HttpResult{result=string}  "id of the new sample item"
// @Failure      400,401,403,500  {object}  util.HttpResult
// @Router       /api/v1/samples [post]
func (h *Handler) Create(c *echo.Context) error {
	var inputDTO CreateInputDto
	err := c.Bind(&inputDTO)
	if err != nil {
		return c.JSON(http.StatusBadRequest, util.HttpError(err, "Unable to parse input values"))
	}
	if err = inputDTO.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, util.HttpError(err, err.Error()))
	}

	loggedInUser := util.CurrentUser(c)
	newlyCreated, err := h.service.Create(c.Request().Context(), inputDTO, loggedInUser.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.HttpError(err, "Error on creating a new sample item"))
	}

	return c.JSON(http.StatusCreated, util.HttpData(newlyCreated))
}

// Update godoc
// @Summary      Update a sample item
// @Tags         Samples
// @Security     ApiKeyAuth
// @Produce      json
// @Description  Admin only. Soft-deleted rows can't be updated (404).
// @Accept       json
// @Param        id       path  string  true  "sample item id (UUID)"
// @Param        request  body  sample.UpdateInputDto  true  "New values"
// @Success      200  {object}  util.HttpResult
// @Failure      400,401,403,404,500  {object}  util.HttpResult
// @Router       /api/v1/samples/{id} [put]
func (h *Handler) Update(c *echo.Context) error {
	id, err := util.ParseUUIDParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, util.HttpError(err, "Sample item id is not a valid UUID"))
	}

	var inputDTO UpdateInputDto
	if err = c.Bind(&inputDTO); err != nil {
		return c.JSON(http.StatusBadRequest, util.HttpError(err, "Unable to parse input values"))
	}
	if err = inputDTO.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, util.HttpError(err, err.Error()))
	}

	loggedInUser := util.CurrentUser(c)
	err = h.service.Update(c.Request().Context(), id, inputDTO, loggedInUser.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, util.HttpErrorMessage("Sample item not found"))
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, util.HttpError(err, "Unable to update sample item"))
	}

	return c.JSON(http.StatusOK, util.HttpSuccessStatus())
}

// Delete godoc
// @Summary      Delete a sample item
// @Tags         Samples
// @Security     ApiKeyAuth
// @Produce      json
// @Description  Admin only. Soft delete - the row and its audit history are kept.
// @Param        id  path  string  true  "sample item id (UUID)"
// @Success      200  {object}  util.HttpResult
// @Failure      400,401,403,404,500  {object}  util.HttpResult
// @Router       /api/v1/samples/{id} [delete]
func (h *Handler) Delete(c *echo.Context) error {
	id, err := util.ParseUUIDParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, util.HttpError(err, "Sample item id is not a valid UUID"))
	}

	loggedInUser := util.CurrentUser(c)
	err = h.service.Delete(c.Request().Context(), id, loggedInUser.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, util.HttpErrorMessage("Sample item not found"))
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, util.HttpError(err, "Unable to delete sample item"))
	}

	return c.JSON(http.StatusOK, util.HttpSuccessStatus())
}
