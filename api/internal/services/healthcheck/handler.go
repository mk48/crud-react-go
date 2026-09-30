package healthcheck

import (
	"context"
	"net/http"
	"time"

	"kfamily/internal/util"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
)

type Handler struct {
	envs *util.AppENV
	db   *sqlx.DB
}

func NewHandler(envs *util.AppENV, db *sqlx.DB) *Handler {
	return &Handler{envs: envs, db: db}
}

// HealthCheckHandler godoc
// @Tags         Health
// @Summary      Liveness check
// @Description  The process is up and serving HTTP. Doesn't touch the database - see /healthcheck/ready.
// @Produce      json
// @Success      200  {object}  HealthCheckResponse
// @Router       /healthcheck [get]
func (h *Handler) HealthCheckHandler(c *echo.Context) error {
	return c.JSON(http.StatusOK, HealthCheckResponse{Status: "ok", Version: h.envs.Version})
}

// ReadinessHandler godoc
// @Tags         Health
// @Summary      Readiness check
// @Description  Ready to serve requests: the database answers a ping within 2 seconds.
// @Produce      json
// @Success      200  {object}  HealthCheckResponse
// @Failure      503  {object}  HealthCheckResponse
// @Router       /healthcheck/ready [get]
func (h *Handler) ReadinessHandler(c *echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 2*time.Second)
	defer cancel()

	if err := h.db.PingContext(ctx); err != nil {
		c.Logger().Error("readiness check: database ping failed", "err", err)
		return c.JSON(http.StatusServiceUnavailable, HealthCheckResponse{Status: "database unavailable", Version: h.envs.Version})
	}

	return c.JSON(http.StatusOK, HealthCheckResponse{Status: "ok", Version: h.envs.Version})
}
