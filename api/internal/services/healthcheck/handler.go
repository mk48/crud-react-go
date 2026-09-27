package healthcheck

import (
	"kfamily/internal/util"
	"net/http"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	envs *util.AppENV
}

func NewHandler(envs *util.AppENV) *Handler {
	return &Handler{envs: envs}
}

// ------------------------------------♙♙♙♙♙♙♙♙-----------------------------------------------
// HealthCheck godoc
// @Tags         health
// @Summary      Health check
// @Description  information about the application status, operating environment and version.
// @Accept       json
// @Produce      json
// @Success      200  {object}  HealthCheckResponse
// @Router       /healthcheck [get]
func (h *Handler) HealthCheckHandler(c *echo.Context) error {

	return c.JSON(http.StatusOK, map[string]string{
		"environment": h.envs.Env,
		"version":     h.envs.Version,
	})
}
