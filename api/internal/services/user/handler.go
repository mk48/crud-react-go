package user

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
// @Summary      Get one user
// @Tags         Users
// @Security     ApiKeyAuth
// @Produce      json
// @Param        id              path   string  true   "user id (UUID)"
// @Param        includeDeleted  query  bool    false  "Return it even if soft-deleted"
// @Success      200  {object}  util.HttpResult{result=user.Dto}
// @Failure      400,401,404,500  {object}  util.HttpResult
// @Router       /api/v1/users/{id} [get]
func (h *Handler) GetOne(c *echo.Context) error {
	id, err := util.ParseUUIDParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, util.HttpError(err, "User Id is not a valid UUID"))
	}

	includeDeleted, _ := strconv.ParseBool(c.QueryParam("includeDeleted"))

	row, err := h.service.GetOne(c.Request().Context(), id, includeDeleted)
	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, util.HttpErrorMessage("User not found"))
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, util.HttpError(err, "Error while retrieving user"))
	}

	return c.JSON(http.StatusOK, util.HttpData(row))
}

// Me godoc
// @Summary      Get the signed-in user
// @Tags         Users
// @Security     ApiKeyAuth
// @Produce      json
// @Description  The caller's own user record (including isAdmin).
// @Success      200  {object}  util.HttpResult{result=user.Dto}
// @Failure      401,500  {object}  util.HttpResult
// @Router       /api/v1/users/me [get]
func (h *Handler) Me(c *echo.Context) error {
	currentUser := util.CurrentUser(c)

	row, err := h.service.GetOne(c.Request().Context(), currentUser.ID, false)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.HttpError(err, "Error while retrieving current user"))
	}

	return c.JSON(http.StatusOK, util.HttpData(row))
}

// List godoc
// @Summary      List users
// @Tags         Users
// @Security     ApiKeyAuth
// @Produce      json
// @Description  One page of users, optionally narrowed by a free-text search.
// @Param        pageIndex       query  int     false  "Page index, from 0"  default(0)
// @Param        recordsPerPage  query  int     false  "Rows per page (1-150)"  default(10)
// @Param        sortBy          query  string  false  "Sort as <column>:asc|desc, e.g. createdAt:desc"
// @Param        searchText      query  string  false  "Free-text search"
// @Param        includeDeleted  query  bool    false  "Include soft-deleted rows"
// @Success      200  {object}  util.HttpResult{result=dto.PaginationResponse[user.Dto]}
// @Failure      400,401,500  {object}  util.HttpResult
// @Router       /api/v1/users [get]
func (h *Handler) List(c *echo.Context) error {
	// ?results_per_page=10 & page=1 & sort=name:asc & search=kumaran & includeDeleted=false
	var filters dto.Filters = util.FiltersFromContext(c)

	result, err := h.service.List(c.Request().Context(), filters)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.HttpError(err, "Error in getting user list"))
	}

	return c.JSON(http.StatusOK, util.HttpData(result))
}

// Metadata godoc
// @Summary      Users table columns
// @Tags         Users
// @Security     ApiKeyAuth
// @Produce      json
// @Description  Column names and Postgres data types - valid sort/filter columns for the query builder.
// @Success      200  {object}  util.HttpResult{result=[]dto.ColumnMeta}
// @Failure      401,500  {object}  util.HttpResult
// @Router       /api/v1/users/meta [get]
func (h *Handler) Metadata(c *echo.Context) error {
	columns, err := h.service.Metadata(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.HttpError(err, "Error while retrieving user metadata"))
	}

	return c.JSON(http.StatusOK, util.HttpData(columns))
}

type QueryInputDto struct {
	WhereCondition               string `query:"whereCondition"`
	WhereConditionParametersJson string `query:"whereConditionParametersJson"`
}

