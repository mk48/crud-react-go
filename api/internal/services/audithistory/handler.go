package audithistory

import (
	"kfamily/internal/util"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

//------------------------------------------------------------------------------

// AuditHistoryList godoc
// @Security     ApiKeyAuth
// @Tags         Audit history
// @Summary      List audit history for a record
// @Description  Every change recorded for a record (any table's primary key), newest first. Each entry holds the full row right after the change, and the operation that caused it.
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Source record id (primary key)"
// @Success      200  {object}  util.HttpResult{result=[]audithistory.Dto}
// @Failure      400,401,500  {object}  util.HttpResult
// @Router       /api/v1/audit-history/{id} [get]
func (h *Handler) List(c *echo.Context) error {
	sourceId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, util.HttpError(err, "Id is not a valid UUID: "+c.Param("id")))
	}

	items, err := h.service.List(c.Request().Context(), sourceId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.HttpError(err, "Error while retrieving audit history"))
	}

	return c.JSON(http.StatusOK, util.HttpData(items))
}
