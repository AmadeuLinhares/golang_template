package client

import (
	"context"

	"github.com/tradersclub/TCTemplateBack/model"
	"github.com/tradersclub/TCUtils/logger"
	"github.com/tradersclub/TCUtils/natstan"

	"github.com/tradersclub/TCTemplateBack/event/item"
)

type Item interface {
	GetItemByID(ctx context.Context, id string) (*model.Item, error)
}

type itemImpl struct {
	qm natstan.QueueMessager
}

func (i *itemImpl) GetItemByID(ctx context.Context, id string) (*model.Item, error) {
	var res item.ResponseGetItemById

	if err := i.qm.Request(ctx, item.TCTEMPLATEBACK_ITEMS_GETBYID, item.RequestGetItemById{ID: id}, &res); err != nil {
		logger.Error("Client.Item.GetItemByID: ", err.Error())
		return nil, err
	}

	if res.Err != nil {
		logger.Error("Client.Item.GetItemByID.Err: ", res.Err.Error())
		return nil, res.Err
	}

	return res.Data, nil
}
