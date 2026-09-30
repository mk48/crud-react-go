package crud

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"kfamily/internal/dto"
	"kfamily/internal/util"

	"github.com/labstack/echo/v5"
)

// Handler is the HTTP side of the generic endpoints (routes: register.go).
type Handler[TRow, TDto, TCreate, TUpdate any] struct {
	svc *Service[TRow, TDto, TCreate, TUpdate]
}

// validator is implemented by input DTOs that trim/check themselves
// (e.g. sample.CreateInputDto).
type validator interface {
	Validate() error
}

type queryInput struct {
	WhereCondition               string `query:"whereCondition"`
	WhereConditionParametersJson string `query:"whereConditionParametersJson"`
}

func (h *Handler[TRow, TDto, TCreate, TUpdate]) GetOne(c *echo.Context) error {
	id, err := util.ParseUUIDParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, util.HttpError(err, h.svc.res.Label+" id is not a valid UUID"))
	}

	includeDeleted, _ := strconv.ParseBool(c.QueryParam("includeDeleted"))

	item, err := h.svc.GetOne(c.Request().Context(), id, includeDeleted)
	if err != nil {
		return h.fail(c, err, "Error while retrieving "+h.svc.name())
	}

	return c.JSON(http.StatusOK, util.HttpData(item))
}

func (h *Handler[TRow, TDto, TCreate, TUpdate]) List(c *echo.Context) error {
	var filters dto.Filters = util.FiltersFromContext(c)

	result, err := h.svc.List(c.Request().Context(), filters)
	if err != nil {
		return h.fail(c, err, "Error in getting "+h.svc.name()+" list")
	}

	return c.JSON(http.StatusOK, util.HttpData(result))
}

func (h *Handler[TRow, TDto, TCreate, TUpdate]) Metadata(c *echo.Context) error {
	columns, err := h.svc.Metadata(c.Request().Context())
	if err != nil {
		return h.fail(c, err, "Error while retrieving "+h.svc.name()+" metadata")
	}

	return c.JSON(http.StatusOK, util.HttpData(columns))
}

func (h *Handler[TRow, TDto, TCreate, TUpdate]) Query(c *echo.Context) error {
	var input queryInput
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

	result, err := h.svc.Query(c.Request().Context(), filters, input.WhereCondition, whereConditionParams)
	if err != nil {
		// Mostly a rejected WHERE fragment (see util.ValidateWhereCondition).
		return c.JSON(http.StatusBadRequest, util.HttpError(err, "Error in running "+h.svc.name()+" query"))
	}

	return c.JSON(http.StatusOK, util.HttpData(result))
}

func (h *Handler[TRow, TDto, TCreate, TUpdate]) Create(c *echo.Context) error {
	var input TCreate
	if responded, err := bindAndValidate(c, &input); responded {
		return err
	}

	id, err := h.svc.Create(c.Request().Context(), &input, util.CurrentUser(c))
	if err != nil {
		return h.fail(c, err, "Error on creating a new "+h.svc.name())
	}

	return c.JSON(http.StatusCreated, util.HttpData(id))
}

func (h *Handler[TRow, TDto, TCreate, TUpdate]) Update(c *echo.Context) error {
	id, err := util.ParseUUIDParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, util.HttpError(err, h.svc.res.Label+" id is not a valid UUID"))
	}

	var input TUpdate
	if responded, err := bindAndValidate(c, &input); responded {
		return err
	}

	if err := h.svc.Update(c.Request().Context(), id, &input, util.CurrentUser(c)); err != nil {
		return h.fail(c, err, "Unable to update "+h.svc.name())
	}

	return c.JSON(http.StatusOK, util.HttpSuccessStatus())
}

func (h *Handler[TRow, TDto, TCreate, TUpdate]) Delete(c *echo.Context) error {
	id, err := util.ParseUUIDParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, util.HttpError(err, h.svc.res.Label+" id is not a valid UUID"))
	}

	if err := h.svc.Delete(c.Request().Context(), id, util.CurrentUser(c)); err != nil {
		return h.fail(c, err, "Unable to delete "+h.svc.name())
	}

	return c.JSON(http.StatusOK, util.HttpSuccessStatus())
}

// bindAndValidate parses the JSON body into input and runs its Validate
// method, if any. responded is true when it has already answered 400 - the
// caller then just returns err (the result of writing that response).
func bindAndValidate[T any](c *echo.Context, input *T) (responded bool, err error) {
	if err := c.Bind(input); err != nil {
		return true, c.JSON(http.StatusBadRequest, util.HttpError(err, "Unable to parse input values"))
	}
	if v, ok := any(input).(validator); ok {
		if err := v.Validate(); err != nil {
			return true, c.JSON(http.StatusBadRequest, util.HttpError(err, err.Error()))
		}
	}
	return false, nil
}

// fail answers a service error: a hook's crud.Error with its status, no
// such row with 404, anything else with 500 and message.
func (h *Handler[TRow, TDto, TCreate, TUpdate]) fail(c *echo.Context, err error, message string) error {
	if he, ok := asHTTPError(err); ok {
		return c.JSON(he.Status, util.HttpError(he.Err, he.Err.Error()))
	}
	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, util.HttpErrorMessage(h.svc.res.Label+" not found"))
	}
	return c.JSON(http.StatusInternalServerError, util.HttpError(err, message))
}