// Query godoc
// @Summary      Advanced query for users
// @Tags         Users
// @Security     ApiKeyAuth
// @Produce      json
// @Description  Rows matching a parameterized SQL WHERE fragment from the query builder: quoted column names and :p1, :p2... placeholders, validated server-side. Includes soft-deleted rows (filter on deleted_at to exclude them).
// @Param        pageIndex       query  int     false  "Page index, from 0"  default(0)
// @Param        recordsPerPage  query  int     false  "Rows per page (1-150)"  default(10)
// @Param        sortBy          query  string  false  "Sort as <column>:asc|desc, e.g. createdAt:desc"
// @Param        whereCondition                query  string  true   "WHERE fragment with double-quoted column names and :pN placeholders"
// @Param        whereConditionParametersJson  query  string  false  "JSON array of the :p1, :p2... values, in order"
// @Success      200  {object}  util.HttpResult{result=dto.PaginationResponse[user.Dto]}
// @Failure      400,401,500  {object}  util.HttpResult
// @Router       /api/v1/users/query [get]
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
		return c.JSON(http.StatusBadRequest, util.HttpError(err, "Error in running user query"))
	}

	return c.JSON(http.StatusOK, util.HttpData(result))
}

// Update godoc
// @Summary      Update a user
// @Tags         Users
// @Security     ApiKeyAuth
// @Produce      json
// @Description  Admin only. Soft-deleted rows can't be updated (404). Omit isAdmin to leave it unchanged; removing the last admin is refused (409).
// @Accept       json
// @Param        id       path  string  true  "user id (UUID)"
// @Param        request  body  user.UpdateInputDto  true  "New values"
// @Success      200  {object}  util.HttpResult
// @Failure      400,401,403,404,409,500  {object}  util.HttpResult
// @Router       /api/v1/users/{id} [put]
func (h *Handler) Update(c *echo.Context) error {
	id, err := util.ParseUUIDParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, util.HttpError(err, "User Id is not a valid UUID"))
	}

	// parse input
	var inputDTO UpdateInputDto
	if err = c.Bind(&inputDTO); err != nil {
		return c.JSON(http.StatusBadRequest, util.HttpError(err, "Unable to parse input values"))
	}
	if err = inputDTO.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, util.HttpError(err, err.Error()))
	}

	loggedInUser := util.CurrentUser(c)
	// An admin removing their own admin flag could leave no admin able to
	// undo it - another admin has to do it.
	if id == loggedInUser.ID && inputDTO.IsAdmin != nil && !*inputDTO.IsAdmin {
		return c.JSON(http.StatusBadRequest, util.HttpErrorMessage("You can't remove your own admin access"))
	}

	err = h.service.Update(c.Request().Context(), id, inputDTO, loggedInUser.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, util.HttpErrorMessage("User not found"))
	} else if errors.Is(err, ErrLastAdmin) {
		return c.JSON(http.StatusConflict, util.HttpErrorMessage(ErrLastAdmin.Error()))
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, util.HttpError(err, "Unable to update user"))
	}

	return c.JSON(http.StatusOK, util.HttpSuccessStatus())
}

// Delete godoc
// @Summary      Delete a user
// @Tags         Users
// @Security     ApiKeyAuth
// @Produce      json
// @Description  Admin only. Soft delete - the row and its audit history are kept. You can't delete yourself or the last admin (409).
// @Param        id  path  string  true  "user id (UUID)"
// @Success      200  {object}  util.HttpResult
// @Failure      400,401,403,404,409,500  {object}  util.HttpResult
// @Router       /api/v1/users/{id} [delete]
func (h *Handler) Delete(c *echo.Context) error {
	id, err := util.ParseUUIDParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, util.HttpError(err, "User Id is not a valid UUID"))
	}

	loggedInUser := util.CurrentUser(c)
	// Deleted users are locked out by AuthMiddleware, so deleting yourself
	// is an irreversible self-lockout.
	if id == loggedInUser.ID {
		return c.JSON(http.StatusBadRequest, util.HttpErrorMessage("You can't delete your own account"))
	}

	err = h.service.Delete(c.Request().Context(), id, loggedInUser.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, util.HttpErrorMessage("User not found"))
	} else if errors.Is(err, ErrLastAdmin) {
		return c.JSON(http.StatusConflict, util.HttpErrorMessage(ErrLastAdmin.Error()))
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, util.HttpError(err, "Unable to delete user"))
	}

	return c.JSON(http.StatusOK, util.HttpSuccessStatus())
}
