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
// @Security ApiKeyAuth
// @Tags         AuditHistory
// @Summary      List audit history for a record
// @Description  every audit_history row recorded for id (a source table's primary key), most recent first
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Source record id (primary key)"
// @Success      200  {object}  util.HttpResult{Result=[]audithistory.Dto}
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
