package item_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/tradersclub/TCTemplateBack/model"
	"github.com/tradersclub/TCTemplateBack/store/health"
	"github.com/tradersclub/TCTemplateBack/test"
	"github.com/tradersclub/TCUtils/tcerr"

	"github.com/google/go-cmp/cmp"
)

func TestPing(t *testing.T) {
	cases := map[string]struct {
		ExpectedErr  error
		ExpectedData interface{}

		PrepareMock func(mock sqlmock.Sqlmock)
	}{
		"deve retornar sucesso": {ExpectedData: &model.Health{DatabaseStatus: "OK"}, PrepareMock: func(mock sqlmock.Sqlmock) {
			mock.ExpectPing()
		}},
		"deve retornar erro com a mensagem: ocorreu um erro": {ExpectedErr: tcerr.New(http.StatusInternalServerError, "ocorreu um erro", nil), PrepareMock: func(mock sqlmock.Sqlmock) {
			mock.ExpectPing().WillReturnError(errors.New("ocorreu um erro"))
		}},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			db, mock := test.GetDB()
			cs.PrepareMock(mock)

			store := health.NewStore(db)
			ctx := context.Background()

			result := <-store.Ping(ctx)

			if diff := cmp.Diff(result.Data, cs.ExpectedData); diff != "" {
				t.Error(diff)
			}

			if diff := cmp.Diff(result.Error, cs.ExpectedErr); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestCheck(t *testing.T) {
	cases := map[string]struct {
		ExpectedErr  error
		ExpectedData interface{}

		PrepareMock func(mock sqlmock.Sqlmock)
	}{
		"deve retornar sucesso": {ExpectedData: &model.Health{DatabaseStatus: "DB OK"}, PrepareMock: func(mock sqlmock.Sqlmock) {
			rows := test.NewRows("database_status").AddRow("DB OK")
			mock.ExpectQuery("SELECT 'DB OK' AS database_status").WillReturnRows(rows)
		}},
		"deve retornar erro com a mensagem: ocorreu um erro": {ExpectedErr: tcerr.New(http.StatusInternalServerError, "ocorreu um erro", nil), PrepareMock: func(mock sqlmock.Sqlmock) {
			mock.ExpectQuery("SELECT 'DB OK' AS database_status").WillReturnError(errors.New("ocorreu um erro"))
		}},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			db, mock := test.GetDB()
			cs.PrepareMock(mock)

			store := health.NewStore(db)
			ctx := context.Background()

			result := <-store.Check(ctx)

			if diff := cmp.Diff(result.Data, cs.ExpectedData); diff != "" {
				t.Error(diff)
			}

			if diff := cmp.Diff(result.Error, cs.ExpectedErr); diff != "" {
				t.Error(diff)
			}
		})
	}
}
