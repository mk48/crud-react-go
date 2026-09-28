package user

import (
	"database/sql"
	"encoding/json"
	"errors"
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

// UserGetOne godoc
// @Security ApiKeyAuth
// @Tags         User
// @Summary      Get one user
// @Accept       json
// @Produce      json
// @Param        id   					path      string  true  	"User id"
// @Param        includeDeleted 	query      bool  	false  	"included deleted record or omit"
// @Success      200  {object}  util.HttpResult{Result=user.Dto}
// @Router       /api/users/{id} [get]
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

// UserMe godoc
// @Security ApiKeyAuth
// @Tags         User
// @Summary      Get the signed-in user
// @Description  Returns the caller's own user record (including isAdmin), so a frontend can tell whether the signed-in account has admin access
// @Produce      json
// @Success      200  {object}  util.HttpResult{Result=user.Dto}
// @Router       /api/users/me [get]
func (h *Handler) Me(c *echo.Context) error {
	currentUser := util.CurrentUser(c)

	row, err := h.service.GetOne(c.Request().Context(), currentUser.ID, false)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.HttpError(err, "Error while retrieving current user"))
	}

	return c.JSON(http.StatusOK, util.HttpData(row))
}

// UserList godoc
// @Security ApiKeyAuth
// @Tags         User
// @Summary      List User
// @Description  get User list with pagination
// @Accept       json
// @Produce      json
// @Param        results_per_page   			query      int  	true  	"Results/page"
// @Param        pageindex 						query      int  	true  	"page index (starts 0)"
// @Param        sort   						query      string  	false  	"sort column, ex: name:asc"
// @Param        search   						query      string  	false  	"search value"
// @Param        includeDeleted 		query      bool  	false  	"included deleted record or omit"
// @Success      200  {object}  util.HttpResult{Result=dto.PaginationResponse[user.Dto]}
// @Router       /api/users [get]
func (h *Handler) List(c *echo.Context) error {
	// ?results_per_page=10 & page=1 & sort=name:asc & search=kumaran & includeDeleted=false
	var filters = util.FiltersFromContext(c)

	result, err := h.service.List(c.Request().Context(), filters)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.HttpError(err, "Error in getting user list"))
	}

	return c.JSON(http.StatusOK, util.HttpData(result))
}

// UserMetadata godoc
// @Security ApiKeyAuth
// @Tags         User
// @Summary      Get user table metadata
// @Description  columns and their data types, so a frontend can discover valid sort/filter column names
// @Produce      json
// @Success      200  {object}  util.HttpResult{Result=[]dto.ColumnMeta}
// @Router       /api/users/meta [get]
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

// UserQuery godoc
// @Security ApiKeyAuth
// @Tags         User
// @Summary      Advanced query for User
// @Description  list User rows matching a dynamic WHERE clause built by the frontend's react-querybuilder (see AdvancedQueryBuilder.tsx): a parameterized SQL fragment referencing quoted column names and :pN placeholders, plus its parameter values as a JSON array in placeholder order
// @Accept       json
// @Produce      json
// @Param        results_per_page   			query      int  	true  	"Results/page"
// @Param        pageindex 						query      int  	true  	"page index (starts 0)"
// @Param        sort   						query      string  	false  	"sort column, ex: name:asc"
// @Param        whereCondition   				query      string  	true  	`parameterized SQL where clause, ex: ("name" = :p1)`
// @Param        whereConditionParametersJson  	query      string  	true  	"JSON array of parameter values, in :p1, :p2... order"
// @Success      200  {object}  util.HttpResult{Result=dto.PaginationResponse[user.Dto]}
// @Router       /api/users/query [get]
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

	var filters = util.FiltersFromContext(c)

	result, err := h.service.Query(c.Request().Context(), filters, input.WhereCondition, whereConditionParams)
	if err != nil {
		return c.JSON(http.StatusBadRequest, util.HttpError(err, "Error in running user query"))
	}

	return c.JSON(http.StatusOK, util.HttpData(result))
}

// UpdateUser godoc
// @Security ApiKeyAuth
// @Tags        User
// @Summary     Update user
// @Description Only name and isAdmin can be changed - sub/email are set once at first sign-in
// @Accept      json
// @Produce     json
// @Param       id   		path    string  			true 	"User Id"
// @Param 		request 	body	user.UpdateInputDto 	true 	"Update User"
// @Success     200  {object}  util.HttpResult{Result=user.Dto}
// @Router      /api/users/{id} [put]
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

// UserDelete godoc
// @Security ApiKeyAuth
// @Tags        User
// @Summary     Delete user
// @Accept      json
// @Produce     json
// @Param       id   		path    string  			true 	"User Id"
// @Router      /api/users/{id} [delete]
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
