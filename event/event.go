package event

import (
	"github.com/tradersclub/TCTemplateBack/app"
	"github.com/tradersclub/TCTemplateBack/event/item"
	"github.com/tradersclub/TCUtils/logger"
	"github.com/tradersclub/TCUtils/natstan"
)

// Options struct de opções para a criação de uma instancia das rotas
type Options struct {
	Apps          *app.Container
	QueueMessager natstan.QueueMessager
}

// Register handler instance
func Register(opts Options) {
	item.Register(opts.Apps, opts.QueueMessager)

	logger.Info("Registered EVENT")

}
