package api

import (
	"github.com/labstack/echo/v4"
	health "github.com/tradersclub/TCTemplateBack/api/health"
	v1 "github.com/tradersclub/TCTemplateBack/api/v1"
	"github.com/tradersclub/TCTemplateBack/app"
	"github.com/tradersclub/TCUtils/logger"
)

// Options struct de opções para a criação de uma instancia das rotas
type Options struct {
	Group *echo.Group
	Apps  *app.Container
}

// Register api instance
func Register(opts Options) {
	v1.Register(opts.Group, opts.Apps)
	health.Register(opts.Group, opts.Apps)

	logger.Info("Registered API")
}
