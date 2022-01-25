package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/tradersclub/TCTemplateBack/api/v1/item"
	"github.com/tradersclub/TCTemplateBack/app"
)

// Register regristra as rotas v1
func Register(g *echo.Group, apps *app.Container) {
	v1 := g.Group("/v1", apps.Session.InjectSession)

	item.Register(v1.Group("/item"), apps)
}
