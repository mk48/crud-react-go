package user

import (
	"net/http"

	"kfamily/internal/crud"
	"kfamily/internal/util"

	"github.com/labstack/echo/v5"
)

// Handler holds the user endpoints that aren't generic crud ones - an
// example of a custom endpoint (see crud.Register for the pattern).
type Handler struct {
	svc *crud.Service[row, Dto, crud.NoInput, UpdateInputDto]
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

	item, err := h.svc.GetOne(c.Request().Context(), currentUser.ID, false)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.HttpError(err, "Error while retrieving current user"))
	}

	return c.JSON(http.StatusOK, util.HttpData(item))
}
