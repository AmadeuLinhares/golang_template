package item_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/tradersclub/TCUtils/do"
	"github.com/tradersclub/TCUtils/tcerr"

	"github.com/tradersclub/TCTemplateBack/store"
	"github.com/tradersclub/TCTemplateBack/test"

	"github.com/golang/mock/gomock"

	"github.com/google/go-cmp/cmp"
	"github.com/tradersclub/TCTemplateBack/app/health"
	"github.com/tradersclub/TCTemplateBack/mocks"
	"github.com/tradersclub/TCTemplateBack/model"
)

func TestPing(t *testing.T) {
	startedAt := time.Now()

	cases := map[string]struct {
		ExpectedErr  error
		ExpectedData *model.Health

		InputVersion  string
		InputDatetime time.Time

		PrepareMock func(mock *mocks.MockHealthStore)
	}{
		"deve retornar sucesso": {InputVersion: "1", InputDatetime: startedAt, ExpectedData: &model.Health{DatabaseStatus: "OK", Version: "1", ServerStatedAt: startedAt.UTC().String()}, PrepareMock: func(mock *mocks.MockHealthStore) {
			mock.EXPECT().Ping(gomock.Any()).Times(1).
				Return(do.Do(func(result *do.Result) {
					result.Data = &model.Health{DatabaseStatus: "OK"}
				}))
		}},
		"deve retornar erro com a mensagem: ocorreu um erro": {InputVersion: "1", InputDatetime: startedAt, ExpectedErr: tcerr.New(http.StatusInternalServerError, "ocorreu um erro", nil), PrepareMock: func(mock *mocks.MockHealthStore) {
			mock.EXPECT().Ping(gomock.Any()).Times(1).
				Return(do.Do(func(result *do.Result) {
					result.Error = tcerr.New(http.StatusInternalServerError, "ocorreu um erro", nil)
				}))
		}},
		"deve retornar erro na conversão da interface para *model.Health": {InputVersion: "1", InputDatetime: startedAt, ExpectedErr: tcerr.New(http.StatusInternalServerError, "não foi possível converter interface{} para *Health", nil), PrepareMock: func(mock *mocks.MockHealthStore) {
			mock.EXPECT().Ping(gomock.Any()).Times(1).
				Return(do.Do(func(result *do.Result) {}))
		}},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := test.NewController(t)
			mock := mocks.NewMockHealthStore(ctrl)

			cs.PrepareMock(mock)

			app := health.NewApp(&store.Container{Health: mock}, cs.InputVersion, cs.InputDatetime)

			data, err := app.Ping(ctx)

			if diff := cmp.Diff(data, cs.ExpectedData); diff != "" {
				t.Error(diff)
			}

			if diff := cmp.Diff(err, cs.ExpectedErr); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestCheck(t *testing.T) {
	startedAt := time.Now()

	cases := map[string]struct {
		ExpectedErr  error
		ExpectedData *model.Health

		InputVersion  string
		InputDatetime time.Time

		PrepareMock func(mock *mocks.MockHealthStore)
	}{
		"deve retornar sucesso": {InputVersion: "1", InputDatetime: startedAt, ExpectedData: &model.Health{DatabaseStatus: "OK", Version: "1", ServerStatedAt: startedAt.UTC().String()}, PrepareMock: func(mock *mocks.MockHealthStore) {
			mock.EXPECT().Check(gomock.Any()).Times(1).
				Return(do.Do(func(result *do.Result) {
					result.Data = &model.Health{DatabaseStatus: "OK"}
				}))
		}},
		"deve retornar erro com a mensagem: ocorreu um erro": {InputVersion: "1", InputDatetime: startedAt, ExpectedErr: tcerr.New(http.StatusInternalServerError, "ocorreu um erro", nil), PrepareMock: func(mock *mocks.MockHealthStore) {
			mock.EXPECT().Check(gomock.Any()).Times(1).
				Return(do.Do(func(result *do.Result) {
					result.Error = tcerr.New(http.StatusInternalServerError, "ocorreu um erro", nil)
				}))
		}},
		"deve retornar erro na conversão da interface para *model.Health": {InputVersion: "1", InputDatetime: startedAt, ExpectedErr: tcerr.New(http.StatusInternalServerError, "não foi possível converter interface{} para *Health", nil), PrepareMock: func(mock *mocks.MockHealthStore) {
			mock.EXPECT().Check(gomock.Any()).Times(1).
				Return(do.Do(func(result *do.Result) {}))
		}},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := test.NewController(t)
			mock := mocks.NewMockHealthStore(ctrl)

			cs.PrepareMock(mock)

			app := health.NewApp(&store.Container{Health: mock}, cs.InputVersion, cs.InputDatetime)

			data, err := app.Check(ctx)

			if diff := cmp.Diff(data, cs.ExpectedData); diff != "" {
				t.Error(diff)
			}

			if diff := cmp.Diff(err, cs.ExpectedErr); diff != "" {
				t.Error(diff)
			}
		})
	}
}
