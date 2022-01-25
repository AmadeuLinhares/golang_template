package test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
)

// NewController inicia um novo controler para o mock
func NewController(t *testing.T) (*gomock.Controller, context.Context) {
	ctx := context.Background()
	return gomock.WithContext(ctx, t)
}
