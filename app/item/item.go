package item

import (
	"context"
	"net/http"
	"time"

	"github.com/tradersclub/TCTemplateBack/model"
	"github.com/tradersclub/TCTemplateBack/store"
	"github.com/tradersclub/TCUtils/cache"
	"github.com/tradersclub/TCUtils/logger"
	"github.com/tradersclub/TCUtils/natstan"
	"github.com/tradersclub/TCUtils/tcerr"
)

const TCSTARTKIT_GET_ITEM_BY_ID = "tcstartkit_get_item_by_id"

type getItemById struct {
	Err  error
	Id   string
	Data *model.Item
}

// App interface de item para implementação
type App interface {
	RequestItemById(ctx context.Context, id string) (*model.Item, error)
	GetItemById(ctx context.Context, id string) (*model.Item, error)
}

// NewApp cria uma nova instancia do serviço de exemplo item
func NewApp(stores *store.Container, qm natstan.QueueMessager, cache cache.Cache) App {
	return &appImpl{
		stores: stores,
		qm:     qm,
		cache:  cache,
	}
}

type appImpl struct {
	stores *store.Container
	qm     natstan.QueueMessager
	cache  cache.Cache
}

// RequestItemById - exemplo de request nats
func (s *appImpl) RequestItemById(ctx context.Context, id string) (*model.Item, error) {

	data := new(getItemById)
	data.Id = id

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	var res getItemById
	err := s.qm.Request(ctx, TCSTARTKIT_GET_ITEM_BY_ID, data, &res)
	if err != nil {
		return nil, tcerr.NewError(http.StatusInternalServerError, "erro ao resgatar o item", nil)
	}

	if res.Err != nil {
		return nil, tcerr.NewError(http.StatusInternalServerError, "erro ao tentar recuperar o item pelo id", nil)
	}

	return res.Data, nil
}

// GetItemById - exemplo de recuperação de dados do cache e store
func (s *appImpl) GetItemById(ctx context.Context, id string) (*model.Item, error) {
	item := new(model.Item)

	// exemplo de consulta em cache caso seja necessário, não esquece de validar junto ao seu líder qual memcached usar
	if err := s.cache.Get(ctx, id, item); err != nil {
		logger.ErrorContext(ctx, "app.item.get_item_by_id", "não encontrei o cache com id: "+id, err.Error())
	}

	// exemplo de consulta em store
	result := <-s.stores.Item.GetItemById(ctx, id)
	if result.Error != nil {
		logger.ErrorContext(ctx, "app.item.get_item_by_id", result.Error.Error())
		return nil, result.Error
	}

	data, err := model.ToItem(result.Data)
	if err != nil {
		logger.ErrorContext(ctx, "app.item.get_item_by_id", err.Error())

		return nil, err
	}

	return data, nil
}
