package health

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tradersclub/TCTemplateBack/app"
	"github.com/tradersclub/TCTemplateBack/model"
)

// Register group health check
func Register(g *echo.Group, apps *app.Container) {
	h := &handler{
		apps: apps,
	}

	grp := g.Group("/health")
	grp.GET("", h.ping)
	grp.GET("/check", h.check)
}

type handler struct {
	apps *app.Container
}

// ping swagger document
// @Description Essa rota é privada com o token valido (Bearer)
// @Tags health
// @Accept  json
// @Produce  json
// @Param item body model.Item true "add Item"
// @Success 200 {object} model.Item
// @Failure 400 {object} string
// @Router /health [get]
func (h *handler) ping(c echo.Context) error {
	ctx := c.Request().Context()

	status, err := h.apps.Health.Ping(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, model.Response{
		Data: status,
	})
}

// check swagger document
// @Tags health
// @Accept  json
// @Produce  json
// @Param item body model.Item true "add Item"
// @Success 200 {object} model.Item
// @Failure 400 {object} string
// @Router /health/check [get]
func (h *handler) check(c echo.Context) error {
	ctx := c.Request().Context()

	status, err := h.apps.Health.Check(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, model.Response{
		Data: status,
	})
}
