package samplechild

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

func (h *Handler) GetOne(c *echo.Context) error {
	id, err := util.ParseUUIDParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, util.HttpError(err, "Sample child item id is not a valid UUID"))
	}

	includeDeleted, _ := strconv.ParseBool(c.QueryParam("includeDeleted"))

	row, err := h.service.GetOne(c.Request().Context(), id, includeDeleted)
	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, util.HttpErrorMessage("Sample child item not found"))
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, util.HttpError(err, "Error while retrieving sample child item"))
	}

	return c.JSON(http.StatusOK, util.HttpData(row))
}

func (h *Handler) List(c *echo.Context) error {
	var filters = util.FiltersFromContext(c)

	result, err := h.service.List(c.Request().Context(), filters)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.HttpError(err, "Error in getting sample child item list"))
	}

	return c.JSON(http.StatusOK, util.HttpData(result))
}

func (h *Handler) Metadata(c *echo.Context) error {
	columns, err := h.service.Metadata(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.HttpError(err, "Error while retrieving sample child item metadata"))
	}

	return c.JSON(http.StatusOK, util.HttpData(columns))
}

type QueryInputDto struct {
	WhereCondition               string `query:"whereCondition"`
	WhereConditionParametersJson string `query:"whereConditionParametersJson"`
}

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
		return c.JSON(http.StatusBadRequest, util.HttpError(err, "Error in running sample child item query"))
	}

	return c.JSON(http.StatusOK, util.HttpData(result))
}

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
	if errors.Is(err, ErrInvalidSampleItem) {
		return c.JSON(http.StatusBadRequest, util.HttpError(err, err.Error()))
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, util.HttpError(err, "Error on creating a new sample child item"))
	}

	return c.JSON(http.StatusCreated, util.HttpData(newlyCreated))
}

func (h *Handler) Update(c *echo.Context) error {
	id, err := util.ParseUUIDParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, util.HttpError(err, "Sample child item id is not a valid UUID"))
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
		return c.JSON(http.StatusNotFound, util.HttpErrorMessage("Sample child item not found"))
	} else if errors.Is(err, ErrInvalidSampleItem) {
		return c.JSON(http.StatusBadRequest, util.HttpError(err, err.Error()))
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, util.HttpError(err, "Unable to update sample child item"))
	}

	return c.JSON(http.StatusOK, util.HttpSuccessStatus())
}

func (h *Handler) Delete(c *echo.Context) error {
	id, err := util.ParseUUIDParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, util.HttpError(err, "Sample child item id is not a valid UUID"))
	}

	loggedInUser := util.CurrentUser(c)
	err = h.service.Delete(c.Request().Context(), id, loggedInUser.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, util.HttpErrorMessage("Sample child item not found"))
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, util.HttpError(err, "Unable to delete sample child item"))
	}

	return c.JSON(http.StatusOK, util.HttpSuccessStatus())
}
