package item

import (
	"net/http"

	"github.com/labstack/echo/v4"
	authModel "github.com/tradersclub/TCAuth/model"
	"github.com/tradersclub/TCTemplateBack/app"
	"github.com/tradersclub/TCTemplateBack/model"
	"github.com/tradersclub/TCUtils/logger"
)

// Register group item check
func Register(g *echo.Group, apps *app.Container) {
	h := &handler{
		apps: apps,
	}

	g.GET("", h.getItem)
	g.POST("", h.postItem)
}

type handler struct {
	apps *app.Container
}

// getItem swagger document
// @Summary Exemplo de como Buscar item por id
// @Description Essa rota é privada com o token valido (Bearer)
// @Tags item
// @Accept  json
// @Produce  json
// @Param id query string true "id"
// @Success 200 {object} model.Item
// @Failure 400 {object} string
// @Security ApiKeyAuth
// @Router /v1/item [get]
func (h *handler) getItem(c echo.Context) error {
	ctx := c.Request().Context()

	// recuperando query param
	id := c.QueryParam("id")

	// Se precisar usar session
	session := authModel.GetSession(c)
	if session == nil {
		return c.JSON(http.StatusForbidden, "Você precisa de autenticação para essa rota")
	}

	logger.Info(session.UserID)

	// Se precisar validar role(s)

	if h.apps.Session.GetClient().Is(&authModel.Roles{
		CurrentRoles: session.Roles,
		ValidRoles:   "system_admin",
	}) {
		logger.Info("Usuário é sys_admin")
	}

	resp, err := h.apps.Item.RequestItemById(ctx, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

// postItem swagger document
// @Summary Exemplo de como postar algum item
// @Description Essa rota é privada com o token valido (Bearer)
// @Tags item
// @Accept  json
// @Produce  json
// @Param item body model.Item true "add Item"
// @Success 200 {object} model.Item
// @Failure 400 {object} string
// @Security ApiKeyAuth
// @Router /v1/item [post]
func (h *handler) postItem(c echo.Context) error {

	ctx := c.Request().Context()

	item := new(model.Item)
	err := c.Bind(item)
	if err != nil {
		return err
	}

	resp, err := h.apps.Item.RequestItemById(ctx, item.Token)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}
