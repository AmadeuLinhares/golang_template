package item

import (
	"context"
	"net/http"

	"github.com/tradersclub/TCTemplateBack/app"
	"github.com/tradersclub/TCUtils/logger"
	"github.com/tradersclub/TCUtils/natstan"
	"github.com/tradersclub/TCUtils/tcerr"
)

const TCTEMPLATEBACK_ITEMS_GETBYID = "tctemplateback.items.getbyid"

// Register group health check
func Register(apps *app.Container, qm natstan.QueueMessager) {
	e := &event{
		apps: apps,
		qm:   qm,
	}

	e.qm.Subscribe(TCTEMPLATEBACK_ITEMS_GETBYID, e.getItemById)
}

type event struct {
	apps *app.Container
	qm   natstan.QueueMessager
}

func (e *event) getItemById(msg *natstan.MsgHandler) {
	ctx := context.Background()

	//recupera o objeto através do metódo criado no TCNatsModel
	req, res := RequestGetItemById{}, ResponseGetItemById{}
	err := msg.Decode(&req)
	if err != nil {
		logger.ErrorContext(ctx, "event.session.get_session_by_id", err.Error())

		res.Err = tcerr.NewError(500, "event.session.get_session_by_id", err.Error())

		if err := msg.Respond(res); err != nil {
			logger.ErrorContext(ctx, "event.session.get_session_by_id.respond", err.Error())
		}

		return
	}

	// pega o item
	session, err := e.apps.Item.GetItemById(ctx, req.ID)
	if err != nil {
		logger.ErrorContext(ctx, "event.session.get_session_by_id", err.Error())
		res.Err = tcerr.NewError(http.StatusInternalServerError, "event.session.get_session_by_id", err.Error())
	} else {
		res.Data = session
	}

	if err := msg.Respond(res); err != nil {
		logger.ErrorContext(ctx, "event.session.get_session_by_id.respond", err.Error())
	}
}
