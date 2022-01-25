package item

import (
	"github.com/tradersclub/TCTemplateBack/model"
	"github.com/tradersclub/TCUtils/tcerr"
)

type RequestGetItemById struct {
	ID string `json:"id"`
}

type ResponseGetItemById struct {
	Err  *tcerr.Error `json:"err"`
	Data *model.Item  `json:"data"`
}
