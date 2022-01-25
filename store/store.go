package store

import (
	"github.com/jmoiron/sqlx"
	"github.com/tradersclub/TCTemplateBack/store/health"
	"github.com/tradersclub/TCTemplateBack/store/item"
	"github.com/tradersclub/TCUtils/logger"
	"github.com/tradersclub/TCUtils/mongodb/v2"
)

// Container modelo para exportação dos repositórios instanciados
type Container struct {
	Health health.Store
	Item   item.Store
}

// Options struct de opções para a criação de uma instancia dos repositórios
type Options struct {
	WriterSQL   *sqlx.DB
	ReaderSQL   *sqlx.DB
	WriterNoSQL mongodb.Database
	ReaderNoSQL mongodb.Database
}

// New cria uma nova instancia dos repositórios
func New(opts Options) *Container {
	container := &Container{
		Health: health.NewStore(opts.ReaderSQL),
		Item:   item.NewStore(opts.WriterSQL, opts.ReaderSQL),
	}

	logger.Info("Registered STORE")

	return container
}
