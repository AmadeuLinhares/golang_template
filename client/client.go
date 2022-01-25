package client

import (
	"github.com/tradersclub/TCUtils/natstan"
)

type Client interface {
	Item() Item
}

func New(qm natstan.QueueMessager) Client {
	return &clientImpl{
		item: &itemImpl{qm},
	}
}

type clientImpl struct {
	item Item
}

func (c *clientImpl) Item() Item {
	return c.item
}
