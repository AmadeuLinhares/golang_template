package item

import (
	"context"
	"net/http"

	"github.com/tradersclub/TCUtils/do"

	"github.com/tradersclub/TCUtils/logger"
	"github.com/tradersclub/TCUtils/tcerr"

	"github.com/jmoiron/sqlx"
	"github.com/tradersclub/TCTemplateBack/model"
)

// Store interface para implementação do item
type Store interface {
	GetItemById(ctx context.Context, id string) do.ChanResult
}

// NewStore cria uma nova instancia do repositorio de item
func NewStore(writer, reader *sqlx.DB) Store {
	return &storeImpl{writer, reader}
}

type storeImpl struct {
	writer *sqlx.DB
	reader *sqlx.DB
}

// GetItemById - pega o item no banco
func (r *storeImpl) GetItemById(ctx context.Context, id string) do.ChanResult {
	return do.Do(func(res *do.Result) {
		data := new(model.Item)
		sql := `SELECT 
					Id,
					Token,
					CreateAt,
					ExpiresAt,
					LastActivityAt,
					UserId,
					DeviceId,
					Roles,
					IsOAuth
				FROM 
					Sessions 
				WHERE 
					Token = ?`
		err := r.reader.GetContext(ctx, data, sql, id)
		if err != nil {
			// Muito importante não retornar o erro do banco, apenas logar e retornar erro personalizado
			logger.ErrorContext(ctx, "store.item.check", err.Error())
			res.Error = tcerr.New(http.StatusInternalServerError, "erro ao buscar item no store", nil)
			return
		}

		res.Data = data
	})
}
