package health

import (
	"context"
	"time"

	"github.com/tradersclub/TCUtils/logger"

	"github.com/tradersclub/TCTemplateBack/model"
	"github.com/tradersclub/TCTemplateBack/store"
)

// App interface de health para implementação
type App interface {
	Ping(ctx context.Context) (*model.Health, error)
	Check(ctx context.Context) (*model.Health, error)
}

// NewApp cria uma nova instancia do serviço de health
func NewApp(stores *store.Container, version string, startedAt time.Time) App {
	return &appImpl{
		stores:    stores,
		version:   version,
		startedAt: startedAt,
	}
}

type appImpl struct {
	stores    *store.Container
	startedAt time.Time
	version   string
}

func (s *appImpl) Ping(ctx context.Context) (*model.Health, error) {
	result := <-s.stores.Health.Ping(ctx)
	if result.Error != nil {
		logger.ErrorContext(ctx, "app.health.ping", result.Error.Error())

		return nil, result.Error
	}

	data, err := model.ToHealth(result.Data)
	if err != nil {
		logger.ErrorContext(ctx, "app.health.ping", err.Error())

		return nil, err
	}
	data.ServerStatedAt = s.startedAt.UTC().String()
	data.Version = s.version

	return data, nil
}

func (s *appImpl) Check(ctx context.Context) (*model.Health, error) {
	result := <-s.stores.Health.Check(ctx)
	if result.Error != nil {
		logger.ErrorContext(ctx, "app.health.check", result.Error.Error())

		return nil, result.Error
	}

	data, err := model.ToHealth(result.Data)
	if err != nil {
		logger.ErrorContext(ctx, "app.health.check", err.Error())

		return nil, err
	}
	data.ServerStatedAt = s.startedAt.UTC().String()
	data.Version = s.version

	return data, nil
}
